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

func TestAccKeycloakRequiredAction_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	requiredActionAlias := "CONFIGURE_TOTP"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRequiredAction_basic(realmName, requiredActionAlias, 37),
				Check:  testAccCheckKeycloakRequiresActionExistsWithCorrectPriority(realmName, requiredActionAlias, 37),
			},
		},
	})
}

func TestAccKeycloakRequiredAction_invalidAlias(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	randomReqActionAlias := "randomRequiredAction-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRequiredAction_basic(realmName, randomReqActionAlias, 37),
				ExpectError: regexp.MustCompile("errors during apply: validation error: required action .+ does not exist on the server, installed providers: .+"),
			},
		},
	})
}

func testKeycloakRequiredAction_basic(realm, requiredActionAlias string, priority int) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_required_action" "custom-terms-and-conditions" {
	realm_name		= "${keycloak_realm.realm.realm}"
	alias			= "%s"
	default_action 	= true
	enabled			= true
	name			= "My required Action"
	priority		= %d

	depends_on = ["keycloak_realm.realm"]
}
	`, realm, requiredActionAlias, priority)
}

func testAccCheckKeycloakRequiresActionExistsWithCorrectPriority(realm, requiredActionAlias string, priority int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)
		action, err := keycloakClient.GetRequiredAction(realm, requiredActionAlias)
		if err != nil {
			return fmt.Errorf("required action not found: %s", requiredActionAlias)
		}

		if action.Priority != priority {
			return fmt.Errorf("expected required action to have priority %d, but got %d", priority, action.Priority)
		}

		return nil
	}
}
