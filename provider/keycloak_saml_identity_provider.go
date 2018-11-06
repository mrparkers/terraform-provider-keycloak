package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakSamlIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlIdentityProviderCreate,
		Read:   resourceKeycloakSamlIdentityProviderRead,
		Update: resourceKeycloakSamlIdentityProviderUpdate,
		Delete: resourceKeycloakSamlIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakSamlIdentityProviderImport,
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The alias uniquely identifies an identity provider and it is also used to build the redirect uri.",
			},
			"realm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Realm Name",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Friendly name for Identity Providers.",
			},
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Provider ID.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable/disable this identity provider.",
			},
			"store_token": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable/disable if tokens must be stored after authenticating users.",
			},
			"add_read_token_role_on_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/disable if new users can read any stored tokens. This assigns the broker.read-token role.",
			},
			"authenticate_by_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable/disable authenticate users by default.",
			},
			"link_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, users cannot log in through this provider.  They can only link to this provider.  This is useful if you don't want to allow login from the provider, but want to integrate with a provider",
			},
			"trust_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If enabled then email provided by this provider is not verified even if verification is enabled for the realm.",
			},
			"first_broker_login_flow_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "first broker login",
				Description: "Alias of authentication flow, which is triggered after first login with this identity provider. Term 'First Login' means that there is not yet existing Keycloak account linked with the authenticated identity provider account.",
			},
			"post_broker_login_flow_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alias of authentication flow, which is triggered after each login with this identity provider. Useful if you want additional verification of each user authenticated with this identity provider (for example OTP). Leave this empty if you don't want any additional authenticators to be triggered after login with this identity provider. Also note, that authenticator implementations must assume that user is already set in ClientSession as identity provider already set it.",
			},
			"backchannel_supported": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Does the external IDP support backchannel logout?",
			},
			"use_jwks_url": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Use JWKS url",
			},
			"validate_signature": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/disable signature validation of SAML responses.",
			},
			"hide_on_login_page": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Hide On Login Page.",
			},
			"name_id_policy_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
				Description: "Name ID Policy Format.",
			},
			"single_logout_service_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Logout URL.",
			},
			"single_sign_on_service_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSO Logout URL.",
			},
			"signing_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Signing Certificate.",
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "RSA_SHA256",
				Description: "Signing Algorithm.",
			},
			"xml_sign_key_info_key_name_transformer": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "KEY_ID",
				Description: "Sign Key Transformer.",
			},
			"post_binding_authn_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Post Binding Authn Request.",
			},
			"post_binding_response": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Post Binding Response.",
			},
			"post_binding_logout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Post Binding Logout.",
			},
			"force_authn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require Force Authn.",
			},
			"want_authn_requests_signed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require Force Authn Requests Sign.",
			},
			"want_assertions_signed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Want Assertions Signed.",
			},
			"want_assertions_encrypted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Want Assertions Encrypted.",
			},
		},
	}
}

func getSamlIdentityProviderFromData(data *schema.ResourceData) (*keycloak.SamlIdentityProvider, error) {
	rec := &keycloak.SamlIdentityProvider{
		InternalId:                data.Id(),
		Realm:                     data.Get("realm").(string),
		Alias:                     data.Get("alias").(string),
		DisplayName:               data.Get("display_name").(string),
		ProviderId:                data.Get("provider_id").(string),
		Enabled:                   data.Get("enabled").(bool),
		StoreToken:                keycloak.KeycloakBool(data.Get("store_token").(bool)),
		AddReadTokenRoleOnCreate:  keycloak.KeycloakBool(data.Get("add_read_token_role_on_create").(bool)),
		AuthenticateByDefault:     data.Get("authenticate_by_default").(bool),
		LinkOnly:                  keycloak.KeycloakBool(data.Get("link_only").(bool)),
		TrustEmail:                keycloak.KeycloakBool(data.Get("trust_email").(bool)),
		FirstBrokerLoginFlowAlias: data.Get("first_broker_login_flow_alias").(string),
		PostBrokerLoginFlowAlias:  data.Get("post_broker_login_flow_alias").(string),
	}
	rec.Config = &keycloak.SamlIdentityProviderConfig{
		UseJwksUrl:                       keycloak.KeycloakBoolQuoted(data.Get("use_jwks_url").(bool)),
		ValidateSignature:                keycloak.KeycloakBoolQuoted(data.Get("validate_signature").(bool)),
		HideOnLoginPage:                  keycloak.KeycloakBoolQuoted(data.Get("hide_on_login_page").(bool)),
		NameIDPolicyFormat:               data.Get("name_id_policy_format").(string),
		SingleLogutServiceUrl:            data.Get("single_logout_service_url").(string),
		SingleSignOnServiceUrl:           data.Get("single_sign_on_service_url").(string),
		SigningCertificate:               data.Get("signing_certificate").(string),
		SignatureAlgorithm:               data.Get("signature_algorithm").(string),
		XmlSignKeyInfoKeyNameTransformer: data.Get("xml_sign_key_info_key_name_transformer").(string),
		PostBindingAuthnRequest:          keycloak.KeycloakBoolQuoted(data.Get("post_binding_authn_request").(bool)),
		PostBindingResponse:              keycloak.KeycloakBoolQuoted(data.Get("post_binding_response").(bool)),
		PostBindingLogout:                keycloak.KeycloakBoolQuoted(data.Get("post_binding_logout").(bool)),
		ForceAuthn:                       keycloak.KeycloakBoolQuoted(data.Get("force_authn").(bool)),
		WantAuthnRequestsSigned:          keycloak.KeycloakBoolQuoted(data.Get("want_authn_requests_signed").(bool)),
		WantAssertionsSigned:             keycloak.KeycloakBoolQuoted(data.Get("want_assertions_signed").(bool)),
		WantAssertionsEncrypted:          keycloak.KeycloakBoolQuoted(data.Get("want_assertions_encrypted").(bool)),
	}
	return rec, nil
}

func setSamlIdentityProviderData(data *schema.ResourceData, samlIdentityProvider *keycloak.SamlIdentityProvider) {
	data.SetId(samlIdentityProvider.Realm + "/" + samlIdentityProvider.Alias)
	data.Set("internal_id", samlIdentityProvider.InternalId)
	data.Set("realm", samlIdentityProvider.Realm)
	data.Set("alias", samlIdentityProvider.Alias)
	data.Set("display_name", samlIdentityProvider.DisplayName)
	data.Set("provider_id", samlIdentityProvider.ProviderId)
	data.Set("enabled", samlIdentityProvider.Enabled)
	data.Set("store_token", samlIdentityProvider.StoreToken)
	data.Set("add_read_token_role_on_create", samlIdentityProvider.AddReadTokenRoleOnCreate)
	data.Set("authenticate_by_default", samlIdentityProvider.AuthenticateByDefault)
	data.Set("link_only", samlIdentityProvider.LinkOnly)
	data.Set("trust_email", samlIdentityProvider.TrustEmail)
	data.Set("first_broker_login_flow_alias", samlIdentityProvider.FirstBrokerLoginFlowAlias)
	data.Set("post_broker_login_flow_alias", samlIdentityProvider.PostBrokerLoginFlowAlias)
	if config := samlIdentityProvider.Config; config != nil {
		data.Set("backchannel_supported", config.BackchannelSupported)
		data.Set("use_jwks_url", config.UseJwksUrl)
		data.Set("validate_signature", config.ValidateSignature)
		data.Set("hide_on_login_page", config.HideOnLoginPage)
		data.Set("name_id_policy_format", config.NameIDPolicyFormat)
		data.Set("single_logout_service_url", config.SingleLogutServiceUrl)
		data.Set("single_sign_on_service_url", config.SingleSignOnServiceUrl)
		data.Set("signing_certificate", config.SigningCertificate)
		data.Set("signature_algorithm", config.SignatureAlgorithm)
		data.Set("xml_sign_key_info_key_name_transformer", config.XmlSignKeyInfoKeyNameTransformer)
		data.Set("post_binding_authn_request", config.PostBindingAuthnRequest)
		data.Set("post_binding_response", config.PostBindingResponse)
		data.Set("post_binding_logout", config.PostBindingLogout)
		data.Set("force_authn", config.ForceAuthn)
		data.Set("want_authn_requests_signed", config.WantAuthnRequestsSigned)
		data.Set("want_assertions_signed", config.WantAssertionsSigned)
		data.Set("want_assertions_encrypted", config.WantAssertionsEncrypted)
	}
}

func resourceKeycloakSamlIdentityProviderCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlIdentityProvider, err := getSamlIdentityProviderFromData(data)

	err = keycloakClient.NewSamlIdentityProvider(samlIdentityProvider)
	if err != nil {
		return err
	}

	setSamlIdentityProviderData(data, samlIdentityProvider)

	return resourceKeycloakSamlIdentityProviderRead(data, meta)
}

func resourceKeycloakSamlIdentityProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	samlIdentityProvider, err := keycloakClient.GetSamlIdentityProvider(realm, alias)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setSamlIdentityProviderData(data, samlIdentityProvider)

	return nil
}

func resourceKeycloakSamlIdentityProviderUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlIdentityProvider, err := getSamlIdentityProviderFromData(data)

	err = keycloakClient.UpdateSamlIdentityProvider(samlIdentityProvider)
	if err != nil {
		return err
	}

	setSamlIdentityProviderData(data, samlIdentityProvider)

	return nil
}

func resourceKeycloakSamlIdentityProviderDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	return keycloakClient.DeleteSamlIdentityProvider(realm, alias)
}

func resourceKeycloakSamlIdentityProviderImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	alias := parts[1]

	d.Set("realm", realm)
	d.Set("alias", alias)

	return []*schema.ResourceData{d}, nil
}
