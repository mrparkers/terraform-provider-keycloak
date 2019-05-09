package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakCustomUserFederation_basic(t *testing.T) {
	skipIfEnvSet(t, "CI") // temporary while I figure out how to load this custom provider in CI

	realmName := "terraform-" + acctest.RandString(10)
	name := "terraform-" + acctest.RandString(10)
	providerId := "custom"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakCustomUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakCustomUserFederation_basic(realmName, name, providerId),
				Check:  testAccCheckKeycloakCustomUserFederationExists("keycloak_custom_user_federation.custom"),
			},
			{
				ResourceName:        "keycloak_custom_user_federation.custom",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakCustomUserFederation_createAfterManualDestroy(t *testing.T) {
	skipIfEnvSet(t, "CI") // temporary while I figure out how to load this custom provider in CI

	var customFederation = &keycloak.CustomUserFederation{}

	realmName := "terraform-" + acctest.RandString(10)
	name := "terraform-" + acctest.RandString(10)
	providerId := "custom"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakCustomUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakCustomUserFederation_basic(realmName, name, providerId),
				Check:  testAccCheckKeycloakCustomUserFederationFetch("keycloak_custom_user_federation.custom", customFederation),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteCustomUserFederation(customFederation.RealmId, customFederation.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakCustomUserFederation_basic(realmName, name, providerId),
				Check:  testAccCheckKeycloakCustomUserFederationExists("keycloak_custom_user_federation.custom"),
			},
		},
	})
}

func TestAccKeycloakCustomUserFederation_validation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	name := "terraform-" + acctest.RandString(10)
	providerId := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakCustomUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakCustomUserFederation_basic(realmName, name, providerId),
				ExpectError: regexp.MustCompile("custom user federation provider with id .+ is not installed on the server"),
			},
		},
	})
}

func testAccCheckKeycloakCustomUserFederationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getCustomUserFederationFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakCustomUserFederationFetch(resourceName string, federation *keycloak.CustomUserFederation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedFederation, err := getCustomUserFederationFromState(s, resourceName)
		if err != nil {
			return err
		}

		federation.Id = fetchedFederation.Id
		federation.RealmId = fetchedFederation.RealmId

		return nil
	}
}

func testAccCheckKeycloakCustomUserFederationDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_custom_user_federation" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			custom, _ := keycloakClient.GetCustomUserFederation(realm, id)
			if custom != nil {
				return fmt.Errorf("custom user federation with id %s still exists", id)
			}
		}

		return nil
	}
}

func getCustomUserFederationFromState(s *terraform.State, resourceName string) (*keycloak.CustomUserFederation, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	custom, err := keycloakClient.GetCustomUserFederation(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting custom user federation with id %s: %s", id, err)
	}

	return custom, nil
}

func testKeycloakCustomUserFederation_basic(realm, name, providerId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_custom_user_federation" "custom" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	provider_id = "%s"

	enabled     = true
}
	`, realm, name, providerId)
}
