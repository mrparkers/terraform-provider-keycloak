package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func TestAccKeycloakDataSourceOpenidClientAuthorizationPolicy_basic(t *testing.T) {
	realm := acctest.RandomWithPrefix("tf-acc-test")
	clientId := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.keycloak_openid_client_authorization_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakOpenidClientAuthorizationPolicyConfig(realm, clientId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "resource_server_id", regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestCheckResourceAttr(dataSourceName, "realm_id", realm),
					resource.TestCheckResourceAttr(dataSourceName, "name", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "decision_strategy", "UNANIMOUS"),
					resource.TestCheckResourceAttr(dataSourceName, "logic", "POSITIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "resource"),
				),
			},
		},
	})
}

func testAccKeycloakOpenidClientAuthorizationPolicyConfig(realm, clientId string) string {
	return fmt.Sprintf(`
resource keycloak_realm test {
  realm                = "%s"
  enabled              = true
  display_name         = "foo"
  account_theme        = "base"
  access_code_lifespan = "30m"
}

resource keycloak_openid_client test {
  client_id                      = "%s"
  name                           = "%s"
  realm_id                       = "${keycloak_realm.test.id}"
  description                    = "a test openid client"
  standard_flow_enabled          = true
  service_accounts_enabled       = true
  access_type                    = "CONFIDENTIAL"
  client_secret                  = "secret"
  valid_redirect_uris            = [
    "http://localhost:5555/callback",
  ]
  authorization {
  	 policy_enforcement_mode = "ENFORCING"
  }
}

data keycloak_openid_client_authorization_policy test {
	resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
	realm_id           = "${keycloak_realm.test.id}"
	name               = "default"
}
`, realm, clientId, clientId)
}
