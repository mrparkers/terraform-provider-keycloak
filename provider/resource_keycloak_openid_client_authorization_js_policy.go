package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationJSPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationJSPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationJSPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationJSPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationJSPolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientAuthorizationJSPolicyImport,
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
				Type:     schema.TypeString,
				Required: true,
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
	var policies []string
	var resources []string
	var scopes []string
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

	resource := keycloak.OpenidClientAuthorizationJSPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		Owner:            data.Get("owner").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "js",
		Policies:         policies,
		Resources:        resources,
		Scopes:           scopes,
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
	data.Set("owner", policy.Owner)
	data.Set("logic", policy.Logic)
	data.Set("policies", policy.Policies)
	data.Set("resources", policy.Resources)
	data.Set("scopes", policy.Scopes)
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

func resourceKeycloakOpenidClientAuthorizationJSPolicyImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
