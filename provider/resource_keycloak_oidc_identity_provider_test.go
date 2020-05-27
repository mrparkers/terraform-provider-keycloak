package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOidcIdentityProvider_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	oidcName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(realmName, oidcName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_custom(t *testing.T) {
	skipIfEnvSet(t, "CI") // temporary while I figure out how to load this custom idp in CI
	//This test does not work in keycloak 10, because the interfaces that our customIdp implements, have changed in the keycloak latest version.
	//We need to decide which keycloak version we going to support and test for the customIdp
	realmName := "terraform-" + acctest.RandString(10)
	oidcName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_custom(realmName, oidcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
				),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_extra_config(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	oidcName := "terraform-" + acctest.RandString(10)
	customConfigValue := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_extra_config(realmName, oidcName, customConfigValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderHasCustomConfigValue("keycloak_oidc_identity_provider.oidc", customConfigValue),
				),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_keyDefaultScopes(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	oidcName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_keyDefaultScopes(realmName, oidcName, "openid random"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
					testAccCheckKeycloakOidcIdentityProviderDefaultScopes("keycloak_oidc_identity_provider.oidc", "openid random"),
				),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_createAfterManualDestroy(t *testing.T) {
	var oidc = &keycloak.IdentityProvider{}

	realmName := "terraform-" + acctest.RandString(10)
	oidcName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(realmName, oidcName),
				Check:  testAccCheckKeycloakOidcIdentityProviderFetch("keycloak_oidc_identity_provider.oidc", oidc),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProvider(oidc.Realm, oidc.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOidcIdentityProvider_basic(realmName, oidcName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	oidcName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(firstRealm, oidcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
					resource.TestCheckResourceAttr("keycloak_oidc_identity_provider.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakOidcIdentityProvider_basic(secondRealm, oidcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
					resource.TestCheckResourceAttr("keycloak_oidc_identity_provider.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	firstEnabled := randomBool()

	firstOidc := &keycloak.IdentityProvider{
		Realm:   realmName,
		Alias:   acctest.RandString(10),
		Enabled: firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			AuthorizationUrl: "https://example.com/auth",
			TokenUrl:         "https://example.com/token",
			ClientId:         acctest.RandString(10),
			ClientSecret:     acctest.RandString(10),
		},
	}

	secondOidc := &keycloak.IdentityProvider{
		Realm:   realmName,
		Alias:   acctest.RandString(10),
		Enabled: !firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			AuthorizationUrl: "https://example.com/auth",
			TokenUrl:         "https://example.com/token",
			ClientId:         acctest.RandString(10),
			ClientSecret:     acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basicFromInterface(firstOidc),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
			{
				Config: testKeycloakOidcIdentityProvider_basicFromInterface(secondOidc),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func testAccCheckKeycloakOidcIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakOidcIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOidcIdentityProviderFetch(resourceName string, oidc *keycloak.IdentityProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedOidc, err := getKeycloakOidcIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		oidc.Alias = fetchedOidc.Alias
		oidc.Realm = fetchedOidc.Realm

		return nil
	}
}

func testAccCheckKeycloakOidcIdentityProviderHasCustomConfigValue(resourceName, customConfigValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedOidc, err := getKeycloakOidcIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		if fetchedOidc.Config.ExtraConfig["dummyConfig"].(string) != customConfigValue {
			return fmt.Errorf("expected custom oidc provider to have config with a custom key 'dummyConfig' with a value %s, but value was %s", customConfigValue, fetchedOidc.Config.ExtraConfig["dummyConfig"].(string))
		}

		return nil
	}
}

func testAccCheckKeycloakOidcIdentityProviderDefaultScopes(resourceName, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedOidc, err := getKeycloakOidcIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		if fetchedOidc.Config.DefaultScope != value {
			return fmt.Errorf("expected oidc provider to have value %s for key 'defaultScope', but value was %s", value, fetchedOidc.Config.DefaultScope)
		}

		return nil
	}
}

func testAccCheckKeycloakOidcIdentityProviderDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_oidc_identity_provider" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			oidc, _ := keycloakClient.GetIdentityProvider(realm, id)
			if oidc != nil {
				return fmt.Errorf("oidc config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakOidcIdentityProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm"]
	alias := rs.Primary.Attributes["alias"]

	oidc, err := keycloakClient.GetIdentityProvider(realm, alias)
	if err != nil {
		return nil, fmt.Errorf("error getting oidc identity provider config with alias %s: %s", alias, err)
	}

	return oidc, nil
}

func testKeycloakOidcIdentityProvider_basic(realm, oidc string) string {
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
	`, realm, oidc)
}

func testKeycloakOidcIdentityProvider_custom(realm, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	provider_id       = "customIdp"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
}
	`, realm, alias)
}

func testKeycloakOidcIdentityProvider_extra_config(realm, alias, customConfigValue string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	provider_id       = "oidc"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
	extra_config      = {
		dummyConfig = "%s"
	}
}
	`, realm, alias, customConfigValue)
}

func testKeycloakOidcIdentityProvider_keyDefaultScopes(realm, alias, value string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	provider_id       = "oidc"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
	default_scopes    = "%s"
}
	`, realm, alias, value)
}

func testKeycloakOidcIdentityProvider_basicFromInterface(oidc *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	alias             = "%s"
	enabled           = %t
	authorization_url = "%s"
	token_url         = "%s"
	client_id         = "%s"
	client_secret     = "%s"
}
	`, oidc.Realm, oidc.Alias, oidc.Enabled, oidc.Config.AuthorizationUrl, oidc.Config.TokenUrl, oidc.Config.ClientId, oidc.Config.ClientSecret)
}
