package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakCustomIdentityProviderMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakCustomIdentityProviderMapperCreate,
		Read:   resourceKeycloakCustomIdentityProviderMapperRead,
		Update: resourceKeycloakCustomIdentityProviderMapperUpdate,
		Delete: resourceKeycloakCustomIdentityProviderMapperDelete,
		Importer: &schema.ResourceImporter{
			// we can use the generic identity provider import func here
			State: resourceKeycloakIdentityProviderMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Realm Name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IDP Mapper Name",
			},
			"identity_provider_alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IDP Alias",
			},
			"identity_provider_mapper": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IDP Mapper Type",
			},
			"extra_config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func getCustomIdentityProviderMapperFromData(data *schema.ResourceData) *keycloak.CustomIdentityProviderMapper {
	return &keycloak.CustomIdentityProviderMapper{
		Id:                     data.Id(),
		Realm:                  data.Get("realm").(string),
		Name:                   data.Get("name").(string),
		IdentityProviderAlias:  data.Get("identity_provider_alias").(string),
		IdentityProviderMapper: data.Get("identity_provider_mapper").(string),
		Config: &keycloak.CustomIdentityProviderMapperConfig{
			ExtraConfig: getExtraConfigFromData(data),
		},
	}
}

func setCustomIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.CustomIdentityProviderMapper) {
	data.SetId(identityProviderMapper.Id)
	data.Set("realm", identityProviderMapper.Realm)
	data.Set("name", identityProviderMapper.Name)
	data.Set("identity_provider_alias", identityProviderMapper.IdentityProviderAlias)
	data.Set("identity_provider_mapper", identityProviderMapper.IdentityProviderMapper)

	setExtraConfigData(data, identityProviderMapper.Config.ExtraConfig)
}

func resourceKeycloakCustomIdentityProviderMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	customIdentityProvider := getCustomIdentityProviderMapperFromData(data)

	err := keycloakClient.NewCustomIdentityProviderMapper(customIdentityProvider)
	if err != nil {
		return err
	}

	setCustomIdentityProviderMapperData(data, customIdentityProvider)

	return resourceKeycloakCustomIdentityProviderMapperRead(data, meta)
}

func resourceKeycloakCustomIdentityProviderMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()

	customIdentityProvider, err := keycloakClient.GetCustomIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setCustomIdentityProviderMapperData(data, customIdentityProvider)

	return nil
}

func resourceKeycloakCustomIdentityProviderMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	customIdentityProvider := getCustomIdentityProviderMapperFromData(data)

	err := keycloakClient.UpdateCustomIdentityProviderMapper(customIdentityProvider)
	if err != nil {
		return err
	}

	setCustomIdentityProviderMapperData(data, customIdentityProvider)

	return nil
}

func resourceKeycloakCustomIdentityProviderMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()

	return keycloakClient.DeleteCustomIdentityProviderMapper(realm, alias, id)
}
