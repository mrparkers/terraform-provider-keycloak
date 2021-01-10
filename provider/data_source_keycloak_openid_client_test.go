package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccKeycloakDataSourceOpenidClient_basic(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.keycloak_openid_client.test"
	resourceName := "keycloak_openid_client.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakOpenidClientConfig(clientId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "client_id", resourceName, "client_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_id", resourceName, "realm_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "access_type", resourceName, "access_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "standard_flow_enabled", resourceName, "standard_flow_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "implicit_flow_enabled", resourceName, "implicit_flow_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "direct_access_grants_enabled", resourceName, "direct_access_grants_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "service_account_user_id", resourceName, "service_account_user_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "service_accounts_enabled", resourceName, "service_accounts_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "resource_server_id", resourceName, "resource_server_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "full_scope_allowed", resourceName, "full_scope_allowed"),
				),
			},
		},
	})
}

func testAccKeycloakOpenidClientConfig(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "test" {
	name                     = "%s"
	client_id                = "%s"
	realm_id                 = data.keycloak_realm.realm.id
	description              = "a test openid client"
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
	full_scope_allowed       = false
}

data "keycloak_openid_client" "test" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = keycloak_openid_client.test.client_id

	depends_on = [
		keycloak_openid_client.test,
	]
}
`, testAccRealm.Realm, clientId, clientId)
}
