package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakHCIdentityImporterIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"proxy_id_user_attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the user attribute where the Proxy Id retrieved from API-C is stored. This will be used to configure User Attribute Mapper in appropriate Client Scope to inject it into JWT",
		},
		"health_cloud_id_user_attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the user attribute where the Health Cloud ID retrieved from Identity API is stored. This will be used to configure User Attribute Mapper in appropriate Client Scope to inject it into JWT",
		},
		"hc_app_id_user_attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute Name for Health Cloud Application Specific ID",
		},
		"proxy_id_api_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter Proxy Id API Url. This is the API used to retrieve the Proxy Id",
		},
		"health_cloud_id_api_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enter Health Cloud Id API Url. This is the API used to retrieve the Health Cloud Id API",
		},
		"tenant": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Tenant",
		},
		"source_code": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Health Cloud Id API Source Code",
		},
		"source_uid_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Source UID Name",
		},
		"idp_id_key_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IDP provided id key name that will be used to retrieve ID Value for the use with Proxy Id API Url. this value will be appended to Proxy Id API Url value before making the API call.",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getHCIdentityImporterIdentityProviderMapperFromData, setHCIdentityImporterIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setHCIdentityImporterIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getHCIdentityImporterIdentityProviderMapperFromData, setHCIdentityImporterIdentityProviderMapperData)
	return genericMapperResource
}

func getHCIdentityImporterIdentityProviderMapperFromData(data *schema.ResourceData, _ interface{}) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "hc-identity-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		ProxyIdUserAttributeName:       data.Get("proxy_id_user_attribute_name").(string),
		HealthCloudIdUserAttributeName: data.Get("health_cloud_id_user_attribute_name").(string),
		HcAppIdUserAttributeName:       data.Get("hc_app_id_user_attribute_name").(string),
		ProxyIdApiUrl:                  data.Get("proxy_id_api_url").(string),
		HealthCloudIdApiUrl:            data.Get("health_cloud_id_api_url").(string),
		Tenant:                         data.Get("tenant").(string),
		SourceCode:                     data.Get("source_code").(string),
		SourceUidName:                  data.Get("source_uid_name").(string),
		IdpIdKeyName:                   data.Get("idp_id_key_name").(string),
	}
	return rec, nil
}

func setHCIdentityImporterIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("proxy_id_user_attribute_name", identityProviderMapper.Config.ProxyIdUserAttributeName)
	data.Set("health_cloud_id_user_attribute_name", identityProviderMapper.Config.HealthCloudIdUserAttributeName)
	data.Set("hc_app_id_user_attribute_name", identityProviderMapper.Config.HcAppIdUserAttributeName)
	data.Set("proxy_id_api_url", identityProviderMapper.Config.ProxyIdApiUrl)
	data.Set("health_cloud_id_api_url", identityProviderMapper.Config.HealthCloudIdApiUrl)
	data.Set("tenant", identityProviderMapper.Config.Tenant)
	data.Set("source_code", identityProviderMapper.Config.SourceCode)
	data.Set("source_uid_name", identityProviderMapper.Config.SourceUidName)
	data.Set("idp_id_key_name", identityProviderMapper.Config.IdpIdKeyName)
	return nil
}
