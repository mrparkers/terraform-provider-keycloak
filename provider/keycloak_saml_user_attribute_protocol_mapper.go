package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var keycloakSamlUserAttributeProtocolMapperNameFormats = []string{"Basic", "URI Reference", "Unspecified"}

func resourceKeycloakSamlUserAttributeProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlUserAttributeProtocolMapperCreate,
		Read:   resourceKeycloakSamlUserAttributeProtocolMapperRead,
		Update: resourceKeycloakSamlUserAttributeProtocolMapperUpdate,
		Delete: resourceKeycloakSamlUserAttributeProtocolMapperDelete,
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
			"user_attribute": {
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

func mapFromDataToSamlUserAttributeProtocolMapper(data *schema.ResourceData) *keycloak.SamlUserAttributeProtocolMapper {
	return &keycloak.SamlUserAttributeProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		UserAttribute:           data.Get("user_attribute").(string),
		FriendlyName:            data.Get("friendly_name").(string),
		SamlAttributeName:       data.Get("saml_attribute_name").(string),
		SamlAttributeNameFormat: data.Get("saml_attribute_name_format").(string),
	}
}

func mapFromSamlUserAttributeMapperToData(mapper *keycloak.SamlUserAttributeProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("user_attribute", mapper.UserAttribute)
	data.Set("friendly_name", mapper.FriendlyName)
	data.Set("saml_attribute_name", mapper.SamlAttributeName)
	data.Set("saml_attribute_name_format", mapper.SamlAttributeNameFormat)
}

func resourceKeycloakSamlUserAttributeProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlUserAttributeMapper := mapFromDataToSamlUserAttributeProtocolMapper(data)

	err := keycloakClient.ValidateSamlUserAttributeProtocolMapper(samlUserAttributeMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewSamlUserAttributeProtocolMapper(samlUserAttributeMapper)
	if err != nil {
		return err
	}

	mapFromSamlUserAttributeMapperToData(samlUserAttributeMapper, data)

	return resourceKeycloakSamlUserAttributeProtocolMapperRead(data, meta)
}

func resourceKeycloakSamlUserAttributeProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	samlUserAttributeMapper, err := keycloakClient.GetSamlUserAttributeProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromSamlUserAttributeMapperToData(samlUserAttributeMapper, data)

	return nil
}

func resourceKeycloakSamlUserAttributeProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	samlUserAttributeMapper := mapFromDataToSamlUserAttributeProtocolMapper(data)

	err := keycloakClient.ValidateSamlUserAttributeProtocolMapper(samlUserAttributeMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateSamlUserAttributeProtocolMapper(samlUserAttributeMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakSamlUserAttributeProtocolMapperRead(data, meta)
}

func resourceKeycloakSamlUserAttributeProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteSamlUserAttributeProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
