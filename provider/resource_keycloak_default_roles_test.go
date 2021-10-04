package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakDefaultRoles_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakDefaultRoles_basic(),
				Check:  testAccCheckDefaultRolesExists("keycloak_default_roles.default_roles"),
			},
		},
	})
}

func TestAccKeycloakDefaultRoles_updateDefaultRoles(t *testing.T) {
	t.Parallel()

	groupDefaultRolesOne := &keycloak.DefaultRoles{
		RealmId:      testAccRealmUserFederation.Realm,
		DefaultRoles: []string{"\"uma_authorization\""},
	}

	groupDefaultRolesTwo := &keycloak.DefaultRoles{
		RealmId:      testAccRealmUserFederation.Realm,
		DefaultRoles: []string{"\"uma_authorization\",", "\"offline_access\""},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakDefaultRoles_basicFromInterface(groupDefaultRolesOne),
				Check:  testAccCheckDefaultRolesExists("keycloak_default_roles.default_roles"),
			},
			{
				Config: testKeycloakDefaultRoles_basicFromInterface(groupDefaultRolesTwo),
				Check:  testAccCheckDefaultRolesExists("keycloak_default_roles.default_roles"),
			},
		},
	})
}

func testAccCheckDefaultRolesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakDefaultRolesFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func getKeycloakDefaultRolesFromState(s *terraform.State, resourceName string) (*keycloak.DefaultRoles, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	composites, err := keycloakClient.GetDefaultRoles(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting defaultRoles with id %s: %s", id, err)
	}

	defaultRoleNamesList, _ := getDefaultRoleNames(composites)

	defaultRoles := &keycloak.DefaultRoles{
		Id:           id,
		RealmId:      realm,
		DefaultRoles: defaultRoleNamesList,
	}

	return defaultRoles, nil
}

func testKeycloakDefaultRoles_basic() string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_default_roles" "default_roles" {
	realm_id  = data.keycloak_realm.realm.id
    default_roles = ["uma_authorization"]
}
	`, testAccRealmUserFederation.Realm)
}

func testKeycloakDefaultRoles_basicFromInterface(defaultRoles *keycloak.DefaultRoles) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_default_roles" "default_roles" {
	realm_id  = data.keycloak_realm.realm.id
    default_roles = %s
}
	`, testAccRealmUserFederation.Realm, defaultRoles.DefaultRoles)
}
