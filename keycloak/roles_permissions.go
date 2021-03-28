package keycloak

import (
	"fmt"
)

type RolePermissionsInput struct {
	Enabled bool `json:"enabled"`
}

type RolePermissions struct {
	RealmId          string                 `json:"-"`
	RoleId           string                 `json:"-"`
	Enabled          bool                   `json:"enabled"`
	Resource         string                 `json:"resource"`
	ScopePermissions map[string]interface{} `json:"scopePermissions"`
}

func (keycloakClient *KeycloakClient) EnableRolePermissions(realmId, clientId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/roles-by-id/%s/management/permissions", realmId, clientId), RolePermissionsInput{Enabled: true})
}

func (keycloakClient *KeycloakClient) DisableRolePermissions(realmId, clientId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/roles-by-id/%s/management/permissions", realmId, clientId), RolePermissionsInput{Enabled: false})
}

func (keycloakClient *KeycloakClient) GetRolePermissions(realmId, roleId string) (*RolePermissions, error) {
	var rolePermissions RolePermissions
	rolePermissions.RealmId = realmId
	rolePermissions.RoleId = roleId

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/roles-by-id/%s/management/permissions", realmId, roleId), &rolePermissions, nil)
	if err != nil {
		return nil, err
	}

	return &rolePermissions, nil
}
