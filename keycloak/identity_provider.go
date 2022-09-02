package keycloak

import (
	"context"
	"fmt"
	"reflect"
)

type IdentityProviderConfig struct {
	Key                             string                 `json:"key,omitempty"`
	HostIp                          string                 `json:"hostIp,omitempty"`
	UseJwksUrl                      KeycloakBoolQuoted     `json:"useJwksUrl,omitempty"`
	JwksUrl                         string                 `json:"jwksUrl,omitempty"`
	ClientId                        string                 `json:"clientId,omitempty"`
	ClientSecret                    string                 `json:"clientSecret,omitempty"`
	DisableUserInfo                 KeycloakBoolQuoted     `json:"disableUserInfo"`
	UserInfoUrl                     string                 `json:"userInfoUrl,omitempty"`
	HideOnLoginPage                 KeycloakBoolQuoted     `json:"hideOnLoginPage"`
	NameIDPolicyFormat              string                 `json:"nameIDPolicyFormat,omitempty"`
	EntityId                        string                 `json:"entityId,omitempty"`
	SingleLogoutServiceUrl          string                 `json:"singleLogoutServiceUrl,omitempty"`
	SingleSignOnServiceUrl          string                 `json:"singleSignOnServiceUrl,omitempty"`
	SigningCertificate              string                 `json:"signingCertificate,omitempty"`
	SignatureAlgorithm              string                 `json:"signatureAlgorithm,omitempty"`
	XmlSigKeyInfoKeyNameTransformer string                 `json:"xmlSigKeyInfoKeyNameTransformer,omitempty"`
	PostBindingAuthnRequest         KeycloakBoolQuoted     `json:"postBindingAuthnRequest,omitempty"`
	PostBindingResponse             KeycloakBoolQuoted     `json:"postBindingResponse,omitempty"`
	PostBindingLogout               KeycloakBoolQuoted     `json:"postBindingLogout,omitempty"`
	ForceAuthn                      KeycloakBoolQuoted     `json:"forceAuthn,omitempty"`
	WantAuthnRequestsSigned         KeycloakBoolQuoted     `json:"wantAuthnRequestsSigned,omitempty"`
	WantAssertionsSigned            KeycloakBoolQuoted     `json:"wantAssertionsSigned,omitempty"`
	WantAssertionsEncrypted         KeycloakBoolQuoted     `json:"wantAssertionsEncrypted,omitempty"`
	BackchannelSupported            KeycloakBoolQuoted     `json:"backchannelSupported,omitempty"`
	ValidateSignature               KeycloakBoolQuoted     `json:"validateSignature,omitempty"`
	AuthorizationUrl                string                 `json:"authorizationUrl,omitempty"`
	TokenUrl                        string                 `json:"tokenUrl,omitempty"`
	LoginHint                       string                 `json:"loginHint,omitempty"`
	UILocales                       KeycloakBoolQuoted     `json:"uiLocales,omitempty"`
	LogoutUrl                       string                 `json:"logoutUrl,omitempty"`
	DefaultScope                    string                 `json:"defaultScope,omitempty"`
	AcceptsPromptNoneForwFrmClt     KeycloakBoolQuoted     `json:"acceptsPromptNoneForwardFromClient,omitempty"`
	HostedDomain                    string                 `json:"hostedDomain,omitempty"`
	UserIp                          KeycloakBoolQuoted     `json:"userIp,omitempty"`
	OfflineAccess                   KeycloakBoolQuoted     `json:"offlineAccess,omitempty"`
	PrincipalType                   string                 `json:"principalType,omitempty"`
	PrincipalAttribute              string                 `json:"principalAttribute,omitempty"`
	GuiOrder                        string                 `json:"guiOrder,omitempty"`
	SyncMode                        string                 `json:"syncMode,omitempty"`
	ExtraConfig                     map[string]interface{} `json:"-"`
	AuthnContextClassRefs           KeycloakSliceQuoted    `json:"authnContextClassRefs,omitempty"`
	AuthnContextComparisonType      string                 `json:"authnContextComparisonType,omitempty"`
	AuthnContextDeclRefs            KeycloakSliceQuoted    `json:"authnContextDeclRefs,omitempty"`
}

type IdentityProvider struct {
	Realm                     string                  `json:"-"`
	InternalId                string                  `json:"internalId,omitempty"`
	Alias                     string                  `json:"alias"`
	DisplayName               string                  `json:"displayName"`
	ProviderId                string                  `json:"providerId"`
	Enabled                   bool                    `json:"enabled"`
	StoreToken                bool                    `json:"storeToken"`
	AddReadTokenRoleOnCreate  bool                    `json:"addReadTokenRoleOnCreate"`
	AuthenticateByDefault     bool                    `json:"authenticateByDefault"`
	LinkOnly                  bool                    `json:"linkOnly"`
	TrustEmail                bool                    `json:"trustEmail"`
	FirstBrokerLoginFlowAlias string                  `json:"firstBrokerLoginFlowAlias"`
	PostBrokerLoginFlowAlias  string                  `json:"postBrokerLoginFlowAlias"`
	Config                    *IdentityProviderConfig `json:"config"`
}

func (keycloakClient *KeycloakClient) NewIdentityProvider(ctx context.Context, identityProvider *IdentityProvider) error {
	_, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances", identityProvider.Realm), identityProvider)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetIdentityProvider(ctx context.Context, realm, alias string) (*IdentityProvider, error) {
	var identityProvider IdentityProvider
	identityProvider.Realm = realm

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias), &identityProvider, nil)
	if err != nil {
		return nil, err
	}

	return &identityProvider, nil
}

func (keycloakClient *KeycloakClient) UpdateIdentityProvider(ctx context.Context, identityProvider *IdentityProvider) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s", identityProvider.Realm, identityProvider.Alias), identityProvider)
}

func (keycloakClient *KeycloakClient) DeleteIdentityProvider(ctx context.Context, realm, alias string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s", realm, alias), nil)
}

func (f *IdentityProviderConfig) UnmarshalJSON(data []byte) error {
	return unmarshalExtraConfig(data, reflect.ValueOf(f).Elem(), &f.ExtraConfig)
}

func (f *IdentityProviderConfig) MarshalJSON() ([]byte, error) {
	return marshalExtraConfig(reflect.ValueOf(f).Elem(), f.ExtraConfig)
}
