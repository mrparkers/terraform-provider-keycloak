package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakUserTemplateImporterIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username For Template Import",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getUserTemplateImporterIdentityProviderMapperFromData, setUserTemplateImporterIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setUserTemplateImporterIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getUserTemplateImporterIdentityProviderMapperFromData, setUserTemplateImporterIdentityProviderMapperData)
	return genericMapperResource
}

func getUserTemplateImporterIdentityProviderMapperFromData(data *schema.ResourceData, meta interface{}) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	rec, _ := getIdentityProviderMapperFromData(data)
	extraConfig := map[string]interface{}{}
	if v, ok := data.GetOk("extra_config"); ok {
		for key, value := range v.(map[string]interface{}) {
			extraConfig[key] = value
		}
	}
	identityProvider, err := keycloakClient.GetIdentityProvider(rec.Realm, rec.IdentityProviderAlias)
	if err != nil {
		return nil, handleNotFoundError(err, data)
	}
	rec.IdentityProviderMapper = fmt.Sprintf("%s-username-idp-mapper", identityProvider.ProviderId)
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Template:    data.Get("template").(string),
		ExtraConfig: extraConfig,
	}
	return rec, nil
}

func setUserTemplateImporterIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("template", identityProviderMapper.Config.Template)
	data.Set("extra_config", identityProviderMapper.Config.ExtraConfig)
	return nil
}
