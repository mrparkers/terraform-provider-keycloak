package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakGroupMemberships_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroupMemberships_basic(realmName, groupName, username),
				Check:  testAccCheckKeycloakUserBelongsToGroup("keycloak_group_memberships.group_members", username),
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
