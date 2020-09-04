package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakAuthenticationFlow_basic(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	authFlowAlias := "terraform-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationFlow_basic(realmName, authFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
			},
			{
				ResourceName:        "keycloak_authentication_flow.flow",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakAuthenticationFlow_createAfterManualDestroy(t *testing.T) {
	var authenticationFlow = &keycloak.AuthenticationFlow{}

	realmName := "terraform-" + acctest.RandString(10)
	authFlowAlias := "terraform-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationFlow_basic(realmName, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					testAccCheckKeycloakAuthenticationFlowFetch("keycloak_authentication_flow.flow", authenticationFlow),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteAuthenticationFlow(authenticationFlow.RealmId, authenticationFlow.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAuthenticationFlow_basic(realmName, authFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
			},
		},
	})
}

func TestAccKeycloakAuthenticationFlow_updateAuthenticationFlow(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)

	authFlowAliasBefore := "terraform-flow-before-" + acctest.RandString(10)
	authFlowAliasAfter := "terraform-flow-after-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationFlow_basic(realmName, authFlowAliasBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					resource.TestCheckResourceAttr("keycloak_authentication_flow.flow", "alias", authFlowAliasBefore),
				),
			},
			{
				Config: testKeycloakAuthenticationFlow_basic(realmName, authFlowAliasAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					resource.TestCheckResourceAttr("keycloak_authentication_flow.flow", "alias", authFlowAliasAfter),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationFlow_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)

	authFlowAlias := "terraform-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationFlow_updateRealmBefore(realmOne, realmTwo, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					testAccCheckKeycloakAuthenticationFlowBelongsToRealm("keycloak_authentication_flow.flow", realmOne),
				),
			},
			{
				Config: testKeycloakAuthenticationFlow_updateRealmAfter(realmOne, realmTwo, authFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					testAccCheckKeycloakAuthenticationFlowBelongsToRealm("keycloak_authentication_flow.flow", realmTwo),
				),
			},
		},
	})
}

func testAccCheckKeycloakAuthenticationFlowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getAuthenticationFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationFlowFetch(resourceName string, authenticationFlow *keycloak.AuthenticationFlow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedAuthenticationFlow, err := getAuthenticationFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		authenticationFlow.Id = fetchedAuthenticationFlow.Id
		authenticationFlow.RealmId = fetchedAuthenticationFlow.RealmId

		return nil
	}
}

func testAccCheckKeycloakAuthenticationFlowBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		authenticationFlow, err := getAuthenticationFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		if authenticationFlow.RealmId != realm {
			return fmt.Errorf("expected authentication flow with id %s to have realm_id of %s, but got %s", authenticationFlow.Id, realm, authenticationFlow.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationFlowDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_authentication_flow" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			authenticationFlow, _ := keycloakClient.GetAuthenticationFlow(realm, id)
			if authenticationFlow != nil {
				return fmt.Errorf("authentication flow with id %s still exists", id)
			}
		}

		return nil
	}
}

func getAuthenticationFlowFromState(s *terraform.State, resourceName string) (*keycloak.AuthenticationFlow, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	authenticationFlow, err := keycloakClient.GetAuthenticationFlow(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting authentication flow with id %s: %s", id, err)
	}

	return authenticationFlow, nil
}

func testKeycloakAuthenticationFlow_basic(realm, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "%s"
}
	`, realm, alias)
}

func testKeycloakAuthenticationFlow_updateRealmBefore(realmOne, realmTwo, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	alias    = "%s"
	realm_id = "${keycloak_realm.realm_1.id}"
}
	`, realmOne, realmTwo, alias)
}

func testKeycloakAuthenticationFlow_updateRealmAfter(realmOne, realmTwo, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	alias    = "%s"
	realm_id = "${keycloak_realm.realm_2.id}"
}
	`, realmOne, realmTwo, alias)
}
