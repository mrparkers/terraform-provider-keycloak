package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakDataSourceRole_basic(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	realmRole := "terraform-role-" + acctest.RandString(10)
	clientRole := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakRole_basic(realm, client, realmRole, clientRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRoleExists("keycloak_role.realm_role"),
					testAccCheckKeycloakRoleExists("keycloak_role.client_role"),
					// realm role
					resource.TestCheckResourceAttrPair("keycloak_role.realm_role", "id", "data.keycloak_role.realm_role", "id"),
					resource.TestCheckResourceAttrPair("keycloak_role.realm_role", "realm_id", "data.keycloak_role.realm_role", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_role.realm_role", "name", "data.keycloak_role.realm_role", "name"),
					resource.TestCheckResourceAttrPair("keycloak_role.realm_role", "description", "data.keycloak_role.realm_role", "description"),
					testAccCheckDataKeycloakRole("data.keycloak_role.realm_role"),
					// client role
					resource.TestCheckResourceAttrPair("keycloak_role.client_role", "id", "data.keycloak_role.client_role", "id"),
					resource.TestCheckResourceAttrPair("keycloak_role.client_role", "realm_id", "data.keycloak_role.client_role", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_role.client_role", "client_id", "data.keycloak_role.client_role", "client_id"),
					resource.TestCheckResourceAttrPair("keycloak_role.client_role", "name", "data.keycloak_role.client_role", "name"),
					resource.TestCheckResourceAttrPair("keycloak_role.client_role", "description", "data.keycloak_role.client_role", "description"),
					testAccCheckDataKeycloakRole("data.keycloak_role.client_role"),
					// offline_access
					resource.TestCheckResourceAttrPair("keycloak_realm.realm", "realm", "data.keycloak_role.realm_offline_access", "realm_id"),
					resource.TestCheckResourceAttr("data.keycloak_role.realm_offline_access", "name", "offline_access"),
					testAccCheckDataKeycloakRole("data.keycloak_role.realm_offline_access"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakRole(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		name := rs.Primary.Attributes["name"]

		role, err := keycloakClient.GetRole(realmId, id)
		if err != nil {
			return err
		}

		if role.Name != name {
			return fmt.Errorf("expected role with ID %s to have name %s, but got %s", id, name, role.Name)
		}

		return nil
	}
}

func testDataSourceKeycloakRole_basic(realm, client, realmRole, clientRole string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_role" "realm_role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.client.id}"
}

data "keycloak_role" "realm_role" {
	realm_id = "${keycloak_realm.realm.id}"
	name     = "${keycloak_role.realm_role.name}"
}

data "keycloak_role" "client_role" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.client.id}"
	name      = "${keycloak_role.client_role.name}"
}

data "keycloak_role" "realm_offline_access" {
	realm_id = "${keycloak_realm.realm.id}"
	name     = "offline_access"
}
	`, realm, client, realmRole, clientRole)
}
