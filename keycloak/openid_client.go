package keycloak

import (
	"fmt"
)

type openidClientSecret struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type OpenidClient struct {
	Id                      string `json:"id,omitempty"`
	ClientId                string `json:"clientId"`
	RealmId                 string `json:"-"`
	Name                    string `json:"name"`
	Protocol                string `json:"protocol"`                // always openid-connect for this resource
	ClientAuthenticatorType string `json:"clientAuthenticatorType"` // always client-secret for now, don't have a need for JWT here
	ClientSecret            string `json:"secret,omitempty"`

	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`

	// Attributes below indicate client access type. If both are false, access type is confidential. Both cannot be true (although the Keycloak API lets you do this)
	PublicClient bool `json:"publicClient"`
	BearerOnly   bool `json:"bearerOnly"`

	ValidRedirectUris []string `json:"redirectUris"`
}

func (keycloakClient *KeycloakClient) ValidateOpenidClient(client *OpenidClient) error {
	if !client.BearerOnly && len(client.ValidRedirectUris) == 0 {
		return fmt.Errorf("validation error: must specify at least one valid redirect uri if access type is PUBLIC or CONFIDENTIAL")
	}

	return nil
}

func (keycloakClient *KeycloakClient) NewOpenidClient(client *OpenidClient) error {
	client.Protocol = "openid-connect"
	client.ClientAuthenticatorType = "client-secret"

	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients", client.RealmId), client)
	if err != nil {
		return err
	}

	client.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClient(realmId, id string) (*OpenidClient, error) {
	var client OpenidClient
	var clientSecret openidClientSecret

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), &client)
	if err != nil {
		return nil, err
	}

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/client-secret", realmId, id), &clientSecret)
	if err != nil {
		return nil, err
	}

	client.RealmId = realmId
	client.ClientSecret = clientSecret.Value

	return &client, nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClient(client *OpenidClient) error {
	client.Protocol = "openid-connect"
	client.ClientAuthenticatorType = "client-secret"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s", client.RealmId, client.Id), client)
}

func (keycloakClient *KeycloakClient) DeleteOpenidClient(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s", realmId, id))
}
