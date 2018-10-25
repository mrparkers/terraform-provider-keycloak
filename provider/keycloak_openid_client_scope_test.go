package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakClientScope_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
			},
			{
				ResourceName:        "keycloak_openid_client_scope.client_scope",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakClientScope_createAfterManualDestroy(t *testing.T) {
	var clientScope = &keycloak.OpenidClientScope{}

	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClientScope_basic(realmName, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
					testAccCheckKeycloakClientScopeFetch("keycloak_openid_client_scope.client_scope", clientScope),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClientScope(clientScope.RealmId, clientScope.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
			},
		},
	})
}

func TestAccKeycloakClientScope_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClientScope_updateRealmBefore(realmOne, realmTwo, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
					testAccCheckKeycloakClientScopeBelongsToRealm("keycloak_openid_client_scope.client_scope", realmOne),
				),
			},
			{
				Config: testKeycloakClientScope_updateRealmAfter(realmOne, realmTwo, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
					testAccCheckKeycloakClientScopeBelongsToRealm("keycloak_openid_client_scope.client_scope", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakClientScope_consentScreenText(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientScopeName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientScopeDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
			},
			{
				Config: testKeycloakClientScope_withConsentText(realmName, clientScopeName, acctest.RandString(10)),
				Check:  testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
			},
			{
				Config: testKeycloakClientScope_basic(realmName, clientScopeName),
				Check:  testAccCheckKeycloakClientScopeExistsWithCorrectProtocol("keycloak_openid_client_scope.client_scope"),
			},
		},
	})
}

func testAccCheckKeycloakClientScopeExistsWithCorrectProtocol(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientScope, err := getClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		if clientScope.Protocol != "openid-connect" {
			return fmt.Errorf("expected openid client scope to have openid-connect protocol, but got %s", clientScope.Protocol)
		}

		return nil
	}
}

func testAccCheckKeycloakClientScopeFetch(resourceName string, clientScope *keycloak.OpenidClientScope) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClientScope, err := getClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		clientScope.Id = fetchedClientScope.Id
		clientScope.RealmId = fetchedClientScope.RealmId

		return nil
	}
}

func testAccCheckKeycloakClientScopeBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clientScope, err := getClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		if clientScope.RealmId != realm {
			return fmt.Errorf("expected openid client scope %s to have realm_id of %s, but got %s", clientScope.Id, realm, clientScope.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakClientScopeDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_scope" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			clientScope, _ := keycloakClient.GetOpenidClientScope(realm, id)
			if clientScope != nil {
				return fmt.Errorf("openid client scope %s still exists", id)
			}
		}

		return nil
	}
}

func getClientScopeFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientScope, error) {
	keycloakClientScope := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	clientScope, err := keycloakClientScope.GetOpenidClientScope(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client scope %s: %s", id, err)
	}

	return clientScope, nil
}

func testKeycloakClientScope_basic(realm, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}
	`, realm, clientScopeName)
}

func testKeycloakClientScope_withConsentText(realm, clientScopeName, consentText string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name                = "%s"
	realm_id            = "${keycloak_realm.realm.id}"

	description         = "test description"

	consent_screen_text = "%s"
}
	`, realm, clientScopeName, consentText)
}

func testKeycloakClientScope_updateRealmBefore(realmOne, realmTwo, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm_1.id}"
}
	`, realmOne, realmTwo, clientScopeName)
}

func testKeycloakClientScope_updateRealmAfter(realmOne, realmTwo, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm_2.id}"
}
	`, realmOne, realmTwo, clientScopeName)
}
