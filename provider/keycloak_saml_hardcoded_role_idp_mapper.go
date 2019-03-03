package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlHardcodedRoleIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"role": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Role To Grant To User",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakSamlHardcodedRoleIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakSamlHardcodedRoleIdpMapperRead
	genericMapperResource.Update = resourceKeycloakSamlHardcodedRoleIdpMapperUpdate
	return genericMapperResource
}

func getSamlHardcodedRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "saml-hardcoded-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role: data.Get("role").(string),
	}
	return rec, nil
}

func setSamlHardcodedRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	return nil
}

func resourceKeycloakSamlHardcodedRoleIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlHardcodedRoleIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlHardcodedRoleIdpMapperData(data, identityProvider)
	return resourceKeycloakSamlHardcodedRoleIdpMapperRead(data, meta)
}

func resourceKeycloakSamlHardcodedRoleIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setSamlHardcodedRoleIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakSamlHardcodedRoleIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlHardcodedRoleIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlHardcodedRoleIdpMapperData(data, identityProvider)
	return nil
}
