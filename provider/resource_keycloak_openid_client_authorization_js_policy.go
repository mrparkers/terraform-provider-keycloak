package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationJSPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationJSPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationJSPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationJSPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationJSPolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: genericResourcePolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"resource_server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"decision_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakPolicyLogicTypes, false),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func getOpenidClientAuthorizationJSPolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationJSPolicy {

	resource := keycloak.OpenidClientAuthorizationJSPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "js",
		Code:             data.Get("code").(string),
		Description:      data.Get("description").(string),
	}
	return &resource
}

func setOpenidClientAuthorizationJSPolicyResourceData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationJSPolicy) {
	data.SetId(policy.Id)

	data.Set("resource_server_id", policy.ResourceServerId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("logic", policy.Logic)
	data.Set("description", policy.Description)
	data.Set("code", policy.Code)
}

func resourceKeycloakOpenidClientAuthorizationJSPolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationJSPolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationJSPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationJSPolicyResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationJSPolicyRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationJSPolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationJSPolicy(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationJSPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationJSPolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationJSPolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationJSPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationJSPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationJSPolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationJSPolicy(realmId, resourceServerId, id)
}
