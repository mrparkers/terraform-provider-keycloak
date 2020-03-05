package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKeycloakUserRoles_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user_roles.userRoles"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserRoles_basic(realmName, roleName, username),
				Check:  testAccCheckKeycloakUserRolesExists(resourceName),
			},
			{
				// we need a separate test for destroy instead of using CheckDestroy because this resource is implicitly
				// destroyed at the end of each test via destroying users or groups they're tied to
				Config: testKeycloakUserRoles_noRole(realmName, roleName, username),
				Check:  testAccCheckUsersDontHaveRole(resourceName, []string{username}),
			},
		},
	})
}

func testAccCheckUsersDontHaveRole(resourceName string, username []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

func testKeycloakUserRoles_noRole(realm, roleName, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}
	`, realm, username)
}

func testKeycloakUserRoles_basic(realm, roleName, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_role" "role" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}

resource "keycloak_user_roles" "userRoles" {
	realm_id = "${keycloak_realm.realm.id}"
	user_id = "${keycloak_user.user.id}"
	roles = [
		"${keycloak_role.role.name}"
	]
}
	`, realm, roleName, username)
}

func testAccCheckKeycloakUserRolesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRealmRolesFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}
