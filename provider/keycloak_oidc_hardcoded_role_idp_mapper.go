package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcHardcodedRoleIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Role To Grant To User",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getOidcHardcodedRoleIdpMapperFromData, setOidcHardcodedRoleIdpMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setOidcHardcodedRoleIdpMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getOidcHardcodedRoleIdpMapperFromData, setOidcHardcodedRoleIdpMapperData)
	return genericMapperResource
}

func getOidcHardcodedRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-hardcoded-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role: data.Get("role").(string),
	}
	return rec, nil
}

func setOidcHardcodedRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	return nil
}
