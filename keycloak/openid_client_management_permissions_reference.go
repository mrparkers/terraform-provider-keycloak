package keycloak

import (
	"fmt"
)

type OpenIdClientManagementPermissionsReference struct {
	RealmId          string
	ClientId         string
	Enabled          bool
	Resource         string
	ScopePermissions map[string]string
}

func (reference *OpenIdClientManagementPermissionsReference) openIdClientManagementPermissionsReferencePath() string {
	return openIdClientManagementPermissionsReferencePath(reference.RealmId, reference.ClientId)
}

func openIdClientManagementPermissionsReferencePath(realmId, clientId string) string {
	return fmt.Sprintf("/realms/%s/clients/%s/management/permissions", realmId, clientId)
}

func (reference *OpenIdClientManagementPermissionsReference) convertToGenericManagementPermissionsReference() *managementPermissionReference {
	return &managementPermissionReference{
		Enabled:          reference.Enabled,
		Resource:         reference.Resource,
		ScopePermissions: reference.ScopePermissions,
	}
}

func (genericReference *managementPermissionReference) convertToOpenIdClientManagementPermissionsReference(realmId, clientId string) *OpenIdClientManagementPermissionsReference {
	return &OpenIdClientManagementPermissionsReference{
		RealmId:          realmId,
		ClientId:         clientId,
		Enabled:          genericReference.Enabled,
		Resource:         genericReference.Resource,
		ScopePermissions: genericReference.ScopePermissions,
	}
}

func (keycloakClient *KeycloakClient) GetOpenIdClientManagementPermissionsReference(realmId, clientId string) (*OpenIdClientManagementPermissionsReference, error) {
	var genericReference *managementPermissionReference

	err := keycloakClient.get(openIdClientManagementPermissionsReferencePath(realmId, clientId), &genericReference, nil)

	if err != nil {
		return nil, err
	}

	return genericReference.convertToOpenIdClientManagementPermissionsReference(realmId, clientId), nil
}

func (keycloakClient *KeycloakClient) CreateOpenIdClientManagementPermissionsReference(realmId, clientId string) error {
	return keycloakClient.put(openIdClientManagementPermissionsReferencePath(realmId, clientId), enableClientManagementPermissionsReference())
}

func (keycloakClient *KeycloakClient) DeleteOpenIdClientManagementPermissionsReference(realmId, clientId string) error {
	return keycloakClient.put(openIdClientManagementPermissionsReferencePath(realmId, clientId), disableClientManagementPermissionsReference())
}

func (keycloakClient *KeycloakClient) UpdateOpenIdClientManagementPermissionsReference(reference *OpenIdClientManagementPermissionsReference) error {
	return keycloakClient.put(reference.openIdClientManagementPermissionsReferencePath(), reference.convertToGenericManagementPermissionsReference())
}
