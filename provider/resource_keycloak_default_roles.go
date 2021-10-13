package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakDefaultRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakDefaultRolesReconcile,
		Read:   resourceKeycloakDefaultRolesRead,
		Delete: resourceKeycloakDefaultRolesDelete,
		Update: resourceKeycloakDefaultRolesReconcile,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakDefaultRolesImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_roles": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Realm level roles assigned to new users.",
				Required:    true,
			},
		},
	}
}

func mapFromDataToDefaultRoles(data *schema.ResourceData) *keycloak.DefaultRoles {
	defaultRolesList := make([]string, 0)
	if v, ok := data.GetOk("default_roles"); ok {
		for _, defaultRole := range v.(*schema.Set).List() {
			defaultRolesList = append(defaultRolesList, defaultRole.(string))
		}
	}

	defaultRoles := &keycloak.DefaultRoles{
		Id:           data.Id(),
		RealmId:      data.Get("realm_id").(string),
		DefaultRoles: defaultRolesList,
	}

	return defaultRoles
}

func mapFromDefaultRolesToData(data *schema.ResourceData, defaultRoles *keycloak.DefaultRoles) {
	data.SetId(defaultRoles.Id)

	data.Set("realm_id", defaultRoles.RealmId)
	data.Set("default_roles", defaultRoles.DefaultRoles)
}

func resourceKeycloakDefaultRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	composites, err := keycloakClient.GetDefaultRoles(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	defaultRoleNamesList := getDefaultRoleNames(composites)

	defaultRoles := &keycloak.DefaultRoles{
		Id:           id,
		RealmId:      realmId,
		DefaultRoles: defaultRoleNamesList,
	}

	mapFromDefaultRolesToData(data, defaultRoles)

	return nil
}

func resourceKeycloakDefaultRolesReconcile(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	defaultRoles := mapFromDataToDefaultRoles(data)

	realm, err := keycloakClient.GetRealm(defaultRoles.RealmId)
	if err != nil {
		return err
	}

	data.SetId(realm.DefaultRole.Id)

	composites, err := keycloakClient.GetDefaultRoles(defaultRoles.RealmId, realm.DefaultRole.Id)
	if err != nil {
		return err
	}

	defaultRoleNamesList := getDefaultRoleNames(composites)
	rolesList, err := keycloakClient.GetRealmRoles(defaultRoles.RealmId)
	if err != nil {
		return err
	}

	// skip if actual default roles in keycloak same as we want
	if roleListsEqual(defaultRoleNamesList, defaultRoles.DefaultRoles) {
		return nil
	}

	var putList, deleteList []*keycloak.Role
	for _, roleName := range defaultRoles.DefaultRoles {
		if !roleListContains(defaultRoleNamesList, roleName) {
			defaultRoles, err := getRoleByNameFromList(rolesList, roleName)
			if err != nil {
				return err
			}
			putList = append(putList, defaultRoles)
		}
	}
	for _, roleName := range defaultRoleNamesList {
		if !roleListContains(defaultRoles.DefaultRoles, roleName) {
			defaultRoles, err := getRoleByNameFromList(rolesList, roleName)
			if err != nil {
				return err
			}
			deleteList = append(deleteList, defaultRoles)
		}
	}

	// apply if not empty
	if len(putList) > 0 {
		role := &keycloak.Role{
			RealmId: defaultRoles.RealmId,
			Id:      realm.DefaultRole.Id,
		}
		err := keycloakClient.AddCompositesToRole(role, putList)
		if err != nil {
			return err
		}
	}
	if len(deleteList) > 0 {
		role := &keycloak.Role{
			RealmId: defaultRoles.RealmId,
			Id:      realm.DefaultRole.Id,
		}
		err := keycloakClient.RemoveCompositesFromRole(role, deleteList)
		if err != nil {
			return err
		}
	}

	return resourceKeycloakDefaultRolesRead(data, meta)
}

// remove all roles from default
func resourceKeycloakDefaultRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	realm, err := keycloakClient.GetRealm(realmId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	defaultRoles, err := keycloakClient.GetDefaultRoles(realmId, realm.DefaultRole.Id)
	if err != nil {
		return err
	}

	if len(defaultRoles) > 0 {
		role := &keycloak.Role{
			RealmId: realmId,
			Id:      realm.DefaultRole.Id,
		}
		err := keycloakClient.RemoveCompositesFromRole(role, defaultRoles)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceKeycloakDefaultRolesImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/{{defaultRoleId}}.")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func getDefaultRoleNames(roles []*keycloak.Role) []string {
	var defaultRolesNames []string
	for _, defaultRoles := range roles {
		defaultRolesNames = append(defaultRolesNames, defaultRoles.Name)
	}
	return defaultRolesNames
}

func getRoleByNameFromList(defaultRoles []*keycloak.Role, name string) (*keycloak.Role, error) {
	for _, element := range defaultRoles {
		if element.Name == name {
			return element, nil
		}
	}
	return nil, fmt.Errorf("defaultRoles not found by name")
}

func roleListContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func roleListsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
