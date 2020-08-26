package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenidClientAuthorizationPermission_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	resourceName := "terraform-" + acctest.RandString(10)
	permissionName := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationPermissionDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationPermission_basic(realmName, clientId, resourceName, permissionName, scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationPermissionExists("keycloak_openid_client_authorization_permission.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationPermission_createAfterManualDestroy(t *testing.T) {
	var authorizationPermission = &keycloak.OpenidClientAuthorizationPermission{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	resourceName := "terraform-" + acctest.RandString(10)
	permissionName := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationPermissionDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationPermission_basic(realmName, clientId, resourceName, permissionName, scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationPermissionFetch("keycloak_openid_client_authorization_permission.test", authorizationPermission),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClientAuthorizationPermission(authorizationPermission.RealmId, authorizationPermission.ResourceServerId, authorizationPermission.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientAuthorizationPermission_basic(realmName, clientId, resourceName, permissionName, scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationPermissionExists("keycloak_openid_client_authorization_permission.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationPermission_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	resourceName := "terraform-" + acctest.RandString(10)
	permissionName := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationPermissionDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationPermission_basic(firstRealm, clientId, resourceName, permissionName, scopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientAuthorizationPermissionExists("keycloak_openid_client_authorization_permission.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_authorization_permission.test", "realm_id", firstRealm),
				),
			},
			{
				Config: testKeycloakOpenidClientAuthorizationPermission_basic(secondRealm, clientId, resourceName, permissionName, scopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientAuthorizationPermissionExists("keycloak_openid_client_authorization_permission.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_authorization_permission.test", "realm_id", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationPermission_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	firstAuthrorizationPermission := &keycloak.OpenidClientAuthorizationPermission{
		RealmId:     realmName,
		Name:        acctest.RandString(10),
		Description: acctest.RandString(10),
	}

	secondAuthrorizationPermission := &keycloak.OpenidClientAuthorizationPermission{
		RealmId:     realmName,
		Name:        acctest.RandString(10),
		Description: acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationPermissionDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationPermission_basicFromInterface(clientId, firstAuthrorizationPermission, acctest.RandString(10), scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationPermissionExists("keycloak_openid_client_authorization_permission.test"),
			},
			{
				Config: testKeycloakOpenidClientAuthorizationPermission_basicFromInterface(clientId, secondAuthrorizationPermission, acctest.RandString(10), scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationPermissionExists("keycloak_openid_client_authorization_permission.test"),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientAuthorizationPermissionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakOpenidClientAuthorizationPermissionFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAuthorizationPermissionFetch(resourceName string, mapper *keycloak.OpenidClientAuthorizationPermission) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakOpenidClientAuthorizationPermissionFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.ResourceServerId = fetchedMapper.ResourceServerId
		mapper.RealmId = fetchedMapper.RealmId
		mapper.Id = fetchedMapper.Id

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAuthorizationPermissionDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_authorization_permission" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			resourceServerId := rs.Primary.Attributes["resource_server_id"]
			id := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			authorizationPermission, _ := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, resourceServerId, id)
			if authorizationPermission != nil {
				return fmt.Errorf("test config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakOpenidClientAuthorizationPermissionFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientAuthorizationPermission, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	resourceServerId := rs.Primary.Attributes["resource_server_id"]
	id := rs.Primary.ID

	authorizationPermission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, resourceServerId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting authorization permission config with id %s: %s", id, err)
	}

	return authorizationPermission, nil
}

func testKeycloakOpenidClientAuthorizationPermission_basic(realm, clientId, resourceName, permissionName, scopeName string) string {
	return fmt.Sprintf(`
resource keycloak_realm test {
	realm = "%s"
}

resource keycloak_openid_client test {
	client_id                = "%s"
	realm_id                 = "${keycloak_realm.test.id}"
	access_type              = "CONFIDENTIAL"
	service_accounts_enabled = true
	authorization {
		policy_enforcement_mode = "ENFORCING"
	}
}

data keycloak_openid_client_authorization_policy default {
  realm_id           = "${keycloak_realm.test.id}"
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "default"
}

resource keycloak_openid_client_authorization_resource test {
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "%s"
  realm_id           = "${keycloak_realm.test.id}"

  uris = [
    "/endpoint/*"
  ]
}

resource keycloak_openid_client_authorization_scope test {
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "%s"
  realm_id           = "${keycloak_realm.test.id}"
}

resource keycloak_openid_client_authorization_permission test {
	resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
	realm_id           = "${keycloak_realm.test.id}"
	name               = "%s"
	policies           = ["${data.keycloak_openid_client_authorization_policy.default.id}"]
	resources          = ["${keycloak_openid_client_authorization_resource.test.id}"]

}
	`, realm, clientId, resourceName, scopeName, permissionName)
}

func testKeycloakOpenidClientAuthorizationPermission_basicFromInterface(clientId string, authorizationPermission *keycloak.OpenidClientAuthorizationPermission, resourceName, scopeName string) string {
	return fmt.Sprintf(`
resource keycloak_realm test {
	realm = "%s"
}

resource keycloak_openid_client test {
	client_id                = "%s"
	realm_id                 = "${keycloak_realm.test.id}"
	access_type              = "CONFIDENTIAL"
	service_accounts_enabled = true
	authorization {
		policy_enforcement_mode = "ENFORCING"
	}
}

data keycloak_openid_client_authorization_policy default {
  realm_id           = "${keycloak_realm.test.id}"
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "default"
}

resource keycloak_openid_client_authorization_resource resource {
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "%s"
	realm_id           = "${keycloak_realm.test.id}"

  uris = [
    "/endpoint/*"
  ]
}

resource keycloak_openid_client_authorization_scope test {
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "%s"
  realm_id           = "${keycloak_realm.test.id}"
}

resource keycloak_openid_client_authorization_permission test {
	resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
	realm_id           = "${keycloak_realm.test.id}"
	name               = "%s"
	policies           = ["${data.keycloak_openid_client_authorization_policy.default.id}"]
   resources          = ["${keycloak_openid_client_authorization_resource.resource.id}"]
	 description        = "%s"
	scopes = ["${keycloak_openid_client_authorization_scope.test.id}"]
}
	`, authorizationPermission.RealmId, clientId, resourceName, scopeName, authorizationPermission.Name, authorizationPermission.Description)
}
