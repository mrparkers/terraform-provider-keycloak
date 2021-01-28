package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlScriptProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlScriptProtocolMapperCreate,
		Read:   resourceKeycloakSamlScriptProtocolMapperRead,
		Update: resourceKeycloakSamlScriptProtocolMapperUpdate,
		Delete: resourceKeycloakSamlScriptProtocolMapperDelete,
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
			"single_value_attribute": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"script": {
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

func mapFromDataToSamlScriptProtocolMapper(data *schema.ResourceData) *keycloak.SamlScriptProtocolMapper {
	return &keycloak.SamlScriptProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		SingleValueAttribute: data.Get("single_value_attribute").(bool),

		SamlScript:              data.Get("script").(string),
		FriendlyName:            data.Get("friendly_name").(string),
		SamlAttributeName:       data.Get("saml_attribute_name").(string),
		SamlAttributeNameFormat: data.Get("saml_attribute_name_format").(string),
	}
}

func mapFromSamlScriptMapperToData(mapper *keycloak.SamlScriptProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("single_value_attribute", mapper.SingleValueAttribute)
	data.Set("script", mapper.SamlScript)
	data.Set("friendly_name", mapper.FriendlyName)
	data.Set("saml_attribute_name", mapper.SamlAttributeName)
	data.Set("saml_attribute_name_format", mapper.SamlAttributeNameFormat)
}

func resourceKeycloakSamlScriptProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlScriptMapper := mapFromDataToSamlScriptProtocolMapper(data)

	err := keycloakClient.ValidateSamlScriptProtocolMapper(samlScriptMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewSamlScriptProtocolMapper(samlScriptMapper)
	if err != nil {
		return err
	}

	mapFromSamlScriptMapperToData(samlScriptMapper, data)

	return resourceKeycloakSamlScriptProtocolMapperRead(data, meta)
}

func resourceKeycloakSamlScriptProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	samlScriptMapper, err := keycloakClient.GetSamlScriptProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromSamlScriptMapperToData(samlScriptMapper, data)

	return nil
}

func resourceKeycloakSamlScriptProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlScriptMapper := mapFromDataToSamlScriptProtocolMapper(data)

	err := keycloakClient.ValidateSamlScriptProtocolMapper(samlScriptMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateSamlScriptProtocolMapper(samlScriptMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakSamlScriptProtocolMapperRead(data, meta)
}

func resourceKeycloakSamlScriptProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteSamlScriptProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
