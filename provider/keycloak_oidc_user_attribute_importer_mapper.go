package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcUserAttributeImporterMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"claim_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "OIDC Claim Name",
		},
		"user_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "User Attribute",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakOidcUserAttributeImporterMapperCreate
	genericMapperResource.Read = resourceKeycloakOidcUserAttributeImporterMapperRead
	genericMapperResource.Update = resourceKeycloakOidcUserAttributeImporterMapperUpdate
	return genericMapperResource
}

func getOidcUserAttributeImporterMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-user-attribute-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Claim:         data.Get("claim_name").(string),
		UserAttribute: data.Get("user_attribute").(string),
	}
	return rec, nil
}

func setOidcUserAttributeImporterMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("claim", identityProviderMapper.Config.Claim)
	data.Set("user_attribute", identityProviderMapper.Config.UserAttribute)
	return nil
}

func resourceKeycloakOidcUserAttributeImporterMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcUserAttributeImporterMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcUserAttributeImporterMapperData(data, identityProvider)
	return resourceKeycloakOidcUserAttributeImporterMapperRead(data, meta)
}

func resourceKeycloakOidcUserAttributeImporterMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setOidcUserAttributeImporterMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakOidcUserAttributeImporterMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcUserAttributeImporterMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcUserAttributeImporterMapperData(data, identityProvider)
	return nil
}
