package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakDataSourceGroup_basic(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	group := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakGroup_basic(realm, group),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group"),
					// realm role
					resource.TestCheckResourceAttrPair("keycloak_group.group", "id", "data.keycloak_group.group", "id"),
					resource.TestCheckResourceAttrPair("keycloak_group.group", "realm_id", "data.keycloak_group.group", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_group.group", "name", "data.keycloak_group.group", "name"),
					testAccCheckDataKeycloakGroup("data.keycloak_group.group"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakGroup(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		name := rs.Primary.Attributes["name"]

		group, err := keycloakClient.GetGroup(realmId, id)
		if err != nil {
			return err
		}

		if group.Name != name {
			return fmt.Errorf("expected group with ID %s to have name %s, but got %s", id, name, group.Name)
		}

		return nil
	}
}

func testDataSourceKeycloakGroup_basic(realm, group string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

data "keycloak_group" "group" {
	realm_id = "${keycloak_realm.realm.id}"
	name     = "${keycloak_group.group.name}"
}
	`, realm, group)
}
