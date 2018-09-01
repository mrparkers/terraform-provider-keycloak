package keycloak

import "net/http"

type KeycloakClient struct {
	url        string
	httpClient *http.Client
}

func NewKeycloakClient(url string) *KeycloakClient {
	httpClient := &http.Client{}

	return &KeycloakClient{
		url:        url,
		httpClient: httpClient,
	}
}
