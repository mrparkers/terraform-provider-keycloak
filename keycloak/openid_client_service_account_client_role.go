package keycloak

import (
	"fmt"
)

type OpenidClientServiceAccountClientRole struct {
	Id                   string `json:"id"`
	RealmId              string `json:"-"`
	ClientId             string `json:"-"`
	ServiceAccountUserId string `json:"-"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientServiceAccountClientRole(serviceAccountRole *OpenidClientServiceAccountClientRole) error {
	serviceAccountRoles := []OpenidClientServiceAccountClientRole{*serviceAccountRole}

	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ClientId), serviceAccountRoles)

	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientServiceAccountClientRole(realm, client, serviceAccountUserId, roleId string) error {
	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountClientRole(realm, client, serviceAccountUserId, roleId)
	if err != nil {
		return err
	}
	serviceAccountRoles := []OpenidClientServiceAccountClientRole{*serviceAccountRole}
	err = keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, serviceAccountRole.ClientId), &serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountClientRole(realm, client, serviceAccountUserId, roleId string) (*OpenidClientServiceAccountClientRole, error) {
	serviceAccountRoles := []OpenidClientServiceAccountClientRole{
		{
			Id:                   roleId,
			RealmId:              realm,
			ClientId:             client,
			ServiceAccountUserId: serviceAccountUserId,
		},
	}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, client), &serviceAccountRoles, nil)
	if err != nil {
		return nil, err
	}
	for _, serviceAccountRole := range serviceAccountRoles {
		if serviceAccountRole.Id == roleId {
			serviceAccountRole.RealmId = realm
			serviceAccountRole.ClientId = client
			serviceAccountRole.ServiceAccountUserId = serviceAccountUserId
			return &serviceAccountRole, nil
		}
	}
	return &OpenidClientServiceAccountClientRole{}, nil
}
