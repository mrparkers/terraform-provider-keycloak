package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"strings"
	"testing"
)

func TestAccKeycloakGroupMemberships_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroupMemberships_basic(realmName, groupName, username),
				Check:  testAccCheckKeycloakUserBelongsToGroup("keycloak_group_memberships.group_members", username),
			},
		},
	})
}

func TestAccKeycloakGroupMemberships_updateGroupForceNew(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	groupOne := "terraform-group-" + acctest.RandString(10)
	groupTwo := "terraform-group-" + acctest.RandString(10)

	username := "terraform-user-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroupMemberships_updateGroupForceNew(realmName, groupOne, groupTwo, username, "group_one"),
				Check:  testAccCheckKeycloakUserBelongsToGroup("keycloak_group_memberships.group_members", username),
			},
			{
				Config: testKeycloakGroupMemberships_updateGroupForceNew(realmName, groupOne, groupTwo, username, "group_two"),
				Check:  testAccCheckKeycloakUserBelongsToGroup("keycloak_group_memberships.group_members", username),
			},
		},
	})
}

func TestAccKeycloakGroupMemberships_updateInPlace(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)

	allUsersForTest := []string{
		"terraform-user-" + acctest.RandString(10),
		"terraform-user-" + acctest.RandString(10),
		"terraform-user-" + acctest.RandString(10),
	}
	randomUserToRemove := acctest.RandIntRange(0, len(allUsersForTest)-1)

	var subsetOfUsers []string
	for index, user := range allUsersForTest {
		if index != randomUserToRemove {
			subsetOfUsers = append(subsetOfUsers, user)
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// init
			{
				Config: testKeycloakGroupMemberships_multipleUsers(realmName, groupName, allUsersForTest, allUsersForTest),
				Check:  testAccCheckKeycloakUsersBelongToGroup("keycloak_group_memberships.group_members", allUsersForTest),
			},
			// remove
			{
				Config: testKeycloakGroupMemberships_multipleUsers(realmName, groupName, allUsersForTest, subsetOfUsers),
				Check:  testAccCheckKeycloakUsersBelongToGroup("keycloak_group_memberships.group_members", subsetOfUsers),
			},
			// add
			{
				Config: testKeycloakGroupMemberships_multipleUsers(realmName, groupName, allUsersForTest, allUsersForTest),
				Check:  testAccCheckKeycloakUsersBelongToGroup("keycloak_group_memberships.group_members", allUsersForTest),
			},
		},
	})
}

func TestAccKeycloakGroupMemberships_userDoesNotExist(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakGroupMemberships_userDoesNotExist(realmName, groupName, username),
				ExpectError: regexp.MustCompile("user with username .+ does not exist"),
			},
		},
	})
}

func testAccCheckKeycloakUserBelongsToGroup(resourceName, user string) resource.TestCheckFunc {
	return testAccCheckKeycloakUsersBelongToGroup(resourceName, []string{user})
}

func testAccCheckKeycloakUsersBelongToGroup(resourceName string, users []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realmId := rs.Primary.Attributes["realm_id"]
		groupId := rs.Primary.Attributes["group_id"]

		usersInGroup, err := keycloakClient.GetGroupMembers(realmId, groupId)
		if err != nil {
			return err
		}

		for _, user := range users {
			userFound := false

			for _, userInGroup := range usersInGroup {
				if user == userInGroup.Username {
					userFound = true

					break
				}
			}

			if !userFound {
				return fmt.Errorf("unable to find user %s in group with id %s", user, groupId)
			}
		}

		return nil
	}
}

func testKeycloakGroupMemberships_basic(realm, group, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_group_memberships" "group_members" {
	realm_id = "${keycloak_realm.realm.id}"
	group_id = "${keycloak_group.group.id}"

	members = [
		"${keycloak_user.user.username}"
	]
}
	`, realm, group, username)
}

func testKeycloakGroupMemberships_updateGroupForceNew(realm, groupOne, groupTwo, username, currentGroup string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group_one" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_group" "group_two" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_group_memberships" "group_members" {
	realm_id = "${keycloak_realm.realm.id}"
	group_id = "${keycloak_group.%s.id}"

	members = [
		"${keycloak_user.user.username}"
	]
}
	`, realm, groupOne, groupTwo, username, currentGroup)
}

// this tf config provides a good way to test users that exist within keycloak but are not necessarily part of a group
func testKeycloakGroupMemberships_multipleUsers(realm, group string, definedUsers, usersInGroup []string) string {
	var userResources strings.Builder
	for _, username := range definedUsers {
		userResources.WriteString(fmt.Sprintf(`
resource "keycloak_user" "user_%s" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}
		`, username, username))
	}

	var usersInGroupInterpolated []string
	for _, userInGroup := range usersInGroup {
		usersInGroupInterpolated = append(usersInGroupInterpolated, fmt.Sprintf("${keycloak_user.user_%s.username}", userInGroup))
	}

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

%s

resource "keycloak_group_memberships" "group_members" {
	realm_id = "${keycloak_realm.realm.id}"
	group_id = "${keycloak_group.group.id}"

	members = %s
}
	`, realm, group, userResources.String(), arrayOfStringsForTerraformResource(usersInGroupInterpolated))
}

func testKeycloakGroupMemberships_userDoesNotExist(realm, group, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_group_memberships" "group_members" {
	realm_id = "${keycloak_realm.realm.id}"
	group_id = "${keycloak_group.group.id}"

	members = [
		"%s"
	]
}
	`, realm, group, username)
}
