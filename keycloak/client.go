package keycloak

import (
	"fmt"
	"strings"
)

type Client struct {
	Id       string `json:"id,omitempty"`
	ClientId string `json:"clientId"`
	RealmId  string `json:"-"`
}

func (keycloakClient *KeycloakClient) NewClient(client *Client) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients", client.RealmId), client)
	if err != nil {
		return err
	}

	client.Id = parseClientLocation(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetClient(realmId, id string) (*Client, error) {
	var client Client

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), &client)
	if err != nil {
		return nil, err
	}

	client.RealmId = realmId

	return &client, nil
}

func (keycloakClient *KeycloakClient) UpdateClient(client *Client) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s", client.RealmId, client.Id), client)
}

func (keycloakClient *KeycloakClient) DeleteClient(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s", realmId, id))
}

func parseClientLocation(location string) string {
	parts := strings.Split(location, "/")

	return parts[len(parts)-1]
}
