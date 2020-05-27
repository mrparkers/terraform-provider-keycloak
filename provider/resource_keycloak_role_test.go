package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakRole_basicRealm(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				ResourceName:        "keycloak_role.role",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakRole_basicRealmUrlRoleName(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-httpfoo.bara1b2" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				ResourceName:        "keycloak_role.role",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakRole_basicClient(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicClient(realmName, clientId, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				ResourceName:        "keycloak_role.role",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakRole_basicSamlClient(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicSamlClient(realmName, clientId, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				ResourceName:        "keycloak_role.role",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakRole_basicRealmUpdate(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)
	descriptionOne := acctest.RandString(50)
	descriptionTwo := acctest.RandString(50)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicRealmWithDescription(realmName, roleName, descriptionOne),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicRealmWithDescription(realmName, roleName, descriptionTwo),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
		},
	})
}

func TestAccKeycloakRole_basicClientUpdate(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)
	descriptionOne := acctest.RandString(50)
	descriptionTwo := acctest.RandString(50)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicClientWithDescription(realmName, clientId, roleName, descriptionOne),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicClientWithDescription(realmName, clientId, roleName, descriptionTwo),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
			{
				Config: testKeycloakRole_basicClient(realmName, clientId, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
		},
	})
}

func TestAccKeycloakRole_createAfterManualDestroy(t *testing.T) {
	var role = &keycloak.Role{}

	realmName := "terraform-" + acctest.RandString(10)
	roleName := "terraform-role-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRoleExists("keycloak_role.role"),
					testAccCheckKeycloakRoleFetch("keycloak_role.role", role),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteRole(role.RealmId, role.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakRole_basicRealm(realmName, roleName),
				Check:  testAccCheckKeycloakRoleExists("keycloak_role.role"),
			},
		},
	})
}

func TestAccKeycloakRole_composites(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientOne := "terraform-client-" + acctest.RandString(10)
	clientTwo := "terraform-client-" + acctest.RandString(10)
	roleOne := "terraform-role-one-" + acctest.RandString(10)
	roleTwo := "terraform-role-two-" + acctest.RandString(10)
	roleThree := "terraform-role-three-" + acctest.RandString(10)
	roleFour := "terraform-role-four-" + acctest.RandString(10)
	roleWithComposites := "terraform-role-with-composites-" + acctest.RandString(10)
	roleWithCompositesResourceName := "keycloak_role.role_with_composites"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			// initial setup - no composites attached
			{
				Config: testKeycloakRole_composites(realmName, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRoleExists("keycloak_role.role_1"),
					testAccCheckKeycloakRoleExists("keycloak_role.role_2"),
					testAccCheckKeycloakRoleExists("keycloak_role.role_3"),
					testAccCheckKeycloakRoleExists("keycloak_role.role_with_composites"),
					testAccCheckKeycloakRoleHasComposites(roleWithCompositesResourceName, []string{}),
				),
			},
			// add all composites
			{
				Config: testKeycloakRole_composites(realmName, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites, []string{
					"${keycloak_role.role_1.id}",
					"${keycloak_role.role_2.id}",
					"${keycloak_role.role_3.id}",
					"${keycloak_role.role_4.id}",
				}),
				Check: testAccCheckKeycloakRoleHasComposites(roleWithCompositesResourceName, []string{
					roleOne,
					roleTwo,
					roleThree,
					roleFour,
				}),
			},
			// remove two composites
			{
				Config: testKeycloakRole_composites(realmName, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites, []string{
					"${keycloak_role.role_1.id}",
					"${keycloak_role.role_2.id}",
				}),
				Check: testAccCheckKeycloakRoleHasComposites(roleWithCompositesResourceName, []string{
					roleOne,
					roleTwo,
				}),
			},
			// add them back and remove the others
			{
				Config: testKeycloakRole_composites(realmName, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites, []string{
					"${keycloak_role.role_3.id}",
					"${keycloak_role.role_4.id}",
				}),
				Check: testAccCheckKeycloakRoleHasComposites(roleWithCompositesResourceName, []string{
					roleThree,
					roleFour,
				}),
			},
			// remove them all
			{
				Config: testKeycloakRole_composites(realmName, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites, []string{}),
				Check:  testAccCheckKeycloakRoleHasComposites(roleWithCompositesResourceName, []string{}),
			},
		},
	})
}

func testAccCheckKeycloakRoleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakRoleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_role" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			role, _ := keycloakClient.GetRole(realm, id)
			if role != nil {
				return fmt.Errorf("role with id %s still exists", id)
			}
		}

		return nil
	}
}

func testAccCheckKeycloakRoleFetch(resourceName string, role *keycloak.Role) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedRole, err := getRoleFromState(state, resourceName)
		if err != nil {
			return err
		}

		role.Id = fetchedRole.Id
		role.Name = fetchedRole.Name
		role.RealmId = fetchedRole.RealmId
		role.ClientId = fetchedRole.ClientId

		return nil
	}
}

func testAccCheckKeycloakRoleHasComposites(resourceName string, compositeRoleNames []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		role, err := getRoleFromState(state, resourceName)
		if err != nil {
			return err
		}

		if len(compositeRoleNames) != 0 && !role.Composite {
			return fmt.Errorf("expected role %s to have composites, but has none", role.Name)
		}

		if len(compositeRoleNames) == 0 && role.Composite {
			return fmt.Errorf("expected role %s to have no composites, but has some", role.Name)
		}

		composites, err := keycloakClient.GetRoleComposites(role)
		if err != nil {
			return err
		}

		for _, compositeRoleName := range compositeRoleNames {
			var found bool

			for _, composite := range composites {
				if composite.Name == compositeRoleName {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("expected role %s to have composite %s", role.Name, compositeRoleName)
			}
		}

		for _, composite := range composites {
			var found bool

			for _, compositeRoleName := range compositeRoleNames {
				if composite.Name == compositeRoleName {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("role %s had unexpected composite %s", role.Name, composite.Name)
			}
		}

		return nil
	}
}

func getRoleFromState(s *terraform.State, resourceName string) (*keycloak.Role, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	role, err := keycloakClient.GetRole(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting role with id %s: %s", id, err)
	}

	return role, nil
}

func testKeycloakRole_basicRealm(realm, role string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
	`, realm, role)
}

func testKeycloakRole_basicRealmWithDescription(realm, role, description string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name        = "%s"
	description = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
}
	`, realm, role, description)
}

func testKeycloakRole_basicClient(realm, clientId, role string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_role" "role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.client.id}"
}
	`, realm, clientId, role)
}

func testKeycloakRole_basicSamlClient(realm, clientId, role string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.client.id}"
}
	`, realm, clientId, role)
}

func testKeycloakRole_basicClientWithDescription(realm, clientId, role, description string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_role" "role" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "${keycloak_openid_client.client.id}"
	description = "%s"
}
	`, realm, clientId, role, description)
}

func testKeycloakRole_composites(realm, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites string, composites []string) string {
	var tfComposites string
	if len(composites) != 0 {
		tfComposites = fmt.Sprintf("composite_roles = %s", arrayOfStringsForTerraformResource(composites))
	}

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client_one" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_openid_client" "client_two" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_role" "role_1" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "role_2" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.client_one.id}"
}

resource "keycloak_role" "role_3" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "role_4" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.client_two.id}"
}

resource "keycloak_role" "role_with_composites" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"

	%s
}
	`, realm, clientOne, clientTwo, roleOne, roleTwo, roleThree, roleFour, roleWithComposites, tfComposites)
}
