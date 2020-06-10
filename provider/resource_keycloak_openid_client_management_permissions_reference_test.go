package provider

import (
	"fmt"
	"testing"

	"github.com/mrparkers/terraform-provider-keycloak/keycloak"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKeycloakOpenIdClientManagementPermissionsReference_basic(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testKeycloakOpenIdClientManagementPermissionReferenceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdClientManagementPermissionsReference(realmName, clientId),
				Check:  testKeycloakOpenIdClientManagementPermissionsReferenceExists("keycloak_openid_client_management_permissions_reference.management_permissions_reference"),
			},
		},
	})
}

func testKeycloakOpenIdClientManagementPermissionsReference(realmId, clientId string) string {
	return fmt.Sprintf(`
	resource "keycloak_realm" "realm" {
		realm = "%s"
	}
	
	resource "keycloak_openid_client" "client" {
		client_id   = "%s"
		realm_id    = "${keycloak_realm.realm.id}"

		access_type = "PUBLIC"
	}
	
	resource "keycloak_openid_client_management_permissions_reference" "management_permissions_reference" {
		realm_id    = "${keycloak_realm.realm.id}"
		client_id 	= "${keycloak_openid_client.client.id}"
	}
	`, realmId, clientId)
}

func getOpenIdClientManagementPermissionsReferenceUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdClientManagementPermissionsReference, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdClientManagementPermissionsReference(realmId, clientId)
}

func testKeycloakOpenIdClientManagementPermissionsReferenceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getOpenIdClientManagementPermissionsReferenceUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdClientManagementPermissionReferenceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_management_permissions_reference" {
				continue
			}

			id := rs.Primary.ID
			realmId := rs.Primary.Attributes["realm_id"]
			clientId := rs.Primary.Attributes["client_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			reference, _ := keycloakClient.GetOpenIdClientManagementPermissionsReference(realmId, clientId)
			if reference != nil {
				return fmt.Errorf("management permission reference with id %s still exists", id)
			}
		}

		return nil
	}
}
