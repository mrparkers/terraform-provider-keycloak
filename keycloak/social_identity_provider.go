package keycloak

import (
	"fmt"
	"log"
)

type SocialIdentityProviderConfig struct {
	Key             string             `json:"key,omitempty"`
	HostIp          string             `json:"hostIp,omitempty"`
	UseJwksUrl      KeycloakBoolQuoted `json:"useJwksUrl,omitempty"`
	ClientId        string             `json:"clientId,omitempty"`
	ClientSecret    string             `json:"clientSecret,omitempty"`
	DisableUserInfo KeycloakBoolQuoted `json:"disableUserInfo"`
	HideOnLoginPage KeycloakBoolQuoted `json:"hideOnLoginPage"`
}

type SocialIdentityProvider struct {
	Realm                       string                        `json:"-"`
	InternalId                  string                        `json:"internalId,omitempty"`
	UpdateProfileFirstLoginMode string                        `json:"updateProfileFirstLoginMode,omitempty"`
	Alias                       string                        `json:"alias,omitempty"`
	DisplayName                 string                        `json:"displayName,omitempty"`
	ProviderId                  string                        `json:"providerId,omitempty"`
	Enabled                     bool                          `json:"enabled,omitempty"`
	StoreToken                  KeycloakBool                  `json:"storeToken"`
	AddReadTokenRoleOnCreate    KeycloakBool                  `json:"addReadTokenRoleOnCreate"`
	AuthenticateByDefault       bool                          `json:"authenticateByDefault"`
	LinkOnly                    KeycloakBool                  `json:"linkOnly"`
	TrustEmail                  KeycloakBool                  `json:"trustEmail"`
	FirstBrokerLoginFlowAlias   string                        `json:"firstBrokerLoginFlowAlias,omitempty"`
	PostBrokerLoginFlowAlias    string                        `json:"postBrokerLoginFlowAlias"`
	Config                      *SocialIdentityProviderConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewSocialIdentityProvider(socialIdentityProvider *SocialIdentityProvider) error {
	log.Printf("[WARN] Realm: %s", socialIdentityProvider.Realm)
	_, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances", socialIdentityProvider.Realm), socialIdentityProvider)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetSocialIdentityProvider(realm, alias string) (*SocialIdentityProvider, error) {
	var socialIdentityProvider SocialIdentityProvider
	socialIdentityProvider.Realm = realm

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias), &socialIdentityProvider)
	if err != nil {
		return nil, err
	}

	return &socialIdentityProvider, nil
}

func (keycloakClient *KeycloakClient) UpdateSocialIdentityProvider(socialIdentityProvider *SocialIdentityProvider) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", socialIdentityProvider.Realm, socialIdentityProvider.Alias), socialIdentityProvider)
}

func (keycloakClient *KeycloakClient) DeleteSocialIdentityProvider(realm, alias string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias))
}
