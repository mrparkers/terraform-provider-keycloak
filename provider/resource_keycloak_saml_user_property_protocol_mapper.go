package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlUserPropertyProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlUserPropertyProtocolMapperCreate,
		Read:   resourceKeycloakSamlUserPropertyProtocolMapperRead,
		Update: resourceKeycloakSamlUserPropertyProtocolMapperUpdate,
		Delete: resourceKeycloakSamlUserPropertyProtocolMapperDelete,
		Importer: &schema.ResourceImporter{
			// import a mapper tied to a client:
			// {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
			// or a client scope:
			// {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}
			State: genericProtocolMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"client_scope_id"},
			},
			"client_scope_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"client_id"},
			},
			"user_property": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"saml_attribute_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"saml_attribute_name_format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakSamlUserAttributeProtocolMapperNameFormats, false),
			},
		},
	}
}

func mapFromDataToSamlUserPropertyProtocolMapper(data *schema.ResourceData) *keycloak.SamlUserPropertyProtocolMapper {
	return &keycloak.SamlUserPropertyProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		UserProperty:            data.Get("user_property").(string),
		FriendlyName:            data.Get("friendly_name").(string),
		SamlAttributeName:       data.Get("saml_attribute_name").(string),
		SamlAttributeNameFormat: data.Get("saml_attribute_name_format").(string),
	}
}

func mapFromSamlUserPropertyProtocolMapperToData(mapper *keycloak.SamlUserPropertyProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("user_property", mapper.UserProperty)
	data.Set("friendly_name", mapper.FriendlyName)
	data.Set("saml_attribute_name", mapper.SamlAttributeName)
	data.Set("saml_attribute_name_format", mapper.SamlAttributeNameFormat)
}

func resourceKeycloakSamlUserPropertyProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlUserPropertyMapper := mapFromDataToSamlUserPropertyProtocolMapper(data)

	err := keycloakClient.ValidateSamlUserPropertyProtocolMapper(samlUserPropertyMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewSamlUserPropertyProtocolMapper(samlUserPropertyMapper)
	if err != nil {
		return err
	}

	mapFromSamlUserPropertyProtocolMapperToData(samlUserPropertyMapper, data)

	return resourceKeycloakSamlUserPropertyProtocolMapperRead(data, meta)
}

func resourceKeycloakSamlUserPropertyProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	samlUserPropertyMapper, err := keycloakClient.GetSamlUserPropertyProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromSamlUserPropertyProtocolMapperToData(samlUserPropertyMapper, data)

	return nil
}

func resourceKeycloakSamlUserPropertyProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlUserPropertyMapper := mapFromDataToSamlUserPropertyProtocolMapper(data)

	err := keycloakClient.ValidateSamlUserPropertyProtocolMapper(samlUserPropertyMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateSamlUserPropertyProtocolMapper(samlUserPropertyMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakSamlUserPropertyProtocolMapperRead(data, meta)
}

func resourceKeycloakSamlUserPropertyProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteSamlUserPropertyProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
