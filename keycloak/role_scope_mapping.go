package keycloak

import (
	"fmt"
)

func roleScopeMappingUrl(realmId, clientId string, role *Role) string {
	return fmt.Sprintf("/realms/%s/clients/%s/scope-mappings/clients/%s", realmId, clientId, role.ClientId)
}

func (keycloakClient *KeycloakClient) CreateRoleScopeMapping(realmId string, clientId string, role *Role) error {
	roleUrl := roleScopeMappingUrl(realmId, clientId, role)

	_, _, err := keycloakClient.post(roleUrl, []Role{*role})
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetRoleScopeMapping(realmId string, clientId string, role *Role) (*Role, error) {
	roleUrl := roleScopeMappingUrl(realmId, clientId, role)
	var roles []Role

	err := keycloakClient.get(roleUrl, &roles, nil)
	if err != nil {
		return nil, err
	}

	for _, mappedRole := range roles {
		if mappedRole.Id == role.Id {
			return role, nil
		}
	}

	return nil, nil
}

func (keycloakClient *KeycloakClient) DeleteRoleScopeMapping(realmId string, clientId string, role *Role) error {
	roleUrl := roleScopeMappingUrl(realmId, clientId, role)
	return keycloakClient.delete(roleUrl, nil)
}
