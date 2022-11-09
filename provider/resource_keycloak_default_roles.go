package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakDefaultRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakDefaultRolesReconcile,
		ReadContext:   resourceKeycloakDefaultRolesRead,
		DeleteContext: resourceKeycloakDefaultRolesDelete,
		UpdateContext: resourceKeycloakDefaultRolesReconcile,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakDefaultRolesImport,
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
				Description: "Realm level roles (name) assigned to new users.",
				Required:    true,
			},
			"default_role_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Realm level roles (ID) assigned to new users.",
				Optional:    true,
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

	var defaultRoleIds []string
	if v, ok := data.GetOk("default_role_ids"); ok {
		for _, defaultRole := range v.(*schema.Set).List() {
			defaultRoleIds = append(defaultRoleIds, defaultRole.(string))
		}
	}

	defaultRoles := &keycloak.DefaultRoles{
		Id:             data.Id(),
		RealmId:        data.Get("realm_id").(string),
		DefaultRoles:   defaultRolesList,
		DefaultRoleIds: defaultRoleIds,
	}

	return defaultRoles
}

func mapFromDefaultRolesToData(data *schema.ResourceData, defaultRoles *keycloak.DefaultRoles) {
	data.SetId(defaultRoles.Id)

	data.Set("realm_id", defaultRoles.RealmId)
	data.Set("default_roles", defaultRoles.DefaultRoles)
	data.Set("default_role_ids", defaultRoles.DefaultRoleIds)
}

func resourceKeycloakDefaultRolesRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	composites, err := keycloakClient.GetDefaultRoles(ctx, realmId, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	defaultRoleNamesList, defaultRoleIds := getDefaultRoles(composites)

	defaultRoles := &keycloak.DefaultRoles{
		Id:             id,
		RealmId:        realmId,
		DefaultRoles:   defaultRoleNamesList,
		DefaultRoleIds: defaultRoleIds,
	}

	mapFromDefaultRolesToData(data, defaultRoles)

	return nil
}

func resourceKeycloakDefaultRolesReconcile(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	if ok, err := keycloakClient.VersionIsGreaterThanOrEqualTo(ctx, keycloak.Version_13); !ok && err == nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "this resource requires Keycloak v13 or higher",
		}}
	} else if err != nil {
		return diag.FromErr(err)
	}

	defaultRoles := mapFromDataToDefaultRoles(data)

	realm, err := keycloakClient.GetRealm(ctx, defaultRoles.RealmId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(realm.DefaultRole.Id)

	composites, err := keycloakClient.GetDefaultRoles(ctx, defaultRoles.RealmId, realm.DefaultRole.Id)
	if err != nil {
		return diag.FromErr(err)
	}

	defaultRoleNamesList, defaultRoleIds := getDefaultRoles(composites)
	rolesList, err := keycloakClient.GetRealmRoles(ctx, defaultRoles.RealmId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, roleId := range append(defaultRoleIds, defaultRoles.DefaultRoleIds...) {
		role, err := keycloakClient.GetRole(ctx, defaultRoles.RealmId, roleId)
		if err != nil {
			return diag.FromErr(err)
		}
		rolesList = append(rolesList, role)
	}

	// skip if actual default roles in keycloak same as we want
	if roleListsEqual(defaultRoleNamesList, defaultRoles.DefaultRoles) && roleListsEqual(defaultRoleIds, defaultRoles.DefaultRoleIds) {
		return nil
	}

	var putList, deleteList []*keycloak.Role
	for _, roleName := range defaultRoles.DefaultRoles {
		if !roleListContains(defaultRoleNamesList, roleName) {
			defaultRoles, err := getRoleByNameFromList(rolesList, roleName)
			if err != nil {
				return diag.FromErr(err)
			}
			putList = append(putList, defaultRoles)
		}
	}
	for _, roleId := range defaultRoles.DefaultRoleIds {
		if !roleListContains(defaultRoleIds, roleId) {
			defaultRoles, err := getRoleByIdFromList(rolesList, roleId)
			if err != nil {
				return diag.FromErr(err)
			}
			putList = append(putList, defaultRoles)
		}
	}
	for _, roleName := range defaultRoleNamesList {
		if !roleListContains(defaultRoles.DefaultRoles, roleName) {
			defaultRoles, err := getRoleByNameFromList(rolesList, roleName)
			if err != nil {
				return diag.FromErr(err)
			}
			deleteList = append(deleteList, defaultRoles)
		}
	}
	for _, roleId := range defaultRoleIds {
		if !roleListContains(defaultRoles.DefaultRoleIds, roleId) {
			defaultRoles, err := getRoleByIdFromList(rolesList, roleId)
			if err != nil {
				return diag.FromErr(err)
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
		err := keycloakClient.AddCompositesToRole(ctx, role, putList)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if len(deleteList) > 0 {
		role := &keycloak.Role{
			RealmId: defaultRoles.RealmId,
			Id:      realm.DefaultRole.Id,
		}
		err := keycloakClient.RemoveCompositesFromRole(ctx, role, deleteList)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceKeycloakDefaultRolesRead(ctx, data, meta)
}

// remove all roles from default
func resourceKeycloakDefaultRolesDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	realm, err := keycloakClient.GetRealm(ctx, realmId)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	defaultRoles, err := keycloakClient.GetDefaultRoles(ctx, realmId, realm.DefaultRole.Id)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(defaultRoles) > 0 {
		role := &keycloak.Role{
			RealmId: realmId,
			Id:      realm.DefaultRole.Id,
		}
		err := keycloakClient.RemoveCompositesFromRole(ctx, role, defaultRoles)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceKeycloakDefaultRolesImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import, supported import format: {{realm}}/{{defaultRoleId}}")
	}

	_, err := keycloakClient.GetDefaultRoles(ctx, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	diagnostics := resourceKeycloakDefaultRolesRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

// getDefaultRoles sorts out default realm role names and client role IDs
func getDefaultRoles(roles []*keycloak.Role) ([]string, []string) {
	var (
		defaultRolesNames []string
		defaultRolesIds   []string
	)
	for _, defaultRoles := range roles {
		if defaultRoles.ClientRole {
			defaultRolesIds = append(defaultRolesIds, defaultRoles.Id)
			continue
		}
		defaultRolesNames = append(defaultRolesNames, defaultRoles.Name)
	}
	return defaultRolesNames, defaultRolesIds
}

func getRoleByNameFromList(defaultRoles []*keycloak.Role, name string) (*keycloak.Role, error) {
	return getRoleBy(defaultRoles, func(r *keycloak.Role) bool { return r.Name == name }, " by name: "+name)
}

func getRoleByIdFromList(defaultRoles []*keycloak.Role, id string) (*keycloak.Role, error) {
	return getRoleBy(defaultRoles, func(r *keycloak.Role) bool { return r.Id == id }, " by id: "+id)
}

func getRoleBy(defaultRoles []*keycloak.Role, grep func(*keycloak.Role) bool, msg string) (*keycloak.Role, error) {
	for _, element := range defaultRoles {
		if grep(element) {
			return element, nil
		}
	}
	return nil, fmt.Errorf("defaultRoles not found%s", msg)
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
