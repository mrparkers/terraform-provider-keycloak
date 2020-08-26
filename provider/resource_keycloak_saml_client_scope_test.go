package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakSamlClientScope_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
			{
				ResourceName:        "keycloak_saml_client_scope.client_scope",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakSamlClientScope_createAfterManualDestroy(t *testing.T) {
	var clientScope = &keycloak.SamlClientScope{}

	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
					testAccCheckKeycloakSamlClientScopeFetch("keycloak_saml_client_scope.client_scope", clientScope),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteSamlClientScope(clientScope.RealmId, clientScope.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
		},
	})
}

func TestAccKeycloakSamlClientScope_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClientScope_updateRealmBefore(realmOne, realmTwo, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
					testAccCheckKeycloakSamlClientScopeBelongsToRealm("keycloak_saml_client_scope.client_scope", realmOne),
				),
			},
			{
				Config: testKeycloakSamlClientScope_updateRealmAfter(realmOne, realmTwo, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
					testAccCheckKeycloakSamlClientScopeBelongsToRealm("keycloak_saml_client_scope.client_scope", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakSamlClientScope_consentScreenText(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
			{
				Config: testKeycloakSamlClientScope_withConsentText(realmName, clientScopeName, acctest.RandString(10)),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
			{
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
		},
	})
}

func TestAccKeycloakSamlClientScope_guiOrder(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)
	guiOrder := acctest.RandIntRange(0, 1000)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
			{
				Config: testKeycloakSamlClientScope_withGuiOrder(realmName, clientScopeName, guiOrder),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
					testAccCheckKeycloakSamlClientScopeExistsWithCorrectGuiOrder("keycloak_saml_client_scope.client_scope", guiOrder),
				),
			},
			{
				Config: testKeycloakSamlClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol("keycloak_saml_client_scope.client_scope"),
			},
		},
	})
}

func testAccCheckKeycloakSamlClientScopeExistsWithCorrectProtocol(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientScope, err := getSamlClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		if clientScope.Protocol != "saml" {
			return fmt.Errorf("expected saml client scope to have saml protocol, but got %s", clientScope.Protocol)
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientScopeExistsWithCorrectGuiOrder(resourceName string, guiOrder int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientScope, err := getSamlClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		if clientScope.Attributes.GuiOrder != strconv.Itoa(guiOrder) {
			return fmt.Errorf("expected saml client guiOrder to have %d, but got %s", guiOrder, clientScope.Attributes.GuiOrder)
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientScopeFetch(resourceName string, clientScope *keycloak.SamlClientScope) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClientScope, err := getSamlClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		clientScope.Id = fetchedClientScope.Id
		clientScope.RealmId = fetchedClientScope.RealmId

		return nil
	}
}

func testAccCheckKeycloakSamlClientScopeBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientScope, err := getSamlClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		if clientScope.RealmId != realm {
			return fmt.Errorf("expected saml client scope %s to have realm_id of %s, but got %s", clientScope.Id, realm, clientScope.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientScopeDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_saml_client_scope" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			clientScope, _ := keycloakClient.GetSamlClientScope(realm, id)
			if clientScope != nil {
				return fmt.Errorf("saml client scope %s still exists", id)
			}
		}

		return nil
	}
}

func getSamlClientScopeFromState(s *terraform.State, resourceName string) (*keycloak.SamlClientScope, error) {
	keycloakClientScope := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	clientScope, err := keycloakClientScope.GetSamlClientScope(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting saml client scope %s: %s", id, err)
	}

	return clientScope, nil
}

func testKeycloakSamlClientScope_basic(realm, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}
	`, realm, clientScopeName)
}

func testKeycloakSamlClientScope_withConsentText(realm, clientScopeName, consentText string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client_scope" "client_scope" {
	name                = "%s"
	realm_id            = "${keycloak_realm.realm.id}"

	description         = "test description"

	consent_screen_text = "%s"
}
	`, realm, clientScopeName, consentText)
}

func testKeycloakSamlClientScope_withGuiOrder(realm, clientScopeName string, guiOrder int) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client_scope" "client_scope" {
	name                = "%s"
	realm_id            = "${keycloak_realm.realm.id}"

	description         = "test description"

	gui_order           = %d
}
	`, realm, clientScopeName, guiOrder)
}

func testKeycloakSamlClientScope_updateRealmBefore(realmOne, realmTwo, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_saml_client_scope" "client_scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm_1.id}"
}
	`, realmOne, realmTwo, clientScopeName)
}

func testKeycloakSamlClientScope_updateRealmAfter(realmOne, realmTwo, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_saml_client_scope" "client_scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm_2.id}"
}
	`, realmOne, realmTwo, clientScopeName)
}
