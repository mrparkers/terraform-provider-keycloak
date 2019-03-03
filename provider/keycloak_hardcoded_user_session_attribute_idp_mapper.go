package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakHardcodedUserSessionAttributeIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"attribute_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Attribute Name",
		},
		"attribute_value": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Attribute Value",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakHardcodedUserSessionAttributeIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakHardcodedUserSessionAttributeIdpMapperRead
	genericMapperResource.Update = resourceKeycloakHardcodedUserSessionAttributeIdpMapperUpdate
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

func resourceKeycloakHardcodedUserSessionAttributeIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getHardcodedUserSessionAttributeIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setHardcodedUserSessionAttributeIdpMapperData(data, identityProvider)
	return resourceKeycloakHardcodedUserSessionAttributeIdpMapperRead(data, meta)
}

func resourceKeycloakHardcodedUserSessionAttributeIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setHardcodedUserSessionAttributeIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakHardcodedUserSessionAttributeIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getHardcodedUserSessionAttributeIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setHardcodedUserSessionAttributeIdpMapperData(data, identityProvider)
	return nil
}
