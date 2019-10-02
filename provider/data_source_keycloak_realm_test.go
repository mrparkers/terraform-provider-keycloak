package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakDataSourceRealm_basic(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakRealm_basic(realm),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakRealm(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		realmId := rs.Primary.Attributes["realm_id"]
		name := rs.Primary.Attributes["display_name"]

		realm, err := keycloakClient.GetRealm(realmId)
		if err != nil {
			return err
		}

		if realm.DisplayName != name {
			return fmt.Errorf("expected realm with ID %s to have display_name %s, but got %s", realmId, name, realm.DisplayName)
		}

		return nil
	}
}

func testDataSourceKeycloakRealm_basic(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	display_name = "foo"
}`, realm)
}
