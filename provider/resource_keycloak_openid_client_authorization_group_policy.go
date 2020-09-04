package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"logic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakPolicyLogicTypes, false),
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
	var groups []keycloak.OpenidClientAuthorizationGroup
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
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "group",
		GroupsClaim:      data.Get("groups_claim").(string),
		Groups:           groups,
		Description:      data.Get("description").(string),
	}

	return &resource
}

func setOpenidClientAuthorizationGroupPolicyResourceData(keycloakClient *keycloak.KeycloakClient, policy *keycloak.OpenidClientAuthorizationGroupPolicy, data *schema.ResourceData) error {
	data.SetId(policy.Id)

	data.Set("resource_server_id", policy.ResourceServerId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("logic", policy.Logic)
	data.Set("description", policy.Description)
	data.Set("groups_claim", policy.GroupsClaim)

	var groups []interface{}
	for _, g := range policy.Groups {
		// the "path" attribute is omitted by Keycloak, so we have to look this group up ourselves to get the path
		group, err := keycloakClient.GetGroup(policy.RealmId, g.Id)
		if err != nil {
			return err
		}

		groups = append(groups, map[string]interface{}{
			"id":              g.Id,
			"path":            group.Path,
			"extend_children": g.ExtendChildren,
		})
	}

	data.Set("groups", groups)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationGroupPolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationGroupPolicy(resource)
	if err != nil {
		return err
	}

	err = setOpenidClientAuthorizationGroupPolicyResourceData(keycloakClient, resource, data)
	if err != nil {
		return err
	}

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

	err = setOpenidClientAuthorizationGroupPolicyResourceData(keycloakClient, resource, data)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationGroupPolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationGroupPolicy(resource)
	if err != nil {
		return err
	}

	err = setOpenidClientAuthorizationGroupPolicyResourceData(keycloakClient, resource, data)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientAuthorizationGroupPolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationGroupPolicy(realmId, resourceServerId, id)
}
