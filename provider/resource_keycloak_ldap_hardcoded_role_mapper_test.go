package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakLdapHardcodedRoleMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapHardcodedRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapHardcodedRoleMapper(realmName, roleMapperName),
				Check:  testAccCheckKeycloakLdapHardcodedRoleMapperExists("keycloak_ldap_hardcoded_role_mapper.hardcoded_role_mapper"),
			},
			{
				ResourceName:      "keycloak_ldap_hardcoded_role_mapper.hardcoded_role_mapper",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapGenericMapperImportId("keycloak_ldap_hardcoded_role_mapper.hardcoded_role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapHardcodedRoleMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.LdapHardcodedRoleMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapHardcodedRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapHardcodedRoleMapper(realmName, roleMapperName),
				Check:  testAccCheckKeycloakLdapHardcodedRoleMapperFetch("keycloak_ldap_hardcoded_role_mapper.hardcoded_role_mapper", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteLdapHardcodedRoleMapper(mapper.RealmId, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapHardcodedRoleMapper(realmName, roleMapperName),
				Check:  testAccCheckKeycloakLdapHardcodedRoleMapperExists("keycloak_ldap_hardcoded_role_mapper.hardcoded_role_mapper"),
			},
		},
	})
}

func testAccCheckKeycloakLdapHardcodedRoleMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapHardcodedRoleMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapHardcodedRoleMapperFetch(resourceName string, mapper *keycloak.LdapHardcodedRoleMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getLdapHardcodedRoleMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckKeycloakLdapHardcodedRoleMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_hardcoded_role_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldapMapper, _ := keycloakClient.GetLdapHardcodedRoleMapper(realm, id)
			if ldapMapper != nil {
				return fmt.Errorf("ldap hardcoded role mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapHardcodedRoleMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapHardcodedRoleMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapMapper, err := keycloakClient.GetLdapHardcodedRoleMapper(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap group mapper with id %s: %s", id, err)
	}

	return ldapMapper, nil
}

func testKeycloakLdapHardcodedRoleMapper(realm, roleMapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "openldap"
	realm_id                = keycloak_realm.realm.id

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "admin"
}

resource "keycloak_role" "hardcoded_role_mapper_test" {
    realm_id    = keycloak_realm.realm.id
    name        = "hardcoded-role-test"
}

resource "keycloak_ldap_hardcoded_role_mapper" "hardcoded_role_mapper" {
	name                        = "%s"
	realm_id                    = keycloak_realm.realm.id
	ldap_user_federation_id     = keycloak_ldap_user_federation.openldap.id
	role                        = keycloak_role.hardcoded_role_mapper_test.name
}
	`, realm, roleMapperName)
}
