package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
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

func getHardcodedRoleIdentityProviderMapperFromData(data *schema.ResourceData, meta interface{}) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	rec, _ := getIdentityProviderMapperFromData(data)
	identityProvider, err := keycloakClient.GetIdentityProvider(rec.Realm, rec.IdentityProviderAlias)
	if err != nil {
		return nil, handleNotFoundError(err, data)
	}
	rec.IdentityProviderMapper = "oidc-hardcoded-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role: data.Get("role").(string),
	}
	return rec, nil
}

func setHardcodedRoleIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	return nil
}
