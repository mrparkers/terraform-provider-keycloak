package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
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
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
			{
				ResourceName:      "keycloak_ldap_group_mapper.group_mapper",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapGenericMapperImportId("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.LdapGroupMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapGroupMapper_basic(realmName, groupMapperName),
				Check:  testAccCheckKeycloakLdapGroupMapperFetch("keycloak_ldap_group_mapper.group_mapper", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteLdapGroupMapper(mapper.RealmId, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapGroupMapper_basic(realmName, groupMapperName),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_modeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)
	mode := randomStringInSlice(keycloakLdapGroupMapperModes)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "mode", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected mode to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "mode", mode),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_membershipAttributeTypeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)
	membershipAttributeType := randomStringInSlice(keycloakLdapGroupMapperMembershipAttributeTypes)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "membership_attribute_type", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected membership_attribute_type to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "membership_attribute_type", membershipAttributeType),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_userRolesRetrieveStrategyValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)
	userRolesRetrieveStrategy := randomStringInSlice(keycloakLdapGroupMapperUserRolesRetrieveStrategies)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "user_roles_retrieve_strategy", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected user_roles_retrieve_strategy to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "user_roles_retrieve_strategy", userRolesRetrieveStrategy),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_groupsLdapFilterValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)
	groupsLdapFilter := "(" + acctest.RandString(10) + ")"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "groups_ldap_filter", acctest.RandString(10)),
				ExpectError: regexp.MustCompile(`validation error: groups ldap filter must start with '\(' and end with '\)'`),
			},
			{
				Config: testKeycloakLdapGroupMapper_basicWithAttrValidation(realmName, groupMapperName, "groups_ldap_filter", groupsLdapFilter),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_groupInheritanceValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapGroupMapper_groupInheritanceValidation(realmName, groupMapperName),
				ExpectError: regexp.MustCompile("validation error: group inheritance cannot be preserved while membership attribute type is UID"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_updateLdapUserFederationForceNew(t *testing.T) {
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
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
			{
				Config: testKeycloakLdapGroupMapper_updateLdapUserFederationAfter(realmOne, realmTwo, groupMapperName),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapGroupMapper_updateLdapUserFederationInPlace(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	preserveGroupInheritance := true
	ignoreMissingGroups := randomBool()
	dropNonExistingGroupsDuringSync := randomBool()

	groupMapperOne := &keycloak.LdapGroupMapper{
		Name:                            acctest.RandString(10),
		RealmId:                         realm,
		LdapGroupsDn:                    acctest.RandString(10),
		GroupNameLdapAttribute:          acctest.RandString(10),
		GroupObjectClasses:              []string{acctest.RandString(10), acctest.RandString(10)},
		PreserveGroupInheritance:        preserveGroupInheritance,
		IgnoreMissingGroups:             ignoreMissingGroups,
		MembershipLdapAttribute:         acctest.RandString(10),
		MembershipAttributeType:         "DN",
		MembershipUserLdapAttribute:     acctest.RandString(10),
		GroupsLdapFilter:                "(" + acctest.RandString(10) + ")",
		Mode:                            randomStringInSlice(keycloakLdapGroupMapperModes),
		UserRolesRetrieveStrategy:       randomStringInSlice(keycloakLdapGroupMapperUserRolesRetrieveStrategies),
		MemberofLdapAttribute:           acctest.RandString(10),
		MappedGroupAttributes:           []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
		DropNonExistingGroupsDuringSync: dropNonExistingGroupsDuringSync,
	}

	groupMapperTwo := &keycloak.LdapGroupMapper{
		Name:                            acctest.RandString(10),
		RealmId:                         realm,
		LdapGroupsDn:                    acctest.RandString(10),
		GroupNameLdapAttribute:          acctest.RandString(10),
		GroupObjectClasses:              []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
		PreserveGroupInheritance:        !preserveGroupInheritance,
		IgnoreMissingGroups:             !ignoreMissingGroups,
		MembershipLdapAttribute:         acctest.RandString(10),
		MembershipAttributeType:         randomStringInSlice(keycloakLdapGroupMapperMembershipAttributeTypes),
		MembershipUserLdapAttribute:     acctest.RandString(10),
		GroupsLdapFilter:                "(" + acctest.RandString(10) + ")",
		Mode:                            randomStringInSlice(keycloakLdapGroupMapperModes),
		UserRolesRetrieveStrategy:       randomStringInSlice(keycloakLdapGroupMapperUserRolesRetrieveStrategies),
		MemberofLdapAttribute:           acctest.RandString(10),
		MappedGroupAttributes:           []string{acctest.RandString(10)},
		DropNonExistingGroupsDuringSync: !dropNonExistingGroupsDuringSync,
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapGroupMapper_basicFromInterface(realm, groupMapperOne),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
			},
			{
				Config: testKeycloakLdapGroupMapper_basicFromInterface(realm, groupMapperTwo),
				Check:  testAccCheckKeycloakLdapGroupMapperExists("keycloak_ldap_group_mapper.group_mapper"),
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

func testAccCheckKeycloakLdapGroupMapperFetch(resourceName string, mapper *keycloak.LdapGroupMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getLdapGroupMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

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

resource "keycloak_ldap_group_mapper" "group_mapper" {
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

func testKeycloakLdapGroupMapper_basicWithAttrValidation(realm, groupMapperName, attr, val string) string {
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

resource "keycloak_ldap_group_mapper" "group_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	%s                          = "%s"

	ldap_groups_dn                 = "dc=example,dc=org"
	group_name_ldap_attribute      = "cn"
	group_object_classes           = [
		"groupOfNames"
	]
	membership_ldap_attribute      = "member"
	membership_user_ldap_attribute = "cn"
	memberof_ldap_attribute        = "memberOf"
}
	`, realm, groupMapperName, attr, val)
}

func testKeycloakLdapGroupMapper_groupInheritanceValidation(realm, groupMapperName string) string {
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

resource "keycloak_ldap_group_mapper" "group_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	membership_attribute_type      = "UID"
	preserve_group_inheritance     = true

	ldap_groups_dn                 = "dc=example,dc=org"
	group_name_ldap_attribute      = "cn"
	group_object_classes           = [
		"groupOfNames"
	]
	membership_ldap_attribute      = "member"
	membership_user_ldap_attribute = "cn"
	memberof_ldap_attribute        = "memberOf"
}
	`, realm, groupMapperName)
}

func testKeycloakLdapGroupMapper_basicFromInterface(realm string, mapper *keycloak.LdapGroupMapper) string {
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

resource "keycloak_ldap_group_mapper" "group_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap.id}"

	ldap_groups_dn                       = "%s"
	group_name_ldap_attribute            = "%s"
	group_object_classes                 = %s
	preserve_group_inheritance           = %t
	ignore_missing_groups                = %t
	membership_ldap_attribute            = "%s"
	membership_attribute_type            = "%s"
	membership_user_ldap_attribute       = "%s"
	groups_ldap_filter                   = "%s"
	mode                                 = "%s"
	user_roles_retrieve_strategy         = "%s"
	memberof_ldap_attribute              = "%s"
	mapped_group_attributes              = %s
	drop_non_existing_groups_during_sync = %t
}
	`, realm, mapper.Name, mapper.LdapGroupsDn, mapper.GroupNameLdapAttribute, arrayOfStringsForTerraformResource(mapper.GroupObjectClasses), mapper.PreserveGroupInheritance, mapper.IgnoreMissingGroups, mapper.MembershipLdapAttribute, mapper.MembershipAttributeType, mapper.MembershipUserLdapAttribute, mapper.GroupsLdapFilter, mapper.Mode, mapper.UserRolesRetrieveStrategy, mapper.MemberofLdapAttribute, arrayOfStringsForTerraformResource(mapper.MappedGroupAttributes), mapper.DropNonExistingGroupsDuringSync)
}

func testKeycloakLdapGroupMapper_updateLdapUserFederationBefore(realmOne, realmTwo, groupMapperName string) string {
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

resource "keycloak_ldap_group_mapper" "group_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm_one.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap_one.id}"

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

resource "keycloak_ldap_group_mapper" "group_mapper" {
	name                        = "%s"
	realm_id                    = "${keycloak_realm.realm_two.id}"
	ldap_user_federation_id     = "${keycloak_ldap_user_federation.openldap_two.id}"

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
