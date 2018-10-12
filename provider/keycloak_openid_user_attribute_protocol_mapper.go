package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenIdUserAttributeProtocolMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdUserAttributeProtocolMapperCreate,
		Read:   resourceKeycloakOpenIdUserAttributeProtocolMapperRead,
		Update: resourceKeycloakOpenIdUserAttributeProtocolMapperUpdate,
		Delete: resourceKeycloakOpenIdUserAttributeProtocolMapperDelete,
		Importer: &schema.ResourceImporter{
			// import a mapper tied to a client:
			// {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
			// or a client scope:
			// {{realmId}}/client-scope/{{clientScopeId}}/protocolMapperId
			State: resourceKeycloakOpenIdUserAttributeProtocolMapperImport,
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
			"add_to_user_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the attribute should appear in the userinfo response body.",
			},
			"multivalued": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether this attribute is a single value or an array of values.",
			},
			"user_attribute": {
				Type:     schema.TypeString,
				Required: true,
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
				ValidateFunc: validation.StringInSlice([]string{"String", "long", "int", "boolean"}, true),
			},
		},
	}
}

func mapFromDataToOpenIdUserAttributeProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdUserAttributeProtocolMapper {
	return &keycloak.OpenIdUserAttributeProtocolMapper{
		Id:               data.Id(),
		Name:             data.Get("name").(string),
		RealmId:          data.Get("realm_id").(string),
		ClientId:         data.Get("client_id").(string),
		ClientScopeId:    data.Get("client_scope_id").(string),
		AddToIdToken:     data.Get("add_to_id_token").(bool),
		AddToAccessToken: data.Get("add_to_access_token").(bool),
		AddToUserInfo:    data.Get("add_to_user_info").(bool),

		UserAttribute:  data.Get("user_attribute").(string),
		ClaimName:      data.Get("claim_name").(string),
		ClaimValueType: data.Get("claim_value_type").(string),
	}
}

func mapFromOpenIdUserAttributeMapperToData(mapper *keycloak.OpenIdUserAttributeProtocolMapper, data *schema.ResourceData) {
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
	data.Set("add_to_user_info", mapper.AddToUserInfo)
	data.Set("user_attribute", mapper.UserAttribute)
	data.Set("claim_name", mapper.ClaimName)
	data.Set("claim_value_type", mapper.ClaimValueType)
}

func resourceKeycloakOpenIdUserAttributeProtocolMapperCreate(data *schema.ResourceData, meta interface{}) error {
	clientId := data.Get("client_id")
	clientScopeId := data.Get("client_scope_id")

	if clientId == "" && clientScopeId == "" {
		return fmt.Errorf("one of client_id or client_scope_id must be set")
	}

	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdUserAttributeMapper := mapFromDataToOpenIdUserAttributeProtocolMapper(data)

	err := keycloakClient.NewOpenIdUserAttributeProtocolMapper(openIdUserAttributeMapper)

	if err != nil {
		return err
	}

	mapFromOpenIdUserAttributeMapperToData(openIdUserAttributeMapper, data)

	return resourceKeycloakOpenIdUserAttributeProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdUserAttributeProtocolMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	var openIdUserAttributeMapper *keycloak.OpenIdUserAttributeProtocolMapper
	var err error
	if clientId != "" {
		openIdUserAttributeMapper, err = keycloakClient.GetOpenIdUserAttributeProtocolMapperForClient(realmId, clientId, data.Id())
	} else {
		openIdUserAttributeMapper, err = keycloakClient.GetOpenIdUserAttributeProtocolMapperForClientScope(realmId, clientScopeId, data.Id())
	}

	if err != nil {
		return err
	}

	mapFromOpenIdUserAttributeMapperToData(openIdUserAttributeMapper, data)

	return nil
}

func resourceKeycloakOpenIdUserAttributeProtocolMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdUserAttributeMapper := mapFromDataToOpenIdUserAttributeProtocolMapper(data)
	err := keycloakClient.UpdateOpenIdUserAttributeProtocolMapper(openIdUserAttributeMapper)

	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdUserAttributeProtocolMapperRead(data, meta)
}

func resourceKeycloakOpenIdUserAttributeProtocolMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	if clientId != "" {
		return keycloakClient.DeleteOpenIdUserAttributeProtocolMapperForClient(realmId, clientId, data.Id())
	} else {
		return keycloakClient.DeleteOpenIdUserAttributeProtocolMapperForClientScope(realmId, clientScopeId, data.Id())
	}
}

func resourceKeycloakOpenIdUserAttributeProtocolMapperImport(data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")

	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid import. supported import formats: {{realmId}}/client/{{clientId}}/protocolMapperId or {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}")
	}

	realmId := parts[0]
	parentResourceType := parts[1]
	parentResourceId := parts[2]
	mapperId := parts[3]

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
