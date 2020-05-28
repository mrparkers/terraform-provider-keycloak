package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

/*
 * Keycloak does not allow two protocol mappers to exist with the same name, but if you create two of them with the same
 * name quickly enough (which Terraform is pretty good at), the API allows it.
 *
 * The following tests check to see if that error is caught by creating a single protocol mapper first (each test uses a
 * different one), then creating another mapper with the same name in the next test step.
 */

func TestAccKeycloakOpenIdFullNameProtocolMapper_clientDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	groupMembershipProtocolMapperResourceName := "keycloak_openid_group_membership_protocol_mapper.group_membership_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientGroupMembershipMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(groupMembershipProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientFullNameAndGroupMembershipMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	groupMembershipProtocolMapperResourceName := "keycloak_openid_group_membership_protocol_mapper.group_membership_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeGroupMembershipMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(groupMembershipProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeFullNameAndGroupMembershipMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_clientDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	fullNameProtocolMapperResourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientFullNameMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(fullNameProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientFullNameAndGroupMembershipMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	fullNameProtocolMapperResourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeFullNameMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(fullNameProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeFullNameAndGroupMembershipMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_clientDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	groupMembershipProtocolMapperResourceName := "keycloak_openid_group_membership_protocol_mapper.group_membership_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientGroupMembershipMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(groupMembershipProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientGroupMembershipAndUserAttributeMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	groupMembershipProtocolMapperResourceName := "keycloak_openid_group_membership_protocol_mapper.group_membership_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeGroupMembershipMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(groupMembershipProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeGroupMembershipAndUserAttributeMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_clientDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	userAttributeProtocolMapperResourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientUserAttributeMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(userAttributeProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientUserAttributeAndUserPropertyMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	userAttributeProtocolMapperResourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeUserAttributeMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(userAttributeProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeUserAttributeAndUserPropertyMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_clientDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	userPropertyProtocolMapperResourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientUserPropertyMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(userPropertyProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientUserPropertyAndHardcodedClaimMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-protocol-mapper-" + acctest.RandString(5)

	userPropertyProtocolMapperResourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeUserPropertyMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(userPropertyProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeUserPropertyAndHardcodedClaimMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-user-realm-role-mapper-" + acctest.RandString(5)

	fullNameProtocolMapperResourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeUserRealmRoleMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(fullNameProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeUserRealmRoleAndHardcodedClaimMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_clientScopeDuplicateNameValidation(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-user-client-role-mapper-" + acctest.RandString(5)

	userClientRoleProtocolMapperResourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGenericProtocolMapperValidation_clientScopeUserClientRoleMapper(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(userClientRoleProtocolMapperResourceName),
			},
			{
				Config:      testGenericProtocolMapperValidation_clientScopeUserClientRoleAndHardcodedClaimMapper(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: a protocol mapper with name .+ already exists for this client"),
			},
		},
	})
}

/*
 * Protocol mappers must be attached to either a client or client scope.  The following tests assert that errors are raised
 * if neither are specified.
 */

func TestAccKeycloakOpenIdFullNameProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-full-name-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdFullNameProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-group-membership-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdGroupMembershipProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserAttributeProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserPropertyProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdHardcodedClaimProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserRealmRoleProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserClientRoleProtocolMapper_parentResourceValidation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func testGenericProtocolMapperValidation_clientGroupMembershipMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_client" {
  	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
  	client_id       = "${keycloak_openid_client.openid_client.id}"

  	claim_name      = "foo"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientFullNameMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client" {
  	name       = "%s"
	realm_id   = "${keycloak_realm.realm.id}"
  	client_id  = "${keycloak_openid_client.openid_client.id}"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientUserAttributeMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client" {
  	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
  	client_id      = "${keycloak_openid_client.openid_client.id}"
  	user_attribute = "foo-attribute"
  	claim_name     = "bar-attribute"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientUserPropertyMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client" {
  	name          = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
  	client_id     = "${keycloak_openid_client.openid_client.id}"
  	user_property = "foo-property"
  	claim_name    = "bar-property"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientFullNameAndGroupMembershipMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_client" {
  	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
  	client_id       = "${keycloak_openid_client.openid_client.id}"

  	claim_name      = "foo"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client" {
  	name       = "%s"
	realm_id   = "${keycloak_realm.realm.id}"
  	client_id  = "${keycloak_openid_client.openid_client.id}"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientGroupMembershipAndUserAttributeMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_client" {
  	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
  	client_id       = "${keycloak_openid_client.openid_client.id}"

  	claim_name      = "foo"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client" {
  	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
  	client_id      = "${keycloak_openid_client.openid_client.id}"
  	user_attribute = "foo"
  	claim_name     = "bar"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientUserAttributeAndUserPropertyMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client" {
  	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
  	client_id      = "${keycloak_openid_client.openid_client.id}"
  	user_attribute = "foo-attribute"
  	claim_name     = "bar-attribute"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client" {
  	name          = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
  	client_id     = "${keycloak_openid_client.openid_client.id}"
  	user_property = "foo-property"
  	claim_name    = "bar-property"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientUserPropertyAndHardcodedClaimMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client" {
  	name          = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
  	client_id     = "${keycloak_openid_client.openid_client.id}"
  	user_property = "foo-property"
  	claim_name    = "bar-property"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeGroupMembershipMapper(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar"
}`, realmName, clientScopeId, mapperName)
}

func testGenericProtocolMapperValidation_clientUserRealmRoleMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientUserClientRoleMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeFullNameMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserAttributeMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_attribute  = "foo-attribute"
	claim_name      = "bar-attribute"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserPropertyMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_property   = "foo-property"
	claim_name      = "bar-property"
}`, realmName, clientId, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserRealmRoleMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar-property"
}`, realmName, clientId, mapperName)

}

func testGenericProtocolMapperValidation_clientScopeUserClientRoleMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar-property"
}`, realmName, clientId, mapperName)

}

func testGenericProtocolMapperValidation_clientScopeFullNameAndGroupMembershipMapper(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
}`, realmName, clientScopeId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeGroupMembershipAndUserAttributeMapper(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_attribute  = "foo"
	claim_name      = "bar"
}`, realmName, clientScopeId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserAttributeAndUserPropertyMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_attribute  = "foo-attribute"
	claim_name      = "bar-attribute"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_property   = "foo-property"
	claim_name      = "bar-property"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserPropertyAndHardcodedClaimMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_property   = "foo-property"
	claim_name      = "bar-property"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserRealmRoleAndHardcodedClaimMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar-property"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, mapperName)
}

func testGenericProtocolMapperValidation_clientScopeUserClientRoleAndHardcodedClaimMapper(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name      = "bar-property"
}
resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, mapperName)
}

func testKeycloakOpenIdFullNameProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_validation" {
	name       = "%s"
	realm_id   = "${keycloak_realm.realm.id}"
}`, realmName, mapperName)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group_membership_mapper_validation" {
	name       = "%s"
	realm_id   = "${keycloak_realm.realm.id}"
	claim_name = "bar"
}`, realmName, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_validation" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	user_attribute = "foo"
	claim_name     = "bar"
}`, realmName, mapperName)
}

func testKeycloakOpenIdUserPropertyProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_validation" {
	name          = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
	user_property = "foo"
	claim_name    = "bar"
}`, realmName, mapperName)
}

func testKeycloakOpenIdHardcodedClaimProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, mapperName)
}

func testKeycloakOpenIdUserRealmRoleProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"

	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, mapperName)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_parentResourceValidation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, mapperName)
}
