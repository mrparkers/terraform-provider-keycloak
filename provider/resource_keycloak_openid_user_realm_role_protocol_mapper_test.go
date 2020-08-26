package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client"
	clientScopeResourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdUserRealmRoleProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdUserRealmRoleProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	updatedClaimName := "claim-name-update-" + acctest.RandString(10)

	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(realmName, clientId, mapperName, updatedClaimName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdUserRealmRoleProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdUserRealmRoleProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserRealmRoleProtocolMapper_validateClaimValueType(realmName, mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(realmName, updatedClientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserRealmRoleProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-realm-role-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_realm_role_protocol_mapper.user_realm_role_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(newRealmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdUserRealmRoleProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_user_realm_role_protocol_mapper" {
				continue
			}

			mapper, _ := getUserRealmRoleMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdUserRealmRoleProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getUserRealmRoleMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdUserRealmRoleProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdUserRealmRoleProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getUserRealmRoleMapperUsingState(state, resourceName)
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

func getUserRealmRoleMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdUserRealmRoleProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdUserRealmRoleProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

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

func testKeycloakOpenIdUserRealmRoleProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserRealmRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "%s"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, claimName)
}

func testKeycloakOpenIdUserRealmRoleProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
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
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserRealmRoleProtocolMapper_validateClaimValueType(realmName, mapperName, claimValueType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "openid-client"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "user_realm_role_mapper_validation" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name      = "foo"
	claim_value_type = "%s"
}`, realmName, mapperName, claimValueType)
}
