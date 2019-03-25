package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakHardcodedUserSessionAttributeIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Attribute Name",
		},
		"attribute_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Attribute Value",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getHardcodedUserSessionAttributeIdpMapperFromData, setHardcodedUserSessionAttributeIdpMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setHardcodedUserSessionAttributeIdpMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getHardcodedUserSessionAttributeIdpMapperFromData, setHardcodedUserSessionAttributeIdpMapperData)
	return genericMapperResource
}

func getHardcodedUserSessionAttributeIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "hardcoded-user-session-attribute-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Attribute:      data.Get("attribute_name").(string),
		AttributeValue: data.Get("attribute_value").(string),
	}
	return rec, nil
}

func setHardcodedUserSessionAttributeIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setHardcodedUserSessionAttributeIdpMapperData(data, identityProviderMapper)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	return nil
}