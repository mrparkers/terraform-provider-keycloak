package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKeycloakDataSourceOpenidClientAuthorizationPolicy_basic(t *testing.T) {
	realm := acctest.RandomWithPrefix("tf-acc-test")
	clientId := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.keycloak_openid_client_authorization_policy.test"
	resourceName := "keycloak_openid_client_authorization_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakOpenidClientAuthorizationPolicyConfig(realm, clientId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "resource_server_id", resourceName, "resource_server_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_id", resourceName, "realm_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "decision_strategy", resourceName, "decision_strategy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "owner", resourceName, "owner"),
					resource.TestCheckResourceAttrPair(dataSourceName, "logic", resourceName, "logic"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "policies", resourceName, "policies"),
					resource.TestCheckResourceAttrPair(dataSourceName, "resources", resourceName, "resources"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scopes", resourceName, "scopes"),
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
	name               = "Default Policy"
}
`, realm, clientId, clientId)
}
