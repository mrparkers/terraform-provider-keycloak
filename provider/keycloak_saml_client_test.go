package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strconv"
	"testing"
)

func TestAccKeycloakSamlClient_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_basic(realmName, clientId),
				Check:  testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
			},
			{
				ResourceName:        "keycloak_saml_client.saml_client",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakSamlClient_createAfterManualDestroy(t *testing.T) {
	var client = &keycloak.SamlClient{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_basic(realmName, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientFetch("keycloak_saml_client.saml_client", client),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteSamlClient(client.RealmId, client.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakSamlClient_basic(realmName, clientId),
				Check:  testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
			},
		},
	})
}

func TestAccKeycloakSamlClient_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_updateRealmBefore(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					resource.TestCheckResourceAttr("keycloak_saml_client.saml_client", "realm_id", realmOne),
				),
			},
			{
				Config: testKeycloakSamlClient_updateRealmAfter(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					resource.TestCheckResourceAttr("keycloak_saml_client.saml_client", "realm_id", realmTwo),
				),
			},
		},
	})
}

// Keycloak typically sets some values as default if they aren't provided
// This test asserts that these default values are present if none are provided
func TestAccKeycloakSamlClient_keycloakDefaults(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_basic(realmName, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientHasDefaultBooleanAttributes("keycloak_saml_client.saml_client"),
					TestCheckResourceAttrNot("keycloak_saml_client.saml_client", "signing_certificate", ""),
					TestCheckResourceAttrNot("keycloak_saml_client.saml_client", "signing_private_key", ""),
				),
			},
		},
	})
}

func testAccCheckKeycloakSamlClientExistsWithCorrectProtocol(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Protocol != "saml" {
			return fmt.Errorf("expected saml client to have saml protocol, but got %s", client.Protocol)
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientFetch(resourceName string, client *keycloak.SamlClient) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClient, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		client.Id = fetchedClient.Id
		client.RealmId = fetchedClient.RealmId

		return nil
	}
}

func testAccCheckKeycloakSamlClientDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_saml_client" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			client, _ := keycloakClient.GetSamlClient(realm, id)
			if client != nil {
				return fmt.Errorf("saml client %s still exists", id)
			}
		}

		return nil
	}
}

func getSamlClientFromState(s *terraform.State, resourceName string) (*keycloak.SamlClient, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetSamlClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting saml client %s: %s", id, err)
	}

	return client, nil
}

func testAccCheckKeycloakSamlClientHasDefaultBooleanAttributes(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		includeAuthnStatement, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["include_authn_statement"])
		if err != nil {
			return err
		}

		signDocuments, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["sign_documents"])
		if err != nil {
			return err
		}

		signAssertions, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["sign_assertions"])
		if err != nil {
			return err
		}

		clientSignatureRequired, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["client_signature_required"])
		if err != nil {
			return err
		}

		forcePostBinding, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["force_post_binding"])
		if err != nil {
			return err
		}

		if !includeAuthnStatement && !signDocuments && !signAssertions && !clientSignatureRequired && !forcePostBinding {
			return fmt.Errorf("expected saml client with id %s to have some defaults set by Keycloak", rs.Primary.ID)
		}

		return nil
	}
}

func parseBoolAndTreatEmptyStringAsFalse(b string) (bool, error) {
	if b == "" {
		return false, nil
	}

	return strconv.ParseBool(b)
}

func testKeycloakSamlClient_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
}
	`, realm, clientId)
}

func testKeycloakSamlClient_updateRealmBefore(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm_1.id}"
}
	`, realmOne, realmTwo, clientId)
}

func testKeycloakSamlClient_updateRealmAfter(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm_2.id}"
}
	`, realmOne, realmTwo, clientId)
}
