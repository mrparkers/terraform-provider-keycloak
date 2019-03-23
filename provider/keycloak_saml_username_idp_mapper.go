package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlUsernameIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getSamlUsernameIdpMapperFromData, setSamlUsernameIdpMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setSamlUsernameIdpMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getSamlUsernameIdpMapperFromData, setSamlUsernameIdpMapperData)
	return genericMapperResource
}

func getSamlUsernameIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "saml-username-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Template: data.Get("template").(string),
	}
	return rec, nil
}

func setSamlUsernameIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("template", identityProviderMapper.Config.Template)
	return nil
}
