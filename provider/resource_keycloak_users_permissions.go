package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakUsersPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUsersPermissionsCreate,
		Read:   resourceKeycloakUsersPermissionsRead,
		Delete: resourceKeycloakUsersPermissionsDelete,
		Update: resourceKeycloakUsersPermissionsUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakUsersPermissionsImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
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
			"view_scope": {
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
			},
			"manage_scope": {
				Type:     schema.TypeSet,
				Optional: true,
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
			},
			"map_roles_scope": {
				Type:     schema.TypeSet,
				Optional: true,
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
			},
			"manage_group_membership_scope": {
				Type:     schema.TypeSet,
				Optional: true,
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
			},
			"impersonate_scope": {
				Type:     schema.TypeSet,
				Optional: true,
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
			},
			"user_impersonated_scope": {
				Type:     schema.TypeSet,
				Optional: true,
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
			},
		},
	}
}

func getUsersScopePermissions(keycloakClient *keycloak.KeycloakClient, realmId string, realmManagementClientId, permissionId string) (map[string]interface{}, error) {

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClientId, permissionId)
	if err != nil {
		return nil, err
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

	if permission.Description == "" && permission.DecisionStrategy == "UNANIMOUS" && len(permission.Policies) == 0 {
		permissionViewSettings = nil
	}

	return permissionViewSettings, nil
}

func setUsersScopePermission(keycloakClient *keycloak.KeycloakClient, realmId string, scopeName string, scopeDataSet *schema.Set) error {

	var policies []string

	scopeData := scopeDataSet.List()[0]
	scopePermission := scopeData.(map[string]interface{})

	if v, ok := scopePermission["policies"]; ok {
		for _, policy := range v.(*schema.Set).List() {
			policies = append(policies, policy.(string))
		}
	}

	usersPermissions, err := keycloakClient.GetUsersPermissions(realmId)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, usersPermissions.ScopePermissions[scopeName].(string))
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

func resourceKeycloakUsersPermissionsCreate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakUsersPermissionsUpdate(data, meta)
}

func resourceKeycloakUsersPermissionsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)

	err := keycloakClient.EnableUsersPermissions(realmId)
	if err != nil {
		return err
	}

	viewScope, ok := data.GetOk("view_scope")
	if ok {
		err := setUsersScopePermission(keycloakClient, realmId, "view", viewScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	manageScope, ok := data.GetOk("manage_scope")
	if ok {
		err := setUsersScopePermission(keycloakClient, realmId, "manage", manageScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	mapRolesScope, ok := data.GetOk("map_roles_scope")
	if ok {
		err := setUsersScopePermission(keycloakClient, realmId, "map-roles", mapRolesScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	manageGroupMembershipScope, ok := data.GetOk("manage_group_membership_scope")
	if ok {
		err := setUsersScopePermission(keycloakClient, realmId, "manage-group-membership", manageGroupMembershipScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	impersonateScope, ok := data.GetOk("impersonate_scope")
	if ok {
		err := setUsersScopePermission(keycloakClient, realmId, "impersonate", impersonateScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	userImpersonatedScope, ok := data.GetOk("user_impersonated_scope")
	if ok {
		err := setUsersScopePermission(keycloakClient, realmId, "user-impersonated", userImpersonatedScope.(*schema.Set))
		if err != nil {
			return err
		}
	}

	return resourceKeycloakUsersPermissionsRead(data, meta)
}

func resourceKeycloakUsersPermissionsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	usersPermissions, err := keycloakClient.GetUsersPermissions(realmId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	data.SetId(usersPermissions.RealmId)
	data.Set("realm_id", usersPermissions.RealmId)
	data.Set("enabled", usersPermissions.Enabled)
	data.Set("authorization_resource_server_id", realmManagementClient.Id)

	permissionView, err := getUsersScopePermissions(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["view"].(string))
	if err != nil {
		return err
	}
	if permissionView != nil {
		data.Set("view_scope", []interface{}{permissionView})
	}

	permissionManage, err := getUsersScopePermissions(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["manage"].(string))
	if err != nil {
		return err
	}
	if permissionManage != nil {
		data.Set("manage_scope", []interface{}{permissionManage})
	}

	permissionMapRoles, err := getUsersScopePermissions(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["map-roles"].(string))
	if err != nil {
		return err
	}
	if permissionMapRoles != nil {
		data.Set("map_roles_scope", []interface{}{permissionMapRoles})
	}

	permissionManageGroupMembership, err := getUsersScopePermissions(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["manage-group-membership"].(string))
	if err != nil {
		return err
	}
	if permissionManageGroupMembership != nil {
		data.Set("manage_group_membership_scope", []interface{}{permissionManageGroupMembership})
	}

	permissionImpersonate, err := getUsersScopePermissions(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["impersonate"].(string))
	if err != nil {
		return err
	}
	if permissionImpersonate != nil {
		data.Set("impersonate_scope", []interface{}{permissionImpersonate})
	}

	permissionUserImpersonated, err := getUsersScopePermissions(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["user-impersonated"].(string))
	if err != nil {
		return err
	}
	if permissionUserImpersonated != nil {
		data.Set("user_impersonated_scope", []interface{}{permissionUserImpersonated})
	}

	return nil
}

func resourceKeycloakUsersPermissionsDelete(data *schema.ResourceData, meta interface{}) error {

	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)

	return keycloakClient.DisableUsersPermissions(realmId)
}

func resourceKeycloakUsersPermissionsImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {

	d.Set("realm_id", d.Id())
	d.SetId(d.Id())

	return []*schema.ResourceData{d}, nil
}
