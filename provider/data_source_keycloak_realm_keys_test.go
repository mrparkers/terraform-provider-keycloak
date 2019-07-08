package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKeycloakDataSourceRealmKeys_basic(t *testing.T) {
	realm := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.keycloak_realm_keys.test_keys"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakRealmKeysConfig(realm),
				Check:  testKeycloakRealmKeysCheck_basic(dataSourceName),
			},
		},
	})
}

func TestAccKeycloakDataSourceRealmKeys_filterByAlgorithms(t *testing.T) {
	realm := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.keycloak_realm_keys.test_keys"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakRealmKeysConfig_filterByAlgorithms(realm),
				Check:  testKeycloakRealmKeysCheck_filterByAlgorithms(dataSourceName),
			},
		},
	})
}

func getRealmKeysUsingState(state *terraform.State, resourceName string) (*terraform.ResourceState, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	return rs, nil
}

func testKeycloakRealmKeysCheck_basic(dataSourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		datasourceState, err := getRealmKeysUsingState(state, dataSourceName)
		if err != nil {
			return err
		}

		if len(datasourceState.Primary.Attributes["keys.#"]) == 0 {
			return fmt.Errorf("no key exists")
		}

		return nil
	}
}

func testKeycloakRealmKeysCheck_filterByAlgorithms(dataSourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		datasourceState, err := getRealmKeysUsingState(state, dataSourceName)
		if err != nil {
			return err
		}

		if len(datasourceState.Primary.Attributes["keys.#"]) == 0 {
			return fmt.Errorf("no key exists")
		}

		algorithm := datasourceState.Primary.Attributes["keys.0.algorithm"]
		if algorithm != "AES" && algorithm != "RS256" {
			return fmt.Errorf("filtering by algorithm returned '%s', but this value was not part of the valid values", algorithm)
		}

		return nil
	}
}

func testAccKeycloakRealmKeysConfig(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "test" {
  realm                = "%s"
  enabled              = true
  display_name         = "test"
}

data "keycloak_realm_keys" "test_keys" {
  realm_id  = keycloak_realm.test.id
}
`, realm)
}

func testAccKeycloakRealmKeysConfig_filterByAlgorithms(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "test" {
  realm                = "%s"
  enabled              = true
  display_name         = "test"
}

data "keycloak_realm_keys" "test_keys" {
  realm_id  = keycloak_realm.test.id
  algorithms = ["RS256", "AES"]
}
`, realm)
}
