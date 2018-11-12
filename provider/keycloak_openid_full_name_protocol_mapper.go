package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenIdFullNameProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdFullNameProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdFullNameProtocolMapperRead,
		Update: resourceKeycloakOpenIdFullNameProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdFullNameProtocolMapperDelete,
		Importer: &schema.ResourceImporter{
			// import a mapper tied to a client:
			// {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
			// or a client scope:
			// {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}
			State: resourceKeycloakOpenIdFullNameProtocolMapperImport,
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
				Optional:    true,
				ForceNew:    true,
				Description: "The realm id where the associated client or client scope exists.",
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The mapper's associated client. Cannot be used at the same time as client_scope_id.",
				ConflictsWith: []string{"client_scope_id"},
			},
			"client_scope_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The mapper's associated client scope. Cannot be used at the same time as client_id.",
				ConflictsWith: []string{"client_id"},
			},
			"id_token_claim": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"access_token_claim": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"userinfo_token_claim": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func mapFromDataToOpenIdFullNameProtocolMapper(data *schema.ResourceData, client *keycloak.KeycloakClient) *keycloak.OpenIdFullNameProtocolMapper {

	var realmId string
	if v, ok := data.GetOk("realm_id"); ok {
		realmId = v.(string)
	} else {
		realmId = client.RealmId
	}

	return &keycloak.OpenIdFullNameProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       realmId,
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		AddToIdToken:     data.Get("id_token_claim").(bool),
		AddToAccessToken: data.Get("access_token_claim").(bool),
		AddToUserInfo:    data.Get("userinfo_token_claim").(bool),
	}
}

func mapFromOpenIdFullNameMapperToData(mapper *keycloak.OpenIdFullNameProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("id_token_claim", mapper.AddToIdToken)
	data.Set("access_token_claim", mapper.AddToAccessToken)
	data.Set("userinfo_token_claim", mapper.AddToUserInfo)
}

func resourceKeycloakOpenIdFullNameProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdFullNameMapper := mapFromDataToOpenIdFullNameProtocolMapper(data, keycloakClient)

	err := keycloakClient.ValidateOpenIdFullNameProtocolMapper(openIdFullNameMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenIdFullNameProtocolMapper(openIdFullNameMapper)
	if err != nil {
		return err
	}

	mapFromOpenIdFullNameMapperToData(openIdFullNameMapper, data)

	return resourceKeycloakOpenIdFullNameProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdFullNameProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdFullNameMapper, err := keycloakClient.GetOpenIdFullNameProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdFullNameMapperToData(openIdFullNameMapper, data)

	return nil
}

func resourceKeycloakOpenIdFullNameProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdFullNameMapper := mapFromDataToOpenIdFullNameProtocolMapper(data, keycloakClient)

	err := keycloakClient.ValidateOpenIdFullNameProtocolMapper(openIdFullNameMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenIdFullNameProtocolMapper(openIdFullNameMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdFullNameProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdFullNameProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteOpenIdFullNameProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}

func resourceKeycloakOpenIdFullNameProtocolMapperImport(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")
	keycloakClient := meta.(*keycloak.KeycloakClient)

	var realmId, parentResourceType, parentResourceId, mapperId string
	switch len(parts) {
	case 3:
		realmId = keycloakClient.RealmId
		parentResourceType = parts[1]
		parentResourceId = parts[2]
		mapperId = parts[3]
	case 4:
		realmId = parts[0]
		parentResourceType = parts[1]
		parentResourceId = parts[2]
		mapperId = parts[3]
	default:
		return nil, fmt.Errorf("invalid import. supported import formats: {{realmId}}/client/{{clientId}}/{{protocolMapperId}}, {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}, client-scope/{{clientScopeId}}/{{protocolMapperId}} or client/{{clientId}}/{{protocolMapperId}}")
	}

	data.Set("realm_id", realmId)
	data.SetId(mapperId)

	if parentResourceType == "client" {
		data.Set("client_id", parentResourceId)
	} else if parentResourceType == "client-scope" {
		data.Set("client_scope_id", parentResourceId)
	} else {
		return nil, fmt.Errorf("the associated parent resource must be either a client or a client-scope")
	}

	return []*schema.ResourceData{data}, nil
}
