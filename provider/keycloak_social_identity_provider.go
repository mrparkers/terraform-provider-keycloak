package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakSocialIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSocialIdentityProviderCreate,
		Read:   resourceKeycloakSocialIdentityProviderRead,
		Update: resourceKeycloakSocialIdentityProviderUpdate,
		Delete: resourceKeycloakSocialIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakSocialIdentityProviderImport,
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
			"host_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Google Host IP",
			},
			"use_jwks_url": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Use JWKS url",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "StackOverFlow key.",
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
				Description: "Disable User Info.",
			},
			"hide_on_login_page": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Hide On Login Page.",
			},
		},
	}
}

func getSocialIdentityProviderFromData(data *schema.ResourceData) (*keycloak.SocialIdentityProvider, error) {
	rec := &keycloak.SocialIdentityProvider{
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
	rec.Config = &keycloak.SocialIdentityProviderConfig{
		UseJwksUrl:      keycloak.KeycloakBoolQuoted(data.Get("use_jwks_url").(bool)),
		ClientId:        data.Get("client_id").(string),
		ClientSecret:    data.Get("client_secret").(string),
		DisableUserInfo: keycloak.KeycloakBoolQuoted(data.Get("disable_user_info").(bool)),
		HideOnLoginPage: keycloak.KeycloakBoolQuoted(data.Get("hide_on_login_page").(bool)),
		Key:             data.Get("key").(string),
	}
	return rec, nil
}

func setSocialIdentityProviderData(data *schema.ResourceData, socialIdentityProvider *keycloak.SocialIdentityProvider) {
	data.SetId(socialIdentityProvider.Realm + "/" + socialIdentityProvider.Alias)
	data.Set("internal_id", socialIdentityProvider.InternalId)
	data.Set("realm", socialIdentityProvider.Realm)
	data.Set("alias", socialIdentityProvider.Alias)
	data.Set("display_name", socialIdentityProvider.DisplayName)
	data.Set("provider_id", socialIdentityProvider.ProviderId)
	data.Set("enabled", socialIdentityProvider.Enabled)
	data.Set("store_token", socialIdentityProvider.StoreToken)
	data.Set("add_read_token_role_on_create", socialIdentityProvider.AddReadTokenRoleOnCreate)
	data.Set("authenticate_by_default", socialIdentityProvider.AuthenticateByDefault)
	data.Set("link_only", socialIdentityProvider.LinkOnly)
	data.Set("trust_email", socialIdentityProvider.TrustEmail)
	data.Set("first_broker_login_flow_alias", socialIdentityProvider.FirstBrokerLoginFlowAlias)
	data.Set("post_broker_login_flow_alias", socialIdentityProvider.PostBrokerLoginFlowAlias)
	if config := socialIdentityProvider.Config; config != nil {
		data.Set("use_jwks_url", config.UseJwksUrl)
		data.Set("client_id", config.ClientId)
		data.Set("client_secret", config.ClientSecret)
		data.Set("disable_user_info", config.DisableUserInfo)
		data.Set("hide_on_login_page", config.HideOnLoginPage)
		data.Set("key", config.Key)
		data.Set("host_ip", config.HostIp)
	}
}

func resourceKeycloakSocialIdentityProviderCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	socialIdentityProvider, err := getSocialIdentityProviderFromData(data)

	err = keycloakClient.NewSocialIdentityProvider(socialIdentityProvider)
	if err != nil {
		return err
	}

	setSocialIdentityProviderData(data, socialIdentityProvider)

	return resourceKeycloakSocialIdentityProviderRead(data, meta)
}

func resourceKeycloakSocialIdentityProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	socialIdentityProvider, err := keycloakClient.GetSocialIdentityProvider(realm, alias)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setSocialIdentityProviderData(data, socialIdentityProvider)

	return nil
}

func resourceKeycloakSocialIdentityProviderUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	socialIdentityProvider, err := getSocialIdentityProviderFromData(data)

	err = keycloakClient.UpdateSocialIdentityProvider(socialIdentityProvider)
	if err != nil {
		return err
	}

	setSocialIdentityProviderData(data, socialIdentityProvider)

	return nil
}

func resourceKeycloakSocialIdentityProviderDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	return keycloakClient.DeleteSocialIdentityProvider(realm, alias)
}

func resourceKeycloakSocialIdentityProviderImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	alias := parts[1]

	d.Set("realm", realm)
	d.Set("alias", alias)

	return []*schema.ResourceData{d}, nil
}
