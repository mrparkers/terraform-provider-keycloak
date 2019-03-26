package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAttributeImporterIdentityProviderMapper() *schema.Resource {
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
		"attribute_friendly_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Attribute Friendly Name",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getAttributeImporterIdentityProviderMapperFromData, setAttributeImporterIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setAttributeImporterIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getAttributeImporterIdentityProviderMapperFromData, setAttributeImporterIdentityProviderMapperData)
	return genericMapperResource
}

func getAttributeImporterIdentityProviderMapperFromData(data *schema.ResourceData, meta interface{}) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	rec, _ := getIdentityProviderMapperFromData(data)
	identityProvider, err := keycloakClient.GetIdentityProvider(rec.Realm, rec.IdentityProviderAlias)
	if err != nil {
		return nil, handleNotFoundError(err, data)
	}
	rec.IdentityProviderMapper = fmt.Sprintf("%s-user-attribute-idp-mapper", identityProvider.ProviderId)
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Attribute:             data.Get("attribute_name").(string),
		AttributeValue:        data.Get("attribute_value").(string),
		AttributeFriendlyName: data.Get("attribute_friendly_name").(string),
	}
	return rec, nil
}

func setAttributeImporterIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	data.Set("attribute_friendly_name", identityProviderMapper.Config.AttributeFriendlyName)
	return nil
}
