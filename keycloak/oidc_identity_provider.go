package keycloak

import (
	"fmt"
	"log"
)

type OidcIdentityProviderConfig struct {
	BackchannelSupported KeycloakBoolQuoted `json:"backchannelSupported,omitempty"`
	UseJwksUrl           KeycloakBoolQuoted `json:"useJwksUrl,omitempty"`
	ValidateSignature    KeycloakBoolQuoted `json:"validateSignature,omitempty"`
	AuthorizationUrl     string             `json:"authorizationUrl,omitempty"`
	ClientId             string             `json:"clientId,omitempty"`
	ClientSecret         string             `json:"clientSecret,omitempty"`
	DisableUserInfo      KeycloakBoolQuoted `json:"disableUserInfo,omitempty"`
	HideOnLoginPage      KeycloakBoolQuoted `json:"hideOnLoginPage,omitempty"`
	TokenUrl             string             `json:"tokeUrl,omitempty"`
	LoginHint            string             `json:"loginHint,omitempty"`
}

type OidcIdentityProvider struct {
	RealmId                     string                      `json:"-"`
	InternalId                  string                      `json:"internalId,omitempty"`
	UpdateProfileFirstLoginMode string                      `json:"updateProfileFirstLoginMode,omitempty"`
	Alias                       string                      `json:"alias,omitempty"`
	DisplayName                 string                      `json:"displayName,omitempty"`
	ProviderId                  string                      `json:"providerId,omitempty"`
	Enabled                     bool                        `json:"enabled,omitempty"`
	StoreToken                  KeycloakBool                `json:"storeToken"`
	AddReadTokenRoleOnCreate    KeycloakBool                `json:"addReadTokenRoleOnCreate"`
	AuthenticateByDefault       bool                        `json:"authenticateByDefault"`
	LinkOnly                    KeycloakBool                `json:"linkOnly"`
	TrustEmail                  KeycloakBool                `json:"trustEmail"`
	FirstBrokerLoginFlowAlias   string                      `json:"firstBrokerLoginFlowAlias,omitempty"`
	PostBrokerLoginFlowAlias    string                      `json:"postBrokerLoginFlowAlias"`
	Config                      *OidcIdentityProviderConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewOidcIdentityProvider(oidcIdentityProvider *OidcIdentityProvider) error {
	log.Printf("[WARN] Realm: %s", oidcIdentityProvider.RealmId)
	_, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances", oidcIdentityProvider.RealmId), oidcIdentityProvider)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetOidcIdentityProvider(realmId, alias string) (*OidcIdentityProvider, error) {
	var oidcIdentityProvider *OidcIdentityProvider
	oidcIdentityProvider.RealmId = realmId

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realmId, alias), oidcIdentityProvider)
	if err != nil {
		return nil, err
	}

	return oidcIdentityProvider, nil
}

func (keycloakClient *KeycloakClient) UpdateOidcIdentityProvider(oidcIdentityProvider *OidcIdentityProvider) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", oidcIdentityProvider.RealmId, oidcIdentityProvider.Alias), oidcIdentityProvider)
}

func (keycloakClient *KeycloakClient) DeleteOidcIdentityProvider(realmId, alias string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realmId, alias))
}
