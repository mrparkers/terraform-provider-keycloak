package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOpenidClient_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
			},
			{
				ResourceName:        "keycloak_openid_client.client",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakOpenidClient_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_updateRealmBefore(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientBelongsToRealm("keycloak_openid_client.client", realmOne),
				),
			},
			{
				Config: testKeycloakOpenidClient_updateRealmAfter(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientBelongsToRealm("keycloak_openid_client.client", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_accessType(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_accessType(realmName, clientId, "CONFIDENTIAL"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", false, false),
			},
			{
				Config: testKeycloakOpenidClient_accessType(realmName, clientId, "PUBLIC"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", true, false),
			},
			{
				Config: testKeycloakOpenidClient_accessType(realmName, clientId, "BEARER-ONLY"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", false, true),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_updateInPlace(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	enabled := randomBool()
	accessTypeBefore := randomStringInSlice(keycloakOpenidClientAccessTypes)
	accessTypeAfter := randomStringInSlice(keycloakOpenidClientAccessTypes)

	openidClientBefore := &keycloak.OpenidClient{
		RealmId:      realm,
		ClientId:     clientId,
		Enabled:      enabled,
		Description:  acctest.RandString(50),
		ClientSecret: acctest.RandString(10),
	}

	openidClientAfter := &keycloak.OpenidClient{
		RealmId:      realm,
		ClientId:     clientId,
		Enabled:      !enabled,
		Description:  acctest.RandString(50),
		ClientSecret: acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_fromInterface(openidClientBefore, accessTypeBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
				),
			},
			{
				Config: testKeycloakOpenidClient_fromInterface(openidClientAfter, accessTypeAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_secret(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	clientSecret := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_secret(realmName, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientHasClientSecret("keycloak_openid_client.client", clientSecret),
				),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Protocol != "openid-connect" {
			return fmt.Errorf("expected openid client to have openid-connect protocol, but got %s", client.Protocol)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAccessType(resourceName string, public, bearer bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.PublicClient != public {
			return fmt.Errorf("expected openid client to have public set to %t, but got %t", public, client.PublicClient)
		}

		if client.BearerOnly != bearer {
			return fmt.Errorf("expected openid client to have bearer set to %t, but got %t", bearer, client.BearerOnly)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.RealmId != realm {
			return fmt.Errorf("expected openid client %s to have realm_id of %s, but got %s", client.ClientId, realm, client.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientHasClientSecret(resourceName, secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.ClientSecret != secret {
			return fmt.Errorf("expected openid client %s to have secret value of %s, but got %s", client.ClientId, secret, client.ClientSecret)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			client, _ := keycloakClient.GetOpenidClient(realm, id)
			if client != nil {
				return fmt.Errorf("openid client %s still exists", id)
			}
		}

		return nil
	}
}

func getOpenidClientFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClient, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetOpenidClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client %s: %s", id, err)
	}

	return client, nil
}

func testKeycloakOpenidClient_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}
	`, realm, clientId)
}

func testKeycloakOpenidClient_accessType(realm, clientId, accessType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "%s"
}
	`, realm, clientId, accessType)
}

func testKeycloakOpenidClient_updateRealmBefore(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-1" {
	realm = "%s"
}

resource "keycloak_realm" "realm-2" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm-1.id}"
	access_type = "CONFIDENTIAL"
}
	`, realmOne, realmTwo, clientId)
}

func testKeycloakOpenidClient_updateRealmAfter(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-1" {
	realm = "%s"
}

resource "keycloak_realm" "realm-2" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm-2.id}"
	access_type = "CONFIDENTIAL"
}
	`, realmOne, realmTwo, clientId)
}

func testKeycloakOpenidClient_fromInterface(openidClient *keycloak.OpenidClient, accessType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id     = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
	access_type   = "%s"
	client_secret = "%s"

	enabled     = %t
	description = "%s"
}
	`, openidClient.RealmId, openidClient.ClientId, accessType, openidClient.ClientSecret, openidClient.Enabled, openidClient.Description)
}

func testKeycloakOpenidClient_secret(realm, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id     = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
	access_type   = "CONFIDENTIAL"
	client_secret = "%s"
}
	`, realm, clientId, clientSecret)
}
