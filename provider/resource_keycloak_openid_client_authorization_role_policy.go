package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationRolePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationRolePolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationRolePolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationRolePolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationRolePolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientAuthorizationRolePolicyImport,
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "POSITIVE",
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"required": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
				},
			},
		},
	}
}

func getOpenidClientAuthorizationRolePolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationRolePolicy {
	var scopes []string
	var rolesList []keycloak.OpenidClientAuthorizationRole
	if v, ok := data.GetOk("scopes"); ok {
		for _, scope := range v.(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}
	}
	if v, ok := data.Get("role").([]interface{}); ok {
		for _, role := range v {
			roleMap := role.(map[string]interface{})
			tempRole := keycloak.OpenidClientAuthorizationRole{
				Id:       roleMap["id"].(string),
				Required: roleMap["required"].(bool),
			}
			rolesList = append(rolesList, tempRole)
		}
	}

	resource := keycloak.OpenidClientAuthorizationRolePolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "role",
		Scopes:           scopes,
		Roles:            rolesList,
		Description:      data.Get("description").(string),
	}

	return &resource
}

func setOpenidClientAuthorizationRolePolicyResourceData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationRolePolicy) {
	data.SetId(policy.Id)

	data.Set("resource_server_id", policy.ResourceServerId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("logic", policy.Logic)
	data.Set("resources", policy.Resources)
	data.Set("scopes", policy.Scopes)
	data.Set("type", policy.Type)
	data.Set("description", policy.Description)
	data.Set("roles", policy.Roles)
}

func resourceKeycloakOpenidClientAuthorizationRolePolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationRolePolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationRolePolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationRolePolicyResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationRolePolicyRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationRolePolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationRolePolicy(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationRolePolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationRolePolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationRolePolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationRolePolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationRolePolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationRolePolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationRolePolicy(realmId, resourceServerId, id)
}

func resourceKeycloakOpenidClientAuthorizationRolePolicyImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
