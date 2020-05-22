package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenIdUserClientRoleProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdUserClientRoleProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdUserClientRoleProtocolMapperRead,
		Update: resourceKeycloakOpenIdUserClientRoleProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdUserClientRoleProtocolMapperDelete,
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
			"claim_value_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Claim type used when serializing tokens.",
				Default:      "String",
				ValidateFunc: validation.StringInSlice([]string{"JSON", "String", "long", "int", "boolean"}, true),
			},
			"multivalued": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether this attribute is a single value or an array of values.",
			},
			"client_id_for_role_mappings": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client ID for role mappings.",
			},
			"client_role_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix that will be added to each client role.",
			},
		},
	}
}

func mapFromDataToOpenIdUserClientRoleProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdUserClientRoleProtocolMapper {
	return &keycloak.OpenIdUserClientRoleProtocolMapper{
		Id:               data.Id(),
		Name:             data.Get("name").(string),
		RealmId:          data.Get("realm_id").(string),
		ClientId:         data.Get("client_id").(string),
		ClientScopeId:    data.Get("client_scope_id").(string),
		AddToIdToken:     data.Get("add_to_id_token").(bool),
		AddToAccessToken: data.Get("add_to_access_token").(bool),
		AddToUserInfo:    data.Get("add_to_userinfo").(bool),

		ClaimName:               data.Get("claim_name").(string),
		ClaimValueType:          data.Get("claim_value_type").(string),
		Multivalued:             data.Get("multivalued").(bool),
		ClientIdForRoleMappings: data.Get("client_id_for_role_mappings").(string),
		ClientRolePrefix:        data.Get("client_role_prefix").(string),
	}
}

func mapFromOpenIdUserClientRoleMapperToData(mapper *keycloak.OpenIdUserClientRoleProtocolMapper, data *schema.ResourceData) {
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
	data.Set("claim_value_type", mapper.ClaimValueType)
	data.Set("multivalued", mapper.Multivalued)
	data.Set("client_id_for_role_mappings", mapper.ClientIdForRoleMappings)
	data.Set("client_role_prefix", mapper.ClientRolePrefix)
}

func resourceKeycloakOpenIdUserClientRoleProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdUserClientRoleMapper := mapFromDataToOpenIdUserClientRoleProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdUserClientRoleProtocolMapper(openIdUserClientRoleMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenIdUserClientRoleProtocolMapper(openIdUserClientRoleMapper)
	if err != nil {
		return err
	}

	mapFromOpenIdUserClientRoleMapperToData(openIdUserClientRoleMapper, data)

	return resourceKeycloakOpenIdUserClientRoleProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdUserClientRoleProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdUserClientRoleMapper, err := keycloakClient.GetOpenIdUserClientRoleProtocolMapper(realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdUserClientRoleMapperToData(openIdUserClientRoleMapper, data)

	return nil
}

func resourceKeycloakOpenIdUserClientRoleProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdUserClientRoleMapper := mapFromDataToOpenIdUserClientRoleProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdUserClientRoleProtocolMapper(openIdUserClientRoleMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenIdUserClientRoleProtocolMapper(openIdUserClientRoleMapper)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdUserClientRoleProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdUserClientRoleProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return keycloakClient.DeleteOpenIdUserClientRoleProtocolMapper(realmId, clientId, clientScopeId, data.Id())
}
