package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func KeycloakProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"keycloak_realm":                                   resourceKeycloakRealm(),
			"keycloak_group":                                   resourceKeycloakGroup(),
			"keycloak_group_memberships":                       resourceKeycloakGroupMemberships(),
			"keycloak_user":                                    resourceKeycloakUser(),
			"keycloak_openid_client":                           resourceKeycloakOpenidClient(),
			"keycloak_openid_client_scope":                     resourceKeycloakOpenidClientScope(),
			"keycloak_ldap_user_federation":                    resourceKeycloakLdapUserFederation(),
			"keycloak_ldap_user_attribute_mapper":              resourceKeycloakLdapUserAttributeMapper(),
			"keycloak_ldap_group_mapper":                       resourceKeycloakLdapGroupMapper(),
			"keycloak_ldap_msad_user_account_control_mapper":   resourceKeycloakLdapMsadUserAccountControlMapper(),
			"keycloak_ldap_full_name_mapper":                   resourceKeycloakLdapFullNameMapper(),
			"keycloak_custom_user_federation":                  resourceKeycloakCustomUserFederation(),
			"keycloak_openid_user_attribute_protocol_mapper":   resourceKeycloakOpenIdUserAttributeProtocolMapper(),
			"keycloak_openid_user_property_protocol_mapper":    resourceKeycloakOpenIdUserPropertyProtocolMapper(),
			"keycloak_openid_group_membership_protocol_mapper": resourceKeycloakOpenIdGroupMembershipProtocolMapper(),
			"keycloak_openid_full_name_protocol_mapper":        resourceKeycloakOpenIdFullNameProtocolMapper(),
			"keycloak_openid_hardcoded_claim_protocol_mapper":  resourceKeycloakOpenIdHardcodedClaimProtocolMapper(),
			"keycloak_openid_audience_protocol_mapper":         resourceKeycloakOpenIdAudienceProtocolMapper(),
			"keycloak_openid_client_default_scopes":            resourceKeycloakOpenidClientDefaultScopes(),
			"keycloak_openid_client_optional_scopes":           resourceKeycloakOpenidClientOptionalScopes(),
			"keycloak_saml_client":                             resourceKeycloakSamlClient(),
			"keycloak_saml_user_attribute_protocol_mapper":     resourceKeycloakSamlUserAttributeProtocolMapper(),
			"keycloak_saml_user_property_protocol_mapper":      resourceKeycloakSamlUserPropertyProtocolMapper(),
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_ID", nil),
			},
			"client_secret": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("KEYCLOAK_CLIENT_SECRET", nil),
				ConflictsWith: []string{"username", "password"},
			},
			"username": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("KEYCLOAK_USER", nil),
				ConflictsWith: []string{"client_secret"},
			},
			"password": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("KEYCLOAK_PASSWORD", nil),
				ConflictsWith: []string{"client_secret"},
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
	return keycloak.NewKeycloakClient(url, clientId, clientSecret, realm, username, password)
}
