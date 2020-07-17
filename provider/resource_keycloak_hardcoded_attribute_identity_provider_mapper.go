package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakHardcodedAttributeIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "OIDC Claim",
		},
		"attribute_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute",
		},
		"user_session": {
			Type:        schema.TypeBool,
			Required:    true,
			ForceNew:    true,
			Description: "Is Attribute Related To a User Session",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getHardcodedAttributeIdentityProviderMapperFromData, setHardcodedAttributeIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setHardcodedAttributeIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getHardcodedAttributeIdentityProviderMapperFromData, setHardcodedAttributeIdentityProviderMapperData)
	return genericMapperResource
}

func getHardcodedAttributeIdentityProviderMapperFromData(data *schema.ResourceData, _ interface{}) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	extraConfig := map[string]interface{}{}
	if v, ok := data.GetOk("extra_config"); ok {
		for key, value := range v.(map[string]interface{}) {
			extraConfig[key] = value
		}
	}
	rec.IdentityProviderMapper = getHardcodedAttributeIdentityProviderMapperType(data.Get("user_session").(bool))
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		HardcodedAttribute: data.Get("attribute_name").(string),
		AttributeValue:     data.Get("attribute_value").(string),
		ExtraConfig:        extraConfig,
	}
	return rec, nil
}

func setHardcodedAttributeIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("attribute_name", identityProviderMapper.Config.HardcodedAttribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	data.Set("extra_config", identityProviderMapper.Config.ExtraConfig)
	mapperType, err := getUserSessionFromHardcodedAttributeIdentityProviderMapperType(identityProviderMapper.IdentityProviderMapper)
	if err != nil {
		return err
	}
	data.Set("user_session", mapperType)
	return nil
}

func getHardcodedAttributeIdentityProviderMapperType(userSession bool) string {
	if userSession {
		return "hardcoded-user-session-attribute-idp-mapper"
	} else {
		return "hardcoded-attribute-idp-mapper"
	}
}

func getUserSessionFromHardcodedAttributeIdentityProviderMapperType(mapperType string) (bool, error) {
	if mapperType == "hardcoded-user-session-attribute-idp-mapper" {
		return true, nil
	} else if mapperType == "hardcoded-attribute-idp-mapper" {
		return false, nil
	} else {
		return false, fmt.Errorf(`provider.keycloak: keycloak_hardcoded_attribute_identity_provider_mapper: mapper type "%s" is not valid`, mapperType)
	}
}
