package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"log"
)

func resourceKeycloakGenericClientProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakGenericClientProtocolMapperCreate,
		Read:   resourceKeycloakGenericClientProtocolMapperRead,
		Delete: resourceKeycloakGenericClientProtocolMapperDelete,
		Update: resourceKeycloakGenericClientProtocolMapperUpdate,
		// This resource can be imported using {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
		Importer: &schema.ResourceImporter{
			State: genericProtocolMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-friendly name that will appear in the Keycloak console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm id where the associated client exists.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The mapper's associated client.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The protocol of the client (openid-connect / saml).",
				ValidateFunc: validation.StringInSlice([]string{"openid-connect", "saml"}, false),
			},
			"protocol_mapper": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the protocol mapper.",
			},
			"config": {
				Type:     schema.TypeMap,
				Required: true,
			},
		},
	}
}

func getGenericClientProtocolMapperFromData(data *schema.ResourceData) *keycloak.GenericClientProtocolMapper {
	config := make(map[string]string)
	if v, ok := data.GetOk("config"); ok {
		for key, value := range v.(map[string]interface{}) {
			config[key] = value.(string)
		}
	}

	return &keycloak.GenericClientProtocolMapper{
		ClientId:       data.Get("client_id").(string),
		Config:         config,
		Id:             data.Id(),
		Name:           data.Get("name").(string),
		Protocol:       data.Get("protocol").(string),
		ProtocolMapper: data.Get("protocol_mapper").(string),
		RealmId:        data.Get("realm_id").(string),
	}
}

func setGenericClientProtocolMapperData(data *schema.ResourceData, resource *keycloak.GenericClientProtocolMapper) {
	data.SetId(resource.Id)
	data.Set("client_id", resource.ClientId)
	data.Set("config", resource.Config)
	data.Set("name", resource.Name)
	data.Set("protocol", resource.Protocol)
	data.Set("protocol_mapper", resource.ProtocolMapper)
	data.Set("realm_id", resource.RealmId)
}

func resourceKeycloakGenericClientProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getGenericClientProtocolMapperFromData(data)

	err := keycloakClient.NewGenericClientProtocolMapper(resource)
	if err != nil {
		return err
	}
	setGenericClientProtocolMapperData(data, resource)

	return resourceKeycloakGenericClientProtocolMapperRead(data, meta)
}

func resourceKeycloakGenericClientProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetGenericClientProtocolMapper(realmId, clientId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setGenericClientProtocolMapperData(data, resource)

	return nil
}

func resourceKeycloakGenericClientProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] updating\n")
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getGenericClientProtocolMapperFromData(data)

	err := keycloakClient.UpdateGenericClientProtocolMapper(resource)
	if err != nil {
		return err
	}

	setGenericClientProtocolMapperData(data, resource)

	return nil
}

func resourceKeycloakGenericClientProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	id := data.Id()

	return keycloakClient.DeleteGenericClientProtocolMapper(realmId, clientId, id)
}
