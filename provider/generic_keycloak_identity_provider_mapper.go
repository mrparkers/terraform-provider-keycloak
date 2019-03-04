package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var getIdentityProviderMapperDataFunctions = map[string]func(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error){
	"hardcoded-attribute-idp-mapper":              getHardcodedAttributeIdpMapperFromData,
	"hardcoded-user-session-attribute-idp-mapper": getHardcodedUserSessionAttributeIdpMapperFromData,
	"oidc-hardcoded-role-idp-mapper":              getOidcHardcodedRoleIdpMapperFromData,
	"oidc-role-idp-mapper":                        getOidcRoleIdpMapperFromData,
	"oidc-user-attribute-idp-mapper":              getOidcUserAttributeIdpMapperFromData,
	"oidc-username-idp-mapper":                    getOidcUsernameIdpMapperFromData,
	"saml-hardcoded-role-idp-mapper":              getSamlHardcodedRoleIdpMapperFromData,
	"saml-role-idp-mapper":                        getSamlRoleIdpMapperFromData,
	"saml-user-attribute-idp-mapper":              getSamlUserAttributeIdpMapperFromData,
	"saml-username-idp-mapper":                    getSamlUsernameIdpMapperFromData,
}

var setIdentityProviderMapperDataFunctions = map[string]func(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error{
	"hardcoded-attribute-idp-mapper":              setHardcodedAttributeIdpMapperData,
	"hardcoded-user-session-attribute-idp-mapper": setHardcodedUserSessionAttributeIdpMapperData,
	"oidc-hardcoded-role-idp-mapper":              setOidcHardcodedRoleIdpMapperData,
	"oidc-role-idp-mapper":                        setOidcRoleIdpMapperData,
	"oidc-user-attribute-idp-mapper":              setOidcUserAttributeIdpMapperData,
	"oidc-username-idp-mapper":                    setOidcUsernameIdpMapperData,
	"saml-hardcoded-role-idp-mapper":              setSamlHardcodedRoleIdpMapperData,
	"saml-role-idp-mapper":                        setSamlRoleIdpMapperData,
	"saml-user-attribute-idp-mapper":              setSamlUserAttributeIdpMapperData,
	"saml-username-idp-mapper":                    setSamlUsernameIdpMapperData,
}

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
	data.Set("id", identityProviderMapper.Id)
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

func resourceKeycloakIdentityProviderMapperCreate(providerId string) func(data *schema.ResourceData, meta interface{}) error {
	setIdentityProviderMapperDataFunction := setIdentityProviderMapperDataFunctions[providerId]
	getIdentityProviderMapperDataFunction := getIdentityProviderMapperDataFunctions[providerId]
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		identityProvider, err := getIdentityProviderMapperDataFunction(data)
		err = keycloakClient.NewIdentityProviderMapper(identityProvider)
		if err != nil {
			return err
		}
		setIdentityProviderMapperDataFunction(data, identityProvider)
		return resourceKeycloakIdentityProviderMapperRead(providerId)(data, meta)
	}
}

func resourceKeycloakIdentityProviderMapperRead(providerId string) func(data *schema.ResourceData, meta interface{}) error {
	setIdentityProviderMapperDataFunction := setIdentityProviderMapperDataFunctions[providerId]
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		realm := data.Get("realm").(string)
		alias := data.Get("identity_provider_alias").(string)
		id := data.Id()
		identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
		if err != nil {
			return handleNotFoundError(err, data)
		}
		setIdentityProviderMapperDataFunction(data, identityProvider)
		return nil
	}
}

func resourceKeycloakIdentityProviderMapperUpdate(providerId string) func(data *schema.ResourceData, meta interface{}) error {
	setIdentityProviderMapperDataFunction := setIdentityProviderMapperDataFunctions[providerId]
	getIdentityProviderMapperDataFunction := getIdentityProviderMapperDataFunctions[providerId]
	return func(data *schema.ResourceData, meta interface{}) error {
		keycloakClient := meta.(*keycloak.KeycloakClient)
		identityProvider, err := getIdentityProviderMapperDataFunction(data)
		err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
		if err != nil {
			return err
		}
		setIdentityProviderMapperDataFunction(data, identityProvider)
		return nil
	}
}
