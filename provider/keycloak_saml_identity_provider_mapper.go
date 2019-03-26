package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var samlIdentityProviderMappers = map[string]string{
	"Hardcoded Role":                   "saml-hardcoded-role-idp-mapper",
	"Hardcoded Attribute":              "hardcoded-attribute-idp-mapper",
	"Hardcoded User Session Attribute": "hardcoded-user-session-attribute-idp-mapper",
	"SAML Attribute To Role":           "saml-role-idp-mapper",
	"Attribute Mapper":                 "saml-user-attribute-idp-mapper",
	"User Template Importer":           "saml-username-idp-mapper",
}

func resourceKeycloakSamlIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Identity Provider Mapper Type",
			ValidateFunc: validation.StringInSlice(keys(samlIdentityProviderMappers), false),
		},
		"attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "OIDC Claim",
		},
		"attribute_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute",
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
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getSamlIdentityProviderMapperFromData, setSamlIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setSamlIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getSamlIdentityProviderMapperFromData, setSamlIdentityProviderMapperData)
	return genericMapperResource
}

func getSamlIdentityProviderMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	mapperType := data.Get("type").(string)
	rec.IdentityProviderMapper = samlIdentityProviderMappers[mapperType]
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role:                  data.Get("role").(string),
		Attribute:             data.Get("attribute_name").(string),
		AttributeValue:        data.Get("attribute_value").(string),
		AttributeFriendlyName: data.Get("attribute_friendly_name").(string),
		Template:              data.Get("template").(string),
	}
	return rec, nil
}

func setSamlIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	data.Set("attribute_friendly_name", identityProviderMapper.Config.AttributeFriendlyName)
	data.Set("template", identityProviderMapper.Config.Template)
	return nil
}
