package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakRolePermission_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRolePermissionAreDisabled("keycloak_role.role"),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRolePermission_basic(realmName, roleName),
				Check:  testAccCheckKeycloakRolePermissionExists("keycloak_role_permissions.my_permission"),
			},
			{
				ResourceName:      "keycloak_role_permissions.my_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKeycloakRolePermissionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		permissions, err := getRolePermissionsFromState(s, resourceName)
		if err != nil {
			return err
		}
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		authorizationResourceServerId := rs.Primary.Attributes["authorization_resource_server_id"]

		var realmManagementId string
		clients, _ := keycloakClient.GetOpenidClients(permissions.RealmId, false)
		for _, client := range clients {
			if client.ClientId == "realm-management" {
				realmManagementId = client.Id
				break
			}
		}

		if authorizationResourceServerId != realmManagementId {
			return fmt.Errorf("computed authorizationResourceServerId %s was not equal to %s (the id of the realm-management client)", authorizationResourceServerId, realmManagementId)
		}

		mapRoleScopePolicyId := rs.Primary.Attributes["map_role_scope.0.policies.0"]
		mapRoleScopeDescription := rs.Primary.Attributes["map_role_scope.0.description"]
		mapRoleScopeDecisionStrategy := rs.Primary.Attributes["map_role_scope.0.decision_strategy"]

		authzClientMapRolesScope, err := keycloakClient.GetOpenidClientAuthorizationPermission(permissions.RealmId, realmManagementId, permissions.ScopePermissions["map-role"].(string))
		if err != nil {
			return err
		}
		policyId := authzClientMapRolesScope.Policies[0]

		if mapRoleScopePolicyId != policyId {
			return fmt.Errorf("computed mapRoleScopePolicyId %s was not equal to policyId %s", mapRoleScopePolicyId, policyId)
		}

		if authzClientMapRolesScope.Description != mapRoleScopeDescription {
			return fmt.Errorf("DecisionStrategy %s was not equal to %s", authzClientMapRolesScope.DecisionStrategy, mapRoleScopeDescription)
		}

		if authzClientMapRolesScope.DecisionStrategy != mapRoleScopeDecisionStrategy {
			return fmt.Errorf("DecisionStrategy %s was not equal to %s", authzClientMapRolesScope.DecisionStrategy, mapRoleScopeDecisionStrategy)
		}

		mapRolesScope := rs.Primary.Attributes["map_role_client_scope_scope"]

		if mapRolesScope != "" {
			return fmt.Errorf("map_role_client_scope_scope found")
		}

		return nil
	}
}

func testAccCheckKeycloakRolePermissionAreDisabled(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realmId := rs.Primary.Attributes["realm_id"]
		roleId := rs.Primary.ID

		permissions, err := keycloakClient.GetRolePermissions(realmId, roleId)
		if err != nil {
			return fmt.Errorf("error getting role permissions with realm id %s and role id %s : %s", realmId, roleId, err)
		}

		if permissions.Enabled != false {
			return fmt.Errorf("Users Permission in Keycloak is not disabled")
		}

		return nil
	}
}

func getRolePermissionsFromState(s *terraform.State, resourceName string) (*keycloak.RolePermissions, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	roleId := rs.Primary.Attributes["role_id"]

	permissions, err := keycloakClient.GetRolePermissions(realmId, roleId)
	if err != nil {
		return nil, fmt.Errorf("error getting role permissions with realm id %s and role id %s : %s", realmId, roleId, err)
	}

	return permissions, nil
}

func testKeycloakRolePermission_basic(realmId, roleName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

data "keycloak_openid_client" "realm_management" {
  realm_id  = keycloak_realm.realm.id
  client_id = "realm-management"  
}

resource "keycloak_openid_client_permissions" "realm-management_permission" {
	realm_id   = keycloak_realm.realm.id
	client_id  = data.keycloak_openid_client.realm_management.id
	enabled    = true
}

resource "keycloak_user" "test" {
	realm_id   = keycloak_realm.realm.id
	username   = "test-user"

	email      = "test-user@fakedomain.com"
	first_name = "Testy"
	last_name  = "Tester"
}

resource "keycloak_openid_client_user_policy" "test" {
	realm_id           = keycloak_realm.realm.id
	resource_server_id = data.keycloak_openid_client.realm_management.id
	name 			   = "client_user_policy_test"

	users = [
		keycloak_user.test.id
	]
	logic = "POSITIVE"
	decision_strategy = "UNANIMOUS"
	depends_on = [
		keycloak_openid_client_permissions.realm-management_permission,
	]
}

resource "keycloak_role" "role" {
	realm_id = keycloak_realm.realm.id
	name     = "%s"
}

resource "keycloak_role_permissions" "my_permission" {
	realm_id                               = keycloak_realm.realm.id
	role_id                                = keycloak_role.role.id

	map_role_scope {
		policies          = [
			keycloak_openid_client_user_policy.test.id
		]
		description       = "map_role_scope"
		decision_strategy = "CONSENSUS"
	}
	
}

	`, realmId, roleName)
}
