package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakAttributeToRoleIdentityProviderMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)
	claimValue := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role, claimName, claimValue),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakAttributeToRoleIdentityProviderMapper_withExtraConfig(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)
	claimValue := "terraform-" + acctest.RandString(10)
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, role, claimName, claimValue, syncMode),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakAttributeToRoleIdentityProviderMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)
	claimValue := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role, claimName, claimValue),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperFetch("keycloak_attribute_to_role_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role, claimName, claimValue),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakAttributeToRoleIdentityProviderMapper_withExtraConfig_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)
	claimValue := "terraform-" + acctest.RandString(10)
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, role, claimName, claimValue, syncMode),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperFetch("keycloak_attribute_to_role_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, role, claimName, claimValue, syncMode),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakAttributeToRoleIdentityProviderMapper_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)
	claimValue := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basic(firstRealm, alias, mapperName, role, claimName, claimValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_attribute_to_role_identity_provider_mapper.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basic(secondRealm, alias, mapperName, role, claimName, claimValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_attribute_to_role_identity_provider_mapper.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakAttributeToRoleIdentityProviderMapper_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	identityProviderAliasName := "terraform-" + acctest.RandString(10)

	firstMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			AttributeValue: acctest.RandString(10),
			Attribute:      acctest.RandString(10),
			Role:           acctest.RandString(10),
		},
	}

	secondMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			AttributeValue: acctest.RandString(10),
			Attribute:      acctest.RandString(10),
			Role:           acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basicFromInterface(firstMapper),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.saml"),
			},
			{
				Config: testKeycloakAttributeToRoleIdentityProviderMapper_basicFromInterface(secondMapper),
				Check:  testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists("keycloak_attribute_to_role_identity_provider_mapper.saml"),
			},
		},
	})
}

func testAccCheckKeycloakAttributeToRoleIdentityProviderMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakAttributeToRoleIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakAttributeToRoleIdentityProviderMapperFetch(resourceName string, mapper *keycloak.IdentityProviderMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakAttributeToRoleIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.IdentityProviderAlias = fetchedMapper.IdentityProviderAlias
		mapper.Realm = fetchedMapper.Realm
		mapper.Id = fetchedMapper.Id

		return nil
	}
}

func testAccCheckKeycloakAttributeToRoleIdentityProviderMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_attribute_to_role_identity_provider_mapper" {
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

func getKeycloakAttributeToRoleIdentityProviderMapperFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProviderMapper, error) {
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

func testKeycloakAttributeToRoleIdentityProviderMapper_basic(realm, alias, name, role, claimName, claimValue string) string {
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

resource keycloak_attribute_to_role_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	role                    = "%s"
	claim_name              = "%s"
	claim_value             = "%s"
}
	`, realm, alias, name, role, claimName, claimValue)
}

func testKeycloakAttributeToRoleIdentityProviderMapper_withExtraConfig(realm, alias, name, role, claimName, claimValue, syncMode string) string {
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

resource keycloak_attribute_to_role_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	role                    = "%s"
	claim_name              = "%s"
	claim_value             = "%s"
	extra_config 			= {
		syncMode = "%s"
	}
}
	`, realm, alias, name, role, claimName, claimValue, syncMode)
}

func testKeycloakAttributeToRoleIdentityProviderMapper_basicFromInterface(mapper *keycloak.IdentityProviderMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm                      = "${keycloak_realm.realm.id}"
	alias                      = "%s"
	single_sign_on_service_url = "https://example.com/auth"
}

resource keycloak_attribute_to_role_identity_provider_mapper saml {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_saml_identity_provider.saml.alias}"
	role                    = "%s"
	attribute_name          = "%s"
	attribute_value         = "%s"
}
	`, mapper.Realm, mapper.IdentityProviderAlias, mapper.Name, mapper.Config.Role, mapper.Config.Attribute, mapper.Config.AttributeValue)
}
