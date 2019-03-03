package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlRoleIdpMapper() *schema.Resource {
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
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Role Name",
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
	genericMapperResource.Create = resourceKeycloakSamlRoleIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakSamlRoleIdpMapperRead
	genericMapperResource.Update = resourceKeycloakSamlRoleIdpMapperUpdate
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

func resourceKeycloakSamlRoleIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlRoleIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlRoleIdpMapperData(data, identityProvider)
	return resourceKeycloakSamlRoleIdpMapperRead(data, meta)
}

func resourceKeycloakSamlRoleIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setSamlRoleIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakSamlRoleIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlRoleIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlRoleIdpMapperData(data, identityProvider)
	return nil
}
