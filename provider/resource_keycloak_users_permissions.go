package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"view_scope":                    scopePermissionsSchema(),
			"manage_scope":                  scopePermissionsSchema(),
			"map_roles_scope":               scopePermissionsSchema(),
			"manage_group_membership_scope": scopePermissionsSchema(),
			"impersonate_scope":             scopePermissionsSchema(),
			"user_impersonated_scope":       scopePermissionsSchema(),
		},
	}
}

func resourceKeycloakUsersPermissionsCreate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakUsersPermissionsUpdate(data, meta)
}

func resourceKeycloakUsersPermissionsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)

	// the existence of this resource implies that it is enabled.
	err := keycloakClient.EnableUsersPermissions(realmId)
	if err != nil {
		return err
	}

	// setting scope permissions requires us to fetch the users permissions details, as well as the realm management client
	usersPermissions, err := keycloakClient.GetUsersPermissions(realmId)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	if viewScope, ok := data.GetOk("view_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["view"].(string), viewScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if manageScope, ok := data.GetOk("manage_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["manage"].(string), manageScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if mapRolesScope, ok := data.GetOk("map_roles_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["map-roles"].(string), mapRolesScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if manageGroupMembershipScope, ok := data.GetOk("manage_group_membership_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["manage-group-membership"].(string), manageGroupMembershipScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if impersonateScope, ok := data.GetOk("impersonate_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["impersonate"].(string), impersonateScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if userImpersonatedScope, ok := data.GetOk("user_impersonated_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["user-impersonated"].(string), userImpersonatedScope.(*schema.Set))
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

	if viewScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["view"].(string)); err == nil && viewScope != nil {
		data.Set("view_scope", []interface{}{viewScope})
	} else if err != nil {
		return err
	}

	if manageScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["manage"].(string)); err == nil && manageScope != nil {
		data.Set("manage_scope", []interface{}{manageScope})
	} else if err != nil {
		return err
	}

	if mapRolesScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["map-roles"].(string)); err == nil && mapRolesScope != nil {
		data.Set("map_roles_scope", []interface{}{mapRolesScope})
	} else if err != nil {
		return err
	}

	if manageGroupMembershipScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["manage-group-membership"].(string)); err == nil && manageGroupMembershipScope != nil {
		data.Set("manage_group_membership_scope", []interface{}{manageGroupMembershipScope})
	} else if err != nil {
		return err
	}

	if impersonateScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["impersonate"].(string)); err == nil && impersonateScope != nil {
		data.Set("impersonate_scope", []interface{}{impersonateScope})
	} else if err != nil {
		return err
	}

	if userImpersonatedScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, usersPermissions.ScopePermissions["user-impersonated"].(string)); err == nil && userImpersonatedScope != nil {
		data.Set("user_impersonated_scope", []interface{}{userImpersonatedScope})
	} else if err != nil {
		return err
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
