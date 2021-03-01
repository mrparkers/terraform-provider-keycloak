package keycloak

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	"golang.org/x/net/publicsuffix"
)

type KeycloakClient struct {
	baseUrl           string
	realm             string
	clientCredentials *ClientCredentials
	httpClient        *http.Client
	initialLogin      bool
	userAgent         string
	version           *version.Version
}

type ClientCredentials struct {
	ClientId     string
	ClientSecret string
	Username     string
	Password     string
	GrantType    string
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

const (
	apiUrl   = "/admin"
	tokenUrl = "%s/realms/%s/protocol/openid-connect/token"
)

func NewKeycloakClient(url, basePath, clientId, clientSecret, realm, username, password string, initialLogin bool, clientTimeout int, caCert string, tlsInsecureSkipVerify bool, userAgent string) (*KeycloakClient, error) {
	cookieJar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
		Proxy:           http.ProxyFromEnvironment,
	}

	httpClient := &http.Client{
		Timeout:   time.Second * time.Duration(clientTimeout),
		Transport: transport,
		Jar:       cookieJar,
	}

	if caCert != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}
	}
	clientCredentials := &ClientCredentials{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
	if password != "" && username != "" {
		clientCredentials.Username = username
		clientCredentials.Password = password
		clientCredentials.GrantType = "password"
	} else if clientSecret != "" {
		clientCredentials.GrantType = "client_credentials"
	} else {
		return nil, fmt.Errorf("must specify client id, username and password for password grant, or client id and secret for client credentials grant")
	}

	keycloakClient := KeycloakClient{
		baseUrl:           url + basePath,
		clientCredentials: clientCredentials,
		httpClient:        httpClient,
		initialLogin:      initialLogin,
		realm:             realm,
		userAgent:         userAgent,
	}

	if keycloakClient.initialLogin {
		err := keycloakClient.login()
		if err != nil {
			return nil, err
		}
	}

	return &keycloakClient, nil
}

func (keycloakClient *KeycloakClient) login() error {
	accessTokenUrl := fmt.Sprintf(tokenUrl, keycloakClient.baseUrl, keycloakClient.realm)
	accessTokenData := keycloakClient.getAuthenticationFormData()

	log.Printf("[DEBUG] Login request: %s", accessTokenData.Encode())

	accessTokenRequest, err := http.NewRequest(http.MethodPost, accessTokenUrl, strings.NewReader(accessTokenData.Encode()))
	if err != nil {
		return err
	}

	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if keycloakClient.userAgent != "" {
		accessTokenRequest.Header.Set("User-Agent", keycloakClient.userAgent)
	}

	accessTokenResponse, err := keycloakClient.httpClient.Do(accessTokenRequest)
	if err != nil {
		return err
	}
	if accessTokenResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending POST request to %s: %s", accessTokenUrl, accessTokenResponse.Status)
	}

	defer accessTokenResponse.Body.Close()

	body, _ := ioutil.ReadAll(accessTokenResponse.Body)

	log.Printf("[DEBUG] Login response: %s", body)

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)
	if err != nil {
		return err
	}

	keycloakClient.clientCredentials.AccessToken = clientCredentials.AccessToken
	keycloakClient.clientCredentials.RefreshToken = clientCredentials.RefreshToken
	keycloakClient.clientCredentials.TokenType = clientCredentials.TokenType

	info, err := keycloakClient.GetServerInfo()
	if err != nil {
		return err
	}

	server_version := info.SystemInfo.ServerVersion
	if strings.Contains(server_version, ".GA") {
		server_version = strings.ReplaceAll(info.SystemInfo.ServerVersion, ".GA", "")
	}

	v, err := version.NewVersion(server_version)
	if err != nil {
		return err
	}

	keycloakClient.version = v

	return nil
}

func (keycloakClient *KeycloakClient) refresh() error {
	refreshTokenUrl := fmt.Sprintf(tokenUrl, keycloakClient.baseUrl, keycloakClient.realm)
	refreshTokenData := keycloakClient.getAuthenticationFormData()

	log.Printf("[DEBUG] Refresh request: %s", refreshTokenData.Encode())

	refreshTokenRequest, err := http.NewRequest(http.MethodPost, refreshTokenUrl, strings.NewReader(refreshTokenData.Encode()))
	if err != nil {
		return err
	}

	refreshTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if keycloakClient.userAgent != "" {
		refreshTokenRequest.Header.Set("User-Agent", keycloakClient.userAgent)
	}

	refreshTokenResponse, err := keycloakClient.httpClient.Do(refreshTokenRequest)
	if err != nil {
		return err
	}

	defer refreshTokenResponse.Body.Close()

	body, _ := ioutil.ReadAll(refreshTokenResponse.Body)

	log.Printf("[DEBUG] Refresh response: %s", body)

	// Handle 401 "User or client no longer has role permissions for client key" until I better understand why that happens in the first place
	if refreshTokenResponse.StatusCode == http.StatusBadRequest {
		log.Printf("[DEBUG] Unexpected 400, attemting to log in again")

		return keycloakClient.login()
	}

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)
	if err != nil {
		return err
	}

	keycloakClient.clientCredentials.AccessToken = clientCredentials.AccessToken
	keycloakClient.clientCredentials.RefreshToken = clientCredentials.RefreshToken
	keycloakClient.clientCredentials.TokenType = clientCredentials.TokenType

	return nil
}

func (keycloakClient *KeycloakClient) getAuthenticationFormData() url.Values {
	authenticationFormData := url.Values{}
	authenticationFormData.Set("client_id", keycloakClient.clientCredentials.ClientId)
	authenticationFormData.Set("grant_type", keycloakClient.clientCredentials.GrantType)

	if keycloakClient.clientCredentials.GrantType == "password" {
		authenticationFormData.Set("username", keycloakClient.clientCredentials.Username)
		authenticationFormData.Set("password", keycloakClient.clientCredentials.Password)

		if keycloakClient.clientCredentials.ClientSecret != "" {
			authenticationFormData.Set("client_secret", keycloakClient.clientCredentials.ClientSecret)
		}

	} else if keycloakClient.clientCredentials.GrantType == "client_credentials" {
		authenticationFormData.Set("client_secret", keycloakClient.clientCredentials.ClientSecret)
	}

	return authenticationFormData
}

func (keycloakClient *KeycloakClient) addRequestHeaders(request *http.Request) {
	tokenType := keycloakClient.clientCredentials.TokenType
	accessToken := keycloakClient.clientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Accept", "application/json")

	if keycloakClient.userAgent != "" {
		request.Header.Set("User-Agent", keycloakClient.userAgent)
	}

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodDelete {
		request.Header.Set("Content-type", "application/json")
	}
}

/**
Sends an HTTP request and refreshes credentials on 403 or 401 errors
*/
func (keycloakClient *KeycloakClient) sendRequest(request *http.Request, body []byte) ([]byte, string, error) {
	if !keycloakClient.initialLogin {
		keycloakClient.initialLogin = true
		err := keycloakClient.login()
		if err != nil {
			return nil, "", fmt.Errorf("error logging in: %s", err)
		}
	}

	requestMethod := request.Method
	requestPath := request.URL.Path

	log.Printf("[DEBUG] Sending %s to %s", requestMethod, requestPath)
	if body != nil {
		request.Body = ioutil.NopCloser(bytes.NewReader(body))
		log.Printf("[DEBUG] Request body: %s", string(body))
	}

	keycloakClient.addRequestHeaders(request)

	response, err := keycloakClient.httpClient.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("error sending request: %v", err)
	}

	// Unauthorized: Token could have expired
	// Forbidden: After creating a realm, following GETs for the realm return 403 until you refresh
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		log.Printf("[DEBUG] Response: %s.  Attempting refresh", response.Status)

		err := keycloakClient.refresh()
		if err != nil {
			return nil, "", fmt.Errorf("error refreshing credentials: %s", err)
		}

		keycloakClient.addRequestHeaders(request)

		if body != nil {
			request.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		response, err = keycloakClient.httpClient.Do(request)
		if err != nil {
			return nil, "", fmt.Errorf("error sending request after refresh: %v", err)
		}
	}

	log.Printf("[DEBUG] Response: %s", response.Status)

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if len(responseBody) != 0 && request.URL.Path != "/auth/admin/serverinfo" {
		log.Printf("[DEBUG] Response body: %s", responseBody)
	}

	if response.StatusCode >= 400 {
		errorMessage := fmt.Sprintf("error sending %s request to %s: %s.", request.Method, request.URL.Path, response.Status)

		if len(responseBody) != 0 {
			errorMessage = fmt.Sprintf("%s Response body: %s", errorMessage, responseBody)
		}

		return nil, "", &ApiError{
			Code:    response.StatusCode,
			Message: errorMessage,
		}
	}

	return responseBody, response.Header.Get("Location"), nil
}

func (keycloakClient *KeycloakClient) get(path string, resource interface{}, params map[string]string) error {
	body, err := keycloakClient.getRaw(path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resource)
}

func (keycloakClient *KeycloakClient) getRaw(path string, params map[string]string) ([]byte, error) {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	request, err := http.NewRequest(http.MethodGet, resourceUrl, nil)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	body, _, err := keycloakClient.sendRequest(request, nil)
	return body, err
}

func (keycloakClient *KeycloakClient) post(path string, requestBody interface{}) ([]byte, string, error) {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, "", err
	}

	request, err := http.NewRequest(http.MethodPost, resourceUrl, nil)
	if err != nil {
		return nil, "", err
	}

	body, location, err := keycloakClient.sendRequest(request, payload)

	return body, location, err
}

func (keycloakClient *KeycloakClient) put(path string, requestBody interface{}) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, resourceUrl, nil)
	if err != nil {
		return err
	}

	_, _, err = keycloakClient.sendRequest(request, payload)

	return err
}

func (keycloakClient *KeycloakClient) delete(path string, requestBody interface{}) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	var (
		payload []byte
		err     error
	)

	if requestBody != nil {
		payload, err = json.Marshal(requestBody)
		if err != nil {
			return err
		}
	}

	request, err := http.NewRequest(http.MethodDelete, resourceUrl, nil)
	if err != nil {
		return err
	}

	_, _, err = keycloakClient.sendRequest(request, payload)

	return err
}
