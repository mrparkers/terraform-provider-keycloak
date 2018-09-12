package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakLdapGroupMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapGroupMapper_basic(realmName, groupMapperName),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group-mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_updateLdapUserFederation(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapGroupMapper_updateLdapUserFederationBefore(realmOne, realmTwo, groupMapperName),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group-mapper"),
			},
			{
				Config: testKeycloakLdapGroupMapper_updateLdapUserFederationAfter(realmOne, realmTwo, groupMapperName),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group-mapper"),
			},
		},
	})
}

func testAccCheckKeycloakLdapGroupMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapGroupMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapGroupMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_group_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldapGroupMapper, _ := keycloakClient.GetLdapGroupMapper(realm, id)
			if ldapGroupMapper != nil {
				return fmt.Errorf("ldap group mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapGroupMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapGroupMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapGroupMapper, err := keycloakClient.GetLdapGroupMapper(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap group mapper with id %s: %s", id, err)
	}

	return ldapGroupMapper, nil
}

func testKeycloakLdapGroupMapper_basic(realm, groupMapperName string) string {
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

resource "keycloak_ldap_group_mapper" "group-mapper" {
  name                        = "%s"
  realm_id                    = "${keycloak_realm.realm.id}"
  ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

  ldap_groups_dn                 = "dc=example,dc=org"
  group_name_ldap_attribute      = "cn"
  group_object_classes           = [
    "groupOfNames"
  ]
  membership_attribute_type      = "DN"
  membership_ldap_attribute      = "member"
  membership_user_ldap_attribute = "cn"
  memberof_ldap_attribute        = "memberOf"
}
	`, realm, groupMapperName)
}

func testKeycloakLdapGroupMapper_updateLdapUserFederationBefore(realmOne, realmTwo, groupMapperName string) string {
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

resource "keycloak_ldap_group_mapper" "group-mapper" {
  name                        = "%s"
  realm_id                    = "${keycloak_realm.realm-one.id}"
  ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap-one.id}"

  ldap_groups_dn                 = "dc=example,dc=org"
  group_name_ldap_attribute      = "cn"
  group_object_classes           = [
    "groupOfNames"
  ]
  membership_attribute_type      = "DN"
  membership_ldap_attribute      = "member"
  membership_user_ldap_attribute = "cn"
  memberof_ldap_attribute        = "memberOf"
}
	`, realmOne, realmTwo, groupMapperName)
}

func testKeycloakLdapGroupMapper_updateLdapUserFederationAfter(realmOne, realmTwo, groupMapperName string) string {
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

resource "keycloak_ldap_group_mapper" "group-mapper" {
  name                        = "%s"
  realm_id                    = "${keycloak_realm.realm-two.id}"
  ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap-two.id}"

  ldap_groups_dn                 = "dc=example,dc=org"
  group_name_ldap_attribute      = "cn"
  group_object_classes           = [
    "groupOfNames"
  ]
  membership_attribute_type      = "DN"
  membership_ldap_attribute      = "member"
  membership_user_ldap_attribute = "cn"
  memberof_ldap_attribute        = "memberOf"
}
	`, realmOne, realmTwo, groupMapperName)
}
