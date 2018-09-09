package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type KeycloakClient struct {
	baseUrl           string
	clientCredentials *ClientCredentials
	httpClient        *http.Client
}

type ClientCredentials struct {
	ClientId     string
	ClientSecret string
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

const (
	apiUrl   = "/auth/admin"
	tokenUrl = "/auth/realms/master/protocol/openid-connect/token"
)

func NewKeycloakClient(baseUrl, clientId, clientSecret string) (*KeycloakClient, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	keycloakClient := KeycloakClient{
		baseUrl: baseUrl,
		clientCredentials: &ClientCredentials{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		},
		httpClient: httpClient,
	}

	err := keycloakClient.login()
	if err != nil {
		return nil, err
	}

	return &keycloakClient, nil
}

func (keycloakClient *KeycloakClient) login() error {
	accessTokenUrl := keycloakClient.baseUrl + tokenUrl

	accessTokenData := url.Values{}

	accessTokenData.Set("client_id", keycloakClient.clientCredentials.ClientId)
	accessTokenData.Set("client_secret", keycloakClient.clientCredentials.ClientSecret)
	accessTokenData.Set("grant_type", "client_credentials")

	log.Printf("[DEBUG] Login request: %s", accessTokenData.Encode())

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, accessTokenUrl, strings.NewReader(accessTokenData.Encode()))
	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	accessTokenResponse, err := keycloakClient.httpClient.Do(accessTokenRequest)
	if err != nil {
		return err
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

	return nil
}

func (keycloakClient *KeycloakClient) refresh() error {
	refreshTokenUrl := keycloakClient.baseUrl + tokenUrl

	refreshTokenData := url.Values{}

	refreshTokenData.Set("grant_type", "refresh_token")
	refreshTokenData.Set("client_id", keycloakClient.clientCredentials.ClientId)
	refreshTokenData.Set("client_secret", keycloakClient.clientCredentials.ClientSecret)
	refreshTokenData.Set("refresh_token", keycloakClient.clientCredentials.RefreshToken)

	log.Printf("[DEBUG] Refresh request: %s", refreshTokenData.Encode())

	accessTokenRequest, _ := http.NewRequest(http.MethodPost, refreshTokenUrl, strings.NewReader(refreshTokenData.Encode()))
	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	refreshTokenResponse, err := keycloakClient.httpClient.Do(accessTokenRequest)
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

func (keycloakClient *KeycloakClient) addRequestHeaders(request *http.Request) {
	tokenType := keycloakClient.clientCredentials.TokenType
	accessToken := keycloakClient.clientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Accept", "application/json")

	if request.Method == http.MethodPost || request.Method == http.MethodPut {
		request.Header.Set("Content-type", "application/json")
	}
}

/**
Sends an HTTP request and refreshes credentials on 403 or 401 errors
*/
func (keycloakClient *KeycloakClient) sendRequest(request *http.Request) ([]byte, string, error) {
	requestMethod := request.Method
	requestPath := request.URL.Path

	log.Printf("[DEBUG] Sending %s to %s", requestMethod, requestPath)
	if request.Body != nil {
		requestBody, err := request.GetBody()
		if err != nil {
			return nil, "", err
		}

		requestBodyBuffer := new(bytes.Buffer)
		requestBodyBuffer.ReadFrom(requestBody)

		log.Printf("[DEBUG] Request body: %s", requestBodyBuffer.String())
	}

	keycloakClient.addRequestHeaders(request)

	response, err := keycloakClient.httpClient.Do(request)
	if err != nil {
		return nil, "", err
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

		response, err = keycloakClient.httpClient.Do(request)
		if err != nil {
			return nil, "", err
		}
	}

	log.Printf("[DEBUG] Response: %s", response.Status)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	if len(body) != 0 {
		log.Printf("[DEBUG] Response body: %s", body)
	}

	if response.StatusCode >= 400 {
		return nil, "", fmt.Errorf("error sending %s request to %s: %s", request.Method, request.URL.Path, response.Status)
	}

	return body, response.Header.Get("Location"), nil
}

func (keycloakClient *KeycloakClient) get(path string, resource interface{}) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	request, err := http.NewRequest(http.MethodGet, resourceUrl, nil)
	if err != nil {
		return err
	}

	body, _, err := keycloakClient.sendRequest(request)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, resource)
}

func (keycloakClient *KeycloakClient) post(path string, requestBody interface{}) (string, error) {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, resourceUrl, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	_, location, err := keycloakClient.sendRequest(request)

	return location, err
}

func (keycloakClient *KeycloakClient) put(path string, requestBody interface{}) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, resourceUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	_, _, err = keycloakClient.sendRequest(request)

	return err
}

func (keycloakClient *KeycloakClient) delete(path string) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	request, err := http.NewRequest(http.MethodDelete, resourceUrl, nil)
	if err != nil {
		return err
	}

	_, _, err = keycloakClient.sendRequest(request)

	return err
}
