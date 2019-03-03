package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOidcIdentityProvider_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
			{
				ResourceName:        "keycloak_oidc_identity_provider.oidc",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_createAfterManualDestroy(t *testing.T) {
	var idp = &keycloak.IdentityProvider{}

	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakOidcIdentityProviderFetch("keycloak_oidc_identity_provider.oidc", idp),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProvider(idp.Realm, idp.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOidcIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basic(firstRealm, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
					resource.TestCheckResourceAttr("keycloak_oidc_identity_provider.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakIdentityProvider_basic(secondRealm, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
					resource.TestCheckResourceAttr("keycloak_oidc_identity_provider.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakIdentityProvider_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	firstAlias := &keycloak.IdentityProvider{
		Realm:      realmName,
		Alias:      "terraform-" + acctest.RandString(10),
		Enabled:    true,
		ProviderId: "saml",
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl: "alias://" + acctest.RandString(10),
		},
	}

	secondAlias := &keycloak.IdentityProvider{
		Realm:      realmName,
		Alias:      "terraform-" + acctest.RandString(10),
		Enabled:    true,
		ProviderId: "saml",
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl: "alias://" + acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basicFromInterface(firstAlias),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
			{
				Config: testKeycloakOidcIdentityProvider_basicFromInterface(secondAlias),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func TestAccKeycloakOidcIdentityProvider_unsetTimeoutDurationStrings(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOidcIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOidcIdentityProvider_basicWithTimeouts(realmName, aliasName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
			{
				Config: testKeycloakOidcIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakOidcIdentityProviderExists("keycloak_oidc_identity_provider.oidc"),
			},
		},
	})
}

func testAccCheckKeycloakOidcIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getOidcIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOidcIdentityProviderFetch(resourceName string, idp *keycloak.IdentityProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedalias, err := getIdentityOidcProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		idp.Alias = fetchedalias.Alias
		idp.Realm = fetchedalias.Realm

		return nil
	}
}

func testAccCheckKeycloakOidcIdentityProviderDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_identity_provider" {
				continue
			}

			realm := rs.Primary.Attributes["realm"]
			alias := rs.Primary.Attributes["alias"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			idp, _ := keycloakClient.GetIdentityProvider(realm, alias)
			if idp != nil {
				return fmt.Errorf("idp config with alias %s still exists", alias)
			}
		}

		return nil
	}
}

func getIdentityOidcProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm"]
	alias := rs.Primary.Attributes["alias"]

	idp, err := keycloakClient.GetIdentityProvider(realm, alias)
	if err != nil {
		return nil, fmt.Errorf("error getting idp config with alias %s: %s", alias, err)
	}

	return idp, nil
}

func testKeycloakOidcIdentityProvider_basic(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_oidc_identity_provider oidc {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  token_url = "https://example.com"
  authorization_url = "https://example.com"
  client_id = "laksjhnlasjdkblakjsbda"
  client_secret = "asdasdjasdklhasldas"
}
   `, realm, alias)
}

func testKeycloakOidcIdentityProvider_basicFromInterface(alias *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_oidc_identity_provider oidc {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  token_url = "https://example.com"
  authorization_url = "https://example.com"
  client_id = "laksjhnlasjdkblakjsbda"
  client_secret = "asdasdjasdklhasldas"
}
   `, alias.Realm, alias.Alias)
}

func testKeycloakOidcIdentityProvider_nobindDnValidation(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_oidc_identity_provider oidc {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  token_url = "https://example.com"
  authorization_url = "https://example.com"
  client_id = "laksjhnlasjdkblakjsbda"
  client_secret = "asdasdjasdklhasldas"
}
   `, realm, alias)
}

func testKeycloakOidcIdentityProvider_noBindCredentialValidation(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_oidc_identity_provider oidc {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  token_url = "https://example.com"
  authorization_url = "https://example.com"
  client_id = "laksjhnlasjdkblakjsbda"
  client_secret = "asdasdjasdklhasldas"
}
   `, realm, alias)
}

func testKeycloakOidcIdentityProvider_basicWithSyncPeriod(realm, alias string, fullSyncPeriod, changedSyncPeriod int) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_oidc_identity_provider oidc {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  token_url = "https://example.com"
  authorization_url = "https://example.com"
  client_id = "laksjhnlasjdkblakjsbda"
  client_secret = "asdasdjasdklhasldas"
}
   `, realm, alias)
}

func testKeycloakOidcIdentityProvider_basicWithTimeouts(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_oidc_identity_provider oidc {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  token_url = "https://example.com"
  authorization_url = "https://example.com"
  client_id = "laksjhnlasjdkblakjsbda"
  client_secret = "asdasdjasdklhasldas"
}
   `, realm, alias)
}
