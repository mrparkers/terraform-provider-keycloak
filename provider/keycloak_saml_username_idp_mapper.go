package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlUsernameIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakSamlUsernameIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakSamlUsernameIdpMapperRead
	genericMapperResource.Update = resourceKeycloakSamlUsernameIdpMapperUpdate
	return genericMapperResource
}

func getSamlUsernameIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "saml-username-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Template: data.Get("template").(string),
	}
	return rec, nil
}

func setSamlUsernameIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("template", identityProviderMapper.Config.Template)
	return nil
}

func resourceKeycloakSamlUsernameIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlUsernameIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlUsernameIdpMapperData(data, identityProvider)
	return resourceKeycloakSamlUsernameIdpMapperRead(data, meta)
}

func resourceKeycloakSamlUsernameIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setSamlUsernameIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakSamlUsernameIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getSamlUsernameIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setSamlUsernameIdpMapperData(data, identityProvider)
	return nil
}
