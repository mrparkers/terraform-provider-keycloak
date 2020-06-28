package keycloak

import "fmt"

// struct for the MappingRepresentation
// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_mappingsrepresentation
type UserRoleMapping struct {
	ClientMappings map[string]*ClientRoleMapping `json:"clientMappings"`
	RealmMappings  []*Role                       `json:"realmMappings"`
}

// struct for the ClientMappingRepresentation
// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_clientmappingsrepresentation
type ClientRoleMapping struct {
	Client   string  `json:"client"`
	Id       string  `json:"id"`
	Mappings []*Role `json:"mappings"`
}

func (keycloakClient *KeycloakClient) GetUserRoleMappings(realmId string, userId string) (*UserRoleMapping, error) {
	var roleMapping *UserRoleMapping
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users/%s/role-mappings", realmId, userId), &roleMapping, nil)
	if err != nil {
		return nil, err
	}

	return roleMapping, nil
}

func (keycloakClient *KeycloakClient) AddRealmRolesToUser(realmId, userId string, roles []*Role) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/users/%s/role-mappings/realm", realmId, userId), roles)

	return err
}

func (keycloakClient *KeycloakClient) AddClientRolesToUser(realmId, userId, clientId string, roles []*Role) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realmId, userId, clientId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveRealmRolesFromUser(realmId, userId string, roles []*Role) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s/role-mappings/realm", realmId, userId), roles)

	return err
}

func (keycloakClient *KeycloakClient) RemoveClientRolesFromUser(realmId, userId, clientId string, roles []*Role) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realmId, userId, clientId), roles)

	return err
}
