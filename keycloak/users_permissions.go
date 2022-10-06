package keycloak

import (
	"context"
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

func (keycloakClient *KeycloakClient) EnableUsersPermissions(ctx context.Context, realmId string) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/users-management-permissions", realmId), UsersPermissionsInput{Enabled: true})
}

func (keycloakClient *KeycloakClient) DisableUsersPermissions(ctx context.Context, realmId string) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/users-management-permissions", realmId), UsersPermissionsInput{Enabled: false})
}

func (keycloakClient *KeycloakClient) GetUsersPermissions(ctx context.Context, realmId string) (*UsersPermissions, error) {
	var openidClientPermissions UsersPermissions

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/users-management-permissions", realmId), &openidClientPermissions, nil)
	if err != nil {
		return nil, err
	}

	openidClientPermissions.RealmId = realmId

	return &openidClientPermissions, nil
}
