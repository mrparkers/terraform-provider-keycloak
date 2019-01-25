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

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client"
	clientScopeResourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdUserAttributeProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdUserAttributeProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	attributeName := "claim-" + acctest.RandString(10)
	updatedAttributeName := "claim-update-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, clientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, clientId, mapperName, updatedAttributeName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdUserAttributeProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdUserAttributeProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserAttributeProtocolMapper_claimValueType(realmName, mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	attributeName := "claim-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, clientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, updatedClientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	attributeName := "claim-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, clientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserAttributeProtocolMapper_claim(newRealmName, clientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_user_attribute_protocol_mapper" {
				continue
			}

			mapper, _ := getUserAttributeMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdUserAttributeProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getUserAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdUserAttributeProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdUserAttributeProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getUserAttributeMapperUsingState(state, resourceName)
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

func getUserAttributeMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdUserAttributeProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdUserAttributeProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.openid_client.id}"
	user_attribute = "foo"
	claim_name     = "bar"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
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
	user_attribute  = "foo"
	claim_name      = "bar"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
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
	user_attribute = "foo"
	claim_name     = "bar"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	user_attribute  = "foo"
	claim_name      = "bar"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, clientId, mapperName, attributeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper" {
	name           = "%s"
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.openid_client.id}"
	user_attribute = "%s"
	claim_name     = "bar"
}`, realmName, clientId, mapperName, attributeName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_claimValueType(realmName, mapperName, claimValueType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user_attribute_mapper_validation" {
	name              = "%s"
	realm_id          = "${keycloak_realm.realm.id}"
	user_attribute    = "foo"
	claim_name        = "bar"
	claim_value_type  = "%s"
}`, realmName, mapperName, claimValueType)
}
