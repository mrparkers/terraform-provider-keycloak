package provider

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakIdentityProviderCreate,
		Read:   resourceKeycloakIdentityProviderRead,
		Update: resourceKeycloakIdentityProviderUpdate,
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
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Friendly name for Identity Providers.",
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
			"oidc": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"saml", "social"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
				Set: configHash,
			},
			"saml": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"oidc", "social"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
				Set: configHash,
			},
			"social": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"oidc", "saml"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "github",
							Description: "IDP name",
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
				},
				Set: configHash,
			},
		},
	}
}

func getIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec := &keycloak.IdentityProvider{
		InternalId:                data.Id(),
		Realm:                     data.Get("realm").(string),
		Alias:                     data.Get("alias").(string),
		DisplayName:               data.Get("display_name").(string),
		Enabled:                   data.Get("enabled").(bool),
		StoreToken:                keycloak.KeycloakBool(data.Get("store_token").(bool)),
		AddReadTokenRoleOnCreate:  keycloak.KeycloakBool(data.Get("add_read_token_role_on_create").(bool)),
		AuthenticateByDefault:     data.Get("authenticate_by_default").(bool),
		LinkOnly:                  keycloak.KeycloakBool(data.Get("link_only").(bool)),
		TrustEmail:                keycloak.KeycloakBool(data.Get("trust_email").(bool)),
		FirstBrokerLoginFlowAlias: data.Get("first_broker_login_flow_alias").(string),
		PostBrokerLoginFlowAlias:  data.Get("post_broker_login_flow_alias").(string),
	}
	if v, ok := data.GetOk("oidc"); ok {
		rec.ProviderId = "oidc"
		configs := v.(*schema.Set).List()
		if len(configs) == 1 {
			config := configs[0].(map[string]interface{})
			rec.Config = &keycloak.IdentityProviderConfig{
				BackchannelSupported: keycloak.KeycloakBoolQuoted(config["backchannel_supported"].(bool)),
				UseJwksUrl:           keycloak.KeycloakBoolQuoted(config["use_jwks_url"].(bool)),
				ValidateSignature:    keycloak.KeycloakBoolQuoted(config["validate_signature"].(bool)),
				AuthorizationUrl:     config["authorization_url"].(string),
				ClientId:             config["client_id"].(string),
				ClientSecret:         config["client_secret"].(string),
				DisableUserInfo:      keycloak.KeycloakBoolQuoted(config["disable_user_info"].(bool)),
				HideOnLoginPage:      keycloak.KeycloakBoolQuoted(config["hide_on_login_page"].(bool)),
				TokenUrl:             config["token_url"].(string),
				LoginHint:            config["login_hint"].(string),
			}
		}
	} else if v, ok := data.GetOk("social"); ok {
		configs := v.(*schema.Set).List()
		if len(configs) == 1 {
			config := configs[0].(map[string]interface{})
			rec.ProviderId = config["provider"].(string)
			rec.Config = &keycloak.IdentityProviderConfig{
				UseJwksUrl:      keycloak.KeycloakBoolQuoted(config["use_jwks_url"].(bool)),
				ClientId:        config["client_id"].(string),
				ClientSecret:    config["client_secret"].(string),
				DisableUserInfo: keycloak.KeycloakBoolQuoted(config["disable_user_info"].(bool)),
				HideOnLoginPage: keycloak.KeycloakBoolQuoted(config["hide_on_login_page"].(bool)),
				Key:             config["key"].(string),
			}
		}
	} else if v, ok := data.GetOk("saml"); ok {
		rec.ProviderId = "saml"
		configs := v.(*schema.Set).List()
		if len(configs) == 1 {
			config := configs[0].(map[string]interface{})
			rec.Config = &keycloak.IdentityProviderConfig{
				UseJwksUrl:                       keycloak.KeycloakBoolQuoted(config["use_jwks_url"].(bool)),
				ValidateSignature:                keycloak.KeycloakBoolQuoted(config["validate_signature"].(bool)),
				HideOnLoginPage:                  keycloak.KeycloakBoolQuoted(config["hide_on_login_page"].(bool)),
				NameIDPolicyFormat:               config["name_id_policy_format"].(string),
				SingleLogutServiceUrl:            config["single_logout_service_url"].(string),
				SingleSignOnServiceUrl:           config["single_sign_on_service_url"].(string),
				SigningCertificate:               config["signing_certificate"].(string),
				SignatureAlgorithm:               config["signature_algorithm"].(string),
				XmlSignKeyInfoKeyNameTransformer: config["xml_sign_key_info_key_name_transformer"].(string),
				PostBindingAuthnRequest:          keycloak.KeycloakBoolQuoted(config["post_binding_authn_request"].(bool)),
				PostBindingResponse:              keycloak.KeycloakBoolQuoted(config["post_binding_response"].(bool)),
				PostBindingLogout:                keycloak.KeycloakBoolQuoted(config["post_binding_logout"].(bool)),
				ForceAuthn:                       keycloak.KeycloakBoolQuoted(config["force_authn"].(bool)),
				WantAuthnRequestsSigned:          keycloak.KeycloakBoolQuoted(config["want_authn_requests_signed"].(bool)),
				WantAssertionsSigned:             keycloak.KeycloakBoolQuoted(config["want_assertions_signed"].(bool)),
				WantAssertionsEncrypted:          keycloak.KeycloakBoolQuoted(config["want_assertions_encrypted"].(bool)),
			}
		}
	} else {
		return nil, fmt.Errorf("No provider config is defined. Please add social, saml or oidc provider")
	}
	return rec, nil
}

func setIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error {
	data.SetId(identityProvider.Realm + "/" + identityProvider.Alias)
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
	if config := identityProvider.Config; config != nil {
		switch identityProvider.ProviderId {
		case "oidc":
			data.Set("provider_id", "oidc")
			data.Set("config", []interface{}{
				map[string]interface{}{
					"backchannel_supported": config.BackchannelSupported,
					"use_jwks_url":          config.UseJwksUrl,
					"validate_signature":    config.ValidateSignature,
					"authorization_url":     config.AuthorizationUrl,
					"client_id":             config.ClientId,
					"client_secret":         config.ClientSecret,
					"disable_user_info":     config.DisableUserInfo,
					"hide_on_login_page":    config.HideOnLoginPage,
					"token_url":             config.TokenUrl,
					"login_hint":            config.LoginHint,
				},
			})
		case "facebook", "stackoverflow", "twitter", "github", "gitlab", "instagram", "bitbucket", "google", "microsoft", "paypal":
			data.Set("provider_id", identityProvider.ProviderId)
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
		case "saml":
			data.Set("provider_id", "saml")
			data.Set("config", []interface{}{
				map[string]interface{}{
					"backchannel_supported":                  config.BackchannelSupported,
					"use_jwks_url":                           config.UseJwksUrl,
					"validate_signature":                     config.ValidateSignature,
					"hide_on_login_page":                     config.HideOnLoginPage,
					"name_id_policy_format":                  config.NameIDPolicyFormat,
					"single_logout_service_url":              config.SingleLogutServiceUrl,
					"single_sign_on_service_url":             config.SingleSignOnServiceUrl,
					"signing_certificate":                    config.SigningCertificate,
					"signature_algorithm":                    config.SignatureAlgorithm,
					"xml_sign_key_info_key_name_transformer": config.XmlSignKeyInfoKeyNameTransformer,
					"post_binding_authn_request":             config.PostBindingAuthnRequest,
					"post_binding_response":                  config.PostBindingResponse,
					"post_binding_logout":                    config.PostBindingLogout,
					"force_authn":                            config.ForceAuthn,
					"want_authn_requests_signed":             config.WantAuthnRequestsSigned,
					"want_assertions_signed":                 config.WantAssertionsSigned,
					"want_assertions_encrypted":              config.WantAssertionsEncrypted,
				},
			})
		default:
			return fmt.Errorf("No provider config is defined. Please add social, saml or oidc provider")
		}
	}
	return nil
}

func resourceKeycloakIdentityProviderCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	identityProvider, err := getIdentityProviderFromData(data)

	err = keycloakClient.NewIdentityProvider(identityProvider)
	if err != nil {
		return err
	}

	setIdentityProviderData(data, identityProvider)

	return resourceKeycloakIdentityProviderRead(data, meta)
}

func resourceKeycloakIdentityProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	identityProvider, err := keycloakClient.GetIdentityProvider(realm, alias)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setIdentityProviderData(data, identityProvider)

	return nil
}

func resourceKeycloakIdentityProviderUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	identityProvider, err := getIdentityProviderFromData(data)

	err = keycloakClient.UpdateIdentityProvider(identityProvider)
	if err != nil {
		return err
	}

	setIdentityProviderData(data, identityProvider)

	return nil
}

func resourceKeycloakIdentityProviderDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)

	return keycloakClient.DeleteIdentityProvider(realm, alias)
}

func configHash(v interface{}) int {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", "config"))
	return hashcode.String(buf.String())
}

func resourceKeycloakIdentityProviderImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	alias := parts[1]

	d.Set("realm", realm)
	d.Set("alias", alias)

	return []*schema.ResourceData{d}, nil
}
