package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcIdentityProvider() *schema.Resource {
	oidcSchema := map[string]*schema.Schema{
		"backchannel_supported": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Does the external IDP support backchannel logout?",
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
			Sensitive:   true,
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
		"ui_locales": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Pass current locale to identity provider",
		},
	}
	oidcResource := resourceKeycloakIdentityProvider()
	oidcResource.Schema = mergeSchemas(oidcResource.Schema, oidcSchema)
	oidcResource.Create = resourceKeycloakIdentityProviderCreate(getOidcIdentityProviderFromData, setOidcIdentityProviderData)
	oidcResource.Read = resourceKeycloakIdentityProviderRead(setOidcIdentityProviderData)
	oidcResource.Update = resourceKeycloakIdentityProviderUpdate(getOidcIdentityProviderFromData, setOidcIdentityProviderData)
	return oidcResource
}

func getOidcIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec, _ := getIdentityProviderFromData(data)
	rec.ProviderId = "oidc"
	rec.Config = &keycloak.IdentityProviderConfig{
		BackchannelSupported: keycloak.KeycloakBoolQuoted(data.Get("backchannel_supported").(bool)),
		UseJwksUrl:           keycloak.KeycloakBoolQuoted(true),
		ValidateSignature:    keycloak.KeycloakBoolQuoted(data.Get("validate_signature").(bool)),
		AuthorizationUrl:     data.Get("authorization_url").(string),
		ClientId:             data.Get("client_id").(string),
		DisableUserInfo:      keycloak.KeycloakBoolQuoted(data.Get("disable_user_info").(bool)),
		HideOnLoginPage:      keycloak.KeycloakBoolQuoted(data.Get("hide_on_login_page").(bool)),
		TokenUrl:             data.Get("token_url").(string),
		UILocales:            keycloak.KeycloakBoolQuoted(data.Get("ui_locales").(bool)),
		LoginHint:            data.Get("login_hint").(string),
	}

	if data.HasChange("client_secret") {
		rec.Config.ClientSecret = data.Get("client_secret").(string)
	}

	return rec, nil
}

func setOidcIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error {
	setIdentityProviderData(data, identityProvider)
	data.Set("backchannel_supported", identityProvider.Config.BackchannelSupported)
	data.Set("use_jwks_url", identityProvider.Config.UseJwksUrl)
	data.Set("validate_signature", identityProvider.Config.ValidateSignature)
	data.Set("authorization_url", identityProvider.Config.AuthorizationUrl)
	data.Set("client_id", identityProvider.Config.ClientId)
	data.Set("disable_user_info", identityProvider.Config.DisableUserInfo)
	data.Set("hide_on_login_page", identityProvider.Config.HideOnLoginPage)
	data.Set("token_url", identityProvider.Config.TokenUrl)
	data.Set("login_hint", identityProvider.Config.LoginHint)
	data.Set("ui_locales", identityProvider.Config.UILocales)
	return nil
}
