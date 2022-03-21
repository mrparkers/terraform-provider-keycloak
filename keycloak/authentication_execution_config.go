package keycloak

import (
	"context"
	"fmt"
)

// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_authenticatorconfigrepresentation
type AuthenticationExecutionConfig struct {
	RealmId     string            `json:"-"`
	ExecutionId string            `json:"-"`
	Id          string            `json:"id"`
	Alias       string            `json:"alias"`
	Config      map[string]string `json:"config"`
}

// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_newexecutionconfig
func (keycloakClient *KeycloakClient) NewAuthenticationExecutionConfig(ctx context.Context, config *AuthenticationExecutionConfig) (string, error) {
	_, location, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/authentication/executions/%s/config", config.RealmId, config.ExecutionId), config)
	if err != nil {
		return "", err
	}
	return getIdFromLocationHeader(location), nil
}

// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_getauthenticatorconfig
func (keycloakClient *KeycloakClient) GetAuthenticationExecutionConfig(ctx context.Context, config *AuthenticationExecutionConfig) error {
	return keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/authentication/config/%s", config.RealmId, config.Id), config, nil)
}

// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_updateauthenticatorconfig
func (keycloakClient *KeycloakClient) UpdateAuthenticationExecutionConfig(ctx context.Context, config *AuthenticationExecutionConfig) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/authentication/config/%s", config.RealmId, config.Id), config)
}

// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_removeauthenticatorconfig
func (keycloakClient *KeycloakClient) DeleteAuthenticationExecutionConfig(ctx context.Context, config *AuthenticationExecutionConfig) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/authentication/config/%s", config.RealmId, config.Id), nil)
}
