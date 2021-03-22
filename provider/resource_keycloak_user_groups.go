package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakUserGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUserGroupsCreate,
		Read:   resourceKeycloakUserGroupsRead,
		Delete: resourceKeycloakUserGroupsDelete,
		Update: resourceKeycloakUserGroupsUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
		},
	}
}

func resourceKeycloakUserGroupsCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	groupIds := data.Get("group_ids").(*schema.Set)

	err := keycloakClient.AddUserToGroups(groupIds, userId, realmId)
	if err != nil {
		return err
	}

	data.SetId(resource.UniqueId())

	return resourceKeycloakUserGroupsRead(data, meta)
}

func resourceKeycloakUserGroupsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	groupIds := data.Get("group_ids").(*schema.Set)

	userGroups, err := keycloakClient.GetUserGroups(realmId, userId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	var groups []string
	for _, group := range userGroups {
		//only add groups that we care about
		if groupIds.Contains(group.Id) {
			groups = append(groups, group.Id)
		}
	}

	data.Set("group_ids", groups)

	return nil
}

func resourceKeycloakUserGroupsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	if data.HasChange("group_ids") {
		realmId := data.Get("realm_id").(string)
		userId := data.Get("user_id").(string)

		o, n := data.GetChange("group_ids")

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns)
		add := ns.Difference(os)

		if err := keycloakClient.RemoveUserFromGroups(remove, userId, realmId); err != nil {
			return err
		}

		if err := keycloakClient.AddUserToGroups(add, userId, realmId); err != nil {
			return err
		}
	}

	return resourceKeycloakUserGroupsRead(data, meta)
}

func resourceKeycloakUserGroupsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	groupIds := data.Get("group_ids").(*schema.Set)

	return keycloakClient.RemoveUserFromGroups(groupIds, userId, realmId)
}
