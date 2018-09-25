package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakClientScope_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakClientScopeExists("keycloak_client_scope.client-scope"),
			},
		},
	})
}

func TestAccKeycloakClientScope_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClientScope_updateRealmBefore(realmOne, realmTwo, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientScopeExists("keycloak_client_scope.client-scope"),
					testAccCheckKeycloakClientScopeBelongsToRealm("keycloak_client_scope.client-scope", realmOne),
				),
			},
			{
				Config: testKeycloakClientScope_updateRealmAfter(realmOne, realmTwo, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientScopeExists("keycloak_client_scope.client-scope"),
					testAccCheckKeycloakClientScopeBelongsToRealm("keycloak_client_scope.client-scope", realmTwo),
				),
			},
		},
	})
}

func testAccCheckKeycloakClientScopeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakClientScopeBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientScope, err := getClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		if clientScope.RealmId != realm {
			return fmt.Errorf("expected client scope %s to have realm_id of %s, but got %s", clientScope.Id, realm, clientScope.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakClientScopeDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_client_scope" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			clientScope, _ := keycloakClient.GetClientScope(realm, id)
			if clientScope != nil {
				return fmt.Errorf("client scope %s still exists", id)
			}
		}

		return nil
	}
}

func getClientScopeFromState(s *terraform.State, resourceName string) (*keycloak.ClientScope, error) {
	keycloakClientScope := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	clientScope, err := keycloakClientScope.GetClientScope(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting client scope %s: %s", id, err)
	}

	return clientScope, nil
}

func testKeycloakClientScope_basic(realm, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_client_scope" "client-scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}
	`, realm, clientScopeName)
}

func testKeycloakClientScope_updateRealmBefore(realmOne, realmTwo, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-1" {
	realm = "%s"
}

resource "keycloak_realm" "realm-2" {
	realm = "%s"
}

resource "keycloak_client_scope" "client-scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm-1.id}"
}
	`, realmOne, realmTwo, clientScopeName)
}

func testKeycloakClientScope_updateRealmAfter(realmOne, realmTwo, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-1" {
	realm = "%s"
}

resource "keycloak_realm" "realm-2" {
	realm = "%s"
}

resource "keycloak_client_scope" "client-scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm-2.id}"
}
	`, realmOne, realmTwo, clientScopeName)
}
