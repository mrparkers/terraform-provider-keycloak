package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakAuthenticationExecutionConfig_basic(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	var config1, config2 keycloak.AuthenticationExecutionConfig

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakAuthenticationExecutionConfig(realmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionConfigExists("keycloak_authentication_execution_config.config", &config1),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "realm_id", realmName),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "alias", "some-config-alias"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.%", "1"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.defaultProvider", "some-config-default-idp"),
				),
			},
			{
				Config: testAccKeycloakAuthenticationExecutionConfigUpdatedConfig(realmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionConfigExists("keycloak_authentication_execution_config.config", &config2),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "realm_id", realmName),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "alias", "some-config-alias"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.%", "1"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.defaultProvider", "some-updated-config-default-idp"),
					testAccCheckKeycloakAuthenticationExecutionConfigForceNew(&config1, &config2, false),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationExecutionConfig_updateForcesNew(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	var config1, config2 keycloak.AuthenticationExecutionConfig

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakAuthenticationExecutionConfig(realmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionConfigExists("keycloak_authentication_execution_config.config", &config1),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "realm_id", realmName),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "alias", "some-config-alias"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.%", "1"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.defaultProvider", "some-config-default-idp"),
				),
			},
			{
				Config: testAccKeycloakAuthenticationExecutionConfigUpdatedAlias(realmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionConfigExists("keycloak_authentication_execution_config.config", &config2),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "realm_id", realmName),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "alias", "some-updated-config-alias"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.%", "1"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution_config.config", "config.defaultProvider", "some-config-default-idp"),
					testAccCheckKeycloakAuthenticationExecutionConfigForceNew(&config1, &config2, true),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationExecutionConfig_import(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakAuthenticationExecutionConfig(realmName),
			},
			{
				ResourceName:      "keycloak_authentication_execution_config.config",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getExecutionConfigImportId("keycloak_authentication_execution_config.config"),
			},
		},
	})
}

func getExecutionConfigImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}

		realmId := rs.Primary.Attributes["realm_id"]
		executionId := rs.Primary.Attributes["execution_id"]
		id := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", realmId, executionId, id), nil
	}
}

func testAccCheckKeycloakAuthenticationExecutionConfigExists(resourceName string, config *keycloak.AuthenticationExecutionConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		config.RealmId = rs.Primary.Attributes["realm_id"]
		config.ExecutionId = rs.Primary.Attributes["execution_id"]
		config.Id = rs.Primary.ID

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)
		if err := keycloakClient.GetAuthenticationExecutionConfig(config); err != nil {
			return fmt.Errorf("error fetching authentication execution config: %v", err)
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationExecutionConfigDestroy(s *terraform.State) error {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "keycloak_authentication_execution_config" {
			continue
		}

		config := &keycloak.AuthenticationExecutionConfig{
			RealmId: rs.Primary.Attributes["realm_id"],
			Id:      rs.Primary.ID,
		}
		if err := keycloakClient.GetAuthenticationExecutionConfig(config); err == nil {
			return fmt.Errorf("authentication execution config still exists")
		} else if !keycloak.ErrorIs404(err) {
			return fmt.Errorf("could not fetch authentication execution config: %v", err)
		}
	}

	return nil
}

func testAccCheckKeycloakAuthenticationExecutionConfigForceNew(old, new *keycloak.AuthenticationExecutionConfig, wantNew bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if wantNew {
			if old.Id == new.Id {
				return fmt.Errorf("expecting authentication execution config ID to differ, got %+v and %+v", old, new)
			}
		} else {
			if old.Id != new.Id {
				return fmt.Errorf("expecting authentication execution config ID to be equal, got %+v and %+v", old, new)
			}
		}
		return nil
	}
}

func testAccKeycloakAuthenticationExecutionConfig(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "some-flow-alias"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	authenticator     = "identity-provider-redirector"
}

resource "keycloak_authentication_execution_config" "config" {
	realm_id     = "${keycloak_realm.realm.id}"
	execution_id = "${keycloak_authentication_execution.execution.id}"
	alias        = "some-config-alias"
	config = {
		defaultProvider = "some-config-default-idp"
	}
}`, realm)
}

func testAccKeycloakAuthenticationExecutionConfigUpdatedConfig(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "some-flow-alias"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	authenticator     = "identity-provider-redirector"
}

resource "keycloak_authentication_execution_config" "config" {
	realm_id     = "${keycloak_realm.realm.id}"
	execution_id = "${keycloak_authentication_execution.execution.id}"
	alias        = "some-config-alias"
	config = {
		defaultProvider = "some-updated-config-default-idp"
	}
}`, realm)
}

func testAccKeycloakAuthenticationExecutionConfigUpdatedAlias(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "some-flow-alias"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	authenticator     = "identity-provider-redirector"
}

resource "keycloak_authentication_execution_config" "config" {
	realm_id     = "${keycloak_realm.realm.id}"
	execution_id = "${keycloak_authentication_execution.execution.id}"
	alias        = "some-updated-config-alias"
	config = {
		defaultProvider = "some-config-default-idp"
	}
}`, realm)
}
