package keycloak

import (
	"fmt"
)

type ClientScope struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Protocol    string `json:"protocol"`
	Attributes  struct {
		DisplayOnConsentScreen string `json:"display.on.consent.screen"` // boolean in string form
		ConsentScreenText      string `json:"consent.screen.text"`
	} `json:"attributes"`
}

func (keycloakClient *KeycloakClient) NewClientScope(clientScope *ClientScope) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/client-scopes", clientScope.RealmId), clientScope)
	if err != nil {
		return err
	}

	clientScope.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetClientScope(realmId, id string) (*ClientScope, error) {
	var clientScope ClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/client-scopes/%s", realmId, id), &clientScope)
	if err != nil {
		return nil, err
	}

	clientScope.RealmId = realmId

	return &clientScope, nil
}

func (keycloakClient *KeycloakClient) UpdateClientScope(clientScope *ClientScope) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/client-scopes/%s", clientScope.RealmId, clientScope.Id), clientScope)
}

func (keycloakClient *KeycloakClient) DeleteClientScope(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/client-scopes/%s", realmId, id))
}
