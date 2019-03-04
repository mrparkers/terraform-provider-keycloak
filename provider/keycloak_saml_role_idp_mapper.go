package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlRoleIdpMapper() *schema.Resource {
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
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Role Name",
		},
		"attribute_friendly_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Attribute Friendly Name",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate("saml-role-idp-mapper")
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead("saml-role-idp-mapper")
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate("saml-role-idp-mapper")
	return genericMapperResource
}

func getSamlRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "saml-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Attribute:             data.Get("attribute_name").(string),
		AttributeValue:        data.Get("attribute_value").(string),
		Role:                  data.Get("role").(string),
		AttributeFriendlyName: data.Get("attribute_friendly_name").(string),
	}
	return rec, nil
}

func setSamlRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	data.Set("role", identityProviderMapper.Config.Role)
	data.Set("attribute_friendly_name", identityProviderMapper.Config.AttributeFriendlyName)
	return nil
}
