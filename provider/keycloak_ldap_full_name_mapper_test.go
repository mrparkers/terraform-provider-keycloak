package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakLdapFullNameMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	fullNameMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapFullNameMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapFullNameMapper_basic(realmName, fullNameMapperName),
				Check:  testAccCheckKeycloakLdapFullNameMapperExists("keycloak_ldap_full_name_mapper.full-name-mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapFullNameMapper_updateLdapUserFederation(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapFullNameMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapFullNameMapper_updateLdapUserFederationBefore(realmOne, realmTwo, mapperName),
				Check:  testAccCheckKeycloakLdapFullNameMapperExists("keycloak_ldap_full_name_mapper.full-name-mapper"),
			},
			{
				Config: testKeycloakLdapFullNameMapper_updateLdapUserFederationAfter(realmOne, realmTwo, mapperName),
				Check:  testAccCheckKeycloakLdapFullNameMapperExists("keycloak_ldap_full_name_mapper.full-name-mapper"),
			},
		},
	})
}

func testAccCheckKeycloakLdapFullNameMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapFullNameMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapFullNameMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_full_name_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldapFullNameMapper, _ := keycloakClient.GetLdapFullNameMapper(realm, id)
			if ldapFullNameMapper != nil {
				return fmt.Errorf("ldap full name mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapFullNameMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapFullNameMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapFullNameMapper, err := keycloakClient.GetLdapFullNameMapper(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap full name mapper with id %s: %s", id, err)
	}

	return ldapFullNameMapper, nil
}

func testKeycloakLdapFullNameMapper_basic(realm, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
  name                    = "openldap"
  realm_id                = "${keycloak_realm.realm.id}"

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

resource "keycloak_ldap_full_name_mapper" "full-name-mapper" {
  name                     = "%s"
  realm_id                 = "${keycloak_realm.realm.id}"
  ldap_user_federation_id  = "${keycloak_ldap_user_federation.openldap.id}"

  ldap_full_name_attribute = "cn"
}
	`, realm, mapperName)
}

func testKeycloakLdapFullNameMapper_updateLdapUserFederationBefore(realmOne, realmTwo, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-one" {
	realm = "%s"
}

resource "keycloak_realm" "realm-two" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap-one" {
  name                    = "openldap"
  realm_id                = "${keycloak_realm.realm-one.id}"

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

resource "keycloak_ldap_user_federation" "openldap-two" {
  name                    = "openldap"
  realm_id                = "${keycloak_realm.realm-two.id}"

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

resource "keycloak_ldap_full_name_mapper" "full-name-mapper" {
  name                     = "%s"
  realm_id                 = "${keycloak_realm.realm-one.id}"
  ldap_user_federation_id  = "${keycloak_ldap_user_federation.openldap-one.id}"

  ldap_full_name_attribute = "cn"
}
	`, realmOne, realmTwo, mapperName)
}

func testKeycloakLdapFullNameMapper_updateLdapUserFederationAfter(realmOne, realmTwo, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm-one" {
	realm = "%s"
}

resource "keycloak_realm" "realm-two" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap-one" {
  name                    = "openldap"
  realm_id                = "${keycloak_realm.realm-one.id}"

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

resource "keycloak_ldap_user_federation" "openldap-two" {
  name                    = "openldap"
  realm_id                = "${keycloak_realm.realm-two.id}"

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

resource "keycloak_ldap_full_name_mapper" "full-name-mapper" {
  name                     = "%s"
  realm_id                 = "${keycloak_realm.realm-two.id}"
  ldap_user_federation_id  = "${keycloak_ldap_user_federation.openldap-two.id}"

  ldap_full_name_attribute = "cn"
}
	`, realmOne, realmTwo, mapperName)
}
