package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationClientPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationClientPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationClientPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationClientPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationClientPolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientAuthorizationClientPolicyImport,
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
	var policies []string
	var resources []string
	var scopes []string
	var clients []string
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
	if v, ok := data.GetOk("clients"); ok {
		for _, client := range v.(*schema.Set).List() {
			clients = append(clients, client.(string))
		}
	}

	resource := keycloak.OpenidClientAuthorizationClientPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		Owner:            data.Get("owner").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "client",
		Policies:         policies,
		Resources:        resources,
		Scopes:           scopes,
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
	data.Set("owner", policy.Owner)
	data.Set("logic", policy.Logic)
	data.Set("policies", policy.Policies)
	data.Set("resources", policy.Resources)
	data.Set("scopes", policy.Scopes)
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

func resourceKeycloakOpenidClientAuthorizationClientPolicyImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
