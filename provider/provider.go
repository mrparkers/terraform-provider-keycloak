package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func KeycloakProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"keycloak_realm":                                   resourceKeycloakRealm(),
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
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_ID", nil),
			},
			"client_secret": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_SECRET", nil),
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

	return keycloak.NewKeycloakClient(url, clientId, clientSecret)
}
