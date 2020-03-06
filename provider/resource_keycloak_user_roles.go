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
				Type:     schema.TypeList,
				Required: true,
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

	realm_id := data.Get("realm_id").(string)
	user_id := data.Get("user_id").(string)
	user, err := keycloakClient.GetUser(realm_id, user_id)

	roles, err := keycloakClient.GetRealmLevelRoleMappings(user)

	for _, realmRole := range data.Get("roles").([]interface{}) {
		role, err := keycloakClient.GetRoleByName(realm_id, "", realmRole.(string))
		if err != nil {
			return err
		}
		roles = append(roles, role)
	}

	keycloakClient.AddRealmLevelRoleMapping(user, roles)
	if err != nil {
		return err
	}

	data.Set("roles", roles)
	data.SetId(user_id)

	return resourceKeycloakUserRolesRead(data, meta)
}

func resourceKeycloakUserRolesRead(data *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceKeycloakUserRolesUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
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
	for _, realmRole := range data.Get("roles").([]interface{}) {
		role, err := keycloakClient.GetRoleByName(realm_id, "", realmRole.(string))
		if err != nil {
			return err
		}
		roles = append(roles, role)
	}
	return keycloakClient.RemoveRealmRolesFromUser(user, roles)
}
