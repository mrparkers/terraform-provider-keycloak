package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakUser_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_basic(realmName, username),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakUser_createAfterManualDestroy(t *testing.T) {
	var user = &keycloak.User{}

	realmName := "terraform-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_basic(realmName, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					testAccCheckKeycloakUserFetch(resourceName, user),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteUser(user.RealmId, user.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakUser_basic(realmName, username),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakUser_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_updateRealmBefore(realmOne, realmTwo, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "realm_id", realmOne),
				),
			},
			{
				Config: testKeycloakUser_updateRealmAfter(realmOne, realmTwo, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "realm_id", realmTwo),
				),
			},
		},
	})
}

func testAccCheckKeycloakUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getUserFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakUserFetch(resourceName string, user *keycloak.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedUser, err := getUserFromState(s, resourceName)
		if err != nil {
			return err
		}

		user.Id = fetchedUser.Id
		user.RealmId = fetchedUser.RealmId

		return nil
	}
}

func testAccCheckKeycloakUserDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_user" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			user, _ := keycloakClient.GetUser(realm, id)
			if user != nil {
				return fmt.Errorf("user with id %s still exists", id)
			}
		}

		return nil
	}
}

func getUserFromState(s *terraform.State, resourceName string) (*keycloak.User, error) {
	keycloakUser := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	user, err := keycloakUser.GetUser(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user with id %s: %s", id, err)
	}

	return user, nil
}

func testKeycloakUser_basic(realm, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}
	`, realm, username)
}

func testKeycloakUser_updateRealmBefore(realmOne, realmTwo, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id  = "${keycloak_realm.realm_1.id}"
	username  = "%s"
}
	`, realmOne, realmTwo, username)
}

func testKeycloakUser_updateRealmAfter(realmOne, realmTwo, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id  = "${keycloak_realm.realm_2.id}"
	username  = "%s"
}
	`, realmOne, realmTwo, username)
}
