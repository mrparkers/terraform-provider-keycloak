package provider

import (
	"fmt"
	"testing"

	"github.com/mrparkers/terraform-provider-keycloak/keycloak"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	TF_RESOURCE_NAME = "client_protocol_mapper"
)

func TestAccKeycloakGenericClientProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-generic-client-protocol-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_generic_client_protocol_mapper." + TF_RESOURCE_NAME

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakGenericClientProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericClientProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakGenericClientProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakGenericClientProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-generic-client-protocol-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_generic_client_protocol_mapper." + TF_RESOURCE_NAME

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakGenericClientProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericClientProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakGenericClientProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakGenericClientProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-generic-client-protocol-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_generic_client_protocol_mapper." + TF_RESOURCE_NAME

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakGenericClientProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericClientProtocolMapper_import(realmName, clientId, mapperName),
				Check:  testKeycloakGenericClientProtocolMapperExists(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(resourceName),
			},
		},
	})
}

func TestGenericClientProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-generic-client-protocol-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_generic_client_protocol_mapper." + TF_RESOURCE_NAME

	oldAttributeName := "attribute-name-" + acctest.RandString(10)
	oldAttributeValue := "attribute-name-" + acctest.RandString(10)
	newAttributeName := "attribute-value-" + acctest.RandString(10)
	newAttributeValue := "attribute-value-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakGenericClientProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericClientProtocolMapper_update(realmName, clientId, mapperName, oldAttributeName, oldAttributeValue),
				Check:  testKeycloakGenericClientProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakGenericClientProtocolMapper_update(realmName, clientId, mapperName, newAttributeName, newAttributeValue),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakGenericClientProtocolMapperExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "config.attribute.name", newAttributeName),
					resource.TestCheckResourceAttr(resourceName, "config.attribute.value", newAttributeValue)),
			},
		},
	})
}

/*  =================================================================================================================
    Helper functions
    ================================================================================================================= */
func testAccKeycloakGenericClientProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_generic_client_protocol_mapper" {
				continue
			}

			mapper, _ := getGenericClientProtocolMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("generic client protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func getGenericClientProtocolMapperUsingState(state *terraform.State, resourceName string) (*keycloak.GenericClientProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	mapperId := rs.Primary.ID
	realmId := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetGenericClientProtocolMapper(realmId, clientId, clientScopeId, mapperId)
}

func testKeycloakGenericClientProtocolMapper_basic_client(realmName string, clientId string, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "%s"
}

resource "keycloak_generic_client_protocol_mapper" "%s" {
  client_id       = "${keycloak_saml_client.saml_client.id}"
  name            = "%s"
  protocol        = "saml"
  protocol_mapper = "saml-hardcode-attribute-mapper"
  realm_id        = "${keycloak_realm.realm.id}"
  config = {
    "attribute.name"       = "name"
    "attribute.nameformat" = "Basic"
    "attribute.value"      = "value"
    "friendly.name"        = "%s"
  }
}`, realmName, clientId, TF_RESOURCE_NAME, mapperName, mapperName)
}

func testKeycloakGenericClientProtocolMapper_basic_clientScope(realmName string, clientScopeId string, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_generic_client_protocol_mapper" "%s" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
	protocol        = "openid-connect"
	protocol_mapper = "oidc-usermodel-property-mapper"
	config = {
		"user.attribute" = "foo"
		"claim.name"     = "bar"
	}
}`, realmName, clientScopeId, TF_RESOURCE_NAME, mapperName)
}

func testKeycloakGenericClientProtocolMapper_import(realmName string, clientId string, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "%s"
}

resource "keycloak_generic_client_protocol_mapper" "%s" {
  client_id       = "${keycloak_saml_client.saml_client.id}"
  name            = "%s"
  protocol        = "saml"
  protocol_mapper = "saml-hardcode-attribute-mapper"
  realm_id        = "${keycloak_realm.realm.id}"
  config = {
    "attribute.name"       = "name"
    "attribute.nameformat" = "Basic"
    "attribute.value"      = "value"
    "friendly.name"        = "%s"
  }
}`, realmName, clientId, TF_RESOURCE_NAME, mapperName, mapperName)
}

func testKeycloakGenericClientProtocolMapper_update(realmName string, clientId string, mapperName string, attributeName string, attributeValue string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
  realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "%s"
}

resource "keycloak_generic_client_protocol_mapper" "%s" {
  client_id       = "${keycloak_saml_client.saml_client.id}"
  name            = "%s"
  protocol        = "saml"
  protocol_mapper = "saml-hardcode-attribute-mapper"
  realm_id        = "${keycloak_realm.realm.id}"
  config = {
    "attribute.name"       = "%s"
    "attribute.nameformat" = "Basic"
    "attribute.value"      = "%s"
    "friendly.name"        = "%s"
  }
}`, realmName, clientId, TF_RESOURCE_NAME, mapperName, attributeName, attributeValue, mapperName)
}

func testKeycloakGenericClientProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getGenericClientProtocolMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}
