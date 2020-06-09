package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakUserRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUserRolesCreate,
		Read:   resourceKeycloakUserRolesRead,
		Update: resourceKeycloakUserRolesUpdate,
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

func resourceKeycloakUserRolesCreate(data *schema.ResourceData, meta interface{}) error {
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

	// get the list of currently assigned roles. Due to default-realm- and client-roles
	// (e.g. roles of the account-client) this is probably not empty upon resource creation
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

func resourceKeycloakUserRolesUpdate(data *schema.ResourceData, meta interface{}) error {
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

	roleMappings, err := keycloakClient.GetUserRoleMappings(realmId, userId)
	if err != nil {
		return err
	}

	updates := calculateRoleMappingUpdates(tfRoles, intoRoleMapping(roleMappings))

	// `tfRoles` contains all roles that need to be added
	// `remoteRoles` contains all roles that need to be removed

	err = addRolesToUser(keycloakClient, updates.clientRolesToAdd, updates.realmRolesToAdd, user)
	if err != nil {
		return err
	}

	err = removeRolesFromUser(keycloakClient, updates.clientRolesToRemove, updates.realmRolesToRemove, user)
	if err != nil {
		return err
	}

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

// a struct that represents the "desired" state configured via terraform
// the key for 'clientRoles' is keycloak's client-id (the uuid, not to be confused with the OAuth Client Id)
type roleMapping struct {
	clientRoles map[string][]*keycloak.Role
	realmRoles  []*keycloak.Role
}

// transform the keycloak response UserRoleMapping into the internal roleMapping
// so we can compare the current state (keycloak) against the desired state (terraform)
func intoRoleMapping(userRoleMapping *keycloak.UserRoleMapping) *roleMapping {
	clientRoles := make(map[string][]*keycloak.Role)
	for _, clientRoleMapping := range userRoleMapping.ClientMappings {
		clientRoles[clientRoleMapping.Id] = clientRoleMapping.Mappings
	}

	mapping := roleMapping{
		clientRoles: clientRoles,
		realmRoles:  userRoleMapping.RealmMappings,
	}

	return &mapping
}

// given a list or roleIds, query keycloak for role-details to find out if a role is a client-role or a
// realm-rule (which is required to POST the role-assignment to the correct API-endpoint)
func getExtendedRoleMapping(keycloakClient *keycloak.KeycloakClient, realmId string, roleIds []string) (*roleMapping, error) {
	clientRoles := make(map[string][]*keycloak.Role)
	var realmRoles []*keycloak.Role

	for _, roleId := range roleIds {
		role, err := keycloakClient.GetRole(realmId, roleId)
		if err != nil {
			return nil, err
		}

		if role.ClientRole {
			clientRoles[role.ClientId] = append(clientRoles[role.ClientId], role)
		} else {
			realmRoles = append(realmRoles, role)
		}
	}

	mapping := roleMapping{
		clientRoles: clientRoles,
		realmRoles:  realmRoles,
	}

	return &mapping, nil
}

type terraformRoleMappingUpdates struct {
	realmRolesToRemove  []*keycloak.Role
	realmRolesToAdd     []*keycloak.Role
	clientRolesToRemove map[string][]*keycloak.Role
	clientRolesToAdd    map[string][]*keycloak.Role
}

// given the existing roles (queried from keycloak) and the requested roles (via tf)
// calculate the required updates, i.e. roles to remove and roles to add
func calculateRoleMappingUpdates(requestedRoles *roleMapping, existingRoles *roleMapping) *terraformRoleMappingUpdates {
	clientRolesToRemove := make(map[string][]*keycloak.Role)
	clientRolesToAdd := make(map[string][]*keycloak.Role)

	realmRolesToRemove := minusRoles(existingRoles.realmRoles, requestedRoles.realmRoles)
	realmRolesToAdd := minusRoles(requestedRoles.realmRoles, existingRoles.realmRoles)

	for clientId, requestedClientRoles := range requestedRoles.clientRoles {
		if existingClientRoles, ok := existingRoles.clientRoles[clientId]; ok {
			clientRolesToAdd[clientId] = minusRoles(requestedClientRoles, existingClientRoles)
			clientRolesToRemove[clientId] = minusRoles(existingClientRoles, requestedClientRoles)
		} else {
			// if no roles for this client exist yet, then of course, all requested roles need to be created
			clientRolesToAdd[clientId] = requestedClientRoles
		}
	}

	// now check all existing roles, if there are even any roles configured for each client
	for clientId, existingClientRoles := range existingRoles.clientRoles {
		if _, ok := requestedRoles.clientRoles[clientId]; !ok {
			// no role requested for this client? -> remove all existing client-roles
			clientRolesToRemove[clientId] = existingClientRoles
		}
	}

	updates := terraformRoleMappingUpdates{
		realmRolesToRemove:  realmRolesToRemove,
		realmRolesToAdd:     realmRolesToAdd,
		clientRolesToRemove: clientRolesToRemove,
		clientRolesToAdd:    clientRolesToAdd,
	}

	return &updates
}

// check if given role exists in a list of roles
func roleExists(roles []*keycloak.Role, role *keycloak.Role) bool {
	for _, r := range roles {
		if r.Id == role.Id {
			return true
		}
	}

	return false
}

// calculate the set difference: returns `a \ b`, i.e. every role that exist in a, but not in b
func minusRoles(a, b []*keycloak.Role) []*keycloak.Role {
	var aWithoutB []*keycloak.Role

	for _, role := range a {
		if !roleExists(b, role) {
			aWithoutB = append(aWithoutB, role)
		}
	}

	return aWithoutB
}
