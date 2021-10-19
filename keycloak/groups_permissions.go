package keycloak

import (
	"fmt"
)

type GroupsPermissionsInput struct {
	Enabled bool `json:"enabled"`
}

type GroupsPermissions struct {
	RealmId          string                 `json:"-"`
	GroupId          string                 `json:"-"`
	Enabled          bool                   `json:"enabled"`
	Resource         string                 `json:"resource"`
	ScopePermissions map[string]interface{} `json:"scopePermissions"`
}

func (keycloakClient *KeycloakClient) EnableGroupsPermissions(realmId, groupId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/groups/%s/management/permissions", realmId, groupId), GroupsPermissionsInput{Enabled: true})
}

func (keycloakClient *KeycloakClient) DisableGroupsPermissions(realmId, groupId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/groups/%s/management/permissions", realmId, groupId), GroupsPermissionsInput{Enabled: false})
}

func (keycloakClient *KeycloakClient) GetGroupsPermissions(realmId, groupId string) (*GroupsPermissions, error) {
	var groupPermissions GroupsPermissions

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/groups/%s/management/permissions", realmId, groupId), &groupPermissions, nil)
	if err != nil {
		return nil, err
	}

	groupPermissions.RealmId = realmId
	groupPermissions.GroupId = groupId

	return &groupPermissions, nil
}
