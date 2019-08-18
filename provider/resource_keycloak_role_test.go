package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakRole_basicRealm(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				ResourceName:        "keycloak_role.role",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakRole_basicClient(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicClient(realmName, clientId, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				ResourceName:        "keycloak_role.role",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakRole_basicRealmUpdate(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)
	descriptionOne := acctest.RandString(50)
	descriptionTwo := acctest.RandString(50)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicRealmWithDescription(realmName, roleName, descriptionOne),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicRealmWithDescription(realmName, roleName, descriptionTwo),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
		},
	})
}

func TestAccKeycloakRole_basicClientUpdate(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)
	descriptionOne := acctest.RandString(50)
	descriptionTwo := acctest.RandString(50)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicClientWithDescription(realmName, clientId, roleName, descriptionOne),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicClientWithDescription(realmName, clientId, roleName, descriptionTwo),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicClient(realmName, clientId, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
		},
	})
}

func testAccCheckKeycloakRoleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakRoleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_role" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			role, _ := keycloakClient.GetRole(realm, id)
			if role != nil {
				return fmt.Errorf("role with id %s still exists", id)
			}
		}

		return nil
	}
}

func getRoleFromState(s *terraform.State, resourceName string) (*keycloak.Role, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	role, err := keycloakClient.GetRole(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting role with id %s: %s", id, err)
	}

	return role, nil
}

func testKeycloakRole_basicRealm(realm, role string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
	`, realm, role)
}

func testKeycloakRole_basicRealmWithDescription(realm, role, description string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name        = "%s"
	description = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
}
	`, realm, role, description)
}

func testKeycloakRole_basicClient(realm, clientId, role string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_role" "role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.client.id}"
}
	`, realm, clientId, role)
}

func testKeycloakRole_basicClientWithDescription(realm, clientId, role, description string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_role" "role" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "${keycloak_openid_client.client.id}"
	description = "%s"
}
	`, realm, clientId, role, description)
}
