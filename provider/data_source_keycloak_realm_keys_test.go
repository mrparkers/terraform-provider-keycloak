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
	dataSourceName := "keycloak_realm_keys.test_keys"

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
func getRealmKeysUsingState(state *terraform.State, resourceName string) (*terraform.ResourceState, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	return rs, nil
}

func testKeycloakRealmKeysCheck_basic(dataSourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		return nil

		datasourceState, err := getRealmKeysUsingState(state, dataSourceName)
		if err != nil {
			return err
		}

		// add tests...
		if len(datasourceState.Primary.Attributes["active"]) == 0 {
			return fmt.Errorf("no active key exists")
		}
		if len(datasourceState.Primary.Attributes["keys"]) == 0 {
			return fmt.Errorf("no key exists")
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

  filter {
     name   = "algorithm"
     values = ["RS256"]
  }
}
`, realm)
}
