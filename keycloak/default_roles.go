package keycloak

import (
	"context"
	"fmt"
)

type DefaultRoles struct {
	Id           string   `json:"id,omitempty"`
	RealmId      string   `json:"-"`
	DefaultRoles []string `json:"-"`
}

func (keycloakClient *KeycloakClient) GetDefaultRoles(ctx context.Context, realmId, id string) ([]*Role, error) {
	var composites []*Role
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/roles-by-id/%s/composites", realmId, id), &composites, nil)
	if err != nil {
		return nil, err
	}
	for _, composite := range composites {
		if composite.ClientRole && composite.ClientId == "" {
			client, err := keycloakClient.GetGenericClient(ctx, realmId, composite.ContainerId)
			if err != nil {
				return nil, err
			}
			composite.ClientId = client.ClientId
		}
	}
	return composites, nil
}

func (keycloakClient *KeycloakClient) GetQualifiedRoleName(ctx context.Context, realmId string, role *Role) (string, error) {
	if !role.ClientRole {
		return role.Name, nil
	}
	if role.ClientId != "" {
		return fmt.Sprintf("%s/%s", role.ClientId, role.Name), nil
	}
	genericClient, err := keycloakClient.GetGenericClient(ctx, realmId, role.ContainerId)
	if err != nil {
		return "", err
	}
	role.ClientId = genericClient.ClientId
	return fmt.Sprintf("%s/%s", role.ClientId, role.Name), nil
}

func (keycloakClient *KeycloakClient) GetRoleFullNames(ctx context.Context, realmId string, roles []*Role) ([]string, error) {
	fullNames := make([]string, len(roles))
	for i, role := range roles {
		if !role.ClientRole {
			fullNames[i] = role.Name
			continue
		}
		fullName, err := keycloakClient.GetQualifiedRoleName(ctx, realmId, role)
		if err != nil {
			return []string{}, err
		}
		fullNames[i] = fullName
	}
	return fullNames, nil
}
