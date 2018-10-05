package keycloak

import (
	"fmt"
)

type OpenidClient struct {
	Id       string `json:"id,omitempty"`
	ClientId string `json:"clientId"`
	RealmId  string `json:"-"`
	Protocol string `json:"protocol"`
}

func (keycloakClient *KeycloakClient) NewOpenidClient(client *OpenidClient) error {
	client.Protocol = "openid-connect"

	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients", client.RealmId), client)
	if err != nil {
		return err
	}

	client.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClient(realmId, id string) (*OpenidClient, error) {
	var client OpenidClient

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), &client)
	if err != nil {
		return nil, err
	}

	client.RealmId = realmId

	return &client, nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClient(client *OpenidClient) error {
	client.Protocol = "openid-connect"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s", client.RealmId, client.Id), client)
}

func (keycloakClient *KeycloakClient) DeleteOpenidClient(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s", realmId, id))
}
