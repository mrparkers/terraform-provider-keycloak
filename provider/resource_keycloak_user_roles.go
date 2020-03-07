package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeyCloakUserRoles() *schema.Resource {
	return &schema.Resource{
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
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceKeycloakUserRolesCreate,
		Read:   resourceKeycloakUserRolesRead,
		Update: resourceKeycloakUserRolesUpdate,
		Delete: resourceKeycloakUserRolesDelete,
	}
}

func resourceKeycloakUserRolesCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	roles := data.Get("roles").(*schema.Set)

	//for _, realmRole := range data.Get("roles").([]interface{}) {
	//	role, err := keycloakClient.GetRoleByName(realmId, "", realmRole.(string))
	//	if err != nil {
	//		return err
	//	}
	//	roles = append(roles, role)
	//}

	err := keycloakClient.AddRealmRolesToUser(realmId, userId, roles.List())
	if err != nil {
		return err
	}

	data.Set("roles", roles)
	data.SetId(userId)

	return resourceKeycloakUserRolesRead(data, meta)
}

func resourceKeycloakUserRolesRead(data *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceKeycloakUserRolesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	tfRoles := data.Get("roles").(*schema.Set)

	user, err := keycloakClient.GetUser(realmId, userId)

	if err != nil {
		return err
	}

	roles, err := keycloakClient.GetRealmRoleMappings(user)

	if err != nil {
		return err
	}

	for _, role := range roles {
		if tfRoles.Contains(role.Name) {
			tfRoles.Remove(role.Name)
		} else {
			err = keycloakClient.RemoveRealmRolesFromUser(user, []*keycloak.Role{role})
			if err != nil {
				return err
			}
		}
	}

	err = keycloakClient.AddRealmRolesToUser(realmId, userId, tfRoles.List())
	if err != nil {
		return err
	}

	data.SetId(userId)

	return resourceKeycloakUserRolesRead(data, meta)
}

func resourceKeycloakUserRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm_id := data.Get("realm_id").(string)
	user_id := data.Get("user_id").(string)

	user, err := keycloakClient.GetUser(realm_id, user_id)
	if err != nil {
		return err
	}

	var roles []*keycloak.Role
	for _, realmRole := range data.Get("roles").(*schema.Set).List() {
		role, err := keycloakClient.GetRoleByName(realm_id, "", realmRole.(string))
		if err != nil {
			return err
		}
		roles = append(roles, role)
	}
	return keycloakClient.RemoveRealmRolesFromUser(user, roles)
}
