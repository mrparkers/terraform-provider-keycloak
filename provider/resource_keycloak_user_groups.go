package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakUserGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUserGroupsReconcile,
		Read:   resourceKeycloakUserGroupsRead,
		Delete: resourceKeycloakUserGroupsDelete,
		Update: resourceKeycloakUserGroupsReconcile,
		// This resource can be imported using {{realm}}/{{userId}}.
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakUserGroupsImport,
		},
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
			"exhaustive": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceKeycloakUserGroupsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	groupIds := data.Get("group_ids").(*schema.Set)
	exhaustive := data.Get("exhaustive").(bool)

	userGroups, err := keycloakClient.GetUserGroups(realmId, userId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	var groups []string
	for _, group := range userGroups {
		//only add groups that we care about
		if exhaustive || groupIds.Contains(group.Id) {
			groups = append(groups, group.Id)
		}
	}

	data.Set("group_ids", groups)
	data.SetId(userGroupsId(realmId, userId))

	return nil
}

func resourceKeycloakUserGroupsReconcile(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	groupIds := interfaceSliceToStringSlice(data.Get("group_ids").(*schema.Set).List())
	exhaustive := data.Get("exhaustive").(bool)

	if data.HasChange("group_ids") {
		o, n := data.GetChange("group_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := interfaceSliceToStringSlice(os.Difference(ns).List())

		if err := keycloakClient.RemoveUserFromGroups(remove, userId, realmId); err != nil {
			return err
		}
	}

	userGroups, err := keycloakClient.GetUserGroups(realmId, userId)
	if err != nil {
		return err
	}

	var userGroupsIds []string
	for _, group := range userGroups {
		userGroupsIds = append(userGroupsIds, group.Id)
	}

	remove := stringArrayDifference(userGroupsIds, groupIds)
	add := stringArrayDifference(groupIds, userGroupsIds)

	if err := keycloakClient.AddUserToGroups(add, userId, realmId); err != nil {
		return err
	}

	if exhaustive {
		if err := keycloakClient.RemoveUserFromGroups(remove, userId, realmId); err != nil {
			return err
		}
	}

	data.SetId(userGroupsId(realmId, userId))
	return resourceKeycloakUserGroupsRead(data, meta)
}

func resourceKeycloakUserGroupsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	groupIds := interfaceSliceToStringSlice(data.Get("group_ids").(*schema.Set).List())

	return keycloakClient.RemoveUserFromGroups(groupIds, userId, realmId)
}

func resourceKeycloakUserGroupsImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/{{userId}}.")
	}

	d.Set("realm_id", parts[0])
	d.Set("user_id", parts[1])

	d.SetId(userGroupsId(parts[0], parts[1]))

	return []*schema.ResourceData{d}, nil
}

func userGroupsId(realmId, userId string) string {
	return fmt.Sprintf("%s/%s", realmId, userId)
}
