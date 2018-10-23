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

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-group-membership-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_group_membership_protocol_mapper.group-membership-mapper-client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{

				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-group-membership-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_group_membership_protocol_mapper.group-membership-mapper-client-scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{

				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_update(t *testing.T) {
	resourceName := "keycloak_openid_group_membership_protocol_mapper.group-membership-mapper"

	mapperOne := &keycloak.OpenIdGroupMembershipProtocolMapper{
		Name:               acctest.RandString(10),
		RealmId:            "terraform-realm-" + acctest.RandString(10),
		ClientId:           "terraform-client-" + acctest.RandString(10),
		ClaimName:          acctest.RandString(10),
		FullPath:           randomBool(),
		IdTokenClaim:       randomBool(),
		AccessTokenClaim:   randomBool(),
		UserinfoTokenClaim: randomBool(),
	}

	mapperTwo := &keycloak.OpenIdGroupMembershipProtocolMapper{
		Name:               mapperOne.Name,
		RealmId:            mapperOne.RealmId,
		ClientId:           mapperOne.ClientId,
		ClaimName:          acctest.RandString(10),
		FullPath:           randomBool(),
		IdTokenClaim:       randomBool(),
		AccessTokenClaim:   randomBool(),
		UserinfoTokenClaim: randomBool(),
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_fromInterface(mapperOne),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_fromInterface(mapperTwo),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-group-membership-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdGroupMembershipProtocolMapper_validation(realmName, mapperName),
				ExpectError: regexp.MustCompile("validation error: one of ClientId or ClientScopeId must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_updateMapperNameForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperNameOne := acctest.RandString(10)
	mapperNameTwo := acctest.RandString(10)

	resourceName := "keycloak_openid_group_membership_protocol_mapper.group-membership-mapper-client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_basic_client(realmName, clientId, mapperNameOne),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_basic_client(realmName, clientId, mapperNameTwo),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientIdOne := "terraform-client-" + acctest.RandString(10)
	clientIdTwo := "terraform-client-" + acctest.RandString(10)

	resourceName := "keycloak_openid_group_membership_protocol_mapper.group-membership-mapper-client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_updateClientForceNew(realmName, clientIdOne, clientIdTwo, "openid-client-one"),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_updateClientForceNew(realmName, clientIdOne, clientIdTwo, "openid-client-two"),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdGroupMembershipProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeOne := "terraform-client-scope-" + acctest.RandString(10)
	clientScopeTwo := "terraform-client-scope-" + acctest.RandString(10)

	resourceName := "keycloak_openid_group_membership_protocol_mapper.group-membership-mapper-client-scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_updateClientScopeForceNew(realmName, clientScopeOne, clientScopeTwo, "client-scope-one"),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdGroupMembershipProtocolMapper_updateClientScopeForceNew(realmName, clientScopeOne, clientScopeTwo, "client-scope-two"),
				Check:  testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdGroupMembershipProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_group_membership_protocol_mapper" {
				continue
			}

			mapper, _ := getGroupMembershipMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdGroupMembershipProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getGroupMembershipMapperUsingState(state, resourceName)

		if err != nil {
			return err
		}

		return nil
	}
}

func getGroupMembershipMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdGroupMembershipProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdGroupMembershipProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid-client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group-membership-mapper-client" {
  	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
  	client_id        = "${keycloak_openid_client.openid-client.id}"
  	claim_name       = "bar"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client-scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group-membership-mapper-client-scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client-scope.id}"
	claim_name      = "bar"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_fromInterface(mapper *keycloak.OpenIdGroupMembershipProtocolMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid-client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group-membership-mapper" {
  	name                 = "%s"
	realm_id             = "${keycloak_realm.realm.id}"
  	client_id            = "${keycloak_openid_client.openid-client.id}"

  	claim_name           = "%s"
	full_path            = %t
	id_token_claim       = %t
	access_token_claim   = %t
	userinfo_token_claim = %t
}`, mapper.RealmId, mapper.ClientId, mapper.Name, mapper.ClaimName, mapper.FullPath, mapper.IdTokenClaim, mapper.AccessTokenClaim, mapper.UserinfoTokenClaim)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_updateClientForceNew(realmId, clientIdOne, clientIdTwo, currentClient string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid-client-one" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_client" "openid-client-two" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group-membership-mapper-client" {
  	name                 = "group-mapper"
	realm_id             = "${keycloak_realm.realm.id}"
  	client_id            = "${keycloak_openid_client.%s.id}"

  	claim_name           = "foo"
}`, realmId, clientIdOne, clientIdTwo, currentClient)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_updateClientScopeForceNew(realmId, clientScopeIdOne, clientScopeIdTwo, currentClientScope string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client-scope-one" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client_scope" "client-scope-two" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group-membership-mapper-client-scope" {
  	name                 = "group-mapper"
	realm_id             = "${keycloak_realm.realm.id}"
  	client_scope_id      = "${keycloak_openid_client_scope.%s.id}"

  	claim_name           = "foo"
}`, realmId, clientScopeIdOne, clientScopeIdTwo, currentClientScope)
}

func testKeycloakOpenIdGroupMembershipProtocolMapper_validation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_group_membership_protocol_mapper" "group-membership-mapper-validation" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	claim_name     = "bar"
}`, realmName, mapperName)
}
