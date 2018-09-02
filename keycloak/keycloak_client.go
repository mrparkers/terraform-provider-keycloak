package keycloak

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	AccessToken string
	TokenType   string
}

func NewKeycloakClient(baseUrl, clientId, clientSecret string) (*KeycloakClient, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	accessTokenUrl := fmt.Sprintf("%s/auth/realms/master/protocol/openid-connect/token", baseUrl)

	accessTokenData := url.Values{}

	accessTokenData.Set("client_id", clientId)
	accessTokenData.Set("client_secret", clientSecret)
	accessTokenData.Set("grant_type", "client_credentials")

	accessTokenRequest, _ := http.NewRequest("POST", accessTokenUrl, strings.NewReader(accessTokenData.Encode()))
	accessTokenResponse, err := httpClient.Do(accessTokenRequest)

	if err != nil {
		return nil, err
	}

	defer accessTokenResponse.Body.Close()

	body, _ := ioutil.ReadAll(accessTokenResponse.Body)

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)

	if err != nil {
		return nil, err
	}

	fmt.Println(clientCredentials.AccessToken)

	return &KeycloakClient{
		baseUrl:           baseUrl,
		clientCredentials: &clientCredentials,
		httpClient:        httpClient,
	}, nil
}
