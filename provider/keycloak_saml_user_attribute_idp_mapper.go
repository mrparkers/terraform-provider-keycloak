package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlUserAttributeIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Attribute Name",
		},
		"user_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute",
		},
		"attribute_friendly_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Attribute Friendly Name",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate("saml-user-attribute-idp-mapper")
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead("saml-user-attribute-idp-mapper")
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate("saml-user-attribute-idp-mapper")
	return genericMapperResource
}

func getSamlUserAttributeIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "saml-user-attribute-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Attribute:             data.Get("attribute_name").(string),
		UserAttribute:         data.Get("user_attribute").(string),
		AttributeFriendlyName: data.Get("attribute_friendly_name").(string),
	}
	return rec, nil
}

func setSamlUserAttributeIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("user_attribute", identityProviderMapper.Config.UserAttribute)
	data.Set("attribute_friendly_name", identityProviderMapper.Config.AttributeFriendlyName)
	return nil
}
