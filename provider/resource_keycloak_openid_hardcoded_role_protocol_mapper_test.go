package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_client(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	role := "terraform-role-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_client(realmName, role, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	role := "terraform-role-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientScope(realmName, role, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdHardcodedRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedRoleProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	role := "terraform-role-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-role-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper_client"
	clientScopeResourceName := "keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_import(realmName, role, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdHardcodedRoleProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdHardcodedRoleProtocolMapperExists(clientScopeResourceName),
				),
			},
			{
				ResourceName:      clientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(clientResourceName),
			},
			{
				ResourceName:      clientScopeResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClientScope(clientScopeResourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedRoleProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	roleOne := "terraform-role-" + acctest.RandString(10)
	roleTwo := "terraform-role-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientUpdateBefore(realmName, roleOne, roleTwo, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientUpdateAfter(realmName, roleOne, roleTwo, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedRoleProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdHardcodedRoleProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	role := "terraform-role-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_client(realmName, role, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedRoleProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdHardcodedRoleProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_client(realmName, role, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdHardcodedRoleProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_hardcoded_role_protocol_mapper" {
				continue
			}

			mapper, _ := getHardcodedRoleMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdHardcodedRoleProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getHardcodedRoleMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdHardcodedRoleProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdHardcodedRoleProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getHardcodedRoleMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.ClientId = fetchedMapper.ClientId
		mapper.ClientScopeId = fetchedMapper.ClientScopeId
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func getHardcodedRoleMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdHardcodedRoleProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdHardcodedRoleProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_client(realmName, role, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper_client" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.openid_client.id}"
	role_id        = "${keycloak_role.role.id}"
}`, realmName, role, clientId, mapperName)
}

func testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientScope(realmName, role, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	role_id         = "${keycloak_role.role.id}"
}`, realmName, role, clientScopeId, mapperName)
}

func testKeycloakOpenIdHardcodedRoleProtocolMapper_import(realmName, role, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	role_id          = "${keycloak_role.role.id}"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	role_id          = "${keycloak_role.role.id}"
}`, realmName, role, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientUpdateBefore(realmName, roleOne, roleTwo, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role_one" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "role_two" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper_client" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.openid_client.id}"
	role_id        = "${keycloak_role.role_one.id}"
}`, realmName, roleOne, roleTwo, clientId, mapperName)
}

func testKeycloakOpenIdHardcodedRoleProtocolMapper_basicRealmRole_clientUpdateAfter(realmName, roleOne, roleTwo, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role_one" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "role_two" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper_client" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.openid_client.id}"
	role_id        = "${keycloak_role.role_two.id}"
}`, realmName, roleOne, roleTwo, clientId, mapperName)
}
