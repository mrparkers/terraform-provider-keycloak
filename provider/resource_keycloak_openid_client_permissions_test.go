package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenidClientPermission_basic(t *testing.T) {
	realmName := "tf_view-" + acctest.RandString(10)
	clientId := "tf-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientPermission_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientPermissionExists("keycloak_openid_client_permissions.my_permission"),
			}, {
				Config: testKeycloakOpenidClientPermissionDelete_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientPermissionDoentExists("keycloak_openid_client_permissions.my_permission"),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientPermissionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		permissions, err := getOpenidClientPermissionsFromState(s, resourceName)
		if err != nil {
			return err
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		authorizationResourceServerId := rs.Primary.Attributes["authorization_resource_server_id"]
		viewScopePolicyId := rs.Primary.Attributes["view_scope_policy_id"]

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

		authzClient, err := keycloakClient.GetOpenidClientAuthorizationPermission(permissions.RealmId, realmManagementId, permissions.ScopePermissions["view"].(string))
		if err != nil {
			return err
		}

		policyId := authzClient.Policies[0]
		if viewScopePolicyId != policyId {
			return fmt.Errorf("computed ViewScopePolicyId %s was not equal to policyId %s", viewScopePolicyId, policyId)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientPermissionDoentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		permissions, err := getOpenidClientPermissionsFromState(s, resourceName)
		if err != nil {
			return err
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if permissions.Enabled != false {
			return fmt.Errorf("Client Permission in Keycloak is not disabled")
		}
		if rs.Primary.Attributes["enabled"] != "false" {
			return fmt.Errorf("Client Permission State is not disabled")
		}

		return nil
	}
}

func getOpenidClientPermissionsFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientPermissions, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]

	permissions, err := keycloakClient.GetOpenidClientPermissions(realmId, clientId)
	if err != nil {
		return nil, fmt.Errorf("error getting openid_client permissions with realm id %s and client id %s: %s", realmId, clientId, err)

	}
	return permissions, nil
}

func testKeycloakOpenidClientPermission_basic(realmId, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

resource keycloak_openid_client "my_openid_client" {
  realm_id              = keycloak_realm.realm.id
  name                  = "my_openid_client"
  client_id             = "%s"
  client_secret         = "secret"
  access_type           = "CONFIDENTIAL"
  standard_flow_enabled = true
  valid_redirect_uris   = [
    "http://localhost:8080/*",
  ]
}

data keycloak_openid_client "realm_management" {
  realm_id  = keycloak_realm.realm.id
  client_id = "realm-management"  
}


resource keycloak_openid_client_permissions "realm-management_permission" {
	realm_id   = keycloak_realm.realm.id
	client_id  = data.keycloak_openid_client.realm_management.id
	enabled = true
}

resource keycloak_user test {
	realm_id = keycloak_realm.realm.id
	username = "test-user"

	email      = "test-user@fakedomain.com"
	first_name = "Testy"
	last_name  = "Tester"
}

resource keycloak_openid_client_user_policy test {
	resource_server_id = "${data.keycloak_openid_client.realm_management.id}"
	realm_id = keycloak_realm.realm.id
	name = "client_user_policy_test"
	users = ["${keycloak_user.test.id}"]
	logic = "POSITIVE"
	decision_strategy = "UNANIMOUS"
	depends_on = [
		keycloak_openid_client_permissions.realm-management_permission,
	]
}

resource "keycloak_openid_client_permissions" "my_permission" {
	realm_id                               = keycloak_realm.realm.id
	client_id                              = keycloak_openid_client.my_openid_client.id

	enabled = true

	view_scope_policy_id                   = keycloak_openid_client_user_policy.test.id
	manage_scope_policy_id                 = keycloak_openid_client_user_policy.test.id
	configure_scope_policy_id              = keycloak_openid_client_user_policy.test.id
	map_roles_scope_policy_id              = keycloak_openid_client_user_policy.test.id
	map_roles_client_scope_scope_policy_id = keycloak_openid_client_user_policy.test.id
	map_roles_composite_scope_policy_id    = keycloak_openid_client_user_policy.test.id
	token_exchange_scope_policy_id         = keycloak_openid_client_user_policy.test.id
}

	`, realmId, clientId)
}

func testKeycloakOpenidClientPermissionDelete_basic(realmId, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

resource keycloak_openid_client "my_openid_client" {
  realm_id              = keycloak_realm.realm.id
  name                  = "my_openid_client"
  client_id             = "%s"
  client_secret         = "secret"
  access_type           = "CONFIDENTIAL"
  standard_flow_enabled = true
  valid_redirect_uris   = [
    "http://localhost:8080/*",
  ]
}

data keycloak_openid_client "realm_management" {
  realm_id  = keycloak_realm.realm.id
  client_id = "realm-management"  
}


resource keycloak_openid_client_permissions "realm-management_permission" {
	realm_id   = keycloak_realm.realm.id
	client_id  = data.keycloak_openid_client.realm_management.id
	enabled = true
}

resource keycloak_user test {
	realm_id = keycloak_realm.realm.id
	username = "test-user"

	email      = "test-user@fakedomain.com"
	first_name = "Testy"
	last_name  = "Tester"
}

resource keycloak_openid_client_user_policy test {
	resource_server_id = "${data.keycloak_openid_client.realm_management.id}"
	realm_id = keycloak_realm.realm.id
	name = "client_user_policy_test"
	users = ["${keycloak_user.test.id}"]
	logic = "POSITIVE"
	decision_strategy = "UNANIMOUS"
	depends_on = [
		keycloak_openid_client_permissions.realm-management_permission,
	]
}

resource "keycloak_openid_client_permissions" "my_permission" {
	realm_id                               = keycloak_realm.realm.id
	client_id                              = keycloak_openid_client.my_openid_client.id

	enabled = false
}

	`, realmId, clientId)
}
