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

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper_client"
	clientScopeResourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdHardcodedClaimProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdHardcodedClaimProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	updatedClaimName := "claim-name-update-" + acctest.RandString(10)
	claimValue := "claim-value-" + acctest.RandString(10)
	updatedClaimValue := "claim-value-update-" + acctest.RandString(10)

	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(realmName, clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(realmName, clientId, mapperName, updatedClaimName, updatedClaimValue),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdHardcodedClaimProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdHardcodedClaimProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdHardcodedClaimProtocolMapper_validateClaimValueType(realmName, mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	claimValue := "claim-value-" + acctest.RandString(10)
	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(realmName, clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(realmName, updatedClientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdHardcodedClaimProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-hardcoded-claim-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	claimValue := "claim-value-" + acctest.RandString(10)
	resourceName := "keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(realmName, clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(newRealmName, clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdHardcodedClaimProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_hardcoded_claim_protocol_mapper" {
				continue
			}

			mapper, _ := getHardcodedClaimMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdHardcodedClaimProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getHardcodedClaimMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdHardcodedClaimProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdHardcodedClaimProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getHardcodedClaimMapperUsingState(state, resourceName)
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

func getHardcodedClaimMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdHardcodedClaimProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdHardcodedClaimProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdHardcodedClaimProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdHardcodedClaimProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdHardcodedClaimProtocolMapper_claimNameAndValue(realmName, clientId, mapperName, claimName, claimValue string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "${keycloak_openid_client.openid_client.id}"

	claim_name  = "%s"
	claim_value = "%s"
}`, realmName, clientId, mapperName, claimName, claimValue)
}

func testKeycloakOpenIdHardcodedClaimProtocolMapper_validateClaimValueType(realmName, mapperName, claimValueType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "openid-client"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper_validation" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_value      = "foo"
	claim_name       = "bar"
	claim_value_type = "%s"
}`, realmName, mapperName, claimValueType)
}
