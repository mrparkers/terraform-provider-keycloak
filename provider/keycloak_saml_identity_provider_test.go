package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakSamlIdentityProvider_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
			{
				ResourceName:        "keycloak_saml_identity_provider.saml",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_createAfterManualDestroy(t *testing.T) {
	var idp = &keycloak.IdentityProvider{}

	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakSamlIdentityProviderFetch("keycloak_saml_identity_provider.saml", idp),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProvider(idp.Realm, idp.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakSamlIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(firstRealm, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
					resource.TestCheckResourceAttr("keycloak_saml_identity_provider.saml", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakIdentityProvider_basic(secondRealm, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
					resource.TestCheckResourceAttr("keycloak_saml_identity_provider.saml", "realm", secondRealm),
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
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basicFromInterface(firstAlias),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
			{
				Config: testKeycloakSamlIdentityProvider_basicFromInterface(secondAlias),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_unsetTimeoutDurationStrings(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	aliasName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basicWithTimeouts(realmName, aliasName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
			{
				Config: testKeycloakSamlIdentityProvider_basic(realmName, aliasName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func testAccCheckKeycloakSamlIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getSamlIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakSamlIdentityProviderFetch(resourceName string, idp *keycloak.IdentityProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedalias, err := getIdentitySamlProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		idp.Alias = fetchedalias.Alias
		idp.Realm = fetchedalias.Realm

		return nil
	}
}

func testAccCheckKeycloakSamlIdentityProviderDestroy() resource.TestCheckFunc {
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

func getIdentitySamlProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
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

func testKeycloakSamlIdentityProvider_basic(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_saml_identity_provider saml {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  single_sign_on_service_url = "https://example.com"
}
   `, realm, alias)
}

func testKeycloakSamlIdentityProvider_basicFromInterface(alias *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_saml_identity_provider saml {
   alias   = "%s"
   realm   = "master"
   single_sign_on_service_url = "https://example.com"
}
   `, alias.Realm, alias.Alias)
}

func testKeycloakSamlIdentityProvider_nobindDnValidation(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_saml_identity_provider saml {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  single_sign_on_service_url = "https://example.com"
}
   `, realm, alias)
}

func testKeycloakSamlIdentityProvider_noBindCredentialValidation(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_saml_identity_provider saml {
  alias   = %s"
  realm   = "${keycloak_realm.realm.realm}"
  single_sign_on_service_url = "https://example.com"
}
   `, realm, alias)
}

func testKeycloakSamlIdentityProvider_basicWithSyncPeriod(realm, alias string, fullSyncPeriod, changedSyncPeriod int) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_saml_identity_provider saml {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  single_sign_on_service_url = "https://example.com"
}
   `, realm, alias)
}

func testKeycloakSamlIdentityProvider_basicWithTimeouts(realm, alias string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_saml_identity_provider saml {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  single_sign_on_service_url = "https://example.com"
}
   `, realm, alias)
}
