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
	samlName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
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

	realmName := "terraform-" + acctest.RandString(10)
	samlName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(realmName, samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderFetch("keycloak_saml_identity_provider.saml", saml),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

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
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	samlName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
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
	realmName := "terraform-" + acctest.RandString(10)
	firstEnabled := randomBool()
	firstBackchannel := randomBool()
	firstValidateSignature := randomBool()
	firstHideOnLogin := randomBool()
	firstForceAuthn := randomBool()
	firstWantAuthnRequests := randomBool()
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
			BackchannelSupported:             firstBackchannel,
			ValidateSignature:                firstValidateSignature,
			HideOnLoginPage:                  firstHideOnLogin,
			NameIdPolicyFormat:               acctest.RandString(10),
			SingleLogoutServiceUrl:           "https://example.com/logout/2",
			SigningCertificate:               acctest.RandString(10),
			SignatureAlgorithm:               acctest.RandString(10),
			XmlSignKeyInfoKeyNameTransformer: acctest.RandString(10),
			PostBindingAuthnRequest:          firstPostBindingRequest,
			PostBindingResponse:              firstPostBindingResponse,
			PostBindingLogout:                firstPostBindingLogout,
			ForceAuthn:                       firstForceAuthn,
			WantAuthnRequestsSigned:          firstWantAuthnRequests,
			WantAssertionsSigned:             firstAssertionsSigned,
			WantAssertionsEncrypted:          firstAssertionsEncrypted,
		},
	}

	secondSaml := &keycloak.IdentityProvider{
		Realm:   realmName,
		Alias:   acctest.RandString(10),
		Enabled: !firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			SingleSignOnServiceUrl:           "https://example.com/signon/2",
			BackchannelSupported:             !firstBackchannel,
			ValidateSignature:                !firstValidateSignature,
			HideOnLoginPage:                  !firstHideOnLogin,
			NameIdPolicyFormat:               acctest.RandString(10),
			SingleLogoutServiceUrl:           "https://example.com/logout/2",
			SigningCertificate:               acctest.RandString(10),
			SignatureAlgorithm:               acctest.RandString(10),
			XmlSignKeyInfoKeyNameTransformer: "KEY_ID_2",
			PostBindingAuthnRequest:          !firstPostBindingRequest,
			PostBindingResponse:              !firstPostBindingResponse,
			PostBindingLogout:                !firstPostBindingLogout,
			ForceAuthn:                       !firstForceAuthn,
			WantAuthnRequestsSigned:          !firstWantAuthnRequests,
			WantAssertionsSigned:             !firstAssertionsSigned,
			WantAssertionsEncrypted:          !firstAssertionsEncrypted,
		},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlIdentityProviderDestroy(),
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
		_, err := getKeycloakIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakSamlIdentityProviderFetch(resourceName string, saml *keycloak.IdentityProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedSaml, err := getKeycloakIdentityProviderFromState(s, resourceName)
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

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			saml, _ := keycloakClient.GetIdentityProvider(realm, id)
			if saml != nil {
				return fmt.Errorf("saml config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakIdentityProviderFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProvider, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

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
	name_id_policy_format      = "urn:oasis:names:tc:SAML:2.0:nameid-format:%s"
	single_logout_service_url  = "%s"
	signing_certificate        = "%s"
	signature_algorithm        = "%s",
	xml_sign_key_info_key_name_transformer = "%s"
	post_binding_authn_request = "%s"
	post_binding_response      = %t
	post_binding_logout        = %t
	force_authn                = %t
	want_authn_requests_signed = %t
	want_assertions_signed     = %t
	want_assertions_encrypted  = %t
}
	`, saml.Realm, saml.Alias, saml.Enabled, saml.Config.SingleSignOnServiceUrl, saml.Config.BackchannelSupported, saml.Config.ValidateSignature, saml.Config.HideOnLoginPage, saml.Config.NameIdPolicyFormat, saml.Config.SingleLogoutServiceUrl, saml.Config.SigningCertificate, saml.Config.SignatureAlgorithm, saml.Config.XmlSignKeyInfoKeyNameTransformer, saml.Config.PostBindingAuthnRequest, saml.Config.PostBindingResponse, saml.Config.PostBindingLogout, saml.Config.ForceAuthn, saml.Config.WantAuthnRequestsSigned, saml.Config.WantAssertionsSigned, saml.Config.WantAssertionsEncrypted)
}
