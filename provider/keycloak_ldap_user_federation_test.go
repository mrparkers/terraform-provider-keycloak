package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakLdapUserFederation_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basic(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_editModeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	validEditModes := []string{"READ_ONLY", "WRITABLE", "UNSYNCED"}
	editMode := randomStringInSlice(validEditModes)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("edit_mode", realmName, ldapName, editMode),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "edit_mode", editMode),
			},
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("edit_mode", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected edit_mode to be one of .+ got .+"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_vendorValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	validVendors := []string{"OTHER", "EDIRECTORY", "AD", "RHDS", "TIVOLI"}
	vendor := randomStringInSlice(validVendors)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("vendor", realmName, ldapName, vendor),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "vendor", vendor),
			},
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("vendor", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected vendor to be one of .+ got .+"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_searchScopeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	validSearchScopes := []string{"ONE_LEVEL", "SUBTREE"}
	searchScope := randomStringInSlice(validSearchScopes)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("search_scope", realmName, ldapName, searchScope),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "search_scope", searchScope),
			},
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("search_scope", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected search_scope to be one of .+ got .+"),
			},
		},
	})
}

func testAccCheckKeycloakLdapUserFederationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapUserFederationFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapUserFederationDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_user_federation" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldap, _ := keycloakClient.GetLdapUserFederation(realm, id)
			if ldap != nil {
				return fmt.Errorf("ldap config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapUserFederationFromState(s *terraform.State, resourceName string) (*keycloak.LdapUserFederation, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldap, err := keycloakClient.GetLdapUserFederation(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap config with id %s: %s", id, err)
	}

	return ldap, nil
}

func testKeycloakLdapUserFederation_basic(realm, ldap string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
  name                    = "%s"
  realm_id                = "master"

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
	`, realm, ldap)
}

func testKeycloakLdapUserFederation_basicWithAttrValidation(attr, realm, ldap, val string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
  name                    = "%s"
  realm_id                = "master"

  enabled                 = true

  %s                      = "%s"

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
	`, realm, ldap, attr, val)
}
