package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcUsernameIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate("oidc-username-idp-mapper")
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead("oidc-username-idp-mapper")
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate("oidc-username-idp-mapper")
	return genericMapperResource
}

func getOidcUsernameIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-username-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Template: data.Get("template").(string),
	}
	return rec, nil
}

func setOidcUsernameIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("template", identityProviderMapper.Config.Template)
	return nil
}
