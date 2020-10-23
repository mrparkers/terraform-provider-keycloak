package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func rolesScopePermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policies": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"decision_strategy": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(keycloakOpenidClientResourcePermissionDecisionStrategies, false),
				},
			},
		},
	}
}

func resourceKeycloakRolePermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRolePermissionsCreate,
		Read:   resourceKeycloakRolePermissionsRead,
		Delete: resourceKeycloakRolePermissionsDelete,
		Update: resourceKeycloakRolePermissionsUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRolePermissionsImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"authorization_resource_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource server id representing the realm management client on which this permission is managed",
			},
			"map_role_scope":              rolesScopePermissionsSchema(),
			"map_role_client_scope_scope": rolesScopePermissionsSchema(),
			"map_role_composite_scope":    rolesScopePermissionsSchema(),
		},
	}
}

func rolePermissionsId(realmId, roleId string) string {
	return fmt.Sprintf("%s/%s", realmId, roleId)
}

func getRoleScopePermissions(keycloakClient *keycloak.KeycloakClient, realmId string, realmManagementClientId, permissionId string) (map[string]interface{}, error) {
	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClientId, permissionId)
	if err != nil {
		return nil, err
	}

	if permission.Description == "" && permission.DecisionStrategy == "AFFIRMATIVE" && len(permission.Policies) == 0 {
		return nil, nil
	}

	permissionViewSettings := make(map[string]interface{})

	if permission.Description != "" {
		permissionViewSettings["description"] = permission.Description
	}

	if permission.DecisionStrategy != "" {
		permissionViewSettings["decision_strategy"] = permission.DecisionStrategy
	}

	if len(permission.Policies) > 0 {
		permissionViewSettings["policies"] = permission.Policies
	}

	return permissionViewSettings, nil
}

func setRoleScopePermission(keycloakClient *keycloak.KeycloakClient, realmId, roleId string, realmManagementClientId string, authorizationPermissionId string, scopeDataSet *schema.Set) error {
	var policies []string

	scopeData := scopeDataSet.List()[0]
	scopePermission := scopeData.(map[string]interface{})

	if v, ok := scopePermission["policies"]; ok {
		for _, policy := range v.(*schema.Set).List() {
			policies = append(policies, policy.(string))
		}
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClientId, authorizationPermissionId)
	if err != nil {
		return err
	}

	if v, ok := scopePermission["description"]; ok {
		permission.Description = v.(string)
	}

	if v, ok := scopePermission["decision_strategy"]; ok {
		permission.DecisionStrategy = v.(string)
	}

	permission.Policies = policies

	return keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
}

func resourceKeycloakRolePermissionsCreate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakRolePermissionsUpdate(data, meta)
}

func resourceKeycloakRolePermissionsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	roleId := data.Get("role_id").(string)

	// the existence of this resource implies that it is enabled.
	err := keycloakClient.EnableRolePermissions(realmId, roleId)
	if err != nil {
		return err
	}

	rolePermissions, err := keycloakClient.GetRolePermissions(realmId, roleId)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	mapRolesScope, ok := data.GetOk("map_role_scope")
	if ok {
		err := setRoleScopePermission(keycloakClient, realmId, roleId, realmManagementClient.Id, rolePermissions.ScopePermissions["map-role"].(string), mapRolesScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	mapRolesClientsScope, ok := data.GetOk("map_role_client_scope_scope")
	if ok {
		err := setRoleScopePermission(keycloakClient, realmId, roleId, realmManagementClient.Id, rolePermissions.ScopePermissions["map-role-client-scope"].(string), mapRolesClientsScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	mapRolesCompositeScope, ok := data.GetOk("map_role_composite_scope")
	if ok {
		err := setRoleScopePermission(keycloakClient, realmId, roleId, realmManagementClient.Id, rolePermissions.ScopePermissions["map-role-composite"].(string), mapRolesCompositeScope.(*schema.Set))
		if err != nil {
			return err
		}
	}

	return resourceKeycloakRolePermissionsRead(data, meta)
}

func resourceKeycloakRolePermissionsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	roleId := data.Get("role_id").(string)

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	rolePermissions, err := keycloakClient.GetRolePermissions(realmId, roleId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	data.SetId(rolePermissionsId(rolePermissions.RealmId, rolePermissions.RoleId))
	data.Set("realm_id", rolePermissions.RealmId)
	data.Set("role_id", rolePermissions.RoleId)
	data.Set("enabled", rolePermissions.Enabled)
	data.Set("authorization_resource_server_id", realmManagementClient.Id)

	permissionMapRoles, err := getRoleScopePermissions(keycloakClient, realmId, realmManagementClient.Id, rolePermissions.ScopePermissions["map-role"].(string))
	if err != nil {
		return err
	}
	if permissionMapRoles != nil {
		data.Set("map_role_scope", []interface{}{permissionMapRoles})
	}

	permissionMapRolesClientScope, err := getRoleScopePermissions(keycloakClient, realmId, realmManagementClient.Id, rolePermissions.ScopePermissions["map-role-client-scope"].(string))
	if err != nil {
		return err
	}
	if permissionMapRolesClientScope != nil {
		data.Set("map_role_client_scope_scope", []interface{}{permissionMapRolesClientScope})
	}

	permissionMapRolesComposite, err := getRoleScopePermissions(keycloakClient, realmId, realmManagementClient.Id, rolePermissions.ScopePermissions["map-role-composite"].(string))
	if err != nil {
		return err
	}
	if permissionMapRolesComposite != nil {
		data.Set("map_role_composite_scope", []interface{}{permissionMapRolesComposite})
	}

	return nil
}

func resourceKeycloakRolePermissionsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	roleId := data.Get("role_id").(string)

	return keycloakClient.DisableRolePermissions(realmId, roleId)
}

func resourceKeycloakRolePermissionsImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{roleId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("role_id", parts[1])

	d.SetId(rolePermissionsId(parts[0], parts[1]))

	return []*schema.ResourceData{d}, nil
}
