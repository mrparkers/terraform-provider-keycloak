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
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/roles-by-id/%s/composites/realm", realmId, id), &composites, nil)
	if err != nil {
		return nil, err
	}

	return composites, nil
}
