package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationClientPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationClientPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationClientPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationClientPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationClientPolicyUpdate,
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
				Optional: true,
			},
			"logic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakPolicyLogicTypes, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clients": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func getOpenidClientAuthorizationClientAuthorizationClientPolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationClientPolicy {
	var clients []string

	if v, ok := data.GetOk("clients"); ok {
		for _, client := range v.(*schema.Set).List() {
			clients = append(clients, client.(string))
		}
	}

	resource := keycloak.OpenidClientAuthorizationClientPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "client",
		Clients:          clients,
		Description:      data.Get("description").(string),
	}
	return &resource
}

func setOpenidClientAuthorizationClientAuthorizationClientPolicyResourceData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationClientPolicy) {
	data.SetId(policy.Id)

	data.Set("resource_server_id", policy.ResourceServerId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("logic", policy.Logic)
	data.Set("description", policy.Description)
	data.Set("clients", policy.Clients)
}

func resourceKeycloakOpenidClientAuthorizationClientPolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationClientAuthorizationClientPolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationClientPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationClientAuthorizationClientPolicyResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationClientPolicyRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationClientPolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationClientPolicy(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationClientAuthorizationClientPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationClientPolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationClientAuthorizationClientPolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationClientPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationClientAuthorizationClientPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationClientPolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationClientPolicy(realmId, resourceServerId, id)
}
