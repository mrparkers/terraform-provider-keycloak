package keycloak

import (
	"context"
	"fmt"
)

type OpenidClientServiceAccountRole struct {
	Id                   string `json:"id"`
	RealmId              string `json:"-"`
	ServiceAccountUserId string `json:"-"`
	Name                 string `json:"name,omitempty"`
	ClientRole           bool   `json:"clientRole"`
	Composite            bool   `json:"composite"`
	ContainerId          string `json:"containerId"`
	Description          string `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientServiceAccountRole(ctx context.Context, serviceAccountRole *OpenidClientServiceAccountRole) error {
	serviceAccountRoles := []OpenidClientServiceAccountRole{*serviceAccountRole}
	_, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ContainerId), serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientServiceAccountRole(ctx context.Context, realm, serviceAccountUserId, clientId, roleId string) error {
	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountRole(ctx, realm, serviceAccountUserId, clientId, roleId)
	if err != nil {
		return err
	}
	serviceAccountRoles := []OpenidClientServiceAccountRole{*serviceAccountRole}
	err = keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountRole(ctx context.Context, realm, serviceAccountUserId, clientId, roleId string) (*OpenidClientServiceAccountRole, error) {
	serviceAccountRoles := []OpenidClientServiceAccountRole{
		{
			Id:                   roleId,
			RealmId:              realm,
			ContainerId:          clientId,
			ServiceAccountUserId: serviceAccountUserId,
		},
	}
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles, nil)
	if err != nil {
		return nil, err
	}
	for _, serviceAccountRole := range serviceAccountRoles {
		if serviceAccountRole.Id == roleId {
			serviceAccountRole.RealmId = realm
			serviceAccountRole.ServiceAccountUserId = serviceAccountUserId
			return &serviceAccountRole, nil
		}
	}
	return &OpenidClientServiceAccountRole{}, nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountRealmRoles(ctx context.Context, realm, serviceAccountUserId string) ([]*OpenidClientServiceAccountRole, error) {
	var serviceAccountRoles []*OpenidClientServiceAccountRole

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/realm/composite", realm, serviceAccountUserId), &serviceAccountRoles, nil)
	if err != nil {
		return nil, err
	}

	for _, serviceAccountRole := range serviceAccountRoles {
		serviceAccountRole.RealmId = realm
		serviceAccountRole.ServiceAccountUserId = serviceAccountUserId
	}

	return serviceAccountRoles, nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountClientRoles(ctx context.Context, realm, serviceAccountUserId, clientId string) ([]*OpenidClientServiceAccountRole, error) {
	var serviceAccountRoles []*OpenidClientServiceAccountRole

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles, nil)
	if err != nil {
		return nil, err
	}

	for _, serviceAccountRole := range serviceAccountRoles {
		serviceAccountRole.RealmId = realm
		serviceAccountRole.ServiceAccountUserId = serviceAccountUserId
	}

	return serviceAccountRoles, nil
}
