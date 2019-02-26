package keycloak

import "fmt"

type GenericClient struct {
	Id       string `json:"id,omitempty"`
	ClientId string `json:"clientId"`
	RealmId  string `json:"-"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`

	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

func (keycloakClient *KeycloakClient) listGenericClients(realmId string) ([]*GenericClient, error) {
	var clients []*GenericClient

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients", realmId), &clients)
	if err != nil {
		return nil, err
	}

	for _, client := range clients {
		client.RealmId = realmId
	}

	return clients, nil
}
