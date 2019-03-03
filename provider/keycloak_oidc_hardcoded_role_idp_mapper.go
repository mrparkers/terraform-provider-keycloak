package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcHardcodedRoleIdpMapper() *schema.Resource {
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
	genericMapperResource.Create = resourceKeycloakOidcHardcodedRoleIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakOidcHardcodedRoleIdpMapperRead
	genericMapperResource.Update = resourceKeycloakOidcHardcodedRoleIdpMapperUpdate
	return genericMapperResource
}

func getOidcHardcodedRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-hardcoded-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Role: data.Get("role").(string),
	}
	return rec, nil
}

func setOidcHardcodedRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	return nil
}

func resourceKeycloakOidcHardcodedRoleIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcHardcodedRoleIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcHardcodedRoleIdpMapperData(data, identityProvider)
	return resourceKeycloakOidcHardcodedRoleIdpMapperRead(data, meta)
}

func resourceKeycloakOidcHardcodedRoleIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setOidcHardcodedRoleIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakOidcHardcodedRoleIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcHardcodedRoleIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcHardcodedRoleIdpMapperData(data, identityProvider)
	return nil
}
