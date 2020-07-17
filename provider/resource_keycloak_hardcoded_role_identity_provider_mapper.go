package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	extraConfig := map[string]interface{}{}
	if v, ok := data.GetOk("extra_config"); ok {
		for key, value := range v.(map[string]interface{}) {
			extraConfig[key] = value
		}
	}
	rec.IdentityProviderMapper = "oidc-hardcoded-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role:        data.Get("role").(string),
		ExtraConfig: extraConfig,
	}
	return rec, nil
}

func setHardcodedRoleIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	data.Set("extra_config", identityProviderMapper.Config.ExtraConfig)
	return nil
}
