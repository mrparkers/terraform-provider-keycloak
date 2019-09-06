package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenIdHardcodedRoleProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdHardcodedRoleProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdHardcodedRoleProtocolMapperRead,
		Update: resourceKeycloakOpenIdHardcodedRoleProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdHardcodedRoleProtocolMapperDelete,
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
				Required:    true,
				ForceNew:    true,
				Description: "A human-friendly name that will appear in the Keycloak console.",
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
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func mapFromDataToOpenIdHardcodedRoleProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdHardcodedRoleProtocolMapper {
	return &keycloak.OpenIdHardcodedRoleProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		RoleId: data.Get("role_id").(string),
	}
}

func mapFromOpenIdHardcodedRoleMapperToData(mapper *keycloak.OpenIdHardcodedRoleProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("role_id", mapper.RoleId)
}

func resourceKeycloakOpenIdHardcodedRoleProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdHardcodedRoleMapper := mapFromDataToOpenIdHardcodedRoleProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdHardcodedRoleProtocolMapper(openIdHardcodedRoleMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenIdHardcodedRoleProtocolMapper(openIdHardcodedRoleMapper)
	if err != nil {
		return err
	}

	mapFromOpenIdHardcodedRoleMapperToData(openIdHardcodedRoleMapper, data)

	return resourceKeycloakOpenIdHardcodedRoleProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdHardcodedRoleProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdHardcodedRoleMapper, err := keycloakClient.GetOpenIdHardcodedRoleProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdHardcodedRoleMapperToData(openIdHardcodedRoleMapper, data)

	return nil
}

func resourceKeycloakOpenIdHardcodedRoleProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdHardcodedRoleMapper := mapFromDataToOpenIdHardcodedRoleProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdHardcodedRoleProtocolMapper(openIdHardcodedRoleMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenIdHardcodedRoleProtocolMapper(openIdHardcodedRoleMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdHardcodedRoleProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdHardcodedRoleProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteOpenIdHardcodedRoleProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
