package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakUserTemplateIdentityProviderMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	template := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_basic(realmName, alias, mapperName, template),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakUserTemplateIdentityProviderMapper_withExtraConfig(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	template := "terraform-" + acctest.RandString(10)
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, template, syncMode),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakUserTemplateIdentityProviderMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	template := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_basic(realmName, alias, mapperName, template),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperFetch("keycloak_user_template_importer_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakUserTemplateIdentityProviderMapper_basic(realmName, alias, mapperName, template),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakUserTemplateIdentityProviderMapper_withExtraConfig_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	template := "terraform-" + acctest.RandString(10)
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, template, syncMode),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperFetch("keycloak_user_template_importer_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakUserTemplateIdentityProviderMapper_basic(realmName, alias, mapperName, template),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakUserTemplateIdentityProviderMapper_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	template := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_basic(firstRealm, alias, mapperName, template),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_user_template_importer_identity_provider_mapper.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_basic(secondRealm, alias, mapperName, template),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_user_template_importer_identity_provider_mapper.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakUserTemplateIdentityProviderMapper_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	identityProviderAliasName := "terraform-" + acctest.RandString(10)

	firstMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			Template: acctest.RandString(10),
		},
	}

	secondMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			Template: acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_basicFromInterface(firstMapper),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.saml"),
			},
			{
				Config: testKeycloakUserTemplateIdentityProviderMapper_basicFromInterface(secondMapper),
				Check:  testAccCheckKeycloakUserTemplateIdentityProviderMapperExists("keycloak_user_template_importer_identity_provider_mapper.saml"),
			},
		},
	})
}

func testAccCheckKeycloakUserTemplateIdentityProviderMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakUserTemplateIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakUserTemplateIdentityProviderMapperFetch(resourceName string, mapper *keycloak.IdentityProviderMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakUserTemplateIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.IdentityProviderAlias = fetchedMapper.IdentityProviderAlias
		mapper.Realm = fetchedMapper.Realm
		mapper.Id = fetchedMapper.Id

		return nil
	}
}

func testAccCheckKeycloakUserTemplateIdentityProviderMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_user_template_importer_identity_provider_mapper" {
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

func getKeycloakUserTemplateIdentityProviderMapperFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProviderMapper, error) {
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

func testKeycloakUserTemplateIdentityProviderMapper_basic(realm, alias, name, template string) string {
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

resource keycloak_user_template_importer_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	template                = "%s"
}
	`, realm, alias, name, template)
}

func testKeycloakUserTemplateIdentityProviderMapper_withExtraConfig(realm, alias, name, template, syncMode string) string {
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

resource keycloak_user_template_importer_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	template                = "%s"
	extra_config 			= {
		syncMode = "%s"
	}
}
	`, realm, alias, name, template, syncMode)
}

func testKeycloakUserTemplateIdentityProviderMapper_basicFromInterface(mapper *keycloak.IdentityProviderMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm                      = "${keycloak_realm.realm.id}"
	alias                      = "%s"
	single_sign_on_service_url = "https://example.com/auth"
}

resource keycloak_user_template_importer_identity_provider_mapper saml {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_saml_identity_provider.saml.alias}"
	template                = "%s"
}
	`, mapper.Realm, mapper.IdentityProviderAlias, mapper.Name, mapper.Config.Template)
}
