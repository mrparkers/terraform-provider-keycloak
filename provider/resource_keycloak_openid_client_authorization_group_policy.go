package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationGroupPolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationGroupPolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationGroupPolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationGroupPolicyUpdate,
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
			"groups_claim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_children": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func getOpenidClientAuthorizationGroupPolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationGroupPolicy {
	var policies []string
	var resources []string
	var scopes []string
	var groups []keycloak.OpenidClientAuthorizationGroup
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
	if v, ok := data.Get("groups").([]interface{}); ok {
		for _, group := range v {
			groupMap := group.(map[string]interface{})
			tempGroup := keycloak.OpenidClientAuthorizationGroup{
				Id:             groupMap["id"].(string),
				Path:           groupMap["path"].(string),
				ExtendChildren: groupMap["extend_children"].(bool),
			}
			groups = append(groups, tempGroup)
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

	return resourceKeycloakOpenidClientAuthorizationGroupPolicyRead(data, meta)
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
