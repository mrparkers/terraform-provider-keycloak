package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcRoleIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"claim_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Claim Name",
		},
		"claim_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Claim Value",
		},
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Role To Grant To User",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate("oidc-role-idp-mapper")
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead("oidc-role-idp-mapper")
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate("oidc-role-idp-mapper")
	return genericMapperResource
}

func getOidcRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Claim:      data.Get("claim_name").(string),
		ClaimValue: data.Get("claim_value").(string),
		Role:       data.Get("role").(string),
	}
	return rec, nil
}

func setOidcRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	data.Set("claim_name", identityProviderMapper.Config.Claim)
	data.Set("claim_value", identityProviderMapper.Config.ClaimValue)
	return nil
}
