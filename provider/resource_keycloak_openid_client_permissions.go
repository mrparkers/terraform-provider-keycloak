package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientPermissionsCreate,
		Read:   resourceKeycloakOpenidClientPermissionsRead,
		Delete: resourceKeycloakOpenidClientPermissionsDelete,
		Update: resourceKeycloakOpenidClientPermissionsUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientPermissionsImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"authorization_resource_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource server id representing the realm management client on which this permission is managed",
			},
			"view_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"manage_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configure_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"map_roles_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"map_roles_client_scope_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"map_roles_composite_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"token_exchange_scope_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func clientPermissionsId(realmId, clientId string) string {
	return fmt.Sprintf("%s/%s", realmId, clientId)
}

func setOpenidClientScopePermissionPolicy(keycloakClient *keycloak.KeycloakClient, realmId, clientId string, scopeName string, policyId string) error {
	openidClientPermissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions[scopeName].(string))
	if err != nil {
		return err
	}

	permission.Policies = []string{policyId}

	return keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
}

func unsetOpenidClientScopePermissionPolicy(keycloakClient *keycloak.KeycloakClient, realmId, clientId, scopeName string) error {
	openidClientPermissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions[scopeName].(string))
	if err != nil {
		return err
	}

	permission.Policies = []string{}
	err = keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientPermissionsCreate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakOpenidClientPermissionsUpdate(data, meta)
}

func resourceKeycloakOpenidClientPermissionsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	if data.Get("enabled").(bool) {
		err := keycloakClient.EnableOpenidClientPermissions(realmId, clientId)
		if err != nil {
			return err
		}
	} else {
		err := keycloakClient.DisableOpenidClientPermissions(realmId, clientId)
		if err != nil {
			return err
		}
	}

	viewScopePolicyId, ok := data.GetOkExists("view_scope_policy_id")
	if ok && viewScopePolicyId != nil {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "view", viewScopePolicyId.(string))
		if err != nil {
			return err
		}
	}
	manageScopePolicyId, ok := data.GetOkExists("manage_scope_policy_id")
	if ok && manageScopePolicyId != "" {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "manage", manageScopePolicyId.(string))
		if err != nil {
			return err
		}
	}
	configureScopePolicyId, ok := data.GetOkExists("configure_scope_policy_id")
	if ok && configureScopePolicyId != "" {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "configure", configureScopePolicyId.(string))
		if err != nil {
			return err
		}
	}
	mapRolesScopePolicyId, ok := data.GetOkExists("map_roles_scope_policy_id")
	if ok && mapRolesScopePolicyId != "" {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "map-roles", mapRolesScopePolicyId.(string))
		if err != nil {
			return err
		}
	}
	mapRolesClientsScopePolicyId, ok := data.GetOkExists("map_roles_client_scope_scope_policy_id")
	if ok && mapRolesClientsScopePolicyId != "" {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "map-roles-client-scope", mapRolesClientsScopePolicyId.(string))
		if err != nil {
			return err
		}
	}
	mapRolesCompositeScopePolicyId, ok := data.GetOkExists("map_roles_composite_scope_policy_id")
	if ok && mapRolesCompositeScopePolicyId != "" {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "map-roles-composite", mapRolesCompositeScopePolicyId.(string))
		if err != nil {
			return err
		}
	}
	tokenExchangeScopePolicyId, ok := data.GetOkExists("token_exchange_scope_policy_id")
	if ok && tokenExchangeScopePolicyId != "" {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "token-exchange", tokenExchangeScopePolicyId.(string))
		if err != nil {
			return err
		}
	}

	return resourceKeycloakOpenidClientPermissionsRead(data, meta)
}

func resourceKeycloakOpenidClientPermissionsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	openidClientPermissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	data.SetId(clientPermissionsId(openidClientPermissions.RealmId, openidClientPermissions.ClientId))
	data.Set("realm_id", openidClientPermissions.RealmId)
	data.Set("client_id", openidClientPermissions.ClientId)

	data.Set("enabled", openidClientPermissions.Enabled)

	if !openidClientPermissions.Enabled {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer enabled", data.Id())
		return nil
	}

	data.Set("view_scope_policy_id", nil)
	data.Set("manage_scope_policy_id", nil)
	data.Set("configure_scope_policy_id", nil)
	data.Set("map_roles_scope_policy_id", nil)
	data.Set("map_roles_client_scope_scope_policy_id", nil)
	data.Set("map_roles_composite_scope_policy_id", nil)
	data.Set("token_exchange_scope_policy_id", nil)

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}
	permissionView, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["view"].(string))
	if err != nil {
		return err
	}
	if permissionView != nil && len(permissionView.Policies) > 0 {
		data.Set("view_scope_policy_id", permissionView.Policies[0])
	}
	permissionManage, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["manage"].(string))
	if err != nil {
		return err
	}
	if permissionManage != nil && len(permissionManage.Policies) > 0 {
		data.Set("manage_scope_policy_id", permissionManage.Policies[0])
	}
	permissionConfigure, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["configure"].(string))
	if err != nil {
		return err
	}
	if permissionConfigure != nil && len(permissionConfigure.Policies) > 0 {
		data.Set("configure_scope_policy_id", permissionConfigure.Policies[0])
	}
	permissionMapRoles, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles"].(string))
	if err != nil {
		return err
	}
	if permissionMapRoles != nil && len(permissionMapRoles.Policies) > 0 {
		data.Set("map_roles_scope_policy_id", permissionMapRoles.Policies[0])
	}
	permissionMapRolesClientScope, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles-client-scope"].(string))
	if err != nil {
		return err
	}
	if permissionMapRolesClientScope != nil && len(permissionMapRolesClientScope.Policies) > 0 {
		data.Set("map_roles_client_scope_scope_policy_id", permissionMapRolesClientScope.Policies[0])
	}
	permissionMapRolesComposite, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles-composite"].(string))
	if err != nil {
		return err
	}
	if permissionMapRolesComposite != nil && len(permissionMapRolesComposite.Policies) > 0 {
		data.Set("map_roles_composite_scope_policy_id", permissionMapRolesComposite.Policies[0])
	}
	permissionTokenExchange, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["token-exchange"].(string))
	if err != nil {
		return err
	}
	if permissionTokenExchange != nil && len(permissionTokenExchange.Policies) > 0 {
		data.Set("token_exchange_scope_policy_id", permissionTokenExchange.Policies[0])
	}
	data.Set("authorization_resource_server_id", realmManagementClient.Id)

	return nil
}

func resourceKeycloakOpenidClientPermissionsDelete(data *schema.ResourceData, meta interface{}) error {

	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	openidClientPermissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
	if err == nil && openidClientPermissions.Enabled {
		_ = unsetOpenidClientScopePermissionPolicy(keycloakClient, realmId, clientId, "view")
	}
	return keycloakClient.DisableOpenidClientPermissions(realmId, clientId)
}

func resourceKeycloakOpenidClientPermissionsImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("client_id", parts[1])

	d.SetId(clientPermissionsId(parts[0], parts[1]))

	return []*schema.ResourceData{d}, nil
}
