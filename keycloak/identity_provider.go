package keycloak

import (
	"fmt"
	"log"
)

type IdentityProviderConfig struct {
	Key                              string             `json:"key,omitempty"`
	HostIp                           string             `json:"hostIp,omitempty"`
	UseJwksUrl                       KeycloakBoolQuoted `json:"useJwksUrl,omitempty"`
	ClientId                         string             `json:"clientId,omitempty"`
	ClientSecret                     string             `json:"clientSecret,omitempty"`
	DisableUserInfo                  KeycloakBoolQuoted `json:"disableUserInfo"`
	HideOnLoginPage                  KeycloakBoolQuoted `json:"hideOnLoginPage"`
	NameIDPolicyFormat               string             `json:"nameIDPolicyFormat,omitempty"`
	SingleLogutServiceUrl            string             `json:"singleLogoutServiceUrl,omitempty"`
	SingleSignOnServiceUrl           string             `json:"singleSignOnServiceUrl,omitempty"`
	SigningCertificate               string             `json:"signingCertificate,omitempty"`
	SignatureAlgorithm               string             `json:"signatureAlgorithm,omitempty"`
	XmlSignKeyInfoKeyNameTransformer string             `json:"xmlSignKeyInfoKeyNameTransformer,omitempty"`
	PostBindingAuthnRequest          KeycloakBoolQuoted `json:"postBindingAuthnRequest,omitempty"`
	PostBindingResponse              KeycloakBoolQuoted `json:"postBindingResponse,omitempty"`
	PostBindingLogout                KeycloakBoolQuoted `json:"postBindingLogout,omitempty"`
	ForceAuthn                       KeycloakBoolQuoted `json:"forceAuthn,omitempty"`
	WantAuthnRequestsSigned          KeycloakBoolQuoted `json:"wantAuthnRequestsSigned,omitempty"`
	WantAssertionsSigned             KeycloakBoolQuoted `json:"wantAssertionsSigned,omitempty"`
	WantAssertionsEncrypted          KeycloakBoolQuoted `json:"wantAssertionsEncrypted,omitempty"`
	BackchannelSupported             KeycloakBoolQuoted `json:"backchannelSupported,omitempty"`
	ValidateSignature                KeycloakBoolQuoted `json:"validateSignature,omitempty"`
	AuthorizationUrl                 string             `json:"authorizationUrl,omitempty"`
	TokenUrl                         string             `json:"tokeUrl,omitempty"`
	LoginHint                        string             `json:"loginHint,omitempty"`
}

type IdentityProvider struct {
	Realm                       string                  `json:"-"`
	InternalId                  string                  `json:"internalId,omitempty"`
	UpdateProfileFirstLoginMode string                  `json:"updateProfileFirstLoginMode,omitempty"`
	Alias                       string                  `json:"alias,omitempty"`
	DisplayName                 string                  `json:"displayName,omitempty"`
	ProviderId                  string                  `json:"providerId,omitempty"`
	Enabled                     bool                    `json:"enabled,omitempty"`
	StoreToken                  KeycloakBool            `json:"storeToken"`
	AddReadTokenRoleOnCreate    KeycloakBool            `json:"addReadTokenRoleOnCreate"`
	AuthenticateByDefault       bool                    `json:"authenticateByDefault"`
	LinkOnly                    KeycloakBool            `json:"linkOnly"`
	TrustEmail                  KeycloakBool            `json:"trustEmail"`
	FirstBrokerLoginFlowAlias   string                  `json:"firstBrokerLoginFlowAlias,omitempty"`
	PostBrokerLoginFlowAlias    string                  `json:"postBrokerLoginFlowAlias"`
	Config                      *IdentityProviderConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewIdentityProvider(identityProvider *IdentityProvider) error {
	log.Printf("[WARN] Realm: %s", identityProvider.Realm)
	_, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances", identityProvider.Realm), identityProvider)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetIdentityProvider(realm, alias string) (*IdentityProvider, error) {
	var identityProvider IdentityProvider
	identityProvider.Realm = realm

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias), &identityProvider)
	if err != nil {
		return nil, err
	}

	return &identityProvider, nil
}

func (keycloakClient *KeycloakClient) UpdateIdentityProvider(identityProvider *IdentityProvider) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", identityProvider.Realm, identityProvider.Alias), identityProvider)
}

func (keycloakClient *KeycloakClient) DeleteIdentityProvider(realm, alias string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias))
}
