package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakGroupRoles_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmRoleName := "terraform-role-" + acctest.RandString(10)
	openIdClientName := "terraform-openid-client-" + acctest.RandString(10)
	openIdRoleName := "terraform-role-" + acctest.RandString(10)
	samlClientName := "terraform-saml-client-" + acctest.RandString(10)
	samlRoleName := "terraform-role-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroupRoles_basic(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName),
				Check:  testAccCheckKeycloakGroupHasRoles("keycloak_group_roles.group_roles"),
			},
			{
				ResourceName:      "keycloak_group_roles.group_roles",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// check destroy
			{
				Config: testKeycloakGroupRoles_noGroupRoles(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName),
				Check:  testAccCheckKeycloakGroupHasNoRoles("keycloak_group.group"),
			},
		},
	})
}

func TestAccKeycloakGroupRoles_update(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmRoleName := "terraform-role-" + acctest.RandString(10)
	openIdClientName := "terraform-openid-client-" + acctest.RandString(10)
	openIdRoleName := "terraform-role-" + acctest.RandString(10)
	samlClientName := "terraform-saml-client-" + acctest.RandString(10)
	samlRoleName := "terraform-role-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// initial setup, resource is defined but no roles are specified
			{
				Config: testKeycloakGroupRoles_update(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName, []string{}),
				Check:  testAccCheckKeycloakGroupHasRoles("keycloak_group_roles.group_roles"),
			},
			// add all roles
			{
				Config: testKeycloakGroupRoles_update(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName, []string{
					"${keycloak_role.realm_role.id}",
					"${keycloak_role.openid_client_role.id}",
					"${keycloak_role.saml_client_role.id}",
					"${data.keycloak_role.offline_access.id}",
				}),
				Check: testAccCheckKeycloakGroupHasRoles("keycloak_group_roles.group_roles"),
			},
			// remove two
			{
				Config: testKeycloakGroupRoles_update(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName, []string{
					"${keycloak_role.openid_client_role.id}",
					"${data.keycloak_role.offline_access.id}",
				}),
				Check: testAccCheckKeycloakGroupHasRoles("keycloak_group_roles.group_roles"),
			},
			// add back and remove others
			{
				Config: testKeycloakGroupRoles_update(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName, []string{
					"${keycloak_role.realm_role.id}",
					"${keycloak_role.saml_client_role.id}",
				}),
				Check: testAccCheckKeycloakGroupHasRoles("keycloak_group_roles.group_roles"),
			},
			// remove all
			{
				Config: testKeycloakGroupRoles_update(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName, []string{}),
				Check:  testAccCheckKeycloakGroupHasRoles("keycloak_group_roles.group_roles"),
			},
		},
	})
}

func flattenGroupRoles(keycloakClient *keycloak.KeycloakClient, group *keycloak.Group) ([]string, error) {
	var roles []string

	for _, realmRole := range group.RealmRoles {
		roles = append(roles, realmRole)
	}

	for clientId, clientRoles := range group.ClientRoles {
		client, err := keycloakClient.GetGenericClientByClientId(group.RealmId, clientId)
		if err != nil {
			return nil, err
		}

		for _, clientRole := range clientRoles {
			roles = append(roles, fmt.Sprintf("%s/%s", client.Id, clientRole))
		}
	}

	return roles, nil
}

func testAccCheckKeycloakGroupHasRoles(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realm := rs.Primary.Attributes["realm_id"]
		groupId := rs.Primary.Attributes["group_id"]

		var roles []*keycloak.Role
		for k, v := range rs.Primary.Attributes {
			if match, _ := regexp.MatchString("role_ids\\.[^#]+", k); !match {
				continue
			}

			role, err := keycloakClient.GetRole(realm, v)
			if err != nil {
				return err
			}

			roles = append(roles, role)
		}

		group, err := keycloakClient.GetGroup(realm, groupId)
		if err != nil {
			return err
		}

		groupRoles, err := flattenGroupRoles(keycloakClient, group)
		if err != nil {
			return err
		}

		if len(groupRoles) != len(roles) {
			return fmt.Errorf("expected number of group roles to be %d, got %d", len(roles), len(groupRoles))
		}

		for _, role := range roles {
			var expectedRoleString string
			if role.ClientRole {
				expectedRoleString = fmt.Sprintf("%s/%s", role.ClientId, role.Name)
			} else {
				expectedRoleString = role.Name
			}

			found := false

			for _, groupRole := range groupRoles {
				if groupRole == expectedRoleString {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("expected to find role %s assigned to group %s", expectedRoleString, group.Name)
			}
		}

		return nil
	}
}

func testAccCheckKeycloakGroupHasNoRoles(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realm := rs.Primary.Attributes["realm_id"]
		id := rs.Primary.ID

		group, err := keycloakClient.GetGroup(realm, id)
		if err != nil {
			return err
		}

		if len(group.RealmRoles) != 0 || len(group.ClientRoles) != 0 {
			return fmt.Errorf("expected group %s to have no roles", group.Name)
		}

		return nil
	}
}

func testKeycloakGroupRoles_basic(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "openid_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "saml_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

data "keycloak_role" "offline_access" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "offline_access"
}

resource "keycloak_group" "group" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}

resource "keycloak_group_roles" "group_roles" {
	realm_id = "${keycloak_realm.realm.id}"
	group_id = "${keycloak_group.group.id}"

	role_ids = [
		"${keycloak_role.realm_role.id}",
		"${keycloak_role.openid_client_role.id}",
		"${keycloak_role.saml_client_role.id}",
		"${data.keycloak_role.offline_access.id}",
	]
}
	`, realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName)
}

func testKeycloakGroupRoles_noGroupRoles(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "openid_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "saml_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

data "keycloak_role" "offline_access" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "offline_access"
}

resource "keycloak_group" "group" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}
	`, realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName)
}

func testKeycloakGroupRoles_update(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName string, roleIds []string) string {
	tfRoleIds := fmt.Sprintf("role_ids = %s", arrayOfStringsForTerraformResource(roleIds))

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "openid_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "saml_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

data "keycloak_role" "offline_access" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "offline_access"
}

resource "keycloak_group" "group" {
	realm_id = "${keycloak_realm.realm.id}"
	name = "%s"
}

resource "keycloak_group_roles" "group_roles" {
	realm_id = "${keycloak_realm.realm.id}"
	group_id = "${keycloak_group.group.id}"

	%s
}
	`, realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, groupName, tfRoleIds)
}
