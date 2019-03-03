package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlUserAttributeIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"attribute_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Attribute Name",
		},
		"user_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "User Attribute",
		},
		"attribute_friendly_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Attribute Friendly Name",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakSamlUserAttributeIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakSamlUserAttributeIdpMapperRead
	genericMapperResource.Update = resourceKeycloakSamlUserAttributeIdpMapperUpdate
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

func resourceKeycloakSamlUserAttributeIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlUserAttributeIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlUserAttributeIdpMapperData(data, identityProvider)
	return resourceKeycloakSamlUserAttributeIdpMapperRead(data, meta)
}

func resourceKeycloakSamlUserAttributeIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setSamlUserAttributeIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakSamlUserAttributeIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlUserAttributeIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlUserAttributeIdpMapperData(data, identityProvider)
	return nil
}
