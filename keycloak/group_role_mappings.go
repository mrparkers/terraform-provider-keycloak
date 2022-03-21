package keycloak

import (
	"context"
	"fmt"
)

func (keycloakClient *KeycloakClient) GetGroupRoleMappings(ctx context.Context, realmId string, userId string) (*RoleMapping, error) {
	var roleMapping *RoleMapping
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/groups/%s/role-mappings", realmId, userId), &roleMapping, nil)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (keycloakClient *KeycloakClient) AddRealmRolesToGroup(ctx context.Context, realmId, groupId string, roles []*Role) error {
	_, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/groups/%s/role-mappings/realm", realmId, groupId), roles)

	return err
}

func (keycloakClient *KeycloakClient) AddClientRolesToGroup(ctx context.Context, realmId, groupId, clientId string, roles []*Role) error {
	_, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/groups/%s/role-mappings/clients/%s", realmId, groupId, clientId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveRealmRolesFromGroup(ctx context.Context, realmId, groupId string, roles []*Role) error {
	err := keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/groups/%s/role-mappings/realm", realmId, groupId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveClientRolesFromGroup(ctx context.Context, realmId, groupId, clientId string, roles []*Role) error {
	err := keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/groups/%s/role-mappings/clients/%s", realmId, groupId, clientId), roles)

	return err
}
