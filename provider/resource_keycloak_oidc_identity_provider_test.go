package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOidcIdentityProvider_basic(t *testing.T) {
	oidcName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(oidcName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_extra_config(t *testing.T) {
	oidcName := acctest.RandomWithPrefix("tf-acc")
	customConfigValue := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_extra_config(oidcName, customConfigValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderHasCustomConfigValue("keycloak_oidc_identity_provider.oidc", customConfigValue),
				),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_keyDefaultScopes(t *testing.T) {
	oidcName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_keyDefaultScopes(oidcName, "openid random"),
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

	oidcName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(oidcName),
				Check:  testAccCheckKeycloakOidcIdentityProviderFetch("keycloak_oidc_identity_provider.oidc", oidc),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteIdentityProvider(oidc.Realm, oidc.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOidcIdentityProvider_basic(oidcName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_basicUpdateAll(t *testing.T) {
	firstEnabled := randomBool()

	firstOidc := &keycloak.IdentityProvider{
		Realm:   testAccRealm.Realm,
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
		Realm:   testAccRealm.Realm,
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
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOidcIdentityProviderDestroy(),
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

			oidc, _ := keycloakClient.GetIdentityProvider(realm, id)
			if oidc != nil {
				return fmt.Errorf("oidc config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakOidcIdentityProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
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

func testKeycloakOidcIdentityProvider_basic(oidc string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = data.keycloak_realm.realm.id
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
}
	`, testAccRealm.Realm, oidc)
}

func testKeycloakOidcIdentityProvider_extra_config(alias, customConfigValue string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = data.keycloak_realm.realm.id
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
	`, testAccRealm.Realm, alias, customConfigValue)
}

func testKeycloakOidcIdentityProvider_keyDefaultScopes(alias, value string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = data.keycloak_realm.realm.id
	provider_id       = "oidc"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
	default_scopes    = "%s"
}
	`, testAccRealm.Realm, alias, value)
}

func testKeycloakOidcIdentityProvider_basicFromInterface(oidc *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = data.keycloak_realm.realm.id
	alias             = "%s"
	enabled           = %t
	authorization_url = "%s"
	token_url         = "%s"
	client_id         = "%s"
	client_secret     = "%s"
}
	`, testAccRealm.Realm, oidc.Alias, oidc.Enabled, oidc.Config.AuthorizationUrl, oidc.Config.TokenUrl, oidc.Config.ClientId, oidc.Config.ClientSecret)
}
