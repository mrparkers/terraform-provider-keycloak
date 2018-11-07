package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOpenidClientDefaultScopes_basic(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	clientScopes := []string{
		"profile",
		"email",
		clientScope,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientDefaultScopes_basic(realm, client, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasDefaultScopes("keycloak_openid_client_default_client_scopes.default_scopes", clientScopes),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientHasDefaultScopes(resourceName string, tfDefaultClientScopes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realm := rs.Primary.Attributes["realm_id"]
		client := rs.Primary.Attributes["client_id"]

		keycloakDefaultClientScopes, err := keycloakClient.GetOpenidClientDefaultScopes(realm, client)
		if err != nil {
			return err
		}

		for _, tfDefaultClientScope := range tfDefaultClientScopes {
			found := false

			for _, keycloakDefaultScope := range keycloakDefaultClientScopes {
				if keycloakDefaultScope.Name == tfDefaultClientScope {
					found = true

					break
				}
			}

			if !found {
				return fmt.Errorf("default scope %s is not assigned to client", tfDefaultClientScope)
			}
		}

		return nil
	}
}

// TODO: don't interpolate client realm_id (implementing delete will fix this)
func testKeycloakOpenidClientDefaultScopes_basic(realm, client, clientScope string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_openid_client_scope.client_scope.realm_id}"
	access_type = "PUBLIC"

	valid_redirect_uris = ["foo"]
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_default_client_scopes" "default_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	default_scopes = [
        "profile",
        "email",
        "${keycloak_openid_client_scope.client_scope.name}"
    ]
}
	`, realm, client, clientScope)
}
