package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakDefaultGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakDefaultGroupsCreate,
		Read:   resourceKeycloakDefaultGroupsRead,
		Update: resourceKeycloakDefaultGroupsUpdate,
		Delete: resourceKeycloakDefaultGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakDefaultGroupsImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
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

func defaultGroupId(realmId string) string {
	return realmId + "/default-groups"
}

func resourceKeycloakDefaultGroupsCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupIds := interfaceSliceToStringSlice(data.Get("group_ids").(*schema.Set).List())

	for _, groupId := range groupIds {
		err := keycloakClient.PutDefaultGroup(realmId, groupId)
		if err != nil {
			return err
		}
	}

	data.SetId(defaultGroupId(realmId))

	return nil
}

func resourceKeycloakDefaultGroupsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)

	groups, err := keycloakClient.GetDefaultGroups(realmId)
	if err != nil {
		return err
	}

	var groupIds []string
	for _, group := range groups {
		groupIds = append(groupIds, group.Id)
	}

	data.SetId(defaultGroupId(realmId))
	data.Set("group_ids", groupIds)

	return nil
}

func resourceKeycloakDefaultGroupsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	newGroupIds := data.Get("group_ids").(*schema.Set)

	originalGroups, err := keycloakClient.GetDefaultGroups(realmId)
	if err != nil {
		return err
	}

	for _, originalGroup := range originalGroups {
		if newGroupIds.Contains(originalGroup.Id) {
			newGroupIds.Remove(originalGroup.Id)
		} else {
			err := keycloakClient.DeleteDefaultGroup(realmId, originalGroup.Id)
			if err != nil {
				return err
			}
		}
	}

	// at this point newGroupIds should contain only users that need to be created
	for _, group := range interfaceSliceToStringSlice(newGroupIds.List()) {
		err := keycloakClient.PutDefaultGroup(realmId, group)
		if err != nil {
			return err
		}
	}

	data.SetId(defaultGroupId(realmId))

	return resourceKeycloakDefaultGroupsRead(data, meta)
}

func resourceKeycloakDefaultGroupsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupIds := interfaceSliceToStringSlice(data.Get("group_ids").(*schema.Set).List())

	for _, groupId := range groupIds {
		err := keycloakClient.DeleteDefaultGroup(realmId, groupId)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceKeycloakDefaultGroupsImport(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	data.Set("realm_id", data.Id())
	data.SetId(data.Id())
	return []*schema.ResourceData{data}, nil
}
