package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOidcRoleIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"claim_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Claim Name",
		},
		"claim_value": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Claim Value",
		},
		"role": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Role To Grant To User",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakOidcRoleIdpMapperCreate
	genericMapperResource.Read = resourceKeycloakOidcRoleIdpMapperRead
	genericMapperResource.Update = resourceKeycloakOidcRoleIdpMapperUpdate
	return genericMapperResource
}

func getOidcRoleIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "oidc-role-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Claim:      data.Get("claim_name").(string),
		ClaimValue: data.Get("claim_value").(string),
		Role:       data.Get("role").(string),
	}
	return rec, nil
}

func setOidcRoleIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("role", identityProviderMapper.Config.Role)
	data.Set("claim_name", identityProviderMapper.Config.Claim)
	data.Set("claim_value", identityProviderMapper.Config.ClaimValue)
	return nil
}

func resourceKeycloakOidcRoleIdpMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcRoleIdpMapperFromData(data)
	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcRoleIdpMapperData(data, identityProvider)
	return resourceKeycloakOidcRoleIdpMapperRead(data, meta)
}

func resourceKeycloakOidcRoleIdpMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()
	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	setOidcRoleIdpMapperData(data, identityProvider)
	return nil
}

func resourceKeycloakOidcRoleIdpMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	identityProvider, err := getOidcRoleIdpMapperFromData(data)
	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}
	setOidcRoleIdpMapperData(data, identityProvider)
	return nil
}
