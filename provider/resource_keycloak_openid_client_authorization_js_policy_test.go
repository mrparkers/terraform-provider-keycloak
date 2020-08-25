package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenidClientAuthorizationJSPolicy(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testResourceKeycloakOpenidClientAuthorizationJSPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceKeycloakOpenidClientAuthorizationJSPolicy_basic(realmName, clientId),
				Check:  testResourceKeycloakOpenidClientAuthorizationJSPolicyExists("keycloak_openid_client_js_policy.test"),
			},
		},
	})
}

func getResourceKeycloakOpenidClientAuthorizationJSPolicyFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientAuthorizationJSPolicy, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]
	resourceServerId := rs.Primary.Attributes["resource_server_id"]
	policyId := rs.Primary.ID

	policy, err := keycloakClient.GetOpenidClientAuthorizationJSPolicy(realm, resourceServerId, policyId)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client auth role policy config with alias %s: %s", resourceServerId, err)
	}

	return policy, nil
}

func testResourceKeycloakOpenidClientAuthorizationJSPolicyDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_js_policy" {
				continue
			}

			realm := rs.Primary.Attributes["realm_id"]
			resourceServerId := rs.Primary.Attributes["resource_server_id"]
			policyId := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			policy, _ := keycloakClient.GetOpenidClientAuthorizationJSPolicy(realm, resourceServerId, policyId)
			if policy != nil {
				return fmt.Errorf("policy config with id %s still exists", policyId)
			}
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationJSPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getResourceKeycloakOpenidClientAuthorizationJSPolicyFromState(s, resourceName)

		if err != nil {
			return err
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationJSPolicy_basic(realm, clientId string) string {
	return fmt.Sprintf(`
	resource keycloak_realm test {
		realm = "%s"
	}

	resource keycloak_openid_client test {
		client_id                = "%s"
		realm_id                 = "${keycloak_realm.test.id}"
		access_type              = "CONFIDENTIAL"
		service_accounts_enabled = true
		authorization {
			policy_enforcement_mode = "ENFORCING"
		}
	}

	resource keycloak_openid_client_js_policy test {
		resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
		realm_id = "${keycloak_realm.test.id}"
		name = "client_js_policy_test"
		logic = "POSITIVE"
		decision_strategy = "UNANIMOUS"
		code = "test"
		description = "description"
	}
	`, realm, clientId)
}
