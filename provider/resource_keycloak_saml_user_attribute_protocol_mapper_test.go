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

// Tests for attaching SAML mappers to SAML client scopes are omitted
// because the keycloak_saml_client_scope resource does not exist yet.

func TestAccKeycloakSamlUserAttributeProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-saml-client-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(clientResourceName),
			},
			{
				ResourceName:      clientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(clientResourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(5)

	userAttribute := "attr-" + acctest.RandString(10)
	updatedUserAttribute := "attr-update-" + acctest.RandString(10)
	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(realmName, clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(realmName, clientId, mapperName, updatedUserAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.SamlUserAttributeProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteSamlUserAttributeProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakSamlUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(10)
	invalidSamlNameFormat := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakSamlUserAttributeProtocolMapper_samlAttributeNameFormat(realmName, clientId, mapperName, invalidSamlNameFormat),
				ExpectError: regexp.MustCompile("expected saml_attribute_name_format to be one of .+ got " + invalidSamlNameFormat),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(5)

	userAttribute := "attr-" + acctest.RandString(10)
	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(realmName, clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(realmName, updatedClientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakSamlUserAttributeProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-saml-user-attribute-mapper-" + acctest.RandString(5)

	userAttribute := "attr-" + acctest.RandString(10)
	resourceName := "keycloak_saml_user_attribute_protocol_mapper.saml_user_attribute_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakSamlUserAttributeProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(realmName, clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakSamlUserAttributeProtocolMapper_userAttribute(newRealmName, clientId, mapperName, userAttribute),
				Check:  testKeycloakSamlUserAttributeProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakSamlUserAttributeProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_saml_user_attribute_protocol_mapper" {
				continue
			}

			mapper, _ := getSamlUserAttributeMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("saml user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakSamlUserAttributeProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getSamlUserAttributeMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakSamlUserAttributeProtocolMapperFetch(resourceName string, mapper *keycloak.SamlUserAttributeProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getSamlUserAttributeMapperUsingState(state, resourceName)
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

func getSamlUserAttributeMapperUsingState(state *terraform.State, resourceName string) (*keycloak.SamlUserAttributeProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetSamlUserAttributeProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakSamlUserAttributeProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
}

resource "keycloak_saml_user_attribute_protocol_mapper" "saml_user_attribute_mapper" {
	name                       = "%s"
	realm_id                   = "${keycloak_realm.realm.id}"
	client_id                  = "${keycloak_saml_client.saml_client.id}"

	user_attribute             = "foo"
	saml_attribute_name        = "bar"
	saml_attribute_name_format = "Unspecified"
}`, realmName, clientId, mapperName)
}

func testKeycloakSamlUserAttributeProtocolMapper_userAttribute(realmName, clientId, mapperName, userAttribute string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
}

resource "keycloak_saml_user_attribute_protocol_mapper" "saml_user_attribute_mapper" {
	name                       = "%s"
	realm_id                   = "${keycloak_realm.realm.id}"
	client_id                  = "${keycloak_saml_client.saml_client.id}"

	user_attribute             = "%s"
	saml_attribute_name        = "bar"
	saml_attribute_name_format = "Unspecified"
}`, realmName, clientId, mapperName, userAttribute)
}

func testKeycloakSamlUserAttributeProtocolMapper_samlAttributeNameFormat(realmName, clientName, mapperName, samlAttributeNameFormat string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
}

resource "keycloak_saml_user_attribute_protocol_mapper" "saml_user_attribute_mapper" {
	name                       = "%s"
	realm_id                   = "${keycloak_realm.realm.id}"
	client_id                  = "${keycloak_saml_client.saml_client.id}"

	user_attribute             = "foo"
	saml_attribute_name        = "bar"
	saml_attribute_name_format = "%s"
}`, realmName, clientName, mapperName, samlAttributeNameFormat)
}
