package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenidClientAuthorizationGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationGroupPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationGroupPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationGroupPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationGroupPolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientAuthorizationGroupPolicyImport,
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
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_claim": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func getOpenidClientAuthorizationGroupPolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationGroupPolicy {
	var policies []string
	var resources []string
	var scopes []string
	var groups []map[string]interface{}
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
	if v, ok := data.GetOk("groups"); ok {
		for _, group := range v.(*schema.Set).List() {
			attributes := map[string]interface{}{}
			for key, value := range group.(map[string]interface{}) {
				attributes[key] = strings.Split(value.(string), ",")
			}
			groups = append(groups, attributes)
		}
	}

	resource := keycloak.OpenidClientAuthorizationGroupPolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		Owner:            data.Get("owner").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "group",
		Policies:         policies,
		Resources:        resources,
		Scopes:           scopes,
		GroupsClaim:      data.Get("groups_claim").(string),
		Groups:           groups,
		Description:      data.Get("description").(string),
	}
	return &resource
}

func setOpenidClientAuthorizationGroupPolicyResourceData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationGroupPolicy) {
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
	data.Set("type", policy.Type)
	data.Set("description", policy.Description)
	data.Set("groups_claim", policy.GroupsClaim)
	data.Set("groups", policy.Groups)
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationGroupPolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationGroupPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationGroupPolicyResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationResourceRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationGroupPolicy(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationGroupPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationGroupPolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationGroupPolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationGroupPolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationGroupPolicy(realmId, resourceServerId, id)
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
