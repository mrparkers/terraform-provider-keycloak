package keycloak

import (
	"fmt"
)

type IdentityProviderPermissionsInput struct {
	Enabled bool `json:"enabled"`
}

type IdentityProviderPermissions struct {
	RealmId          string                 `json:"-"`
	ProviderAlias    string                 `json:"-"`
	Enabled          bool                   `json:"enabled"`
	Resource         string                 `json:"resource"`
	ScopePermissions map[string]interface{} `json:"scopePermissions"`
}

func (keycloakClient *KeycloakClient) EnableIdentityProviderPermissions(realmId, providerAlias string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/management/permissions", realmId, providerAlias), IdentityProviderPermissionsInput{Enabled: true})
}

func (keycloakClient *KeycloakClient) DisableIdentityProviderPermissions(realmId, providerAlias string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/management/permissions", realmId, providerAlias), IdentityProviderPermissionsInput{Enabled: false})
}

func (keycloakClient *KeycloakClient) GetIdentityProviderPermissions(realmId, providerAlias string) (*IdentityProviderPermissions, error) {
	var identityProviderPermissions IdentityProviderPermissions
	identityProviderPermissions.RealmId = realmId
	identityProviderPermissions.ProviderAlias = providerAlias

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/management/permissions", realmId, providerAlias), &identityProviderPermissions, nil)
	if err != nil {
		return nil, err
	}

	return &identityProviderPermissions, nil
}

func (identityProviderPermissions *IdentityProviderPermissions) GetTokenExchangeScopedPermissionId() (string, error) {
	if identityProviderPermissions.Enabled {
		return identityProviderPermissions.ScopePermissions["token-exchange"].(string), nil
	} else {
		return "", fmt.Errorf("identity provider permissions are not enabled, thus can not return the linked 'token-exchange' scope based permission")
	}
}
