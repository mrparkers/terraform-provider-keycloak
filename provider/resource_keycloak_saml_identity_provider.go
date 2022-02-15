package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/imdario/mergo"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var nameIdPolicyFormats = map[string]string{
	"Windows Domain Qualified Name": "urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName",
	"Persistent":                    "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
	"Email":                         "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
	"Kerberos":                      "urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos",
	"X.509 Subject Name":            "urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName",
	"Unspecified":                   "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
}

var signatureAlgorithms = []string{
	"RSA_SHA1",
	"RSA_SHA256",
	"RSA_SHA512",
	"DSA_SHA1",
}

var keyNameTransformers = []string{
	"NONE",
	"KEY_ID",
	"CERT_SUBJECT",
}

var principalTypes = []string{
	"SUBJECT",
	"ATTRIBUTE",
	"FRIENDLY_ATTRIBUTE",
}

var authnComparisonTypes = []string{
	"exact",
	"minimum",
	"maximum",
	"better",
}

func resourceKeycloakSamlIdentityProvider() *schema.Resource {
	samlSchema := map[string]*schema.Schema{
		"backchannel_supported": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Does the external IDP support backchannel logout?",
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
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ValidateFunc: validation.StringInSlice(keys(nameIdPolicyFormats), false),
			StateFunc: func(value interface{}) string {
				return nameIdPolicyFormats[value.(string)]
			},
			Description: "Name ID Policy Format.",
		},
		"single_logout_service_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Logout URL.",
		},
		"entity_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Entity ID that will be used to uniquely identify this SAML Service Provider.",
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
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ValidateFunc: validation.StringInSlice(signatureAlgorithms, false),
			Description:  "Signing Algorithm.",
		},
		"xml_sign_key_info_key_name_transformer": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ValidateFunc: validation.StringInSlice(keyNameTransformers, false),
			Description:  "Sign Key Transformer.",
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
		"principal_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ValidateFunc: validation.StringInSlice(principalTypes, false),
			Description:  "Principal Type",
		},
		"principal_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Principal Attribute",
		},
		"authn_context_class_refs": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "AuthnContext ClassRefs",
		},
		"authn_context_decl_refs": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "AuthnContext DeclRefs",
		},
		"authn_context_comparison_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(authnComparisonTypes, false),
			Description:  "AuthnContext Comparison",
		},
	}
	samlResource := resourceKeycloakIdentityProvider()
	samlResource.Schema = mergeSchemas(samlResource.Schema, samlSchema)
	samlResource.Create = resourceKeycloakIdentityProviderCreate(getSamlIdentityProviderFromData, setSamlIdentityProviderData)
	samlResource.Read = resourceKeycloakIdentityProviderRead(setSamlIdentityProviderData)
	samlResource.Update = resourceKeycloakIdentityProviderUpdate(getSamlIdentityProviderFromData, setSamlIdentityProviderData)
	return samlResource
}

func getSamlIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec, defaultConfig := getIdentityProviderFromData(data)
	rec.ProviderId = "saml"

	samlIdentityProviderConfig := &keycloak.IdentityProviderConfig{
		ValidateSignature:                keycloak.KeycloakBoolQuoted(data.Get("validate_signature").(bool)),
		HideOnLoginPage:                  keycloak.KeycloakBoolQuoted(data.Get("hide_on_login_page").(bool)),
		BackchannelSupported:             keycloak.KeycloakBoolQuoted(data.Get("backchannel_supported").(bool)),
		NameIDPolicyFormat:               nameIdPolicyFormats[data.Get("name_id_policy_format").(string)],
		EntityId:                         data.Get("entity_id").(string),
		SingleLogoutServiceUrl:           data.Get("single_logout_service_url").(string),
		SingleSignOnServiceUrl:           data.Get("single_sign_on_service_url").(string),
		SigningCertificate:               data.Get("signing_certificate").(string),
		SignatureAlgorithm:               data.Get("signature_algorithm").(string),
		XmlSignKeyInfoKeyNameTransformer: data.Get("xml_sign_key_info_key_name_transformer").(string),
		PostBindingAuthnRequest:          keycloak.KeycloakBoolQuoted(data.Get("post_binding_authn_request").(bool)),
		PostBindingResponse:              keycloak.KeycloakBoolQuoted(data.Get("post_binding_response").(bool)),
		PostBindingLogout:                keycloak.KeycloakBoolQuoted(data.Get("post_binding_logout").(bool)),
		ForceAuthn:                       keycloak.KeycloakBoolQuoted(data.Get("force_authn").(bool)),
		WantAssertionsSigned:             keycloak.KeycloakBoolQuoted(data.Get("want_assertions_signed").(bool)),
		WantAssertionsEncrypted:          keycloak.KeycloakBoolQuoted(data.Get("want_assertions_encrypted").(bool)),
		PrincipalType:                    data.Get("principal_type").(string),
		PrincipalAttribute:               data.Get("principal_attribute").(string),
		AuthnContextClassRefs:            data.Get("authn_context_class_refs").([]interface{}),
		AuthnContextComparisonType:       data.Get("authn_context_comparison_type").(string),
		AuthnContextDeclRefs:             data.Get("authn_context_decl_refs").([]interface{}),
	}

	if _, ok := data.GetOk("signature_algorithm"); ok {
		samlIdentityProviderConfig.WantAuthnRequestsSigned = true
	}

	if err := mergo.Merge(samlIdentityProviderConfig, defaultConfig); err != nil {
		return nil, err
	}

	rec.Config = samlIdentityProviderConfig

	return rec, nil
}

func setSamlIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error {
	setIdentityProviderData(data, identityProvider)

	data.Set("backchannel_supported", identityProvider.Config.BackchannelSupported)
	data.Set("validate_signature", identityProvider.Config.ValidateSignature)
	data.Set("hide_on_login_page", identityProvider.Config.HideOnLoginPage)
	data.Set("name_id_policy_format", identityProvider.Config.NameIDPolicyFormat)
	data.Set("entity_id", identityProvider.Config.EntityId)
	data.Set("single_logout_service_url", identityProvider.Config.SingleLogoutServiceUrl)
	data.Set("single_sign_on_service_url", identityProvider.Config.SingleSignOnServiceUrl)
	data.Set("signing_certificate", identityProvider.Config.SigningCertificate)
	data.Set("signature_algorithm", identityProvider.Config.SignatureAlgorithm)
	data.Set("xml_sign_key_info_key_name_transformer", identityProvider.Config.XmlSignKeyInfoKeyNameTransformer)
	data.Set("post_binding_authn_request", identityProvider.Config.PostBindingAuthnRequest)
	data.Set("post_binding_response", identityProvider.Config.PostBindingResponse)
	data.Set("post_binding_logout", identityProvider.Config.PostBindingLogout)
	data.Set("force_authn", identityProvider.Config.ForceAuthn)
	data.Set("want_assertions_signed", identityProvider.Config.WantAssertionsSigned)
	data.Set("want_assertions_encrypted", identityProvider.Config.WantAssertionsEncrypted)
	data.Set("principal_type", identityProvider.Config.PrincipalType)
	data.Set("principal_attribute", identityProvider.Config.PrincipalAttribute)
	data.Set("authn_context_class_refs", identityProvider.Config.AuthnContextClassRefs)
	data.Set("authn_context_comparison_type", identityProvider.Config.AuthnContextComparisonType)
	data.Set("authn_context_decl_refs", identityProvider.Config.AuthnContextDeclRefs)

	return nil
}
