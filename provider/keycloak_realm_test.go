package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"log"
	"testing"
)

func TestAccKeycloakRealm_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
		},
	})
}

func testAccCheckKeycloakRealmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Print("[DEBUG] testAccCheckKeycloakRealmExists")

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		realmName := rs.Primary.Attributes["realm"]

		_, err := keycloakClient.GetRealm(realmName)
		if err != nil {
			return fmt.Errorf("Error getting realm %s: %s", realmName, err)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm" {
				continue
			}

			realmName := rs.Primary.ID
			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			realm, _ := keycloakClient.GetRealm(realmName)
			if realm != nil {
				return fmt.Errorf("Realm %s still exists", realmName)
			}
		}

		return nil
	}
}

func testKeycloakRealm_basic(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
	`, realm)
}
