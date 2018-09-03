package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakRealm_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmDisplayName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				Config: testKeycloakRealm_notEnabled(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmEnabled("keycloak_realm.realm", false),
			},
			{
				Config: testKeycloakRealm_basic(realmName, fmt.Sprintf("%s-changed", realmDisplayName)),
				Check:  testAccCheckKeycloakRealmDisplayName("keycloak_realm.realm", fmt.Sprintf("%s-changed", realmDisplayName)),
			},
		},
	})
}

func testAccCheckKeycloakRealmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakRealmEnabled(resourceName string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realm.Enabled != enabled {
			return fmt.Errorf("expected realm %s to have enabled set to %t, but was %t", realm.Realm, enabled, realm.Enabled)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmDisplayName(resourceName string, displayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realm.DisplayName != displayName {
			return fmt.Errorf("expected realm %s to have display name set to %s, but was %s", realm.Realm, displayName, realm.DisplayName)
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
				return fmt.Errorf("realm %s still exists", realmName)
			}
		}

		return nil
	}
}

func getRealmFromState(s *terraform.State, resourceName string) (*keycloak.Realm, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmName := rs.Primary.Attributes["realm"]

	realm, err := keycloakClient.GetRealm(realmName)
	if err != nil {
		return nil, fmt.Errorf("error getting realm %s: %s", realmName, err)
	}

	return realm, nil
}

func testKeycloakRealm_basic(realm, realmDisplayName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = true
	display_name = "%s"
}
	`, realm, realmDisplayName)
}

func testKeycloakRealm_notEnabled(realm, realmDisplayName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = false
	display_name = "%s"
}
	`, realm, realmDisplayName)
}
