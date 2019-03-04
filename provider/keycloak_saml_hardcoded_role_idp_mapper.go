package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlHardcodedRoleIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Role To Grant To User",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate("saml-hardcoded-role-idp-mapper")
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead("saml-hardcoded-role-idp-mapper")
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate("saml-hardcoded-role-idp-mapper")
	return genericMapperResource
}

func getSamlHardcodedRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "saml-hardcoded-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role: data.Get("role").(string),
	}
	return rec, nil
}

func setSamlHardcodedRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	return nil
}
