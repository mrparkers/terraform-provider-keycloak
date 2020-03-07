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

func TestAccKeycloakUserRoles_moreThan100roles(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	userName := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserRoles_moreThan100roles(realmName, userName),
			},
		},
	})
}

func TestAccKeycloakUserRoles_updateUserForceNew(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	userOne := "terraform-user-" + acctest.RandString(10)
	userTwo := "terraform-user-" + acctest.RandString(10)

	role := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserRoles_updateGroupForceNew(realmName, userOne, userTwo, role, "user_one"),
				Check:  testAccCheckRoleBelongsToUser("keycloak_user_roles.userRoles", role),
			},
			{
				Config: testKeycloakUserRoles_updateGroupForceNew(realmName, userOne, userTwo, role, "user_two"),
				Check:  testAccCheckRoleBelongsToUser("keycloak_user_roles.userRoles", role),
			},
		},
	})
}

func TestAccKeycloakUserRoles_updateInPlace(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	userName := "terraform-user-" + acctest.RandString(10)

	allRolesForTest := []string{
		"terraform-role-" + acctest.RandString(10),
		"terraform-role-" + acctest.RandString(10),
		"terraform-role-" + acctest.RandString(10),
	}
	indexOfRandomRoleToRemove := acctest.RandIntRange(0, len(allRolesForTest)-1)
	//randomRoleToRemove := allRolesForTest[indexOfRandomRoleToRemove]

	var subsetOfRoles []string
	for index, user := range allRolesForTest {
		if index != indexOfRandomRoleToRemove {
			subsetOfRoles = append(subsetOfRoles, user)
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// init
			{
				Config: testKeycloakUserRoles_multipleRoles(realmName, userName, allRolesForTest, allRolesForTest),
				Check:  testAccCheckUsersHaveRole("keycloak_user_roles.userRoles", allRolesForTest),
			},
			// remove
			//{
			//	Config: testKeycloakUserRoles_multipleRoles(realmName, userName, allRolesForTest, subsetOfRoles),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckUsersHaveRole("keycloak_user_roles.userRoles", subsetOfRoles),
			//		testAccCheckUsersDontHaveRole("keycloak_user_roles.userRoles", []string{randomRoleToRemove}),
			//	),
			//},
			//// add
			//{
			//	Config: testKeycloakUserRoles_multipleRoles(realmName, userName, allRolesForTest, allRolesForTest),
			//	Check:  testAccCheckUsersHaveRole("keycloak_user_roles.userRoles", allRolesForTest),
			//},
		},
	})
}

func testKeycloakUserRoles_multipleRoles(realmName, userName string, definiedRoles, rolesinUser []string) string {
	var roleResources strings.Builder
	for _, role := range definiedRoles {
		roleResources.WriteString(fmt.Sprintf(`
resource "keycloak_role" "role_%s" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}
		`, role, role))
	}

	var rolesInUserInterpolated []string
	for _, roleInUser := range rolesinUser {
		rolesInUserInterpolated = append(rolesInUserInterpolated, fmt.Sprintf("${keycloak_role.role_%s.name}", roleInUser))
	}

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

%s

resource "keycloak_user_roles" "userRoles" {
	realm_id = "${keycloak_realm.realm.id}"
	user_id = "${keycloak_user.user.id}"
	roles = %s
}
	`, realmName, userName, roleResources.String(), arrayOfStringsForTerraformResource(rolesInUserInterpolated))
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

	return keycloakClient.GetRealmRoleMappings(user)
}

func testAccCheckRoleBelongsToUser(resourceName, role string) resource.TestCheckFunc {
	return testAccCheckUsersHaveRole(resourceName, []string{role})
}

func testAccCheckUsersHaveRole(resourceName string, roles []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rolesInUser, err := testAccGetRolesFromUserFromUserState(resourceName, s)
		if err != nil {
			return err
		}

		for _, role := range roles {
			roleFound := false

			for _, roleInUser := range rolesInUser {
				if role == roleInUser.Name {
					roleFound = true

					break
				}
			}
			if !roleFound {
				return fmt.Errorf("unable to find role %s in user", role)
			}
		}

		return nil
	}
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

func testAccCheckKeycloakUserRolesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRealmRolesFromState(s, resourceName)
		if err != nil {
			return err
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

func testKeycloakUserRoles_moreThan100roles(realmName, userName string) string {
	count := 100
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_role" "role" {
	count = %d
		realm_id = "${keycloak_realm.realm.id}"
		name = "terraform-user-${count.index}"
}

resource "keycloak_user_roles" "userRoles" {
	realm_id = "${keycloak_realm.realm.id}"
	user_id = "${keycloak_user.user.id}"
	roles = "${keycloak_role.role.*.name}"
}
	`, realmName, userName, count)
}

func testKeycloakUserRoles_updateGroupForceNew(realm, userOne, userTwo, role, currentRole string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user_one" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_user" "user_two" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_role" "role" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}

resource "keycloak_user_roles" "userRoles" {
	realm_id = "${keycloak_realm.realm.id}"
	user_id = "${keycloak_user.%s.id}"
	roles = [
		"${keycloak_role.role.name}"
	]
}
	`, realm, userOne, userTwo, role, currentRole)
}
