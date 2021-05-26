package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakCustomIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"identity_provider_mapper": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IDP Mapper Type",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getCustomIdentityProviderMapperFromData, setCustomIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setCustomIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getCustomIdentityProviderMapperFromData, setCustomIdentityProviderMapperData)
	return genericMapperResource
}

func getCustomIdentityProviderMapperFromData(data *schema.ResourceData, meta interface{}) (*keycloak.IdentityProviderMapper, error) {
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
	identityProviderMapper := data.Get("identity_provider_mapper").(string)
	if strings.Contains(identityProviderMapper, "%s") {
		rec.IdentityProviderMapper = fmt.Sprintf(identityProviderMapper, identityProvider.ProviderId)
	} else {
		rec.IdentityProviderMapper = identityProviderMapper
	}
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		ExtraConfig: extraConfig,
	}
	return rec, nil
}

func setCustomIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("identity_provider_mapper", identityProviderMapper.IdentityProviderMapper)
	data.Set("extra_config", identityProviderMapper.Config.ExtraConfig)
	return nil
}
