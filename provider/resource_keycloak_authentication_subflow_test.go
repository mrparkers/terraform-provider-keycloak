package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakAuthenticationSubFlow_basic(t *testing.T) {
	t.Parallel()

	parentAuthFlowAlias := acctest.RandomWithPrefix("tf-acc")
	authFlowAlias := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(parentAuthFlowAlias, authFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
			},
			{
				ResourceName:      "keycloak_authentication_subflow.subflow",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSubFlowImportId("keycloak_authentication_subflow.subflow"),
			},
		},
	})
}

func TestAccKeycloakAuthenticationSubFlow_createAfterManualDestroy(t *testing.T) {
	t.Parallel()

	var authenticationSubFlow = &keycloak.AuthenticationSubFlow{}

	authParentFlowAlias := acctest.RandomWithPrefix("tf-acc")
	authFlowAlias := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(authParentFlowAlias, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					testAccCheckKeycloakAuthenticationSubFlowFetch("keycloak_authentication_subflow.subflow", authenticationSubFlow),
				),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteAuthenticationSubFlow(testCtx, authenticationSubFlow.RealmId, authenticationSubFlow.ParentFlowAlias, authenticationSubFlow.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAuthenticationSubFlow_basic(authParentFlowAlias, authFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
			},
		},
	})
}

func TestAccKeycloakAuthenticationSubFlow_updateAuthenticationSubFlow(t *testing.T) {
	t.Parallel()

	authParentFlowAlias := acctest.RandomWithPrefix("tf-acc")
	authFlowAliasBefore := acctest.RandomWithPrefix("tf-acc")
	authFlowAliasAfter := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(authParentFlowAlias, authFlowAliasBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "alias", authFlowAliasBefore),
				),
			},
			{
				Config: testKeycloakAuthenticationSubFlow_basic(authParentFlowAlias, authFlowAliasAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "alias", authFlowAliasAfter),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationSubFlow_updateAuthenticationSubFlowRequirement(t *testing.T) {
	t.Parallel()

	authParentFlowAlias := acctest.RandomWithPrefix("tf-acc")
	authFlowAlias := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(authParentFlowAlias, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "requirement", "DISABLED"),
				),
			},
			{
				Config: testKeycloakAuthenticationSubFlow_basicWithRequirement(authParentFlowAlias, authFlowAlias, "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "requirement", "REQUIRED"),
				),
			},
			{
				Config: testKeycloakAuthenticationSubFlow_basicWithRequirement(authParentFlowAlias, authFlowAlias, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "requirement", "DISABLED"),
				),
			},
		},
	})
}

func testAccCheckKeycloakAuthenticationSubFlowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getAuthenticationSubFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationSubFlowExistsWithDescription(resourceName, expectedDescr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		subflow, err := getAuthenticationSubFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		if subflow.Description != expectedDescr {
			return fmt.Errorf("expected authentication subflow's description to be %s, but was %s", expectedDescr, subflow.Description)
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationSubFlowFetch(resourceName string, authenticationSubFlow *keycloak.AuthenticationSubFlow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedAuthenticationSubFlow, err := getAuthenticationSubFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		authenticationSubFlow.Id = fetchedAuthenticationSubFlow.Id
		authenticationSubFlow.ParentFlowAlias = fetchedAuthenticationSubFlow.ParentFlowAlias
		authenticationSubFlow.RealmId = fetchedAuthenticationSubFlow.RealmId
		authenticationSubFlow.Alias = fetchedAuthenticationSubFlow.Alias

		return nil
	}
}

func testAccCheckKeycloakAuthenticationSubFlowDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_authentication_subflow" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]
			parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]

			authenticationSubFlow, _ := keycloakClient.GetAuthenticationSubFlow(testCtx, realm, parentFlowAlias, id)
			if authenticationSubFlow != nil {
				return fmt.Errorf("authentication flow with id %s still exists", id)
			}
		}

		return nil
	}
}
func testAccCheckKeycloakAuthenticationSubFlowNotDestroyed() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_authentication_subflow" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]
			parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]

			client, _ := keycloakClient.GetAuthenticationSubFlow(testCtx, realm, parentFlowAlias, id)
			if client == nil {
				return fmt.Errorf("authentication flow %s does not exists", id)
			}
		}

		return nil
	}
}

func TestAccKeycloakAuthenticationSubFlowImport(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowNotDestroyed(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakAuthenticationSubFlow_import("non-existing-flow", "non-existing-auth-flow", "new-descr"),
				ExpectError: regexp.MustCompile("no parent flow with alias non-existing-flow exists"),
			},
			{
				Config:      testKeycloakAuthenticationSubFlow_import("browser", "non-existing-auth-flow", "new-descr"),
				ExpectError: regexp.MustCompile("no authentication execution under parent flow alias browser with alias non-existing-auth-flow found"),
			},
			{
				// use existing browser flow and change the description to be "descr" (instead of "basic-flow")
				Config: testKeycloakAuthenticationSubFlow_import("browser", "forms", "new-descr"),
				Check:  testAccCheckKeycloakAuthenticationSubFlowExistsWithDescription("keycloak_authentication_subflow.imported-subflow", "new-descr"),
			},
		},
	})
}

func getAuthenticationSubFlowFromState(s *terraform.State, resourceName string) (*keycloak.AuthenticationSubFlow, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]

	authenticationSubFlow, err := keycloakClient.GetAuthenticationSubFlow(testCtx, realm, parentFlowAlias, id)

	if err != nil {
		return nil, fmt.Errorf("error getting authentication subflow with id %s: %s", id, err)
	}

	return authenticationSubFlow, nil
}

func getSubFlowImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]
		realmId := rs.Primary.Attributes["realm_id"]

		return fmt.Sprintf("%s/%s/%s", realmId, parentFlowAlias, id), nil
	}
}

func testKeycloakAuthenticationSubFlow_basic(parentAlias, alias string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = data.keycloak_realm.realm.id
	alias    = "%s"
}

resource "keycloak_authentication_subflow" "subflow" {
	realm_id          = data.keycloak_realm.realm.id
	parent_flow_alias = keycloak_authentication_flow.flow.alias

	alias       = "%s"
	provider_id = "basic-flow"
}
	`, testAccRealm.Realm, parentAlias, alias)
}

func testKeycloakAuthenticationSubFlow_basicWithRequirement(parentAlias, alias, requirement string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = data.keycloak_realm.realm.id
	alias    = "%s"
}

resource "keycloak_authentication_subflow" "subflow" {
	realm_id          = data.keycloak_realm.realm.id
	parent_flow_alias = keycloak_authentication_flow.flow.alias

	alias       = "%s"
	provider_id = "basic-flow"
	requirement = "%s"
}
	`, testAccRealm.Realm, parentAlias, alias, requirement)
}

func testKeycloakAuthenticationSubFlow_import(parentAlias, subflowAlias, description string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_subflow" "imported-subflow" {
	realm_id          = data.keycloak_realm.realm.id
	parent_flow_alias = "%s"

	alias       = "%s"
	description = "%s"
	import = true
}
	`, testAccRealm.Realm, parentAlias, subflowAlias, description)
}
