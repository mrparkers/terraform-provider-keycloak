package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcUsernameTemplateImporterMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakOidcUsernameTemplateImporterMapperCreate
	genericMapperResource.Read = resourceKeycloakOidcUsernameTemplateImporterMapperRead
	genericMapperResource.Update = resourceKeycloakOidcUsernameTemplateImporterMapperUpdate
	return genericMapperResource
}

func getOidcUsernameTemplateImporterMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-username-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Template: data.Get("template").(string),
	}
	return rec, nil
}

func setOidcUsernameTemplateImporterMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("template", identityProviderMapper.Config.Template)
	return nil
}

func resourceKeycloakOidcUsernameTemplateImporterMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcUsernameTemplateImporterMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcUsernameTemplateImporterMapperData(data, identityProvider)
	return resourceKeycloakOidcUsernameTemplateImporterMapperRead(data, meta)
}

func resourceKeycloakOidcUsernameTemplateImporterMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setOidcUsernameTemplateImporterMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakOidcUsernameTemplateImporterMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcUsernameTemplateImporterMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcUsernameTemplateImporterMapperData(data, identityProvider)
	return nil
}
