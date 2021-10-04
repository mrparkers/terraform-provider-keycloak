package keycloak

import (
	"fmt"
)

type DefaultRoles struct {
	Id           string   `json:"id,omitempty"`
	RealmId      string   `json:"-"`
	DefaultRoles []string `json:"-"`
}

func (keycloakClient *KeycloakClient) GetDefaultRoles(realmId, id string) ([]*Role, error) {
	var composites []*Role
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/roles-by-id/%s/composites/realm", realmId, id),
		&composites, nil)
	if err != nil {
		return nil, err
	}

	return composites, nil
}

func (keycloakClient *KeycloakClient) AddDefaultRoles(realmId, id string, compositeRoles []*Role) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/roles-by-id/%s/composites", realmId, id), compositeRoles)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) RemoveDefaultRoles(realmId, id string, compositeRoles []*Role) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/roles-by-id/%s/composites", realmId, id), compositeRoles)
	if err != nil {
		return err
	}

	return nil
}
