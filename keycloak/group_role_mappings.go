package keycloak

import "fmt"

func (keycloakClient *KeycloakClient) GetGroupRoleMappings(realmId string, userId string) (*RoleMapping, error) {
	var roleMapping *RoleMapping
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/groups/%s/role-mappings", realmId, userId), &roleMapping, nil)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (keycloakClient *KeycloakClient) AddRealmRolesToGroup(realmId, groupId string, roles []*Role) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/groups/%s/role-mappings/realm", realmId, groupId), roles)

	return err
}

func (keycloakClient *KeycloakClient) AddClientRolesToGroup(realmId, groupId, clientId string, roles []*Role) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/groups/%s/role-mappings/clients/%s", realmId, groupId, clientId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveRealmRolesFromGroup(realmId, groupId string, roles []*Role) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/groups/%s/role-mappings/realm", realmId, groupId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveClientRolesFromGroup(realmId, groupId, clientId string, roles []*Role) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/groups/%s/role-mappings/clients/%s", realmId, groupId, clientId), roles)

	return err
}
