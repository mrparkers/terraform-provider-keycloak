package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakGroupRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakGroupRolesCreate,
		Read:   resourceKeycloakGroupRolesRead,
		Update: resourceKeycloakGroupRolesUpdate,
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
		},
	}
}

func groupRolesId(realmId, groupId string) string {
	return fmt.Sprintf("%s/%s", realmId, groupId)
}

func getMapOfRealmAndClientRoles(keycloakClient *keycloak.KeycloakClient, realmId string, roleIds []string) (map[string][]*keycloak.Role, error) {
	roles := make(map[string][]*keycloak.Role)

	for _, roleId := range roleIds {
		role, err := keycloakClient.GetRole(realmId, roleId)
		if err != nil {
			return nil, err
		}

		if role.ClientRole {
			roles[role.ClientId] = append(roles[role.ClientId], role)
		} else {
			roles["realm"] = append(roles["realm"], role)
		}
	}

	return roles, nil
}

// given a group and a map of roles we already know about, fetch the roles we don't know about
// `localRoles` is used as a cache to avoid unnecessary http requests
func getMapOfRealmAndClientRolesFromGroup(keycloakClient *keycloak.KeycloakClient, group *keycloak.Group, localRoles map[string][]*keycloak.Role) (map[string][]*keycloak.Role, error) {
	roles := make(map[string][]*keycloak.Role)

	// realm roles
	if len(group.RealmRoles) != 0 {
		var realmRoles []*keycloak.Role

		for _, realmRoleName := range group.RealmRoles {
			found := false

			for _, localRealmRole := range localRoles["realm"] {
				if localRealmRole.Name == realmRoleName {
					found = true
					realmRoles = append(realmRoles, localRealmRole)

					break
				}
			}

			if !found {
				realmRole, err := keycloakClient.GetRoleByName(group.RealmId, "", realmRoleName)
				if err != nil {
					return nil, err
				}

				realmRoles = append(realmRoles, realmRole)
			}
		}

		roles["realm"] = realmRoles
	}

	// client roles
	if len(group.ClientRoles) != 0 {
		for clientName, clientRoleNames := range group.ClientRoles {
			client, err := keycloakClient.GetGenericClientByClientId(group.RealmId, clientName)
			if err != nil {
				return nil, err
			}

			var clientRoles []*keycloak.Role
			for _, clientRoleName := range clientRoleNames {
				found := false

				for _, localClientRole := range localRoles[client.Id] {
					if localClientRole.Name == clientRoleName {
						found = true
						clientRoles = append(clientRoles, localClientRole)

						break
					}
				}

				if !found {
					clientRole, err := keycloakClient.GetRoleByName(group.RealmId, client.Id, clientRoleName)
					if err != nil {
						return nil, err
					}

					clientRoles = append(clientRoles, clientRole)
				}
			}

			roles[client.Id] = clientRoles
		}
	}

	return roles, nil
}

func addRolesToGroup(keycloakClient *keycloak.KeycloakClient, rolesToAdd map[string][]*keycloak.Role, group *keycloak.Group) error {
	if realmRoles, ok := rolesToAdd["realm"]; ok && len(realmRoles) != 0 {
		err := keycloakClient.AddRealmRolesToGroup(group.RealmId, group.Id, realmRoles)
		if err != nil {
			return err
		}
	}

	for k, roles := range rolesToAdd {
		if k == "realm" {
			continue
		}

		err := keycloakClient.AddClientRolesToGroup(group.RealmId, group.Id, k, roles)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeRolesFromGroup(keycloakClient *keycloak.KeycloakClient, rolesToRemove map[string][]*keycloak.Role, group *keycloak.Group) error {
	if realmRoles, ok := rolesToRemove["realm"]; ok && len(realmRoles) != 0 {
		err := keycloakClient.RemoveRealmRolesFromGroup(group.RealmId, group.Id, realmRoles)
		if err != nil {
			return err
		}
	}

	for k, roles := range rolesToRemove {
		if k == "realm" {
			continue
		}

		err := keycloakClient.RemoveClientRolesFromGroup(group.RealmId, group.Id, k, roles)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceKeycloakGroupRolesCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	group, err := keycloakClient.GetGroup(realmId, groupId)
	if err != nil {
		return err
	}

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	rolesToAdd, err := getMapOfRealmAndClientRoles(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	err = addRolesToGroup(keycloakClient, rolesToAdd, group)
	if err != nil {
		return err
	}

	data.SetId(groupRolesId(realmId, groupId))

	return resourceKeycloakGroupRolesRead(data, meta)
}

func resourceKeycloakGroupRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	group, err := keycloakClient.GetGroup(realmId, groupId)
	if err != nil {
		return err
	}

	var roleIds []string

	if len(group.RealmRoles) != 0 {
		for _, realmRole := range group.RealmRoles {
			role, err := keycloakClient.GetRoleByName(realmId, "", realmRole)
			if err != nil {
				return err
			}

			roleIds = append(roleIds, role.Id)
		}
	}

	if len(group.ClientRoles) != 0 {
		for clientName, clientRoles := range group.ClientRoles {
			client, err := keycloakClient.GetGenericClientByClientId(realmId, clientName)
			if err != nil {
				return err
			}

			for _, clientRole := range clientRoles {
				role, err := keycloakClient.GetRoleByName(realmId, client.Id, clientRole)
				if err != nil {
					return err
				}

				roleIds = append(roleIds, role.Id)
			}
		}
	}

	data.Set("role_ids", roleIds)
	data.SetId(groupRolesId(realmId, groupId))

	return nil
}

func resourceKeycloakGroupRolesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	group, err := keycloakClient.GetGroup(realmId, groupId)
	if err != nil {
		return err
	}

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())

	tfRoles, err := getMapOfRealmAndClientRoles(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	remoteRoles, err := getMapOfRealmAndClientRolesFromGroup(keycloakClient, group, tfRoles)
	if err != nil {
		return err
	}

	removeDuplicateRoles(&tfRoles, &remoteRoles)

	// `tfRoles` contains all roles that need to be added
	// `remoteRoles` contains all roles that need to be removed

	err = addRolesToGroup(keycloakClient, tfRoles, group)
	if err != nil {
		return err
	}

	err = removeRolesFromGroup(keycloakClient, remoteRoles, group)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakGroupRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupId := data.Get("group_id").(string)

	group, err := keycloakClient.GetGroup(realmId, groupId)

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	rolesToRemove, err := getMapOfRealmAndClientRoles(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	err = removeRolesFromGroup(keycloakClient, rolesToRemove, group)
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

func removeRoleFromSlice(slice []*keycloak.Role, index int) []*keycloak.Role {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func removeDuplicateRoles(one, two *map[string][]*keycloak.Role) {
	for k := range *one {
		for i1 := 0; i1 < len((*one)[k]); i1++ {
			s1 := (*one)[k][i1]

			for i2 := 0; i2 < len((*two)[k]); i2++ {
				s2 := (*two)[k][i2]

				if s1.Id == s2.Id {
					(*one)[k] = removeRoleFromSlice((*one)[k], i1)
					(*two)[k] = removeRoleFromSlice((*two)[k], i2)

					i1--
					break
				}
			}
		}
	}
}
