package keycloak

import (
	"fmt"
)

type Role struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	ClientId    string `json:"-"`
	RoleId      string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClientRole  bool   `json:"clientRole"`
	ContainerId string `json:"containerId"`
	Composite   bool   `json:"composite"`
}

/*
 * Realm roles: /realms/${realm_id}/roles
 * Client roles: /realms/${realm_id}/clients/${client_id}/roles
 */
func roleByNameUrl(realmId, clientId string) string {
	if clientId == "" {
		return fmt.Sprintf("/realms/%s/roles", realmId)
	}

	return fmt.Sprintf("/realms/%s/clients/%s/roles", realmId, clientId)
}

func (keycloakClient *KeycloakClient) CreateRole(role *Role) error {
	url := roleByNameUrl(role.RealmId, role.ClientId)

	if role.ClientId != "" {
		role.ContainerId = role.ClientId
		role.ClientRole = true
	}

	_, _, err := keycloakClient.post(url, role)
	if err != nil {
		return err
	}

	var createdRole Role
	err = keycloakClient.get(fmt.Sprintf("%s/%s", url, role.Name), &createdRole, nil)
	if err != nil {
		return err
	}

	role.Id = createdRole.Id

	return nil
}

func (keycloakClient *KeycloakClient) GetRole(realmId, id string) (*Role, error) {
	var role Role

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/roles-by-id/%s", realmId, id), &role, nil)
	if err != nil {
		return nil, err
	}

	role.RealmId = realmId

	if role.ClientRole {
		role.ClientId = role.ContainerId
	}

	return &role, nil
}

func (keycloakClient *KeycloakClient) GetRoleByName(realmId, clientId, name string) (*Role, error) {
	var role Role

	err := keycloakClient.get(fmt.Sprintf("%s/%s", roleByNameUrl(realmId, clientId), name), &role, nil)
	if err != nil {
		return nil, err
	}

	role.RealmId = realmId

	if role.ClientRole {
		role.ClientId = role.ContainerId
	}

	return &role, nil
}

func (keycloakClient *KeycloakClient) UpdateRole(role *Role) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/roles-by-id/%s", role.RealmId, role.Id), role)
}

func (keycloakClient *KeycloakClient) DeleteRole(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/roles-by-id/%s", realmId, id), nil)
}

func (keycloakClient *KeycloakClient) AddCompositesToRole(role *Role, compositeRoles []*Role) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/roles-by-id/%s/composites", role.RealmId, role.Id), compositeRoles)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) RemoveCompositesFromRole(role *Role, compositeRoles []*Role) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/roles-by-id/%s/composites", role.RealmId, role.Id), compositeRoles)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetRoleComposites(role *Role) ([]*Role, error) {
	var composites []*Role

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/roles-by-id/%s/composites", role.RealmId, role.Id), &composites, nil)
	if err != nil {
		return nil, err
	}

	return composites, nil
}
