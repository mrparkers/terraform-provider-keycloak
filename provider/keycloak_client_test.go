package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakClient_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClient_basic(realmName, clientId),
				Check:  testAccCheckKeycloakClientExists("keycloak_client.client"),
			},
		},
	})
}

func TestAccKeycloakClient_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakClient_updateRealmBefore(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientExists("keycloak_client.client"),
					testAccCheckKeycloakClientBelongsToRealm("keycloak_client.client", realmOne),
				),
			},
			{
				Config: testKeycloakClient_updateRealmAfter(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakClientExists("keycloak_client.client"),
					testAccCheckKeycloakClientBelongsToRealm("keycloak_client.client", realmTwo),
				),
			},
		},
	})
}

func testAccCheckKeycloakClientExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakClientBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.RealmId != realm {
			return fmt.Errorf("expected client %s to have realm_id of %s, but got %s", client.ClientId, realm, client.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakClientDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_client" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			client, _ := keycloakClient.GetClient(realm, id)
			if client != nil {
				return fmt.Errorf("client %s still exists", id)
			}
		}

		return nil
	}
}

func getClientFromState(s *terraform.State, resourceName string) (*keycloak.Client, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting client %s: %s", id, err)
	}

	return client, nil
}

func testKeycloakClient_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_client" "client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}
	`, realm, clientId)
}

func testKeycloakClient_updateRealmBefore(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-1" {
	realm = "%s"
}

resource "keycloak_realm" "realm-2" {
	realm = "%s"
}

resource "keycloak_client" "client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm-1.id}"
}
	`, realmOne, realmTwo, clientId)
}

func testKeycloakClient_updateRealmAfter(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-1" {
	realm = "%s"
}

resource "keycloak_realm" "realm-2" {
	realm = "%s"
}

resource "keycloak_client" "client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm-2.id}"
}
	`, realmOne, realmTwo, clientId)
}
