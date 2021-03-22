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
	groupIds := data.Get("group_ids").(*schema.Set).List()

	for id := range groupIds {
		var user keycloak.User
		user.Id = userId
		user.RealmId = realmId
		err := keycloakClient.AddUserToGroup(&user, id)

		if err != nil {
			return err
		}
	}

	data.SetId(resource.UniqueId())

	return resourceKeycloakGroupMembershipsRead(data, meta)
}

func resourceKeycloakUserGroupsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)

	// nur relevante Gruppen
	userGroups, err := keycloakClient.GetUserGroups(realmId, userId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	var groups []string
	for _, group := range userGroups {
		groups = append(groups, group.Id)
	}

	data.Set("group_ids", groups)

	return nil
}

func resourceKeycloakUserGroupsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_ids").(string)
	tfMembers := data.Get("members").(*schema.Set).List()

	err := keycloakClient.ValidateGroupMembers(tfMembers.List())
	if err != nil {
		return err
	}

	keycloakMembers, err := keycloakClient.GetGroupMembers(realmId, groupId)
	if err != nil {
		return err
	}

	for _, keycloakMember := range keycloakMembers {
		if tfMembers.Contains(keycloakMember.Username) {
			// if the user exists in keycloak and tf state, no update is required for this member
			// remove them from the set so we can look at members that need to be added later
			tfMembers.Remove(keycloakMember.Username)
		} else {
			// if the user exists in keycloak and not in tf state, they need to be removed from the group
			err = keycloakClient.RemoveUserFromGroup(keycloakMember, groupId)
			if err != nil {
				return err
			}
		}
	}

	// at this point, `tfMembers` should only contain users that exist in tf state but not keycloak. these users need to be added
	err = keycloakClient.AddUsersToGroup(realmId, groupId, tfMembers.List())
	if err != nil {
		return err
	}

	data.SetId(groupMembershipsId(realmId, groupId))

	return resourceKeycloakGroupMembershipsRead(data, meta)
}

func resourceKeycloakUserGroupsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	return keycloakClient.RemoveUsersFromGroup(realmId, groupId, data.Get("members").(*schema.Set).List())
}

