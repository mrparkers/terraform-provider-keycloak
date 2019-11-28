package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func TestAccKeycloakDataSourceOpenidClientServiceAccountUser_basic(t *testing.T) {
	realm := acctest.RandomWithPrefix("tf-acc-test")
	clientId := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.keycloak_openid_client_service_account_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakOpenidClientServiceAccountUserConfig(realm, clientId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "realm_id", realm),
					resource.TestMatchResourceAttr(dataSourceName, "client_id", regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestCheckResourceAttr(dataSourceName, "username", "service-account-"+clientId),
					resource.TestCheckResourceAttr(dataSourceName, "email", "service-account-"+clientId+"@placeholder.org"),
				),
			},
		},
	})
}

func testAccKeycloakOpenidClientServiceAccountUserConfig(realm, clientId string) string {
	return fmt.Sprintf(`
resource keycloak_realm test {
  realm                = "%s"
  enabled              = true
  display_name         = "foo"
  account_theme        = "base"
  access_code_lifespan = "30m"
}

resource keycloak_openid_client test {
  name                  	= "%s"
  client_id 					= "%s"
  realm_id              	= "${keycloak_realm.test.id}"
  description           	= "a test openid client"
  standard_flow_enabled    = true
  access_type              = "CONFIDENTIAL"
  service_accounts_enabled = true
  client_secret            = "secret"
  valid_redirect_uris      = [
   	"http://localhost:5555/callback",
  ]
  authorization {
  		policy_enforcement_mode = "ENFORCING"
  }
  web_origins              = [
		"http://localhost"
  ]
}

data keycloak_openid_client_service_account_user test {
  client_id = "${keycloak_openid_client.test.id}"
  realm_id  = "${keycloak_realm.test.id}"
}
`, realm, clientId, clientId)
}
