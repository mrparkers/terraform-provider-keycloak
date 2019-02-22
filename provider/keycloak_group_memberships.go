package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakGroupMemberships() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakGroupMembershipsCreate,
		Read:   resourceKeycloakGroupMembershipsRead,
		Delete: resourceKeycloakGroupMembershipsDelete,
		Update: resourceKeycloakGroupMembershipsUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
		},
	}
}

func resourceKeycloakGroupMembershipsCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	groupId := data.Get("group_id").(string)
	members := data.Get("members").(*schema.Set).List()
	realmId := data.Get("realm_id").(string)

	err := keycloakClient.ValidateGroupMembers(members)
	if err != nil {
		return err
	}

	err = keycloakClient.AddUsersToGroup(realmId, groupId, members)
	if err != nil {
		return err
	}

	data.SetId(groupMembershipsId(realmId, groupId))

	return resourceKeycloakGroupMembershipsRead(data, meta)
}

func resourceKeycloakGroupMembershipsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	usersInGroup, err := keycloakClient.GetGroupMembers(realmId, groupId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	var members []string
	for _, userInGroup := range usersInGroup {
		members = append(members, userInGroup.Username)
	}

	data.Set("members", members)
	data.SetId(groupMembershipsId(realmId, groupId))

	return nil
}

func resourceKeycloakGroupMembershipsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)
	tfMembers := data.Get("members").(*schema.Set)

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
				return nil
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

func resourceKeycloakGroupMembershipsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	return keycloakClient.RemoveUsersFromGroup(realmId, groupId, data.Get("members").(*schema.Set).List())
}

func groupMembershipsId(realmId, groupId string) string {
	return fmt.Sprintf("%s/group-memberships/%s", realmId, groupId)
}
