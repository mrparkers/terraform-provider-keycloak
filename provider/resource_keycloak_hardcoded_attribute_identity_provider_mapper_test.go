package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakHardcodedAttributeIdentityProviderMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	attributeName := "terraform-" + acctest.RandString(10)
	attributeValue := "terraform-" + acctest.RandString(10)
	userSession := randomBool()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basic(realmName, alias, mapperName, attributeName, attributeValue, userSession),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedAttributeIdentityProviderMapper_withExtraConfig(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	attributeName := "terraform-" + acctest.RandString(10)
	attributeValue := "terraform-" + acctest.RandString(10)
	userSession := randomBool()
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, attributeName, attributeValue, userSession, syncMode),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedAttributeIdentityProviderMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	attributeName := "terraform-" + acctest.RandString(10)
	attributeValue := "terraform-" + acctest.RandString(10)
	userSession := randomBool()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basic(realmName, alias, mapperName, attributeName, attributeValue, userSession),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperFetch("keycloak_hardcoded_attribute_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basic(realmName, alias, mapperName, attributeName, attributeValue, userSession),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedAttributeIdentityProviderMapper_withExtraConfig_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	attributeName := "terraform-" + acctest.RandString(10)
	attributeValue := "terraform-" + acctest.RandString(10)
	userSession := randomBool()
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, attributeName, attributeValue, userSession, syncMode),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperFetch("keycloak_hardcoded_attribute_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basic(realmName, alias, mapperName, attributeName, attributeValue, userSession),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.oidc"),
			},
		},
	})
}
func TestAccKeycloakHardcodedAttributeIdentityProviderMapper_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	attributeName := "terraform-" + acctest.RandString(10)
	attributeValue := "terraform-" + acctest.RandString(10)
	userSession := randomBool()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basic(firstRealm, alias, mapperName, attributeName, attributeValue, userSession),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_hardcoded_attribute_identity_provider_mapper.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basic(secondRealm, alias, mapperName, attributeName, attributeValue, userSession),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_hardcoded_attribute_identity_provider_mapper.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakHardcodedAttributeIdentityProviderMapper_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	identityProviderAliasName := "terraform-" + acctest.RandString(10)
	userSession := randomBool()

	firstMapper := &keycloak.IdentityProviderMapper{
		Realm:                  realmName,
		IdentityProviderAlias:  identityProviderAliasName,
		Name:                   acctest.RandString(10),
		IdentityProviderMapper: getHardcodedAttributeIdentityProviderMapperType(userSession),
		Config: &keycloak.IdentityProviderMapperConfig{
			Attribute:      acctest.RandString(10),
			AttributeValue: acctest.RandString(10),
		},
	}

	secondMapper := &keycloak.IdentityProviderMapper{
		Realm:                  realmName,
		IdentityProviderAlias:  identityProviderAliasName,
		Name:                   acctest.RandString(10),
		IdentityProviderMapper: getHardcodedAttributeIdentityProviderMapperType(!userSession),
		Config: &keycloak.IdentityProviderMapperConfig{
			Attribute:      acctest.RandString(10),
			AttributeValue: acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basicFromInterface(firstMapper, userSession),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.saml"),
			},
			{
				Config: testKeycloakHardcodedAttributeIdentityProviderMapper_basicFromInterface(secondMapper, !userSession),
				Check:  testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists("keycloak_hardcoded_attribute_identity_provider_mapper.saml"),
			},
		},
	})
}

func testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakHardcodedAttributeIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperFetch(resourceName string, mapper *keycloak.IdentityProviderMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakHardcodedAttributeIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.IdentityProviderAlias = fetchedMapper.IdentityProviderAlias
		mapper.Realm = fetchedMapper.Realm
		mapper.Id = fetchedMapper.Id

		return nil
	}
}

func testAccCheckKeycloakHardcodedAttributeIdentityProviderMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_hardcoded_attribute_identity_provider_mapper" {
				continue
			}

			realm := rs.Primary.Attributes["realm"]
			alias := rs.Primary.Attributes["identity_provider_alias"]
			id := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			mapper, _ := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
			if mapper != nil {
				return fmt.Errorf("oidc config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakHardcodedAttributeIdentityProviderMapperFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm"]
	alias := rs.Primary.Attributes["identity_provider_alias"]
	id := rs.Primary.ID

	mapper, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return nil, fmt.Errorf("error getting identity provider mapper config with id %s: %s", id, err)
	}

	return mapper, nil
}

func testKeycloakHardcodedAttributeIdentityProviderMapper_basic(realm, alias, name, attributeName, attributeValue string, userSession bool) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
}

resource keycloak_hardcoded_attribute_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	attribute_name          = "%s"
	attribute_value         = "%s"
	user_session            = %t
}
	`, realm, alias, name, attributeName, attributeValue, userSession)
}

func testKeycloakHardcodedAttributeIdentityProviderMapper_withExtraConfig(realm, alias, name, attributeName, attributeValue string, userSession bool, syncMode string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
}

resource keycloak_hardcoded_attribute_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	attribute_name          = "%s"
	attribute_value         = "%s"
	user_session            = %t
	extra_config 			= {
		syncMode = "%s"
	}
}
	`, realm, alias, name, attributeName, attributeValue, userSession, syncMode)
}

func testKeycloakHardcodedAttributeIdentityProviderMapper_basicFromInterface(mapper *keycloak.IdentityProviderMapper, userSession bool) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm                      = "${keycloak_realm.realm.id}"
	alias                      = "%s"
	single_sign_on_service_url = "https://example.com/auth"
}

resource keycloak_hardcoded_attribute_identity_provider_mapper saml {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_saml_identity_provider.saml.alias}"
	attribute_name          = "%s"
	attribute_value         = "%s"
	user_session            = %t
}
	`, mapper.Realm, mapper.IdentityProviderAlias, mapper.Name, mapper.Config.Attribute, mapper.Config.AttributeValue, userSession)
}
