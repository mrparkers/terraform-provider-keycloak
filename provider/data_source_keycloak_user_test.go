package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakDataSourceUser(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	username := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakUser(realm, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists("keycloak_user.user"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "id", "data.keycloak_user.user", "id"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "realm_id", "data.keycloak_user.user", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "username", "data.keycloak_user.user", "username"),
					testAccCheckDataKeycloakUser("data.keycloak_user.user"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakUser(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		id := rs.Primary.ID
		realmID := rs.Primary.Attributes["realm_id"]
		username := rs.Primary.Attributes["username"]

		user, err := keycloakClient.GetUser(realmID, id)
		if err != nil {
			return err
		}

		if user.Username != username {
			return fmt.Errorf("expected user with ID %s to have username %s, but got %s", id, username, user.Username)
		}

		return nil
	}
}

func testDataSourceKeycloakUser(realm, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm 		= "%s"
}

resource "keycloak_user" "user" {
	username    = "%s"
	realm_id 	= "${keycloak_realm.realm.id}"
	enabled    	= true

    email      	= "bob@domain.com"
    first_name 	= "Bob"
    last_name  	= "Bobson"
}

data "keycloak_user" "user" {
	realm_id 	= "${keycloak_realm.realm.id}"
	username    = "${keycloak_user.user.username}"
}
	`, realm, username)
}
