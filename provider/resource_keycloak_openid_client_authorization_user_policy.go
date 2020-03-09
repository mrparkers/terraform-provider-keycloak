package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakPolicyLogicTypes, false),
			},
			"policies": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
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
	var policies []string
	var resources []string
	var scopes []string
	var users []string
	if v, ok := data.GetOk("resources"); ok {
		for _, resource := range v.(*schema.Set).List() {
			resources = append(resources, resource.(string))
		}
	}
	if v, ok := data.GetOk("policies"); ok {
		for _, policy := range v.(*schema.Set).List() {
			policies = append(policies, policy.(string))
		}
	}
	if v, ok := data.GetOk("scopes"); ok {
		for _, scope := range v.(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}
	}
	if v, ok := data.GetOk("users"); ok {
		for _, user := range v.(*schema.Set).List() {
			users = append(users, user.(string))
		}
	}

	resource := keycloak.OpenidClientAuthorizationUserPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		Owner:            data.Get("owner").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "user",
		Policies:         policies,
		Resources:        resources,
		Scopes:           scopes,
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
	data.Set("owner", policy.Owner)
	data.Set("logic", policy.Logic)
	data.Set("policies", policy.Policies)
	data.Set("resources", policy.Resources)
	data.Set("scopes", policy.Scopes)
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
