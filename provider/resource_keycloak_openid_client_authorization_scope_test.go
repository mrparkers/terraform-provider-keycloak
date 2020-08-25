package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOpenidClientAuthorizationScope_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationScope_basic(realmName, clientId, scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationScopeExists("keycloak_openid_client_authorization_scope.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationScope_createAfterManualDestroy(t *testing.T) {
	var authorizationScope = &keycloak.OpenidClientAuthorizationScope{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationScope_basic(realmName, clientId, scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationScopeFetch("keycloak_openid_client_authorization_scope.test", authorizationScope),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClientAuthorizationScope(authorizationScope.RealmId, authorizationScope.ResourceServerId, authorizationScope.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientAuthorizationScope_basic(realmName, clientId, scopeName),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationScopeExists("keycloak_openid_client_authorization_scope.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationScope_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	scopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationScope_basic(firstRealm, clientId, scopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientAuthorizationScopeExists("keycloak_openid_client_authorization_scope.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_authorization_scope.test", "realm_id", firstRealm),
				),
			},
			{
				Config: testKeycloakOpenidClientAuthorizationScope_basic(secondRealm, clientId, scopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientAuthorizationScopeExists("keycloak_openid_client_authorization_scope.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_authorization_scope.test", "realm_id", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationScope_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	firstAuthrorizationScope := &keycloak.OpenidClientAuthorizationScope{
		RealmId:     realmName,
		Name:        acctest.RandString(10),
		DisplayName: acctest.RandString(10),
		IconUri:     acctest.RandString(10),
	}

	secondAuthrorizationScope := &keycloak.OpenidClientAuthorizationScope{
		RealmId:     realmName,
		Name:        acctest.RandString(10),
		DisplayName: acctest.RandString(10),
		IconUri:     acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientAuthorizationScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientAuthorizationScope_basicFromInterface(clientId, firstAuthrorizationScope),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationScopeExists("keycloak_openid_client_authorization_scope.test"),
			},
			{
				Config: testKeycloakOpenidClientAuthorizationScope_basicFromInterface(clientId, secondAuthrorizationScope),
				Check:  testAccCheckKeycloakOpenidClientAuthorizationScopeExists("keycloak_openid_client_authorization_scope.test"),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientAuthorizationScopeExists(scopeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakOpenidClientAuthorizationScopeFromState(s, scopeName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAuthorizationScopeFetch(scopeName string, authorizationScope *keycloak.OpenidClientAuthorizationScope) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedAuthorizationScope, err := getKeycloakOpenidClientAuthorizationScopeFromState(s, scopeName)
		if err != nil {
			return err
		}

		authorizationScope.ResourceServerId = fetchedAuthorizationScope.ResourceServerId
		authorizationScope.RealmId = fetchedAuthorizationScope.RealmId
		authorizationScope.Id = fetchedAuthorizationScope.Id

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAuthorizationScopeDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_authorization_scope" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			resourceServerId := rs.Primary.Attributes["resource_server_id"]
			id := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			authorizationScope, _ := keycloakClient.GetOpenidClientAuthorizationScope(realmId, resourceServerId, id)
			if authorizationScope != nil {
				return fmt.Errorf("test config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakOpenidClientAuthorizationScopeFromState(s *terraform.State, scopeName string) (*keycloak.OpenidClientAuthorizationScope, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[scopeName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", scopeName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	resourceServerId := rs.Primary.Attributes["resource_server_id"]
	id := rs.Primary.ID

	authorizationScope, err := keycloakClient.GetOpenidClientAuthorizationScope(realmId, resourceServerId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting authorization scope config with id %s: %s", id, err)
	}

	return authorizationScope, nil
}

func testKeycloakOpenidClientAuthorizationScope_basic(realm, clientId, scopeName string) string {
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

resource keycloak_openid_client_authorization_scope test {
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name               = "%s"
  realm_id           = "${keycloak_realm.test.id}"
}
	`, realm, clientId, scopeName)
}

func testKeycloakOpenidClientAuthorizationScope_basicFromInterface(clientId string, authorizationScope *keycloak.OpenidClientAuthorizationScope) string {
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

resource keycloak_openid_client_authorization_scope test {
  resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
  name                 = "%s"
  realm_id             = "${keycloak_realm.test.id}"
  display_name         = "%s"
  icon_uri             = "%s"
}
	`, authorizationScope.RealmId, clientId, authorizationScope.Name, authorizationScope.DisplayName, authorizationScope.IconUri)
}
