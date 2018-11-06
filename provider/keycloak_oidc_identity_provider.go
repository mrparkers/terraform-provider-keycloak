package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOidcIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOidcIdentityProviderCreate,
		Read:   resourceKeycloakOidcIdentityProviderRead,
		Update: resourceKeycloakOidcIdentityProviderUpdate,
		Delete: resourceKeycloakOidcIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOidcIdentityProviderImport,
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
				Default:     false,
				ForceNew:    true,
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
				Default:     false,
				Description: "If true, users cannot log in through this provider.  They can only link to this provider.  This is useful if you don't want to allow login from the provider, but want to integrate with a provider",
			},
			"trust_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
				Default:     "",
				Description: "Alias of authentication flow, which is triggered after each login with this identity provider. Useful if you want additional verification of each user authenticated with this identity provider (for example OTP). Leave this empty if you don't want any additional authenticators to be triggered after login with this identity provider. Also note, that authenticator implementations must assume that user is already set in ClientSession as identity provider already set it.",
			},
			"backchannel_supported": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
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
				Default:     false,
				Description: "Enable/disable signature validation of SAML responses.",
			},
			"authorization_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OIDC authorization URL.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client ID.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client Secret.",
			},
			"disable_user_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable User Info.",
			},
			"hide_on_login_page": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Hide On Login Page.",
			},
			"token_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Token URL.",
			},
			"login_hint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login Hint.",
			},
		},
	}
}

func getOidcIdentityProviderFromData(data *schema.ResourceData) (*keycloak.OidcIdentityProvider, error) {
	rec := &keycloak.OidcIdentityProvider{
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
	rec.Config = &keycloak.OidcIdentityProviderConfig{
		BackchannelSupported: keycloak.KeycloakBoolQuoted(data.Get("backchannel_supported").(bool)),
		UseJwksUrl:           keycloak.KeycloakBoolQuoted(data.Get("use_jwks_url").(bool)),
		ValidateSignature:    keycloak.KeycloakBoolQuoted(data.Get("validate_signature").(bool)),
		AuthorizationUrl:     data.Get("authorization_url").(string),
		ClientId:             data.Get("client_id").(string),
		ClientSecret:         data.Get("client_secret").(string),
		DisableUserInfo:      keycloak.KeycloakBoolQuoted(data.Get("disable_user_info").(bool)),
		HideOnLoginPage:      keycloak.KeycloakBoolQuoted(data.Get("hide_on_login_page").(bool)),
		TokenUrl:             data.Get("token_url").(string),
		LoginHint:            data.Get("login_hint").(string),
	}
	return rec, nil
}

func setOidcIdentityProviderData(data *schema.ResourceData, oidcIdentityProvider *keycloak.OidcIdentityProvider) {
	data.SetId(oidcIdentityProvider.Realm + "/" + oidcIdentityProvider.Alias)
	data.Set("internal_id", oidcIdentityProvider.InternalId)
	data.Set("realm", oidcIdentityProvider.Realm)
	data.Set("alias", oidcIdentityProvider.Alias)
	data.Set("display_name", oidcIdentityProvider.DisplayName)
	data.Set("provider_id", oidcIdentityProvider.ProviderId)
	data.Set("enabled", oidcIdentityProvider.Enabled)
	data.Set("store_token", oidcIdentityProvider.StoreToken)
	data.Set("add_read_token_role_on_create", oidcIdentityProvider.AddReadTokenRoleOnCreate)
	data.Set("authenticate_by_default", oidcIdentityProvider.AuthenticateByDefault)
	data.Set("link_only", oidcIdentityProvider.LinkOnly)
	data.Set("trust_email", oidcIdentityProvider.TrustEmail)
	data.Set("first_broker_login_flow_alias", oidcIdentityProvider.FirstBrokerLoginFlowAlias)
	data.Set("post_broker_login_flow_alias", oidcIdentityProvider.PostBrokerLoginFlowAlias)
	if config := oidcIdentityProvider.Config; config != nil {
		data.Set("backchannel_supported", config.BackchannelSupported)
		data.Set("use_jwks_url", config.UseJwksUrl)
		data.Set("validate_signature", config.ValidateSignature)
		data.Set("authorization_url", config.AuthorizationUrl)
		data.Set("client_id", config.ClientId)
		data.Set("client_secret", config.ClientSecret)
		data.Set("disable_user_info", config.DisableUserInfo)
		data.Set("hide_on_login_page", config.HideOnLoginPage)
		data.Set("token_url", config.TokenUrl)
		data.Set("login_hint", config.LoginHint)
	}
}

func resourceKeycloakOidcIdentityProviderCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	oidcIdentityProvider, err := getOidcIdentityProviderFromData(data)

	err = keycloakClient.NewOidcIdentityProvider(oidcIdentityProvider)
	if err != nil {
		return err
	}

	setOidcIdentityProviderData(data, oidcIdentityProvider)

	return resourceKeycloakOidcIdentityProviderRead(data, meta)
}

func resourceKeycloakOidcIdentityProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	oidcIdentityProvider, err := keycloakClient.GetOidcIdentityProvider(realm, alias)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOidcIdentityProviderData(data, oidcIdentityProvider)

	return nil
}

func resourceKeycloakOidcIdentityProviderUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	oidcIdentityProvider, err := getOidcIdentityProviderFromData(data)

	err = keycloakClient.UpdateOidcIdentityProvider(oidcIdentityProvider)
	if err != nil {
		return err
	}

	setOidcIdentityProviderData(data, oidcIdentityProvider)

	return nil
}

func resourceKeycloakOidcIdentityProviderDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	return keycloakClient.DeleteOidcIdentityProvider(realm, alias)
}

func resourceKeycloakOidcIdentityProviderImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	alias := parts[1]

	d.Set("realm", realm)
	d.Set("alias", alias)

	return []*schema.ResourceData{d}, nil
}
