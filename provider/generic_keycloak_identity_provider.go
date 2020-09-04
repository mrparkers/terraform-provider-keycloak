package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

type identityProviderDataGetterFunc func(data *schema.ResourceData) (*keycloak.IdentityProvider, error)
type identityProviderDataSetterFunc func(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error

func resourceKeycloakIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Delete: resourceKeycloakIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakIdentityProviderImport,
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
			"internal_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal Identity Provider Id",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Friendly name for Identity Providers.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
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
		},
	}
}

func getIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec := &keycloak.IdentityProvider{
		Realm:                     data.Get("realm").(string),
		Alias:                     data.Get("alias").(string),
		DisplayName:               data.Get("display_name").(string),
		Enabled:                   data.Get("enabled").(bool),
		StoreToken:                data.Get("store_token").(bool),
		AddReadTokenRoleOnCreate:  data.Get("add_read_token_role_on_create").(bool),
		AuthenticateByDefault:     data.Get("authenticate_by_default").(bool),
		LinkOnly:                  data.Get("link_only").(bool),
		TrustEmail:                data.Get("trust_email").(bool),
		FirstBrokerLoginFlowAlias: data.Get("first_broker_login_flow_alias").(string),
		PostBrokerLoginFlowAlias:  data.Get("post_broker_login_flow_alias").(string),
		InternalId:                data.Get("internal_id").(string),
	}
	return rec, nil
}

func setIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error {
	data.SetId(identityProvider.Alias)
	data.Set("internal_id", identityProvider.InternalId)
	data.Set("realm", identityProvider.Realm)
	data.Set("alias", identityProvider.Alias)
	data.Set("display_name", identityProvider.DisplayName)
	data.Set("enabled", identityProvider.Enabled)
	data.Set("store_token", identityProvider.StoreToken)
	data.Set("add_read_token_role_on_create", identityProvider.AddReadTokenRoleOnCreate)
	data.Set("authenticate_by_default", identityProvider.AuthenticateByDefault)
	data.Set("link_only", identityProvider.LinkOnly)
	data.Set("trust_email", identityProvider.TrustEmail)
	data.Set("first_broker_login_flow_alias", identityProvider.FirstBrokerLoginFlowAlias)
	data.Set("post_broker_login_flow_alias", identityProvider.PostBrokerLoginFlowAlias)
	return nil
}

func resourceKeycloakIdentityProviderDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	return keycloakClient.DeleteIdentityProvider(realm, alias)
}

func resourceKeycloakIdentityProviderImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realm}}/{{identityProviderAlias}}")
	}

	d.Set("realm", parts[0])
	d.Set("alias", parts[1])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceKeycloakIdentityProviderCreate(getIdentityProviderFromData identityProviderDataGetterFunc, setDataFromIdentityProvider identityProviderDataSetterFunc) func(data *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		identityProvider, err := getIdentityProviderFromData(data)
		if err = keycloakClient.NewIdentityProvider(identityProvider); err != nil {
			return err
		}
		if err = setDataFromIdentityProvider(data, identityProvider); err != nil {
			return err
		}
		return resourceKeycloakIdentityProviderRead(setDataFromIdentityProvider)(data, meta)
	}
}

func resourceKeycloakIdentityProviderRead(setDataFromIdentityProvider identityProviderDataSetterFunc) func(data *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		realm := data.Get("realm").(string)
		alias := data.Get("alias").(string)
		identityProvider, err := keycloakClient.GetIdentityProvider(realm, alias)
		if err != nil {
			return handleNotFoundError(err, data)
		}
		if err = setDataFromIdentityProvider(data, identityProvider); err != nil {
			return err
		}
		return nil
	}
}

func resourceKeycloakIdentityProviderUpdate(getIdentityProviderFromData identityProviderDataGetterFunc, setDataFromIdentityProvider identityProviderDataSetterFunc) func(data *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		identityProvider, err := getIdentityProviderFromData(data)
		if err = keycloakClient.UpdateIdentityProvider(identityProvider); err != nil {
			return err
		}
		if err = setDataFromIdentityProvider(data, identityProvider); err != nil {
			return err
		}
		return nil
	}
}
