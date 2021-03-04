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
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authorization_resource_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource server id representing the realm management client on which this permission is managed",
			},
			"view_scope":                   scopePermissionsSchema(),
			"manage_scope":                 scopePermissionsSchema(),
			"configure_scope":              scopePermissionsSchema(),
			"map_roles_scope":              scopePermissionsSchema(),
			"map_roles_client_scope_scope": scopePermissionsSchema(),
			"map_roles_composite_scope":    scopePermissionsSchema(),
			"token_exchange_scope":         scopePermissionsSchema(),
		},
	}
}

func clientPermissionsId(realmId, clientId string) string {
	return fmt.Sprintf("%s/%s", realmId, clientId)
}

// TODO: is this needed?
//func unsetOpenidClientScopePermissionPolicy(keycloakClient *keycloak.KeycloakClient, realmId, clientId, scopeName string) error {
//	openidClientPermissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
//	if err != nil {
//		return err
//	}
//
//	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
//	if err != nil {
//		return err
//	}
//
//	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions[scopeName].(string))
//	if err != nil {
//		return err
//	}
//
//	permission.Policies = []string{}
//	err = keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func resourceKeycloakOpenidClientPermissionsCreate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakOpenidClientPermissionsUpdate(data, meta)
}

func resourceKeycloakOpenidClientPermissionsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	// the existence of this resource implies that permissions are enabled for this client.
	err := keycloakClient.EnableOpenidClientPermissions(realmId, clientId)
	if err != nil {
		return err
	}

	openidClientPermissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	if viewScope, ok := data.GetOk("view_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["view"].(string), viewScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if manageScope, ok := data.GetOk("manage_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["manage"].(string), manageScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if configureScope, ok := data.GetOk("configure_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["configure"].(string), configureScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if mapRolesScope, ok := data.GetOk("map_roles_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles"].(string), mapRolesScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if mapRolesClientsScope, ok := data.GetOk("map_roles_client_scope_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles-client-scope"].(string), mapRolesClientsScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if mapRolesCompositeScope, ok := data.GetOk("map_roles_composite_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles-composite"].(string), mapRolesCompositeScope.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if tokenExchangeScope, ok := data.GetOk("token_exchange_scope"); ok {
		err := setOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["token-exchange"].(string), tokenExchangeScope.(*schema.Set))
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

	if !openidClientPermissions.Enabled {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer enabled", data.Id())
		data.SetId("")
		return nil
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	data.SetId(clientPermissionsId(openidClientPermissions.RealmId, openidClientPermissions.ClientId))
	data.Set("realm_id", openidClientPermissions.RealmId)
	data.Set("client_id", openidClientPermissions.ClientId)
	data.Set("enabled", openidClientPermissions.Enabled)
	data.Set("authorization_resource_server_id", realmManagementClient.Id)

	if viewScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["view"].(string)); err == nil && viewScope != nil {
		data.Set("view_scope", []interface{}{viewScope})
	} else if err != nil {
		return err
	}

	if manageScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["manage"].(string)); err == nil && manageScope != nil {
		data.Set("manage_scope", []interface{}{manageScope})
	} else if err != nil {
		return err
	}

	if mapRolesScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["configure"].(string)); err == nil && mapRolesScope != nil {
		data.Set("configure_scope", []interface{}{mapRolesScope})
	} else if err != nil {
		return err
	}

	if manageGroupMembershipScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles"].(string)); err == nil && manageGroupMembershipScope != nil {
		data.Set("map_roles_scope", []interface{}{manageGroupMembershipScope})
	} else if err != nil {
		return err
	}

	if impersonateScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles-client-scope"].(string)); err == nil && impersonateScope != nil {
		data.Set("map_roles_client_scope_scope", []interface{}{impersonateScope})
	} else if err != nil {
		return err
	}

	if userImpersonatedScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["map-roles-composite"].(string)); err == nil && userImpersonatedScope != nil {
		data.Set("map_roles_composite_scope", []interface{}{userImpersonatedScope})
	} else if err != nil {
		return err
	}

	if tokenExchangeScope, err := getOpenidClientScopePermissionPolicy(keycloakClient, realmId, realmManagementClient.Id, openidClientPermissions.ScopePermissions["token-exchange"].(string)); err == nil && tokenExchangeScope != nil {
		data.Set("token_exchange_scope", []interface{}{tokenExchangeScope})
	} else if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientPermissionsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

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
