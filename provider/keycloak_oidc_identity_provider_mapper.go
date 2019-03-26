package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var oidcIdentityProviderMappers = map[string]string{
	"Hardcoded Role":                   "oidc-hardcoded-role-idp-mapper",
	"Claim to Role":                    "oidc-role-idp-mapper",
	"Attribute Importer":               "oidc-user-attribute-idp-mapper",
	"Hardcoded User Session Attribute": "hardcoded-user-session-attribute-idp-mapper",
	"Username Template Importer":       "oidc-username-idp-mapper",
	"Hardcoded Attribute":              "hardcoded-attribute-idp-mapper",
}

func resourceKeycloakOidcIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Role To Grant To User",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Identity Provider Mapper Type",
			ValidateFunc: validation.StringInSlice(keys(oidcIdentityProviderMappers), false),
		},
		"claim_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Claim Name",
		},
		"claim_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Claim Value",
		},
		"user_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute",
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
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getOidcIdentityProviderMapperFromData, setOidcIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setOidcIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getOidcIdentityProviderMapperFromData, setOidcIdentityProviderMapperData)
	return genericMapperResource
}

func getOidcIdentityProviderMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	mapperType := data.Get("type").(string)
	rec.IdentityProviderMapper = oidcIdentityProviderMappers[mapperType]
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role:           data.Get("role").(string),
		Claim:          data.Get("claim_name").(string),
		ClaimValue:     data.Get("claim_value").(string),
		UserAttribute:  data.Get("user_attribute").(string),
		Template:       data.Get("template").(string),
		Attribute:      data.Get("attribute_name").(string),
		AttributeValue: data.Get("attribute_value").(string),
	}
	return rec, nil
}

func setOidcIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	data.Set("claim_name", identityProviderMapper.Config.Claim)
	data.Set("claim_value", identityProviderMapper.Config.ClaimValue)
	data.Set("user_attribute", identityProviderMapper.Config.UserAttribute)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	data.Set("template", identityProviderMapper.Config.Template)
	return nil
}
