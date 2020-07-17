package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakHardcodedRoleIdentityProviderMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedRoleIdentityProviderMapper_withExtraConfig(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, role, syncMode),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedRoleIdentityProviderMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperFetch("keycloak_hardcoded_role_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedRoleIdentityProviderMapper_withExtraConfig_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)
	syncMode := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_withExtraConfig(realmName, alias, mapperName, role, syncMode),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperFetch("keycloak_hardcoded_role_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basic(realmName, alias, mapperName, role),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakHardcodedRoleIdentityProviderMapper_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	role := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basic(firstRealm, alias, mapperName, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_hardcoded_role_identity_provider_mapper.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basic(secondRealm, alias, mapperName, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_hardcoded_role_identity_provider_mapper.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakHardcodedRoleIdentityProviderMapper_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	identityProviderAliasName := "terraform-" + acctest.RandString(10)

	firstMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			Role: acctest.RandString(10),
		},
	}

	secondMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			Role: acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basicFromInterface(firstMapper),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.saml"),
			},
			{
				Config: testKeycloakHardcodedRoleIdentityProviderMapper_basicFromInterface(secondMapper),
				Check:  testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists("keycloak_hardcoded_role_identity_provider_mapper.saml"),
			},
		},
	})
}

func testAccCheckKeycloakHardcodedRoleIdentityProviderMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakHardcodedRoleIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakHardcodedRoleIdentityProviderMapperFetch(resourceName string, mapper *keycloak.IdentityProviderMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakHardcodedRoleIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.IdentityProviderAlias = fetchedMapper.IdentityProviderAlias
		mapper.Realm = fetchedMapper.Realm
		mapper.Id = fetchedMapper.Id

		return nil
	}
}

func testAccCheckKeycloakHardcodedRoleIdentityProviderMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_hardcoded_role_identity_provider_mapper" {
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

func getKeycloakHardcodedRoleIdentityProviderMapperFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProviderMapper, error) {
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

func testKeycloakHardcodedRoleIdentityProviderMapper_basic(realm, alias, name, role string) string {
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

resource keycloak_hardcoded_role_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	role                    = "%s"
}
	`, realm, alias, name, role)
}

func testKeycloakHardcodedRoleIdentityProviderMapper_withExtraConfig(realm, alias, name, role, syncMode string) string {
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

resource keycloak_hardcoded_role_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	role                    = "%s"
	extra_config 			= {
		syncMode = "%s"
	}
}
	`, realm, alias, name, role, syncMode)
}

func testKeycloakHardcodedRoleIdentityProviderMapper_basicFromInterface(mapper *keycloak.IdentityProviderMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm                      = "${keycloak_realm.realm.id}"
	alias                      = "%s"
	single_sign_on_service_url = "https://example.com/auth"
}

resource keycloak_hardcoded_role_identity_provider_mapper saml {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_saml_identity_provider.saml.alias}"
	role                    = "%s"
}
	`, mapper.Realm, mapper.IdentityProviderAlias, mapper.Name, mapper.Config.Role)
}
