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

	accessTokenUrl := baseUrl + tokenUrl

	accessTokenData := url.Values{}

	accessTokenData.Set("client_id", clientId)
	accessTokenData.Set("client_secret", clientSecret)
	accessTokenData.Set("grant_type", "client_credentials")

	log.Printf("[DEBUG] Login request: %s", accessTokenData.Encode())

	accessTokenRequest, _ := http.NewRequest("POST", accessTokenUrl, strings.NewReader(accessTokenData.Encode()))
	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	accessTokenResponse, err := httpClient.Do(accessTokenRequest)
	if err != nil {
		return nil, err
	}

	defer accessTokenResponse.Body.Close()

	body, _ := ioutil.ReadAll(accessTokenResponse.Body)

	log.Printf("[DEBUG] Login response: %s", body)

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)

	if err != nil {
		return nil, err
	}

	clientCredentials.ClientId = clientId
	clientCredentials.ClientSecret = clientSecret

	return &KeycloakClient{
		baseUrl:           baseUrl,
		clientCredentials: &clientCredentials,
		httpClient:        httpClient,
	}, nil
}

func (keycloakClient *KeycloakClient) refresh() error {
	refreshTokenUrl := keycloakClient.baseUrl + tokenUrl

	refreshTokenData := url.Values{}

	refreshTokenData.Set("grant_type", "refresh_token")
	refreshTokenData.Set("client_id", keycloakClient.clientCredentials.ClientId)
	refreshTokenData.Set("client_secret", keycloakClient.clientCredentials.ClientSecret)
	refreshTokenData.Set("refresh_token", keycloakClient.clientCredentials.RefreshToken)

	log.Printf("[DEBUG] Refresh request: %s", refreshTokenData.Encode())

	accessTokenRequest, _ := http.NewRequest("POST", refreshTokenUrl, strings.NewReader(refreshTokenData.Encode()))
	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	refreshTokenResponse, err := keycloakClient.httpClient.Do(accessTokenRequest)
	if err != nil {
		return err
	}

	defer refreshTokenResponse.Body.Close()

	body, _ := ioutil.ReadAll(refreshTokenResponse.Body)

	log.Printf("[DEBUG] Refresh response: %s", body)

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

func (keycloakClient *KeycloakClient) get(path string, resource interface{}) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	log.Printf("[DEBUG] Sending GET to %s", resourceUrl)

	request, err := http.NewRequest("GET", resourceUrl, nil)
	if err != nil {
		return err
	}

	tokenType := keycloakClient.clientCredentials.TokenType
	accessToken := keycloakClient.clientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Accept", "application/json")

	response, err := keycloakClient.httpClient.Do(request)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Response: %s", response.Status)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(body) != 0 {
		log.Printf("[DEBUG] GET response body: %s", body)
	}

	return json.Unmarshal(body, resource)
}

func (keycloakClient *KeycloakClient) post(path string, requestBody interface{}) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Sending POST to %s", resourceUrl)
	log.Printf("[DEBUG] Request body: %s", payload)

	request, err := http.NewRequest("POST", resourceUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	tokenType := keycloakClient.clientCredentials.TokenType
	accessToken := keycloakClient.clientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := keycloakClient.httpClient.Do(request)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Response: %s", response.Status)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(body) != 0 {
		log.Printf("[DEBUG] POST response body: %s", body)
	}

	return keycloakClient.refresh()
}

func (keycloakClient *KeycloakClient) delete(path string) error {
	resourceUrl := keycloakClient.baseUrl + apiUrl + path

	request, err := http.NewRequest("DELETE", resourceUrl, nil)
	if err != nil {
		return err
	}

	tokenType := keycloakClient.clientCredentials.TokenType
	accessToken := keycloakClient.clientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Accept", "application/json")

	_, err = keycloakClient.httpClient.Do(request)

	return err
}
