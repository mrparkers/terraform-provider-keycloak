package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenIdPropertyMapperClaimProtocolMapper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperCreate,
		ReadContext:   resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperRead,
		UpdateContext: resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperUpdate,
		DeleteContext: resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperDelete,
		Importer: &schema.ResourceImporter{
			// import a mapper tied to a client:
			// {{realmId}}/client/{{clientId}}/{{protocolMapperId}}
			// or a client scope:
			// {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}
			StateContext: genericProtocolMapperImport,
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
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol type for the expected extra parameters.",
			},
			"protocol_mapper": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol property mapper type.",
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
			"add_to_introspection_token": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the attribute should be a claim in the introspect token.",
			},
			"add_to_lightweight_claim": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if the attribute should appear in the lightweight claim.",
			},
			"claim_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The claim name to display in the token.",
			},
			"json_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Claim type used when serializing tokens.",
				Default:      "String",
				ValidateFunc: validation.StringInSlice([]string{"JSON", "String", "long", "int", "boolean"}, true),
			},
			"set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Mapper values to be merged with the other attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func mapFromDataToOpenIdPropertyMapperClaimProtocolMapper(data *schema.ResourceData) *keycloak.OpenIdPropertyMapperClaimProtocolMapper {

	additionalConfig := map[string]string{}
	setlist := data.Get("set").(*schema.Set).List()

	for _, raw := range setlist {
		set := raw.(map[string]interface{})
		additionalConfig[set["name"].(string)] = set["value"].(string)
	}

	return &keycloak.OpenIdPropertyMapperClaimProtocolMapper{
		Id:                      data.Id(),
		Name:                    data.Get("name").(string),
		RealmId:                 data.Get("realm_id").(string),
		ClientId:                data.Get("client_id").(string),
		ClientScopeId:           data.Get("client_scope_id").(string),
		AddToIdToken:            data.Get("add_to_id_token").(bool),
		AddToAccessToken:        data.Get("add_to_access_token").(bool),
		AddToUserInfo:           data.Get("add_to_userinfo").(bool),
		AddToIntrospectionToken: data.Get("add_to_introspection_token").(bool),
		AddToLightweightClaim:   data.Get("add_to_lightweight_claim").(bool),

		Protocol:       data.Get("protocol").(string),
		ProtocolMapper: data.Get("protocol_mapper").(string),
		ClaimName:      data.Get("claim_name").(string),
		JsonType:       data.Get("json_type").(string),

		AdditionalConfig: additionalConfig,
	}
}

func mapFromOpenIdPropertyMapperClaimMapperToData(mapper *keycloak.OpenIdPropertyMapperClaimProtocolMapper, data *schema.ResourceData) {
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
	data.Set("add_to_introspection_token", mapper.AddToIntrospectionToken)
	data.Set("add_to_lightweight_claim", mapper.AddToLightweightClaim)

	data.Set("protocol", mapper.Protocol)
	data.Set("protocol_mapper", mapper.ProtocolMapper)
	data.Set("claim_name", mapper.ClaimName)
	data.Set("json_type", mapper.JsonType)

	additionalConfig := make([]interface{}, 0, len(mapper.AdditionalConfig))
	for k, v := range mapper.AdditionalConfig {
		item := map[string]interface{}{
			"name":  k,
			"value": v,
		}
		additionalConfig = append(additionalConfig, item)
	}
	data.Set("set", additionalConfig)
}

func resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdPropertyMapperClaimMapper := mapFromDataToOpenIdPropertyMapperClaimProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdPropertyMapperClaimProtocolMapper(ctx, openIdPropertyMapperClaimMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.NewOpenIdPropertyMapperClaimProtocolMapper(ctx, openIdPropertyMapperClaimMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperCreate")

	mapFromOpenIdPropertyMapperClaimMapperToData(openIdPropertyMapperClaimMapper, data)

	return resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperRead(ctx, data, meta)
}

func resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	openIdPropertyMapperClaimMapper, err := keycloakClient.GetOpenIdPropertyMapperClaimProtocolMapper(ctx, realmId, clientId, clientScopeId, data.Id())
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}
	mapFromOpenIdPropertyMapperClaimMapperToData(openIdPropertyMapperClaimMapper, data)

	return nil
}

func resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	openIdPropertyMapperClaimMapper := mapFromDataToOpenIdPropertyMapperClaimProtocolMapper(data)

	err := keycloakClient.ValidateOpenIdPropertyMapperClaimProtocolMapper(ctx, openIdPropertyMapperClaimMapper)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.UpdateOpenIdPropertyMapperClaimProtocolMapper(ctx, openIdPropertyMapperClaimMapper)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperRead(ctx, data, meta)
}

func resourceKeycloakOpenIdPropertyMapperClaimProtocolMapperDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.DeleteOpenIdPropertyMapperClaimProtocolMapper(ctx, realmId, clientId, clientScopeId, data.Id()))
}
