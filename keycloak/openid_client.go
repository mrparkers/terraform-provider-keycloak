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

	StandardFlowEnabled       bool `json:"standardFlowEnabled"`
	ImplicitFlowEnabled       bool `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled bool `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled    bool `json:"serviceAccountsEnabled"`

	ValidRedirectUris []string `json:"redirectUris"`
	WebOrigins        []string `json:"webOrigins"`
}

func (keycloakClient *KeycloakClient) ValidateOpenidClient(client *OpenidClient) error {
	if client.BearerOnly && (client.StandardFlowEnabled || client.ImplicitFlowEnabled || client.DirectAccessGrantsEnabled || client.ServiceAccountsEnabled) {
		return fmt.Errorf("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client")
	}

	if (client.StandardFlowEnabled || client.ImplicitFlowEnabled) && len(client.ValidRedirectUris) == 0 {
		return fmt.Errorf("validation error: standard (authorization code) and implicit flows require at least one valid redirect uri")
	}

	if client.ServiceAccountsEnabled && client.PublicClient {
		return fmt.Errorf("validation error: service accounts (client credentials flow) cannot be enabled on public clients")
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

func (keycloakClient *KeycloakClient) GetOpenidClientByClientId(realmId, clientId string) (*OpenidClient, error) {
	var clients []OpenidClient
	var clientSecret openidClientSecret

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients?clientId=%s", realmId, clientId), &clients)
	if err != nil {
		return nil, err
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("openid client with name %s does not exist", clientId)
	}

	client := clients[0]

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/client-secret", realmId, client.Id), &clientSecret)
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

func (keycloakClient *KeycloakClient) GetOpenidClientDefaultScopes(realmId, clientId string) ([]*OpenidClientScope, error) {
	var scopes []*OpenidClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/default-client-scopes", realmId, clientId), &scopes)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (keycloakClient *KeycloakClient) AttachOpenidClientDefaultScopes(realmId, clientId string, scopeNames []string) error {
	openidClient, err := keycloakClient.GetOpenidClient(realmId, clientId)
	if err != nil && ErrorIs404(err) {
		return fmt.Errorf("validation error: client with id %s does not exist", clientId)
	} else if err != nil {
		return err
	}

	if openidClient.BearerOnly {
		return fmt.Errorf("validation error: client with id %s uses access type BEARER-ONLY which does not use scopes", clientId)
	}

	allOpenidClientScopes, err := keycloakClient.listOpenidClientScopesWithFilter(realmId, includeOpenidClientScopesMatchingNames(scopeNames))
	if err != nil {
		return err
	}

	for _, openidClientScope := range allOpenidClientScopes {
		err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/default-client-scopes/%s", realmId, clientId, openidClientScope.Id), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) DetachOpenidClientDefaultScopes(realmId, clientId string, scopeNames []string) error {
	allOpenidClientScopes, err := keycloakClient.listOpenidClientScopesWithFilter(realmId, includeOpenidClientScopesMatchingNames(scopeNames))
	if err != nil {
		return err
	}

	for _, openidClientScope := range allOpenidClientScopes {
		err := keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/default-client-scopes/%s", realmId, clientId, openidClientScope.Id))
		if err != nil {
			return err
		}
	}

	return nil
}
