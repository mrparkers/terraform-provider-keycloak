package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakHardcodedRoleIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Role Name",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getHardcodedRoleIdentityProviderMapperFromData, setHardcodedRoleIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setHardcodedRoleIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getHardcodedRoleIdentityProviderMapperFromData, setHardcodedRoleIdentityProviderMapperData)
	return genericMapperResource
}

func getHardcodedRoleIdentityProviderMapperFromData(data *schema.ResourceData, _ interface{}) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)

	rec.IdentityProviderMapper = "oidc-hardcoded-role-idp-mapper"
	rec.Config.Role = data.Get("role").(string)

	return rec, nil
}

func setHardcodedRoleIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)

	return nil
}
