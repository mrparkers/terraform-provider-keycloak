package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationUserPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationUserPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationUserPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationUserPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationUserPolicyUpdate,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func getOpenidClientAuthorizationUserPolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationUserPolicy {
	var users []string
	if v, ok := data.GetOk("users"); ok {
		for _, user := range v.(*schema.Set).List() {
			users = append(users, user.(string))
		}
	}

	resource := keycloak.OpenidClientAuthorizationUserPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "user",
		Users:            users,
		Description:      data.Get("description").(string),
	}
	return &resource
}

func setOpenidClientAuthorizationUserPolicyResourceData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationUserPolicy) {
	data.SetId(policy.Id)

	data.Set("resource_server_id", policy.ResourceServerId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("logic", policy.Logic)
	data.Set("description", policy.Description)
	data.Set("users", policy.Users)
}

func resourceKeycloakOpenidClientAuthorizationUserPolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationUserPolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationUserPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationUserPolicyResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationUserPolicyRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationUserPolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationUserPolicy(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationUserPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationUserPolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationUserPolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationUserPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationUserPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationUserPolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationUserPolicy(realmId, resourceServerId, id)
}
