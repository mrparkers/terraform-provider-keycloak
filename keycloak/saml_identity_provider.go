package keycloak

import (
	"fmt"
	"log"
)

type SamlIdentityProviderConfig struct {
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
	UseJwksUrl                       KeycloakBoolQuoted `json:"useJwksUrl,omitempty"`
	ValidateSignature                KeycloakBoolQuoted `json:"validateSignature,omitempty"`
	HideOnLoginPage                  KeycloakBoolQuoted `json:"hideOnLoginPage,omitempty"`
}

type SamlIdentityProvider struct {
	Realm                       string                      `json:"-"`
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
	Config                      *SamlIdentityProviderConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewSamlIdentityProvider(samlIdentityProvider *SamlIdentityProvider) error {
	log.Printf("[WARN] Realm: %s", samlIdentityProvider.Realm)
	_, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances", samlIdentityProvider.Realm), samlIdentityProvider)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetSamlIdentityProvider(realm, alias string) (*SamlIdentityProvider, error) {
	var samlIdentityProvider *SamlIdentityProvider
	samlIdentityProvider.Realm = realm

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias), samlIdentityProvider)
	if err != nil {
		return nil, err
	}

	return samlIdentityProvider, nil
}

func (keycloakClient *KeycloakClient) UpdateSamlIdentityProvider(samlIdentityProvider *SamlIdentityProvider) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", samlIdentityProvider.Realm, samlIdentityProvider.Alias), samlIdentityProvider)
}

func (keycloakClient *KeycloakClient) DeleteSamlIdentityProvider(realm, alias string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias))
}
