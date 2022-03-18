package keycloak

import (
	"context"
	"fmt"
)

func (keycloakClient *KeycloakClient) GetUserRoleMappings(ctx context.Context, realmId string, userId string) (*RoleMapping, error) {
	var roleMapping *RoleMapping
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings", realmId, userId), &roleMapping, nil)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (keycloakClient *KeycloakClient) AddRealmRolesToUser(ctx context.Context, realmId, userId string, roles []*Role) error {
	_, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/realm", realmId, userId), roles)

	return err
}

func (keycloakClient *KeycloakClient) AddClientRolesToUser(ctx context.Context, realmId, userId, clientId string, roles []*Role) error {
	_, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realmId, userId, clientId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveRealmRolesFromUser(ctx context.Context, realmId, userId string, roles []*Role) error {
	err := keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/realm", realmId, userId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveClientRolesFromUser(ctx context.Context, realmId, userId, clientId string, roles []*Role) error {
	err := keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realmId, userId, clientId), roles)

	return err
}
