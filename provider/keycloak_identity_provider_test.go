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

func TestAccKeycloakIdentityProvider_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
			},
			{
				ResourceName:        "keycloak_identity_provider.saml",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakIdentityProvider_createAfterManualDestroy(t *testing.T) {
	var idp = &keycloak.IdentityProvider{}

	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakIdentityProviderFetch("keycloak_identity_provider.saml", idp),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProvider(idp.Realm, idp.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakIdentityProvider_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProvider_basic(firstRealm, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
					resource.TestCheckResourceAttr("keycloak_identity_provider.saml", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakIdentityProvider_basic(secondRealm, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
					resource.TestCheckResourceAttr("keycloak_identity_provider.saml", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakIdentityProvider_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	firstEnabled := randomBool()

	firstalias := &keycloak.IdentityProvider{
		Realm:      realmName,
		Alias:      "terraform-" + acctest.RandString(10),
		ProviderId: "saml",
		Enabled:    firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl: "alias://" + acctest.RandString(10),
		},
	}

	secondalias := &keycloak.IdentityProvider{
		Realm:      realmName,
		Alias:      "terraform-" + acctest.RandString(10),
		ProviderId: "saml",
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl: "alias://" + acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProvider_basicFromInterface(firstalias),
				Check:  testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
			},
			{
				Config: testKeycloakIdentityProvider_basicFromInterface(secondalias),
				Check:  testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakIdentityProvider_unsetTimeoutDurationStrings(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProvider_basicWithTimeouts(realmName, aliasName),
				Check:  testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
			},
			{
				Config: testKeycloakIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakIdentityProviderExists("keycloak_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakIdentityProvider_displayNameValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)
	displayName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakIdentityProvider_basicWithAttrValidation("display_name", realmName, aliasName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected edit_mode to be one of .+ got .+"),
			},
			{
				Config: testKeycloakIdentityProvider_basicWithAttrValidation("display_name", realmName, aliasName, displayName),
				Check:  resource.TestCheckResourceAttr("keycloak_identity_provider.saml", "display_name", displayName),
			},
		},
	})
}

func testAccCheckKeycloakIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakIdentityProviderFetch(resourceName string, idp *keycloak.IdentityProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedalias, err := getIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		idp.Alias = fetchedalias.Alias
		idp.Realm = fetchedalias.Realm

		return nil
	}
}

func testAccCheckKeycloakIdentityProviderDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_identity_provider" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			alias, _ := keycloakClient.GetIdentityProvider(realm, id)
			if alias != nil {
				return fmt.Errorf("alias config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getIdentityProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm"]

	alias, err := keycloakClient.GetIdentityProvider(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting alias config with id %s: %s", id, err)
	}

	return alias, nil
}

func testKeycloakIdentityProvider_basic(realm, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource keycloak_identity_provider saml {
  alias   = "%s"
  realm   = "master"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}
   `, realm, alias)
}

func testKeycloakIdentityProvider_basicFromInterface(alias *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml" {
   alias   = "%s"
   realm   = "master"
   enabled = %t

   saml {
      single_sign_on_service_url = "https://example.com"
   }
}
   `, alias.Realm, alias.Alias, alias.Enabled)
}

func testKeycloakIdentityProvider_basicWithAttrValidation(attr, realm, alias, val string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml" {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  %s      = "%s"

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}
   `, realm, alias, attr, val)
}

func testKeycloakIdentityProvider_nobindDnValidation(realm, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml" {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}
   `, realm, alias)
}

func testKeycloakIdentityProvider_noBindCredentialValidation(realm, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml" {
  alias   = %s"
  realm   = "${keycloak_realm.realm.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}
   `, realm, alias)
}

func testKeycloakIdentityProvider_basicWithSyncPeriod(realm, alias string, fullSyncPeriod, changedSyncPeriod int) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml" {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}
   `, realm, alias)
}

func testKeycloakIdentityProvider_basicWithTimeouts(realm, alias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml" {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}
   `, realm, alias)
}
