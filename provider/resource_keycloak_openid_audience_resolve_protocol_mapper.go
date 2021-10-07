package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenIdAudienceResolveProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdAudienceResolveProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdAudienceResolveProtocolMapperRead,
		//Update: resourceKeycloakOpenIdAudienceResolveProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdAudienceResolveProtocolMapperDelete,
		Importer: &schema.ResourceImporter{
			// import a mapper tied to a client:
			// {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
			// or a client scope:
			// {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}
			State: genericProtocolMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A human-friendly name that will appear in the Keycloak console.",
				Default:     "audience resolve",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
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
		},
	}
}

func mapFromDataToOpenIdAudienceResolveProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdAudienceResolveProtocolMapper {
	return &keycloak.OpenIdAudienceResolveProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),
	}
}

func mapFromOpenIdAudienceResolveMapperToData(mapper *keycloak.OpenIdAudienceResolveProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}
}

func resourceKeycloakOpenIdAudienceResolveProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdAudienceResolveMapper := mapFromDataToOpenIdAudienceResolveProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdAudienceResolveProtocolMapper(openIdAudienceResolveMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenIdAudienceResolveProtocolMapper(openIdAudienceResolveMapper)
	if err != nil {
		return err
	}

	mapFromOpenIdAudienceResolveMapperToData(openIdAudienceResolveMapper, data)

	return resourceKeycloakOpenIdAudienceResolveProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdAudienceResolveProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdAudienceResolveMapper, err := keycloakClient.GetOpenIdAudienceResolveProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdAudienceResolveMapperToData(openIdAudienceResolveMapper, data)

	return nil
}

func resourceKeycloakOpenIdAudienceResolveProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteOpenIdAudienceResolveProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
