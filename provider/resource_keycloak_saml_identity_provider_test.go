package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakSamlIdentityProvider_basic(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")
	samlName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(realmName, samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_createAfterManualDestroy(t *testing.T) {
	var saml = &keycloak.IdentityProvider{}

	realmName := acctest.RandomWithPrefix("tf-acc")
	samlName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(realmName, samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderFetch("keycloak_saml_identity_provider.saml", saml),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteIdentityProvider(saml.Realm, saml.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakSamlIdentityProvider_basic(realmName, samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_basicUpdateRealm(t *testing.T) {
	firstRealm := acctest.RandomWithPrefix("tf-acc")
	secondRealm := acctest.RandomWithPrefix("tf-acc")
	samlName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(firstRealm, samlName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
					resource.TestCheckResourceAttr("keycloak_saml_identity_provider.saml", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakSamlIdentityProvider_basic(secondRealm, samlName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
					resource.TestCheckResourceAttr("keycloak_saml_identity_provider.saml", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_basicUpdateAll(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")
	firstEnabled := randomBool()
	firstBackchannel := randomBool()
	firstValidateSignature := randomBool()
	firstHideOnLogin := randomBool()
	firstForceAuthn := randomBool()
	firstAssertionsEncrypted := randomBool()
	firstAssertionsSigned := randomBool()
	firstPostBindingLogout := randomBool()
	firstPostBindingResponse := randomBool()
	firstPostBindingRequest := randomBool()

	firstSaml := &keycloak.IdentityProvider{
		Realm:   realmName,
		Alias:   acctest.RandString(10),
		Enabled: firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl:           "https://example.com/signon/2",
			BackchannelSupported:             keycloak.KeycloakBoolQuoted(firstBackchannel),
			ValidateSignature:                keycloak.KeycloakBoolQuoted(firstValidateSignature),
			HideOnLoginPage:                  keycloak.KeycloakBoolQuoted(firstHideOnLogin),
			NameIDPolicyFormat:               "Email",
			SingleLogoutServiceUrl:           "https://example.com/logout/2",
			SigningCertificate:               acctest.RandString(10),
			SignatureAlgorithm:               "RSA_SHA512",
			XmlSignKeyInfoKeyNameTransformer: "KEY_ID",
			PostBindingAuthnRequest:          keycloak.KeycloakBoolQuoted(firstPostBindingRequest),
			PostBindingResponse:              keycloak.KeycloakBoolQuoted(firstPostBindingResponse),
			PostBindingLogout:                keycloak.KeycloakBoolQuoted(firstPostBindingLogout),
			ForceAuthn:                       keycloak.KeycloakBoolQuoted(firstForceAuthn),
			WantAssertionsSigned:             keycloak.KeycloakBoolQuoted(firstAssertionsSigned),
			WantAssertionsEncrypted:          keycloak.KeycloakBoolQuoted(firstAssertionsEncrypted),
		},
	}

	secondSaml := &keycloak.IdentityProvider{
		Realm:   realmName,
		Alias:   acctest.RandString(10),
		Enabled: !firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl:           "https://example.com/signon/2",
			BackchannelSupported:             keycloak.KeycloakBoolQuoted(!firstBackchannel),
			ValidateSignature:                keycloak.KeycloakBoolQuoted(!firstValidateSignature),
			HideOnLoginPage:                  keycloak.KeycloakBoolQuoted(!firstHideOnLogin),
			NameIDPolicyFormat:               "Persistent",
			SingleLogoutServiceUrl:           "https://example.com/logout/2",
			SigningCertificate:               acctest.RandString(10),
			SignatureAlgorithm:               "RSA_SHA256",
			XmlSignKeyInfoKeyNameTransformer: "NONE",
			PostBindingAuthnRequest:          keycloak.KeycloakBoolQuoted(!firstPostBindingRequest),
			PostBindingResponse:              keycloak.KeycloakBoolQuoted(!firstPostBindingResponse),
			PostBindingLogout:                keycloak.KeycloakBoolQuoted(!firstPostBindingLogout),
			ForceAuthn:                       keycloak.KeycloakBoolQuoted(!firstForceAuthn),
			WantAssertionsSigned:             keycloak.KeycloakBoolQuoted(!firstAssertionsSigned),
			WantAssertionsEncrypted:          keycloak.KeycloakBoolQuoted(!firstAssertionsEncrypted),
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basicFromInterface(firstSaml),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
			{
				Config: testKeycloakSamlIdentityProvider_basicFromInterface(secondSaml),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func testAccCheckKeycloakSamlIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakSamlIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakSamlIdentityProviderFetch(resourceName string, saml *keycloak.IdentityProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedSaml, err := getKeycloakSamlIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		saml.Alias = fetchedSaml.Alias
		saml.Realm = fetchedSaml.Realm

		return nil
	}
}

func testAccCheckKeycloakSamlIdentityProviderDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_saml_identity_provider" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm"]

			saml, _ := keycloakClient.GetIdentityProvider(realm, id)
			if saml != nil {
				return fmt.Errorf("saml config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakSamlIdentityProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm"]
	alias := rs.Primary.Attributes["alias"]

	saml, err := keycloakClient.GetIdentityProvider(realm, alias)
	if err != nil {
		return nil, fmt.Errorf("error getting saml identity provider config with alias %s: %s", alias, err)
	}

	return saml, nil
}

func testKeycloakSamlIdentityProvider_basic(realm, saml string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm             			= "${keycloak_realm.realm.id}"
	alias             			= "%s"
	single_sign_on_service_url = "https://example.com/auth"
}
	`, realm, saml)
}

func testKeycloakSamlIdentityProvider_basicFromInterface(saml *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm             			= "${keycloak_realm.realm.id}"
	alias             			= "%s"
	enabled           			= %t
	single_sign_on_service_url = "%s"
	backchannel_supported      = %t
	validate_signature         = %t
	hide_on_login_page         = %t
	name_id_policy_format      = "%s"
	single_logout_service_url  = "%s"
	signing_certificate        = "%s"
	signature_algorithm        = "%s"
	xml_sign_key_info_key_name_transformer = "%s"
	post_binding_authn_request = %t
	post_binding_response      = %t
	post_binding_logout        = %t
	force_authn                = %t
	want_assertions_signed     = %t
	want_assertions_encrypted  = %t
}
	`, saml.Realm, saml.Alias, saml.Enabled, saml.Config.SingleSignOnServiceUrl, bool(saml.Config.BackchannelSupported), bool(saml.Config.ValidateSignature), bool(saml.Config.HideOnLoginPage), saml.Config.NameIDPolicyFormat, saml.Config.SingleLogoutServiceUrl, saml.Config.SigningCertificate, saml.Config.SignatureAlgorithm, saml.Config.XmlSignKeyInfoKeyNameTransformer, bool(saml.Config.PostBindingAuthnRequest), bool(saml.Config.PostBindingResponse), bool(saml.Config.PostBindingLogout), bool(saml.Config.ForceAuthn), bool(saml.Config.WantAssertionsSigned), bool(saml.Config.WantAssertionsEncrypted))
}
