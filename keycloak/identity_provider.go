package keycloak

import (
	"fmt"
)

type IdentityProviderConfig struct {
	BaseUrl                          string `json:"baseUrl,omitempty"`
	BackchannelSupported             bool   `json:"backchannelSupported,omitempty"`
	UseJwksUrl                       bool   `json:"useJwksUrl,omitempty"`
	ValidateSignature                bool   `json:"validateSignature,omitempty"`
	AuthorizationUrl                 string `json:"authorizationUrl,omitempty"`
	ClientId                         string `json:"clientId,omitempty"`
	ClientSecret                     string `json:"clientSecret,omitempty"`
	DisableUserInfo                  string `json:"disableUserInfo,omitempty"`
	HideOnLoginPage                  string `json:"hideOnLoginPage,omitempty"`
	TokenUrl                         string `json:"tokeUrl,omitempty"`
	LoginHint                        string `json:"loginHint,omitempty"`
	NameIDPolicyFormat               string `json:"nameIDPolicyFormat,omitempty"`
	SingleLogutServiceUrl            string `json:"singleLogoutServiceUrl,omitempty"`
	SingleSignOnServiceUrl           string `json:"singleSignOnServiceUrl,omitempty"`
	SigningCertificate               string `json:"signingCertificate,omitempty"`
	SignatureAlgorithm               string `json:"signatureAlgorithm,omitempty"`
	XmlSignKeyInfoKeyNameTransformer string `json:"xmlSignKeyInfoKeyNameTransformer,omitempty"`
	PostBindingAuthnRequest          string `json:"postBindingAuthnRequest,omitempty"`
	PostBindingResponse              string `json:"postBindingResponse,omitempty"`
	PostBindingLogout                string `json:"postBindingLogout,omitempty"`
	ForceAuthn                       bool   `json:"forceAuthn,omitempty"`
	WantAuthnRequestsSigned          bool   `json:"wantAuthnRequestsSigned,omitempty"`
	WantAssertionsSigned             bool   `json:"wantAssertionsSigned,omitempty"`
	WantAssertionsEncrypted          bool   `json:"wantAssertionsEncrypted,omitempty"`
}

type IdentityProvider struct {
	Id                        string                  `json:"-"`
	RealmId                   string                  `json:"-"`
	Alias                     string                  `json:"alias,omitempty"`
	DisplayName               string                  `json:"displayName,omitempty"`
	ProviderId                string                  `json:"providerId,omitempty"`
	Enabled                   bool                    `json:"enabled,omitempty"`
	StoreToken                bool                    `json:"storeToken,omitempty"`
	AddReadTokenRoleOnCreate  bool                    `json:"addReadTokenRoleOnCreate,omitempty"`
	AuthenticateByDefault     bool                    `json:"authenticateByDefault,omitempty"`
	LinkOnly                  bool                    `json:"linkOnly,omitempty"`
	TrustEmail                bool                    `json:"trustEmail,omitempty"`
	FirstBrokerLoginFlowAlias string                  `json:"firstBrokerLoginFlowAlias,omitempty"`
	PostBrokerLoginFlowAlias  string                  `json:"postBrokerLoginFlowAlias,omitempty"`
	Config                    *IdentityProviderConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewIdentityProvider(identityProvider *IdentityProvider) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances", identityProvider.RealmId), identityProvider)
	if err != nil {
		return err
	}

	identityProvider.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetIdentityProvider(realmId, alias string) (*IdentityProvider, error) {
	var identityProvider *IdentityProvider

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realmId, alias), identityProvider)
	if err != nil {
		return nil, err
	}

	return identityProvider, nil
}

func (keycloakClient *KeycloakClient) UpdateIdentityProvider(identityProvider *IdentityProvider) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", identityProvider.RealmId, identityProvider.Alias), identityProvider)
}

func (keycloakClient *KeycloakClient) DeleteIdentityProvider(realmId, alias string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realmId, alias))
}
