package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

type identityProviderMapperDataGetterFunc func(data *schema.ResourceData, meta interface{}) (*keycloak.IdentityProviderMapper, error)
type identityProviderMapperDataSetterFunc func(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error

func resourceKeycloakIdentityProviderMapper() *schema.Resource {
	return &schema.Resource{
		Delete: resourceKeycloakIdentityProviderMapperDelete,
		Importer: &schema.ResourceImporter{
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
			"extra_config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func getIdentityProviderMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec := &keycloak.IdentityProviderMapper{
		Id:                    data.Id(),
		Realm:                 data.Get("realm").(string),
		Name:                  data.Get("name").(string),
		IdentityProviderAlias: data.Get("identity_provider_alias").(string),
	}
	return rec, nil
}

func setIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	data.SetId(identityProviderMapper.Id)
	data.Set("realm", identityProviderMapper.Realm)
	data.Set("name", identityProviderMapper.Name)
	data.Set("identity_provider_alias", identityProviderMapper.IdentityProviderAlias)
	return nil
}

func resourceKeycloakIdentityProviderMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()

	return keycloakClient.DeleteIdentityProviderMapper(realm, alias, id)
}

func resourceKeycloakIdentityProviderMapperImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realm}}/{{identityProviderAlias}}/{{identityProviderMapperId}}")
	}

	d.Set("realm", parts[0])
	d.Set("identity_provider_alias", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func resourceKeycloakIdentityProviderMapperCreate(getIdentityProviderMapperFromData identityProviderMapperDataGetterFunc, setDataFromIdentityProviderMapper identityProviderMapperDataSetterFunc) func(data *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		identityProvider, err := getIdentityProviderMapperFromData(data, meta)
		if err != nil {
			return err
		}
		if err = keycloakClient.NewIdentityProviderMapper(identityProvider); err != nil {
			return err
		}
		if err = setDataFromIdentityProviderMapper(data, identityProvider); err != nil {
			return err
		}
		return resourceKeycloakIdentityProviderMapperRead(setDataFromIdentityProviderMapper)(data, meta)
	}
}

func resourceKeycloakIdentityProviderMapperRead(setDataFromIdentityProviderMapper identityProviderMapperDataSetterFunc) func(data *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		realm := data.Get("realm").(string)
		alias := data.Get("identity_provider_alias").(string)
		id := data.Id()
		identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
		if err != nil {
			return handleNotFoundError(err, data)
		}
		if err = setDataFromIdentityProviderMapper(data, identityProvider); err != nil {
			return err
		}
		return nil
	}
}

func resourceKeycloakIdentityProviderMapperUpdate(getIdentityProviderMapperFromData identityProviderMapperDataGetterFunc, setDataFromIdentityProviderMapper identityProviderMapperDataSetterFunc) func(data *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		identityProvider, err := getIdentityProviderMapperFromData(data, meta)
		if err = keycloakClient.UpdateIdentityProviderMapper(identityProvider); err != nil {
			return err
		}
		if err = setDataFromIdentityProviderMapper(data, identityProvider); err != nil {
			return err
		}
		return nil
	}
}
