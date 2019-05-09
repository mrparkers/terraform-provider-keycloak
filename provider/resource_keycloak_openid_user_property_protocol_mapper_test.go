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

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client"
	clientScopeResourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdUserPropertyProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdUserPropertyProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	propertyName := "claim-" + acctest.RandString(10)
	updatedPropertyName := "claim-update-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_claim(realmName, clientId, mapperName, propertyName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_claim(realmName, clientId, mapperName, updatedPropertyName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdUserPropertyProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdUserPropertyProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserPropertyProtocolMapper_claimValueType(realmName, mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	propertyName := "claim-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_claim(realmName, clientId, mapperName, propertyName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_claim(realmName, updatedClientId, mapperName, propertyName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserPropertyProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-property-mapper-" + acctest.RandString(5)

	propertyName := "claim-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_property_protocol_mapper.user_property_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_claim(realmName, clientId, mapperName, propertyName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserPropertyProtocolMapper_claim(newRealmName, clientId, mapperName, propertyName),
				Check:  testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdUserPropertyProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_user_property_protocol_mapper" {
				continue
			}

			mapper, _ := getUserPropertyMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user property protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdUserPropertyProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getUserPropertyMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdUserPropertyProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdUserPropertyProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getUserPropertyMapperUsingState(state, resourceName)
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

func getUserPropertyMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdUserPropertyProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdUserPropertyProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdUserPropertyProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
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
	user_property = "foo"
	claim_name    = "bar"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdUserPropertyProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
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
	user_property   = "foo"
	claim_name      = "bar"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserPropertyProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
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
	user_property = "foo"
	claim_name    = "bar"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_property   = "foo"
	claim_name      = "bar"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserPropertyProtocolMapper_claim(realmName, clientId, mapperName, propertyName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper" {
	name          = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
	client_id     = "${keycloak_openid_client.openid_client.id}"
	user_property = "%s"
	claim_name    = "bar"
}`, realmName, clientId, mapperName, propertyName)
}

func testKeycloakOpenIdUserPropertyProtocolMapper_claimValueType(realmName, mapperName, claimValueType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper_validation" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	user_property    = "foo"
	claim_name       = "bar"
	claim_value_type = "%s"
}`, realmName, mapperName, claimValueType)
}
