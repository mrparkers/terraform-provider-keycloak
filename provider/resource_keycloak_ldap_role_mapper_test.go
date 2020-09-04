package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakLdapRoleMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapRoleMapper_basic(realmName, roleMapperName),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
			{
				ResourceName:      "keycloak_ldap_role_mapper.role_mapper",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapGenericMapperImportId("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.LdapRoleMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapRoleMapper_basic(realmName, roleMapperName),
				Check:  testAccCheckKeycloakLdapRoleMapperFetch("keycloak_ldap_role_mapper.role_mapper", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteLdapRoleMapper(mapper.RealmId, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapRoleMapper_basic(realmName, roleMapperName),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_modeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)
	mode := randomStringInSlice(keycloakLdapRoleMapperModes)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "mode", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected mode to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "mode", mode),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_membershipAttributeTypeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)
	membershipAttributeType := randomStringInSlice(keycloakLdapRoleMapperMembershipAttributeTypes)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "membership_attribute_type", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected membership_attribute_type to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "membership_attribute_type", membershipAttributeType),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_userRolesRetrieveStrategyValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)
	userRolesRetrieveStrategy := randomStringInSlice(keycloakLdapRoleMapperUserRolesRetrieveStrategies)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "user_roles_retrieve_strategy", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected user_roles_retrieve_strategy to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "user_roles_retrieve_strategy", userRolesRetrieveStrategy),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_rolesLdapFilterValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)
	rolesLdapFilter := "(" + acctest.RandString(10) + ")"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "roles_ldap_filter", acctest.RandString(10)),
				ExpectError: regexp.MustCompile(`validation error: roles ldap filter must start with '\(' and end with '\)'`),
			},
			{
				Config: testKeycloakLdapRoleMapper_basicWithAttrValidation(realmName, roleMapperName, "roles_ldap_filter", rolesLdapFilter),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_updateLdapUserFederationForceNew(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	roleMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapRoleMapper_updateLdapUserFederationBefore(realmOne, realmTwo, roleMapperName),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
			{
				Config: testKeycloakLdapRoleMapper_updateLdapUserFederationAfter(realmOne, realmTwo, roleMapperName),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapRoleMapper_updateLdapUserFederationInPlace(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	useRealmRolesMapping := randomBool()

	roleMapperOne := &keycloak.LdapRoleMapper{
		Name:                        acctest.RandString(10),
		RealmId:                     realm,
		LdapRolesDn:                 acctest.RandString(10),
		RoleNameLdapAttribute:       acctest.RandString(10),
		RoleObjectClasses:           []string{acctest.RandString(10), acctest.RandString(10)},
		MembershipLdapAttribute:     acctest.RandString(10),
		MembershipAttributeType:     "DN",
		MembershipUserLdapAttribute: acctest.RandString(10),
		RolesLdapFilter:             "(" + acctest.RandString(10) + ")",
		Mode:                        randomStringInSlice(keycloakLdapRoleMapperModes),
		UserRolesRetrieveStrategy:   randomStringInSlice(keycloakLdapRoleMapperUserRolesRetrieveStrategies),
		MemberofLdapAttribute:       acctest.RandString(10),
		UseRealmRolesMapping:        useRealmRolesMapping,
		ClientId:                    clientId,
	}

	roleMapperTwo := &keycloak.LdapRoleMapper{
		Name:                        acctest.RandString(10),
		RealmId:                     realm,
		LdapRolesDn:                 acctest.RandString(10),
		RoleNameLdapAttribute:       acctest.RandString(10),
		RoleObjectClasses:           []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
		MembershipLdapAttribute:     acctest.RandString(10),
		MembershipAttributeType:     randomStringInSlice(keycloakLdapRoleMapperMembershipAttributeTypes),
		MembershipUserLdapAttribute: acctest.RandString(10),
		RolesLdapFilter:             "(" + acctest.RandString(10) + ")",
		Mode:                        randomStringInSlice(keycloakLdapRoleMapperModes),
		UserRolesRetrieveStrategy:   randomStringInSlice(keycloakLdapRoleMapperUserRolesRetrieveStrategies),
		MemberofLdapAttribute:       acctest.RandString(10),
		UseRealmRolesMapping:        !useRealmRolesMapping,
		ClientId:                    clientId,
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapRoleMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapRoleMapper_basicFromInterface(realm, roleMapperOne),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
			{
				Config: testKeycloakLdapRoleMapper_basicFromInterface(realm, roleMapperTwo),
				Check:  testAccCheckKeycloakLdapRoleMapperExists("keycloak_ldap_role_mapper.role_mapper"),
			},
		},
	})
}

func testAccCheckKeycloakLdapRoleMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapRoleMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapRoleMapperFetch(resourceName string, mapper *keycloak.LdapRoleMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getLdapRoleMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckKeycloakLdapRoleMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_role_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldapRoleMapper, _ := keycloakClient.GetLdapRoleMapper(realm, id)
			if ldapRoleMapper != nil {
				return fmt.Errorf("ldap role mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapRoleMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapRoleMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapRoleMapper, err := keycloakClient.GetLdapRoleMapper(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap role mapper with id %s: %s", id, err)
	}

	return ldapRoleMapper, nil
}

func testKeycloakLdapRoleMapper_basic(realm, roleMapperName string) string {
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

resource "keycloak_ldap_role_mapper" "role_mapper" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"
	ldap_user_federation_id = "${keycloak_ldap_user_federation.openldap.id}"

	ldap_roles_dn                 = "dc=example,dc=org"
        role_name_ldap_attribute      = "cn"
        role_object_classes           = [
                "group"
        ]
        membership_ldap_attribute      = "member"
        membership_user_ldap_attribute = "sAMAccountName"
        memberof_ldap_attribute        = "memberOf"
}
	`, realm, roleMapperName)
}

func testKeycloakLdapRoleMapper_basicWithAttrValidation(realm, roleMapperName, attr, val string) string {
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

resource "keycloak_ldap_role_mapper" "role_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	%s                          = "%s"

	ldap_roles_dn                 = "dc=example,dc=org"
	role_name_ldap_attribute      = "cn"
	role_object_classes           = [
		"group"
	]
	membership_ldap_attribute      = "member"
	membership_user_ldap_attribute = "sAMAccountName"
	memberof_ldap_attribute        = "memberOf"
}
	`, realm, roleMapperName, attr, val)
}

func testKeycloakLdapRoleMapper_basicFromInterface(realm string, mapper *keycloak.LdapRoleMapper) string {
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

resource "keycloak_ldap_role_mapper" "role_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	ldap_roles_dn                  = "%s"
	role_name_ldap_attribute       = "%s"
	role_object_classes            = %s
	membership_ldap_attribute      = "%s"
	membership_attribute_type      = "%s"
	membership_user_ldap_attribute = "%s"
	roles_ldap_filter              = "%s"
	mode                           = "%s"
	user_roles_retrieve_strategy   = "%s"
	memberof_ldap_attribute        = "%s"
	use_realm_roles_mapping        = %t
	client_id                      = "%s"
}
	`, realm, mapper.Name, mapper.LdapRolesDn, mapper.RoleNameLdapAttribute, arrayOfStringsForTerraformResource(mapper.RoleObjectClasses), mapper.MembershipLdapAttribute, mapper.MembershipAttributeType, mapper.MembershipUserLdapAttribute, mapper.RolesLdapFilter, mapper.Mode, mapper.UserRolesRetrieveStrategy, mapper.MemberofLdapAttribute, mapper.UseRealmRolesMapping, mapper.ClientId)
}

func testKeycloakLdapRoleMapper_updateLdapUserFederationBefore(realmOne, realmTwo, roleMapperName string) string {
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

resource "keycloak_ldap_role_mapper" "role_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm_one.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap_one.id}"

	ldap_roles_dn                 = "dc=example,dc=org"
	role_name_ldap_attribute      = "cn"
	role_object_classes           = [
		"group"
	]
	membership_attribute_type      = "DN"
	membership_ldap_attribute      = "member"
	membership_user_ldap_attribute = "sAMAccountName"
	memberof_ldap_attribute        = "memberOf"
}
	`, realmOne, realmTwo, roleMapperName)
}

func testKeycloakLdapRoleMapper_updateLdapUserFederationAfter(realmOne, realmTwo, roleMapperName string) string {
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

resource "keycloak_ldap_role_mapper" "role_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm_two.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap_two.id}"

	ldap_roles_dn                 = "dc=example,dc=org"
	role_name_ldap_attribute      = "cn"
	role_object_classes           = [
		"group"
	]
	membership_attribute_type      = "DN"
	membership_ldap_attribute      = "member"
	membership_user_ldap_attribute = "sAMAccountName"
	memberof_ldap_attribute        = "memberOf"
}
	`, realmOne, realmTwo, roleMapperName)
}
