package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var (
	keycloakOpenidClientResourcePermissionDecisionStrategies = []string{"UNANIMOUS", "AFFIRMATIVE", "CONSENSUS"}
	keycloakOpenidClientPermissionTypes                      = []string{"resource", "scope"}
)

func resourceKeycloakOpenidClientAuthorizationPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationPermissionCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationPermissionRead,
		Delete: resourceKeycloakOpenidClientAuthorizationPermissionDelete,
		Update: resourceKeycloakOpenidClientAuthorizationPermissionUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientAuthorizationPermissionImport,
		},
		Schema: map[string]*schema.Schema{
			"resource_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"decision_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientResourcePermissionDecisionStrategies, false),
				Default:      "UNANIMOUS",
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "resource",
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientPermissionTypes, false),
			},
		},
	}
}

func getOpenidClientAuthorizationPermissionFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationPermission {
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

	permission := keycloak.OpenidClientAuthorizationPermission{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		Description:      data.Get("description").(string),
		Name:             data.Get("name").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Type:             data.Get("type").(string),
		Policies:         policies,
		Scopes:           scopes,
		Resources:        resources,
	}
	return &permission
}

func setOpenidClientAuthorizationPermissionData(data *schema.ResourceData, permission *keycloak.OpenidClientAuthorizationPermission) {
	data.SetId(permission.Id)
	data.Set("resource_server_id", permission.ResourceServerId)
	data.Set("realm_id", permission.RealmId)
	data.Set("description", permission.Description)
	data.Set("name", permission.Name)
	data.Set("decision_strategy", permission.DecisionStrategy)
	data.Set("type", permission.Type)
	data.Set("policies", permission.Policies)
	data.Set("scopes", permission.Scopes)
	data.Set("resources", permission.Resources)
}

func resourceKeycloakOpenidClientAuthorizationPermissionCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	permission := getOpenidClientAuthorizationPermissionFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationPermission(permission)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationPermissionData(data, permission)

	return resourceKeycloakOpenidClientAuthorizationPermissionRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationPermissionRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationPermissionData(data, permission)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationPermissionUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	permission := getOpenidClientAuthorizationPermissionFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationPermissionData(data, permission)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationPermissionDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationPermission(realmId, resourceServerId, id)
}

func resourceKeycloakOpenidClientAuthorizationPermissionImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{permissionId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
