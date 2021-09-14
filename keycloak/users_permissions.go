package keycloak

import (
	"fmt"
)

type UsersPermissionsInput struct {
	Enabled bool `json:"enabled"`
}

type UsersPermissions struct {
	RealmId          string            `json:"-"`
	Enabled          bool              `json:"enabled"`
	Resource         string            `json:"resource"`
	ScopePermissions map[string]string `json:"scopePermissions"`
}

func (keycloakClient *KeycloakClient) EnableUsersPermissions(realmId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/users-management-permissions", realmId), UsersPermissionsInput{Enabled: true})
}

func (keycloakClient *KeycloakClient) DisableUsersPermissions(realmId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/users-management-permissions", realmId), UsersPermissionsInput{Enabled: false})
}

func (keycloakClient *KeycloakClient) GetUsersPermissions(realmId string) (*UsersPermissions, error) {
	var openidClientPermissions UsersPermissions

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users-management-permissions", realmId), &openidClientPermissions, nil)
	if err != nil {
		return nil, err
	}

	openidClientPermissions.RealmId = realmId

	return &openidClientPermissions, nil
}
