package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcUserAttributeIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"claim_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "OIDC Claim Name",
		},
		"user_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getOidcUserAttributeIdpMapperFromData, setOidcUserAttributeIdpMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setOidcUserAttributeIdpMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getOidcUserAttributeIdpMapperFromData, setOidcUserAttributeIdpMapperData)
	return genericMapperResource
}

func getOidcUserAttributeIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-user-attribute-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Claim:         data.Get("claim_name").(string),
		UserAttribute: data.Get("user_attribute").(string),
	}
	return rec, nil
}

func setOidcUserAttributeIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("claim", identityProviderMapper.Config.Claim)
	data.Set("user_attribute", identityProviderMapper.Config.UserAttribute)
	return nil
}
