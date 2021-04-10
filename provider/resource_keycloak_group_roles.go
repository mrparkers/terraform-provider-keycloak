package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakGroupRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakGroupRolesReconcile,
		Read:   resourceKeycloakGroupRolesRead,
		Update: resourceKeycloakGroupRolesReconcile,
		Delete: resourceKeycloakGroupRolesDelete,
		// This resource can be imported using {{realm}}/{{groupId}}.
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakGroupRolesImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
			"exhaustive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func groupRolesId(realmId, groupId string) string {
	return fmt.Sprintf("%s/%s", realmId, groupId)
}

func addRolesToGroup(keycloakClient *keycloak.KeycloakClient, clientRolesToAdd map[string][]*keycloak.Role, realmRolesToAdd []*keycloak.Role, group *keycloak.Group) error {
	if len(realmRolesToAdd) != 0 {
		err := keycloakClient.AddRealmRolesToGroup(group.RealmId, group.Id, realmRolesToAdd)
		if err != nil {
			return err
		}
	}

	for k, roles := range clientRolesToAdd {
		if len(roles) != 0 {
			err := keycloakClient.AddClientRolesToGroup(group.RealmId, group.Id, k, roles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeRolesFromGroup(keycloakClient *keycloak.KeycloakClient, clientRolesToRemove map[string][]*keycloak.Role, realmRolesToRemove []*keycloak.Role, group *keycloak.Group) error {
	if len(realmRolesToRemove) != 0 {
		err := keycloakClient.RemoveRealmRolesFromGroup(group.RealmId, group.Id, realmRolesToRemove)
		if err != nil {
			return err
		}
	}

	for k, roles := range clientRolesToRemove {
		if len(roles) != 0 {
			err := keycloakClient.RemoveClientRolesFromGroup(group.RealmId, group.Id, k, roles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceKeycloakGroupRolesReconcile(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)
	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	exhaustive := data.Get("exhaustive").(bool)

	group, err := keycloakClient.GetGroup(realmId, groupId)
	if err != nil {
		return err
	}

	if data.HasChange("role_ids") {
		o, n := data.GetChange("role_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := interfaceSliceToStringSlice(os.Difference(ns).List())

		tfRolesToRemove, err := getExtendedRoleMapping(keycloakClient, realmId, remove)
		if err != nil {
			return err
		}

		if err = removeRolesFromGroup(keycloakClient, tfRolesToRemove.clientRoles, tfRolesToRemove.realmRoles, group); err != nil {
			return err
		}
	}

	tfRoles, err := getExtendedRoleMapping(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	// get the list of currently assigned roles. Due to default realm and client roles
	roleMappings, err := keycloakClient.GetGroupRoleMappings(realmId, groupId)

	// sort into roles we need to add and roles we need to remove
	updates := calculateRoleMappingUpdates(tfRoles, intoRoleMapping(roleMappings))

	// add roles
	err = addRolesToGroup(keycloakClient, updates.clientRolesToAdd, updates.realmRolesToAdd, group)
	if err != nil {
		return err
	}

	// remove roles if exhaustive (authoritative)
	if exhaustive {
		err = removeRolesFromGroup(keycloakClient, updates.clientRolesToRemove, updates.realmRolesToRemove, group)
		if err != nil {
			return err
		}
	}

	data.SetId(groupRolesId(realmId, groupId))
	return resourceKeycloakGroupRolesRead(data, meta)
}

func resourceKeycloakGroupRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)
	sortedRoleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	exhaustive := data.Get("exhaustive").(bool)

	// check if group exists, remove from state if not found
	if _, err := keycloakClient.GetGroup(realmId, groupId); err != nil {
		return handleNotFoundError(err, data)
	}

	roles, err := keycloakClient.GetGroupRoleMappings(realmId, groupId)
	if err != nil {
		return err
	}

	var roleIds []string

	for _, realmRole := range roles.RealmMappings {
		if exhaustive || stringSliceContains(sortedRoleIds, realmRole.Id) {
			roleIds = append(roleIds, realmRole.Id)
		}
	}

	for _, clientRoleMapping := range roles.ClientMappings {
		for _, clientRole := range clientRoleMapping.Mappings {
			if exhaustive || stringSliceContains(sortedRoleIds, clientRole.Id) {
				roleIds = append(roleIds, clientRole.Id)
			}
		}
	}

	data.Set("role_ids", roleIds)
	data.SetId(groupRolesId(realmId, groupId))

	return nil
}

func resourceKeycloakGroupRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	group, err := keycloakClient.GetGroup(realmId, groupId)

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	rolesToRemove, err := getExtendedRoleMapping(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	err = removeRolesFromGroup(keycloakClient, rolesToRemove.clientRoles, rolesToRemove.realmRoles, group)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakGroupRolesImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/{{groupId}}.")
	}

	d.Set("realm_id", parts[0])
	d.Set("group_id", parts[1])

	d.SetId(groupRolesId(parts[0], parts[1]))

	return []*schema.ResourceData{d}, nil
}
