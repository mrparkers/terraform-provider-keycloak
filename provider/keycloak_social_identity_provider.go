package provider

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"log"
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
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Realm ID",
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
				Optional:    true,
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
		},
	}
}

func getSocialIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec := &keycloak.IdentityProvider{
		Id:                        data.Id(),
		RealmId:                   data.Get("realm_id").(string),
		Alias:                     data.Get("alias").(string),
		DisplayName:               data.Get("display_name").(string),
		ProviderId:                data.Get("provider_id").(string),
		Enabled:                   data.Get("enabled").(bool),
		StoreToken:                data.Get("store_token").(bool),
		AddReadTokenRoleOnCreate:  data.Get("add_read_token_role_on_create").(bool),
		AuthenticateByDefault:     data.Get("authenticate_by_default").(bool),
		LinkOnly:                  data.Get("link_only").(bool),
		TrustEmail:                data.Get("trust_email").(bool),
		FirstBrokerLoginFlowAlias: data.Get("first_broker_login_flow_alias").(string),
		PostBrokerLoginFlowAlias:  data.Get("post_broker_login_flow_alias").(string),
	}
	if v, ok := data.GetOk("config"); ok {
		configs := v.(*schema.Set).List()
		if len(configs) > 1 {
			return nil, fmt.Errorf("You can only define a single alias target per record")
		}
		config := configs[0].(map[string]interface{})
		rec.Config = &keycloak.IdentityProviderConfig{
			UseJwksUrl:      config["use_jwks_url"].(bool),
			ClientId:        config["client_id"].(string),
			ClientSecret:    config["client_secret"].(string),
			DisableUserInfo: config["disable_user_info"].(bool),
			HideOnLoginPage: config["hide_on_login_page"].(bool),
			Key:             config["key"].(string),
		}
		log.Printf("[DEBUG] Creating config: %#v", config)
	} else {
		return nil, fmt.Errorf("No config is defined")
	}
	return rec, nil
}

func setSocialIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) {
	data.SetId(identityProvider.Id)
	data.Set("realm_id", identityProvider.RealmId)
	data.Set("alias", identityProvider.Alias)
	data.Set("display_name", identityProvider.DisplayName)
	data.Set("provider_id", identityProvider.ProviderId)
	data.Set("enabled", identityProvider.Enabled)
	data.Set("store_token", identityProvider.StoreToken)
	data.Set("add_read_token_role_on_create", identityProvider.AddReadTokenRoleOnCreate)
	data.Set("authenticate_by_default", identityProvider.AuthenticateByDefault)
	data.Set("link_only", identityProvider.LinkOnly)
	data.Set("trust_email", identityProvider.TrustEmail)
	data.Set("first_broker_login_flow_alias", identityProvider.FirstBrokerLoginFlowAlias)
	data.Set("post_broker_login_flow_alias", identityProvider.PostBrokerLoginFlowAlias)
	if config := identityProvider.Config; config != nil {
		data.Set("config", []interface{}{
			map[string]interface{}{
				"use_jwks_url":       config.UseJwksUrl,
				"client_id":          config.ClientId,
				"client_secret":      config.ClientSecret,
				"disable_user_info":  config.DisableUserInfo,
				"hide_on_login_page": config.HideOnLoginPage,
				"key":                config.Key,
				"host_ip":            config.HostIp,
			},
		})
	}
}

func resourceKeycloakSocialIdentityProviderCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	identityProvider, err := getSocialIdentityProviderFromData(data)

	err = keycloakClient.NewIdentityProvider(identityProvider)
	if err != nil {
		return err
	}

	setSocialIdentityProviderData(data, identityProvider)

	return resourceKeycloakSocialIdentityProviderRead(data, meta)
}

func resourceKeycloakSocialIdentityProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	identityProvider, err := keycloakClient.GetIdentityProvider(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setSocialIdentityProviderData(data, identityProvider)

	return nil
}

func resourceKeycloakSocialIdentityProviderUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	identityProvider, err := getSocialIdentityProviderFromData(data)

	err = keycloakClient.UpdateIdentityProvider(identityProvider)
	if err != nil {
		return err
	}

	setSocialIdentityProviderData(data, identityProvider)

	return nil
}

func resourceKeycloakSocialIdentityProviderDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteIdentityProvider(realmId, id)
}

func resourceKeycloakSocialIdentityProviderImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	id := parts[1]

	d.Set("realm_id", realm)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
