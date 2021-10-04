package provider

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakDefaultRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakDefaultRolesCreate,
		Read:   resourceKeycloakDefaultRolesRead,
		Delete: resourceKeycloakDefaultRolesDelete,
		Update: resourceKeycloakDefaultRolesUpdate,
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

func resourceKeycloakDefaultRolesCreate(data *schema.ResourceData, meta interface{}) error {
	err := resourceKeycloakDefaultRolesUpdate(data, meta)
	if err != nil {
		return err
	}
	return resourceKeycloakDefaultRolesRead(data, meta)
}

func resourceKeycloakDefaultRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	composites, _ := keycloakClient.GetDefaultRoles(realmId, id)
	defaultRoleNamesList, _ := getDefaultRoleNames(composites)

	defaultRoles := &keycloak.DefaultRoles{
		Id:           id,
		RealmId:      realmId,
		DefaultRoles: defaultRoleNamesList,
	}

	mapFromDefaultRolesToData(data, defaultRoles)

	return nil
}

func resourceKeycloakDefaultRolesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	defaultRoles := mapFromDataToDefaultRoles(data)

	realm, _ := keycloakClient.GetRealm(defaultRoles.RealmId)
	data.SetId(realm.DefaultRole.Id)

	composites, _ := keycloakClient.GetDefaultRoles(defaultRoles.RealmId, realm.DefaultRole.Id)
	defaultRoleNamesList, _ := getDefaultRoleNames(composites)
	rolesList, _ := keycloakClient.GetRealmRoles(defaultRoles.RealmId)

	// skip if actual default roles in keycloak same as we want
	if reflect.DeepEqual(defaultRoleNamesList, defaultRoles.DefaultRoles) {
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
		err := keycloakClient.AddDefaultRoles(defaultRoles.RealmId, realm.DefaultRole.Id, putList)
		if err != nil {
			return err
		}
	}
	if len(deleteList) > 0 {
		err := keycloakClient.RemoveDefaultRoles(defaultRoles.RealmId, realm.DefaultRole.Id, deleteList)
		if err != nil {
			return err
		}
	}

	return resourceKeycloakDefaultRolesRead(data, meta)
}

// remove all roles from default
func resourceKeycloakDefaultRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	defaultRoles := mapFromDataToDefaultRoles(data)

	realm, _ := keycloakClient.GetRealm(defaultRoles.RealmId)
	defaultRoles.Id = realm.DefaultRole.Id
	data.SetId(defaultRoles.Id)

	composites, _ := keycloakClient.GetDefaultRoles(defaultRoles.RealmId, defaultRoles.Id)
	defaultRoleNamesList, _ := getDefaultRoleNames(composites)
	rolesList, _ := keycloakClient.GetRealmRoles(defaultRoles.RealmId)

	if reflect.DeepEqual(defaultRoleNamesList, defaultRoles.DefaultRoles) {
		return nil
	}

	var deleteList []*keycloak.Role
	for _, roleName := range defaultRoleNamesList {
		defaultRoles, err := getRoleByNameFromList(rolesList, roleName)
		if err != nil {
			return err
		}
		deleteList = append(deleteList, defaultRoles)
	}
	if len(deleteList) > 0 {
		err := keycloakClient.RemoveDefaultRoles(defaultRoles.RealmId, defaultRoles.Id, deleteList)
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

func getDefaultRoleNames(roles []*keycloak.Role) ([]string, error) {
	var defaultRolesNames []string
	for _, defaultRoles := range roles {
		defaultRolesNames = append(defaultRolesNames, defaultRoles.Name)
	}
	return defaultRolesNames, nil
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
