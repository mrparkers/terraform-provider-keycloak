package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakDataSourceAuthenticationExecution(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	parentFlowAlias := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakAuthenticationExecution(realm, parentFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "id", "data.keycloak_authentication_execution.execution", "id"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "realm_id", "data.keycloak_authentication_execution.execution", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "parent_flow_alias", "data.keycloak_authentication_execution.execution", "parent_flow_alias"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "authenticator", "data.keycloak_authentication_execution.execution", "provider_id"),
					testAccCheckDataKeycloakAuthenticationExecution("data.keycloak_authentication_execution.execution"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakAuthenticationExecution(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		id := rs.Primary.ID
		realmID := rs.Primary.Attributes["realm_id"]
		parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]
		providerID := rs.Primary.Attributes["provider_id"]

		authenticationExecutionInfo, err := keycloakClient.GetAuthenticationExecutionInfoFromProviderId(realmID, parentFlowAlias, providerID)
		if err != nil {
			return err
		}

		if authenticationExecutionInfo.Id != id {
			return fmt.Errorf("expected authenticationExecutionInfo with ID %s but got %s", id, authenticationExecutionInfo.Id)
		}

		return nil
	}
}

func testDataSourceKeycloakAuthenticationExecution(realm, parentFlowAlias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = keycloak_realm.realm.id
	alias    = "%s"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = keycloak_realm.realm.id
	parent_flow_alias = keycloak_authentication_flow.flow.alias
	authenticator     = "identity-provider-redirector"
	requirement       = "REQUIRED"
}

data "keycloak_authentication_execution" "execution" {
	realm_id 			= keycloak_realm.realm.id
	parent_flow_alias   = keycloak_authentication_flow.flow.alias
	provider_id     	= "identity-provider-redirector"
}
	`, realm, parentFlowAlias)
}
