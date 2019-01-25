package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenIdGroupMembershipProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdGroupMembershipProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdGroupMembershipProtocolMapperRead,
		Update: resourceKeycloakOpenIdGroupMembershipProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdGroupMembershipProtocolMapperDelete,
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
			"claim_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_path": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"add_to_id_token": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"add_to_access_token": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"add_to_userinfo": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func mapFromDataToOpenIdGroupMembershipProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdGroupMembershipProtocolMapper {
	return &keycloak.OpenIdGroupMembershipProtocolMapper{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		RealmId:       data.Get("realm_id").(string),
		ClientId:      data.Get("client_id").(string),
		ClientScopeId: data.Get("client_scope_id").(string),

		ClaimName:        data.Get("claim_name").(string),
		FullPath:         data.Get("full_path").(bool),
		AddToIdToken:     data.Get("add_to_id_token").(bool),
		AddToAccessToken: data.Get("add_to_access_token").(bool),
		AddToUserinfo:    data.Get("add_to_userinfo").(bool),
	}
}

func mapFromOpenIdGroupMembershipMapperToData(mapper *keycloak.OpenIdGroupMembershipProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("claim_name", mapper.ClaimName)
	data.Set("full_path", mapper.FullPath)
	data.Set("add_to_id_token", mapper.AddToIdToken)
	data.Set("add_to_access_token", mapper.AddToAccessToken)
	data.Set("add_to_userinfo", mapper.AddToUserinfo)
}

func resourceKeycloakOpenIdGroupMembershipProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdGroupMembershipMapper := mapFromDataToOpenIdGroupMembershipProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdGroupMembershipProtocolMapper(openIdGroupMembershipMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenIdGroupMembershipProtocolMapper(openIdGroupMembershipMapper)
	if err != nil {
		return err
	}

	mapFromOpenIdGroupMembershipMapperToData(openIdGroupMembershipMapper, data)

	return resourceKeycloakOpenIdGroupMembershipProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdGroupMembershipProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdGroupMembershipMapper, err := keycloakClient.GetOpenIdGroupMembershipProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdGroupMembershipMapperToData(openIdGroupMembershipMapper, data)

	return nil
}

func resourceKeycloakOpenIdGroupMembershipProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdGroupMembershipMapper := mapFromDataToOpenIdGroupMembershipProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdGroupMembershipProtocolMapper(openIdGroupMembershipMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenIdGroupMembershipProtocolMapper(openIdGroupMembershipMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdGroupMembershipProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdGroupMembershipProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteOpenIdGroupMembershipProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
