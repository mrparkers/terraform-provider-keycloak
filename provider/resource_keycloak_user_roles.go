package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakUserRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUserRolesReconcile,
		Read:   resourceKeycloakUserRolesRead,
		Update: resourceKeycloakUserRolesReconcile,
		Delete: resourceKeycloakUserRolesDelete,
		// This resource can be imported using {{realm}}/{{userId}}.
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakUserRolesImport,
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
			"role_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
		},
	}
}

func userRolesId(realmId, userId string) string {
	return fmt.Sprintf("%s/%s", realmId, userId)
}

func addRolesToUser(keycloakClient *keycloak.KeycloakClient, clientRolesToAdd map[string][]*keycloak.Role, realmRolesToAdd []*keycloak.Role, user *keycloak.User) error {
	if len(realmRolesToAdd) != 0 {
		err := keycloakClient.AddRealmRolesToUser(user.RealmId, user.Id, realmRolesToAdd)
		if err != nil {
			return err
		}
	}

	for k, roles := range clientRolesToAdd {
		if len(roles) != 0 {
			err := keycloakClient.AddClientRolesToUser(user.RealmId, user.Id, k, roles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeRolesFromUser(keycloakClient *keycloak.KeycloakClient, clientRolesToRemove map[string][]*keycloak.Role, realmRolesToRemove []*keycloak.Role, user *keycloak.User) error {
	if len(realmRolesToRemove) != 0 {
		err := keycloakClient.RemoveRealmRolesFromUser(user.RealmId, user.Id, realmRolesToRemove)
		if err != nil {
			return err
		}
	}

	for k, roles := range clientRolesToRemove {
		if len(roles) != 0 {
			err := keycloakClient.RemoveClientRolesFromUser(user.RealmId, user.Id, k, roles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceKeycloakUserRolesReconcile(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)

	user, err := keycloakClient.GetUser(realmId, userId)
	if err != nil {
		return err
	}

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	tfRoles, err := getExtendedRoleMapping(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	// get the list of currently assigned roles. Due to default realm and client roles
	// (e.g. roles of the account client) this is probably not empty upon resource creation
	roleMappings, err := keycloakClient.GetUserRoleMappings(realmId, userId)

	// sort into roles we need to add and roles we need to remove
	updates := calculateRoleMappingUpdates(tfRoles, intoRoleMapping(roleMappings))

	// add roles
	err = addRolesToUser(keycloakClient, updates.clientRolesToAdd, updates.realmRolesToAdd, user)
	if err != nil {
		return err
	}

	// remove roles
	err = removeRolesFromUser(keycloakClient, updates.clientRolesToRemove, updates.realmRolesToRemove, user)
	if err != nil {
		return err
	}

	data.SetId(userRolesId(realmId, userId))
	return resourceKeycloakUserRolesRead(data, meta)
}

func resourceKeycloakUserRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)

	roles, err := keycloakClient.GetUserRoleMappings(realmId, userId)
	if err != nil {
		return err
	}

	var roleIds []string

	for _, realmRole := range roles.RealmMappings {
		roleIds = append(roleIds, realmRole.Id)
	}

	for _, clientRoleMapping := range roles.ClientMappings {
		for _, clientRole := range clientRoleMapping.Mappings {
			roleIds = append(roleIds, clientRole.Id)
		}
	}

	data.Set("role_ids", roleIds)
	data.SetId(userRolesId(realmId, userId))

	return nil
}

func resourceKeycloakUserRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)

	user, err := keycloakClient.GetUser(realmId, userId)

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	rolesToRemove, err := getExtendedRoleMapping(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	err = removeRolesFromUser(keycloakClient, rolesToRemove.clientRoles, rolesToRemove.realmRoles, user)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakUserRolesImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/{{userId}}.")
	}

	d.Set("realm_id", parts[0])
	d.Set("user_id", parts[1])

	d.SetId(userRolesId(parts[0], parts[1]))

	return []*schema.ResourceData{d}, nil
}
