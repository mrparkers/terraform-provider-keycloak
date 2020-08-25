package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakAuthenticationSubFlow_basic(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	parentAuthFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)
	authFlowAlias := "terraform-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(realmName, parentAuthFlowAlias, authFlowAlias),
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
	var authenticationSubFlow = &keycloak.AuthenticationSubFlow{}

	realmName := "terraform-" + acctest.RandString(10)
	authParentFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)
	authFlowAlias := "terraform-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(realmName, authParentFlowAlias, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					testAccCheckKeycloakAuthenticationSubFlowFetch("keycloak_authentication_subflow.subflow", authenticationSubFlow),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteAuthenticationSubFlow(authenticationSubFlow.RealmId, authenticationSubFlow.ParentFlowAlias, authenticationSubFlow.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAuthenticationSubFlow_basic(realmName, authParentFlowAlias, authFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
			},
		},
	})
}

func TestAccKeycloakAuthenticationSubFlow_updateAuthenticationSubFlow(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	authParentFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)
	authFlowAliasBefore := "terraform-flow-before-" + acctest.RandString(10)
	authFlowAliasAfter := "terraform-flow-after-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(realmName, authParentFlowAlias, authFlowAliasBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "alias", authFlowAliasBefore),
				),
			},
			{
				Config: testKeycloakAuthenticationSubFlow_basic(realmName, authParentFlowAlias, authFlowAliasAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "alias", authFlowAliasAfter),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationSubFlow_updateAuthenticationSubFlowRequirement(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	authParentFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)
	authFlowAlias := "terraform-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationSubFlow_basic(realmName, authParentFlowAlias, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "requirement", "DISABLED"),
				),
			},
			{
				Config: testKeycloakAuthenticationSubFlow_basicWithRequirement(realmName, authParentFlowAlias, authFlowAlias, "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationSubFlowExists("keycloak_authentication_subflow.subflow"),
					resource.TestCheckResourceAttr("keycloak_authentication_subflow.subflow", "requirement", "REQUIRED"),
				),
			},
			{
				Config: testKeycloakAuthenticationSubFlow_basicWithRequirement(realmName, authParentFlowAlias, authFlowAlias, "DISABLED"),
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

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			authenticationSubFlow, _ := keycloakClient.GetAuthenticationSubFlow(realm, parentFlowAlias, id)
			if authenticationSubFlow != nil {
				return fmt.Errorf("authentication flow with id %s still exists", id)
			}
		}

		return nil
	}
}

func getAuthenticationSubFlowFromState(s *terraform.State, resourceName string) (*keycloak.AuthenticationSubFlow, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]

	authenticationSubFlow, err := keycloakClient.GetAuthenticationSubFlow(realm, parentFlowAlias, id)

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

func testKeycloakAuthenticationSubFlow_basic(realm, parentAlias, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "%s"
}

resource "keycloak_authentication_subflow" "subflow" {
	realm_id = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	alias    = "%s"
	provider_id = "basic-flow"
}
	`, realm, parentAlias, alias)
}

func testKeycloakAuthenticationSubFlow_basicWithRequirement(realm, parentAlias, alias, requirement string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "%s"
}

resource "keycloak_authentication_subflow" "subflow" {
	realm_id = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	alias    = "%s"
	provider_id = "basic-flow"
	requirement = "%s"
}
	`, realm, parentAlias, alias, requirement)
}
