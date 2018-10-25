package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakLdapUserAttributeMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	userAttributeMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserAttributeMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserAttributeMapper_basic(realmName, userAttributeMapperName),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperExists("keycloak_ldap_user_attribute_mapper.username"),
			},
			{
				ResourceName:      "keycloak_ldap_user_attribute_mapper.username",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapGenericMapperImportId("keycloak_ldap_user_attribute_mapper.username"),
			},
		},
	})
}

func TestAccKeycloakLdapUserAttributeMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.LdapUserAttributeMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	userAttributeMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserAttributeMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserAttributeMapper_basic(realmName, userAttributeMapperName),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperFetch("keycloak_ldap_user_attribute_mapper.username", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteLdapUserAttributeMapper(mapper.RealmId, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapUserAttributeMapper_basic(realmName, userAttributeMapperName),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperExists("keycloak_ldap_user_attribute_mapper.username"),
			},
		},
	})
}

func TestAccKeycloakLdapUserAttributeMapper_updateLdapUserFederation(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	userAttributeMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserAttributeMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserAttributeMapper_updateLdapUserFederationBefore(realmOne, realmTwo, userAttributeMapperName),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperExists("keycloak_ldap_user_attribute_mapper.username"),
			},
			{
				Config: testKeycloakLdapUserAttributeMapper_updateLdapUserFederationAfter(realmOne, realmTwo, userAttributeMapperName),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperExists("keycloak_ldap_user_attribute_mapper.username"),
			},
		},
	})
}

func TestAccKeycloakLdapUserAttributeMapper_updateInPlace(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	userAttributeMapperBefore := &keycloak.LdapUserAttributeMapper{
		Name:                    acctest.RandString(10),
		UserModelAttribute:      acctest.RandString(10),
		LdapAttribute:           acctest.RandString(10),
		IsMandatoryInLdap:       randomBool(),
		ReadOnly:                randomBool(),
		AlwaysReadValueFromLdap: randomBool(),
	}
	userAttributeMapperAfter := &keycloak.LdapUserAttributeMapper{
		Name:                    acctest.RandString(10),
		UserModelAttribute:      acctest.RandString(10),
		LdapAttribute:           acctest.RandString(10),
		IsMandatoryInLdap:       randomBool(),
		ReadOnly:                randomBool(),
		AlwaysReadValueFromLdap: randomBool(),
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserAttributeMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserAttributeMapper_basicFromInterface(realm, userAttributeMapperBefore),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperExists("keycloak_ldap_user_attribute_mapper.username"),
			},
			{
				Config: testKeycloakLdapUserAttributeMapper_basicFromInterface(realm, userAttributeMapperAfter),
				Check:  testAccCheckKeycloakLdapUserAttributeMapperExists("keycloak_ldap_user_attribute_mapper.username"),
			},
		},
	})
}

func testAccCheckKeycloakLdapUserAttributeMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapUserAttributeMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapUserAttributeMapperFetch(resourceName string, mapper *keycloak.LdapUserAttributeMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getLdapUserAttributeMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckKeycloakLdapUserAttributeMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_user_attribute_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldapUserAttributeMapper, _ := keycloakClient.GetLdapUserAttributeMapper(realm, id)
			if ldapUserAttributeMapper != nil {
				return fmt.Errorf("ldap user attribute mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapUserAttributeMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapUserAttributeMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapUserAttributeMapper, err := keycloakClient.GetLdapUserAttributeMapper(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap user attribute mapper with id %s: %s", id, err)
	}

	return ldapUserAttributeMapper, nil
}

func testKeycloakLdapUserAttributeMapper_basic(realm, userAttributeMapperName string) string {
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

resource "keycloak_ldap_user_attribute_mapper" "username" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	user_model_attribute        = "username"
	ldap_attribute              = "cn"
}
	`, realm, userAttributeMapperName)
}

func testKeycloakLdapUserAttributeMapper_basicFromInterface(realm string, mapper *keycloak.LdapUserAttributeMapper) string {
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

resource "keycloak_ldap_user_attribute_mapper" "username" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	user_model_attribute        = "%s"
	ldap_attribute              = "%s"

	read_only                   = %t
	always_read_value_from_ldap = %t
	is_mandatory_in_ldap        = %t
}
	`, realm, mapper.Name, mapper.UserModelAttribute, mapper.LdapAttribute, mapper.ReadOnly, mapper.AlwaysReadValueFromLdap, mapper.IsMandatoryInLdap)
}

func testKeycloakLdapUserAttributeMapper_updateLdapUserFederationBefore(realmOne, realmTwo, userAttributeMapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_one" {
	realm = "%s"
}

resource "keycloak_realm" "realm_two" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap_one" {
	name                    = "openldap"
	realm_id                = "${keycloak_realm.realm_one.id}"

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

resource "keycloak_ldap_user_federation" "openldap_two" {
	name                    = "openldap"
	realm_id                = "${keycloak_realm.realm_two.id}"

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

resource "keycloak_ldap_user_attribute_mapper" "username" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm_one.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap_one.id}"

	user_model_attribute        = "username"
	ldap_attribute              = "cn"
}
	`, realmOne, realmTwo, userAttributeMapperName)
}

func testKeycloakLdapUserAttributeMapper_updateLdapUserFederationAfter(realmOne, realmTwo, userAttributeMapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_one" {
	realm = "%s"
}

resource "keycloak_realm" "realm_two" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap_one" {
	name                    = "openldap"
	realm_id                = "${keycloak_realm.realm_one.id}"

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

resource "keycloak_ldap_user_federation" "openldap_two" {
	name                    = "openldap"
	realm_id                = "${keycloak_realm.realm_two.id}"

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

resource "keycloak_ldap_user_attribute_mapper" "username" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm_two.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap_two.id}"

	user_model_attribute        = "username"
	ldap_attribute              = "cn"
}
	`, realmOne, realmTwo, userAttributeMapperName)
}
