package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlIdentityProvider() *schema.Resource {
	samlSchema := map[string]*schema.Schema{
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
	}
	samlResource := resourceKeycloakIdentityProvider()
	samlResource.Schema = mergeSchemas(samlResource.Schema, samlSchema)
	samlResource.Create = resourceKeycloakSamlIdentityProviderCreate
	samlResource.Read = resourceKeycloakSamlIdentityProviderRead
	samlResource.Update = resourceKeycloakSamlIdentityProviderUpdate
	return samlResource
}

func getSamlIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec, _ := getIdentityProviderFromData(data)
	rec.ProviderId = "saml"
	rec.Config = &keycloak.IdentityProviderConfig{
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

func setSamlIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error {
	setIdentityProviderData(data, identityProvider)
	data.Set("backchannel_supported", identityProvider.Config.BackchannelSupported)
	data.Set("use_jwks_url", identityProvider.Config.UseJwksUrl)
	data.Set("validate_signature", identityProvider.Config.ValidateSignature)
	data.Set("hide_on_login_page", identityProvider.Config.HideOnLoginPage)
	data.Set("name_id_policy_format", identityProvider.Config.NameIDPolicyFormat)
	data.Set("single_logout_service_url", identityProvider.Config.SingleLogutServiceUrl)
	data.Set("single_sign_on_service_url", identityProvider.Config.SingleSignOnServiceUrl)
	data.Set("signing_certificate", identityProvider.Config.SigningCertificate)
	data.Set("signature_algorithm", identityProvider.Config.SignatureAlgorithm)
	data.Set("xml_sign_key_info_key_name_transformer", identityProvider.Config.XmlSignKeyInfoKeyNameTransformer)
	data.Set("post_binding_authn_request", identityProvider.Config.PostBindingAuthnRequest)
	data.Set("post_binding_response", identityProvider.Config.PostBindingResponse)
	data.Set("post_binding_logout", identityProvider.Config.PostBindingLogout)
	data.Set("force_authn", identityProvider.Config.ForceAuthn)
	data.Set("want_authn_requests_signed", identityProvider.Config.WantAuthnRequestsSigned)
	data.Set("want_assertions_signed", identityProvider.Config.WantAssertionsSigned)
	data.Set("want_assertions_encrypted", identityProvider.Config.WantAssertionsEncrypted)
	return nil
}

func resourceKeycloakSamlIdentityProviderCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlIdentityProviderFromData(data)
	err = keycloakClient.NewIdentityProvider(identityProvider)
	if err != nil {
		return err
	}
	setSamlIdentityProviderData(data, identityProvider)
	return resourceKeycloakSamlIdentityProviderRead(data, meta)
}

func resourceKeycloakSamlIdentityProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("alias").(string)
	identityProvider, err := keycloakClient.GetIdentityProvider(realm, alias)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setSamlIdentityProviderData(data, identityProvider)
	return nil
}

func resourceKeycloakSamlIdentityProviderUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlIdentityProviderFromData(data)
	err = keycloakClient.UpdateIdentityProvider(identityProvider)
	if err != nil {
		return err
	}
	setSamlIdentityProviderData(data, identityProvider)
	return nil
}
