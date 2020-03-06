package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
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
				Check:  testAccCheckUsersDontHaveRole("keycloak_user.user", []string{roleName}),
			},
		},
	})
}

func testAccGetRolesFromUserFromUserState(resourceName string, s *terraform.State) ([]*keycloak.Role, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]

	var userId string
	if strings.HasPrefix(resourceName, "keycloak_user_roles") {
		userId = rs.Primary.Attributes["user_id"]
	} else {
		userId = rs.Primary.ID
	}

	user, err := keycloakClient.GetUser(realmId, userId)

	if err != nil {
		return nil, err
	}

	return keycloakClient.GetRealmLevelRoleMappings(user)
}

func testAccCheckUsersDontHaveRole(resourceName string, roles []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rolesInUser, err := testAccGetRolesFromUserFromUserState(resourceName, s)
		if err != nil {
			return err
		}

		for _, role := range roles {
			for _, roleInUser := range rolesInUser {
				if role == roleInUser.Name {
					return fmt.Errorf("expected role %s to not belong to user", role)
				}
			}
		}

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

resource "keycloak_role" "role" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}
	`, realm, username, roleName)
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
	`, realm, username, roleName)
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
