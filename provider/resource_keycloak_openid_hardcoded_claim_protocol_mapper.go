package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenIdHardcodedClaimProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdHardcodedClaimProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdHardcodedClaimProtocolMapperRead,
		Update: resourceKeycloakOpenIdHardcodedClaimProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdHardcodedClaimProtocolMapperDelete,
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
			"add_to_id_token": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the attribute should be a claim in the id token.",
			},
			"add_to_access_token": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the attribute should be a claim in the access token.",
			},
			"add_to_userinfo": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the attribute should appear in the userinfo response body.",
			},
			"claim_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"claim_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"claim_value_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Claim type used when serializing tokens.",
				Default:      "String",
				ValidateFunc: validation.StringInSlice([]string{"String", "long", "int", "boolean"}, true),
			},
		},
	}
}

func mapFromDataToOpenIdHardcodedClaimProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdHardcodedClaimProtocolMapper {
	return &keycloak.OpenIdHardcodedClaimProtocolMapper{
		Id:               data.Id(),
		Name:             data.Get("name").(string),
		RealmId:          data.Get("realm_id").(string),
		ClientId:         data.Get("client_id").(string),
		ClientScopeId:    data.Get("client_scope_id").(string),
		AddToIdToken:     data.Get("add_to_id_token").(bool),
		AddToAccessToken: data.Get("add_to_access_token").(bool),
		AddToUserInfo:    data.Get("add_to_userinfo").(bool),

		ClaimName:      data.Get("claim_name").(string),
		ClaimValue:     data.Get("claim_value").(string),
		ClaimValueType: data.Get("claim_value_type").(string),
	}
}

func mapFromOpenIdHardcodedClaimMapperToData(mapper *keycloak.OpenIdHardcodedClaimProtocolMapper, data *schema.ResourceData) {
	data.SetId(mapper.Id)
	data.Set("name", mapper.Name)
	data.Set("realm_id", mapper.RealmId)

	if mapper.ClientId != "" {
		data.Set("client_id", mapper.ClientId)
	} else {
		data.Set("client_scope_id", mapper.ClientScopeId)
	}

	data.Set("add_to_id_token", mapper.AddToIdToken)
	data.Set("add_to_access_token", mapper.AddToAccessToken)
	data.Set("add_to_userinfo", mapper.AddToUserInfo)
	data.Set("claim_name", mapper.ClaimName)
	data.Set("claim_value", mapper.ClaimValue)
	data.Set("claim_value_type", mapper.ClaimValueType)
}

func resourceKeycloakOpenIdHardcodedClaimProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdHardcodedClaimMapper := mapFromDataToOpenIdHardcodedClaimProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdHardcodedClaimProtocolMapper(openIdHardcodedClaimMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenIdHardcodedClaimProtocolMapper(openIdHardcodedClaimMapper)
	if err != nil {
		return err
	}

	mapFromOpenIdHardcodedClaimMapperToData(openIdHardcodedClaimMapper, data)

	return resourceKeycloakOpenIdHardcodedClaimProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdHardcodedClaimProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdHardcodedClaimMapper, err := keycloakClient.GetOpenIdHardcodedClaimProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdHardcodedClaimMapperToData(openIdHardcodedClaimMapper, data)

	return nil
}

func resourceKeycloakOpenIdHardcodedClaimProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdHardcodedClaimMapper := mapFromDataToOpenIdHardcodedClaimProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdHardcodedClaimProtocolMapper(openIdHardcodedClaimMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenIdHardcodedClaimProtocolMapper(openIdHardcodedClaimMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdHardcodedClaimProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdHardcodedClaimProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteOpenIdHardcodedClaimProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
