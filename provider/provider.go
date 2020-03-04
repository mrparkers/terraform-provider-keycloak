package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func KeycloakProvider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"keycloak_group":                              dataSourceKeycloakGroup(),
			"keycloak_openid_client":                      dataSourceKeycloakOpenidClient(),
			"keycloak_openid_client_authorization_policy": dataSourceKeycloakOpenidClientAuthorizationPolicy(),
			"keycloak_openid_client_service_account_user": dataSourceKeycloakOpenidClientServiceAccountUser(),
			"keycloak_realm":                              dataSourceKeycloakRealm(),
			"keycloak_realm_keys":                         dataSourceKeycloakRealmKeys(),
			"keycloak_role":                               dataSourceKeycloakRole(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"keycloak_realm":                                           resourceKeycloakRealm(),
			"keycloak_realm_events":                                    resourceKeycloakRealmEvents(),
			"keycloak_required_action":                                 resourceKeycloakRequiredAction(),
			"keycloak_group":                                           resourceKeycloakGroup(),
			"keycloak_group_memberships":                               resourceKeycloakGroupMemberships(),
			"keycloak_default_groups":                                  resourceKeycloakDefaultGroups(),
			"keycloak_group_roles":                                     resourceKeycloakGroupRoles(),
			"keycloak_user":                                            resourceKeycloakUser(),
			"keycloak_openid_client":                                   resourceKeycloakOpenidClient(),
			"keycloak_openid_client_scope":                             resourceKeycloakOpenidClientScope(),
			"keycloak_ldap_user_federation":                            resourceKeycloakLdapUserFederation(),
			"keycloak_ldap_user_attribute_mapper":                      resourceKeycloakLdapUserAttributeMapper(),
			"keycloak_ldap_group_mapper":                               resourceKeycloakLdapGroupMapper(),
			"keycloak_ldap_hardcoded_role_mapper":                      resourceKeycloakLdapHardcodedRoleMapper(),
			"keycloak_ldap_msad_user_account_control_mapper":           resourceKeycloakLdapMsadUserAccountControlMapper(),
			"keycloak_ldap_msad_lds_user_account_control_mapper":       resourceKeycloakLdapMsadLdsUserAccountControlMapper(),
			"keycloak_ldap_full_name_mapper":                           resourceKeycloakLdapFullNameMapper(),
			"keycloak_custom_user_federation":                          resourceKeycloakCustomUserFederation(),
			"keycloak_openid_user_attribute_protocol_mapper":           resourceKeycloakOpenIdUserAttributeProtocolMapper(),
			"keycloak_openid_user_property_protocol_mapper":            resourceKeycloakOpenIdUserPropertyProtocolMapper(),
			"keycloak_openid_group_membership_protocol_mapper":         resourceKeycloakOpenIdGroupMembershipProtocolMapper(),
			"keycloak_openid_full_name_protocol_mapper":                resourceKeycloakOpenIdFullNameProtocolMapper(),
			"keycloak_openid_hardcoded_claim_protocol_mapper":          resourceKeycloakOpenIdHardcodedClaimProtocolMapper(),
			"keycloak_openid_audience_protocol_mapper":                 resourceKeycloakOpenIdAudienceProtocolMapper(),
			"keycloak_openid_hardcoded_role_protocol_mapper":           resourceKeycloakOpenIdHardcodedRoleProtocolMapper(),
			"keycloak_openid_user_realm_role_protocol_mapper":          resourceKeycloakOpenIdUserRealmRoleProtocolMapper(),
			"keycloak_openid_client_default_scopes":                    resourceKeycloakOpenidClientDefaultScopes(),
			"keycloak_openid_client_optional_scopes":                   resourceKeycloakOpenidClientOptionalScopes(),
			"keycloak_saml_client":                                     resourceKeycloakSamlClient(),
			"keycloak_generic_client_protocol_mapper":                  resourceKeycloakGenericClientProtocolMapper(),
			"keycloak_saml_user_attribute_protocol_mapper":             resourceKeycloakSamlUserAttributeProtocolMapper(),
			"keycloak_saml_user_property_protocol_mapper":              resourceKeycloakSamlUserPropertyProtocolMapper(),
			"keycloak_hardcoded_attribute_identity_provider_mapper":    resourceKeycloakHardcodedAttributeIdentityProviderMapper(),
			"keycloak_hardcoded_role_identity_provider_mapper":         resourceKeycloakHardcodedRoleIdentityProviderMapper(),
			"keycloak_attribute_importer_identity_provider_mapper":     resourceKeycloakAttributeImporterIdentityProviderMapper(),
			"keycloak_attribute_to_role_identity_provider_mapper":      resourceKeycloakAttributeToRoleIdentityProviderMapper(),
			"keycloak_user_template_importer_identity_provider_mapper": resourceKeycloakUserTemplateImporterIdentityProviderMapper(),
			"keycloak_saml_identity_provider":                          resourceKeycloakSamlIdentityProvider(),
			"keycloak_oidc_google_identity_provider":                   resourceKeycloakOidcGoogleIdentityProvider(),
			"keycloak_oidc_identity_provider":                          resourceKeycloakOidcIdentityProvider(),
			"keycloak_openid_client_authorization_resource":            resourceKeycloakOpenidClientAuthorizationResource(),
			"keycloak_openid_client_authorization_scope":               resourceKeycloakOpenidClientAuthorizationScope(),
			"keycloak_openid_client_authorization_permission":          resourceKeycloakOpenidClientAuthorizationPermission(),
			"keycloak_openid_client_service_account_role":              resourceKeycloakOpenidClientServiceAccountRole(),
			"keycloak_openid_client_service_account_realm_role":        resourceKeycloakOpenidClientServiceAccountRealmRole(),
			"keycloak_role":                                            resourceKeycloakRole(),
			"keycloak_authentication_flow":                             resourceKeycloakAuthenticationFlow(),
			"keycloak_authentication_subflow":                          resourceKeycloakAuthenticationSubFlow(),
			"keycloak_authentication_execution":                        resourceKeycloakAuthenticationExecution(),
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_ID", nil),
			},
			"client_secret": {
				Optional:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_SECRET", nil),
			},
			"username": {
				Optional:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_USER", nil),
			},
			"password": {
				Optional:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_PASSWORD", nil),
			},
			"realm": {
				Optional:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_REALM", "master"),
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL of the Keycloak instance, before `/auth`",
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_URL", nil),
			},
			"initial_login": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether or not to login to Keycloak instance on provider initialization",
				Default:     true,
			},
			"client_timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Timeout (in seconds) of the Keycloak client",
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_TIMEOUT", 5),
			},
			"root_ca_certificate": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Allows x509 calls using an unknown CA certificate (for development purposes)",
				Default:     "",
			},
		},
		ConfigureFunc: configureKeycloakProvider,
	}
}

func configureKeycloakProvider(data *schema.ResourceData) (interface{}, error) {
	url := data.Get("url").(string)
	clientId := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	username := data.Get("username").(string)
	password := data.Get("password").(string)
	realm := data.Get("realm").(string)
	initialLogin := data.Get("initial_login").(bool)
	clientTimeout := data.Get("client_timeout").(int)
	rootCaCertificate := data.Get("root_ca_certificate").(string)

	return keycloak.NewKeycloakClient(url, clientId, clientSecret, realm, username, password, initialLogin, clientTimeout, rootCaCertificate)
}
