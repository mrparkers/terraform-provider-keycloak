package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakLdapMsadUserAccountControlMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	msadUacMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapMsadUserAccountControlMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapMsadUserAccountControlMapper_basic(realmName, msadUacMapperName, randomBool()),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperExists("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
			{
				ResourceName:      "keycloak_ldap_msad_user_account_control_mapper.uac_mapper",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapGenericMapperImportId("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapMsadUserAccountControlMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.LdapMsadUserAccountControlMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	msadUacMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapMsadUserAccountControlMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapMsadUserAccountControlMapper_basic(realmName, msadUacMapperName, randomBool()),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperFetch("keycloak_ldap_msad_user_account_control_mapper.uac_mapper", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteLdapMsadUserAccountControlMapper(mapper.RealmId, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapMsadUserAccountControlMapper_basic(realmName, msadUacMapperName, randomBool()),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperExists("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapMsadUserAccountControlMapper_updateLdapUserFederation(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	msadUacMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapMsadUserAccountControlMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapMsadUserAccountControlMapper_updateLdapUserFederationBefore(realmOne, realmTwo, msadUacMapperName),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperExists("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
			{
				Config: testKeycloakLdapMsadUserAccountControlMapper_updateLdapUserFederationAfter(realmOne, realmTwo, msadUacMapperName),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperExists("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapMsadUserAccountControlMapper_updateInPlace(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	passwordHintsEnabled := randomBool()

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapMsadUserAccountControlMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapMsadUserAccountControlMapper_basic(realm, acctest.RandString(10), passwordHintsEnabled),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperExists("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
			{
				Config: testKeycloakLdapMsadUserAccountControlMapper_basic(realm, acctest.RandString(10), !passwordHintsEnabled),
				Check:  testAccCheckKeycloakLdapMsadUserAccountControlMapperExists("keycloak_ldap_msad_user_account_control_mapper.uac_mapper"),
			},
		},
	})
}

func testAccCheckKeycloakLdapMsadUserAccountControlMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapMsadUserAccountControlMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapMsadUserAccountControlMapperFetch(resourceName string, mapper *keycloak.LdapMsadUserAccountControlMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getLdapMsadUserAccountControlMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckKeycloakLdapMsadUserAccountControlMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_msad_user_account_control_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldapMsadUserAccountControlMapper, _ := keycloakClient.GetLdapMsadUserAccountControlMapper(realm, id)
			if ldapMsadUserAccountControlMapper != nil {
				return fmt.Errorf("ldap msad uac mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapMsadUserAccountControlMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapMsadUserAccountControlMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapMsadUserAccountControlMapper, err := keycloakClient.GetLdapMsadUserAccountControlMapper(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap msad uac mapper with id %s: %s", id, err)
	}

	return ldapMsadUserAccountControlMapper, nil
}

func testKeycloakLdapMsadUserAccountControlMapper_basic(realm, msadUacMapperName string, passwordHintsEnabled bool) string {
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

resource "keycloak_ldap_msad_user_account_control_mapper" "uac_mapper" {
	name                               = "%s"
	realm_id                           = "${keycloak_realm.realm.id}"
	ldap_user_federation_id            = "${keycloak_ldap_user_federation.openldap.id}"

	ldap_password_policy_hints_enabled = %t
}
	`, realm, msadUacMapperName, passwordHintsEnabled)
}

func testKeycloakLdapMsadUserAccountControlMapper_updateLdapUserFederationBefore(realmOne, realmTwo, msadUacMapperName string) string {
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

resource "keycloak_ldap_msad_user_account_control_mapper" "uac_mapper" {
	name                               = "%s"
	realm_id                           = "${keycloak_realm.realm_one.id}"
	ldap_user_federation_id            = "${keycloak_ldap_user_federation.openldap_one.id}"
}
	`, realmOne, realmTwo, msadUacMapperName)
}

func testKeycloakLdapMsadUserAccountControlMapper_updateLdapUserFederationAfter(realmOne, realmTwo, msadUacMapperName string) string {
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

resource "keycloak_ldap_msad_user_account_control_mapper" "uac_mapper" {
	name                               = "%s"
	realm_id                           = "${keycloak_realm.realm_two.id}"
	ldap_user_federation_id            = "${keycloak_ldap_user_federation.openldap_two.id}"
}
	`, realmOne, realmTwo, msadUacMapperName)
}
