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

	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user-attribute-mapper-client"

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

	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user-attribute-mapper-client-scope"

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

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(5)

	attributeName := "claim-" + acctest.RandString(10)
	updatedAttributeName := "claim-update-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user-attribute-mapper"

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

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_validateClientOrClientScopeSet(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserAttributeProtocolMapper_validation(realmName, mapperName),
				ExpectError: regexp.MustCompile("one of client_id or client_scope_id must be set"),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserAttributeProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-attribute-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)
	config := fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user-attribute-mapper-validation" {
  	name = "%s"
	realm_id = "${keycloak_realm.realm.id}"
  	user_attribute = "foo"
  	claim_name = "bar"
	claim_value_type = "%s"
}
`, realmName, mapperName, invalidClaimValueType)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(" expected claim_value_type to be one of .+ got " + invalidClaimValueType),
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
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user-attribute-mapper"

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
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user-attribute-mapper-client-scope"

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
	resourceName := "keycloak_openid_user_attribute_protocol_mapper.user-attribute-mapper"

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

	if clientId != "" {
		return keycloakClient.GetOpenIdUserAttributeProtocolMapperForClient(realm, clientId, id)
	} else {
		return keycloakClient.GetOpenIdUserAttributeProtocolMapperForClientScope(realm, clientScopeId, id)
	}
}

func testKeycloakOpenIdUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid-client" {
	realm_id = "${keycloak_realm.realm.id}"
	client_id = "%s"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user-attribute-mapper-client" {
  	name = "%s"
	realm_id = "${keycloak_realm.realm.id}"
  	client_id = "${keycloak_openid_client.openid-client.id}"
  	user_attribute = "foo"
  	claim_name = "bar"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client-scope" {
  name                = "%s"
  realm_id            = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user-attribute-mapper-client-scope" {
  	name = "%s"
	realm_id = "${keycloak_realm.realm.id}"
  	client_scope_id = "${keycloak_openid_client_scope.client-scope.id}"
  	user_attribute = "foo"
  	claim_name = "bar"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_validation(realmName, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user-attribute-mapper-validation" {
  	name = "%s"
	realm_id = "${keycloak_realm.realm.id}"
  	user_attribute = "foo"
  	claim_name = "bar"
}
`, realmName, mapperName)
}

func testKeycloakOpenIdUserAttributeProtocolMapper_claim(realmName, clientId, mapperName, attributeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid-client" {
	realm_id = "${keycloak_realm.realm.id}"
	client_id = "%s"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "user-attribute-mapper" {
  	name = "%s"
	realm_id = "${keycloak_realm.realm.id}"
  	client_id = "${keycloak_openid_client.openid-client.id}"
  	user_attribute = "%s"
  	claim_name = "bar"
}`, realmName, clientId, mapperName, attributeName)
}
