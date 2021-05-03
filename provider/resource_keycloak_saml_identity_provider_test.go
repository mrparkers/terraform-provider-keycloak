package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakSamlIdentityProvider_basic(t *testing.T) {
	samlName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_extraConfig(t *testing.T) {
	samlName := acctest.RandomWithPrefix("tf-acc")
	customConfigValue := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_extra_config(samlName, "dummyConfig", customConfigValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlIdentityProviderHasCustomConfigValue("keycloak_saml_identity_provider.saml", customConfigValue),
				),
			},
		},
	})
}

// ensure that extra_config keys which are covered by top-level attributes are not allowed
func TestAccKeycloakSamlIdentityProvider_extraConfigInvalid(t *testing.T) {
	samlName := acctest.RandomWithPrefix("tf-acc")
	customConfigValue := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakSamlIdentityProvider_extra_config(samlName, "syncMode", customConfigValue),
				ExpectError: regexp.MustCompile("extra_config key \"syncMode\" is not allowed"),
			},
		},
	})
}

func TestAccKeycloakSamlIdentityProvider_createAfterManualDestroy(t *testing.T) {
	var saml = &keycloak.IdentityProvider{}

	samlName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlIdentityProviderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlIdentityProvider_basic(samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderFetch("keycloak_saml_identity_provider.saml", saml),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteIdentityProvider(saml.Realm, saml.Alias)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakSamlIdentityProvider_basic(samlName),
				Check:  testAccCheckKeycloakSamlIdentityProviderExists("keycloak_saml_identity_provider.saml"),
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
			EntityId:                         "https://example.com/entity_id/1",
			SingleSignOnServiceUrl:           "https://example.com/signon/1",
			BackchannelSupported:             keycloak.KeycloakBoolQuoted(firstBackchannel),
			ValidateSignature:                keycloak.KeycloakBoolQuoted(firstValidateSignature),
			HideOnLoginPage:                  keycloak.KeycloakBoolQuoted(firstHideOnLogin),
			NameIDPolicyFormat:               "Email",
			SingleLogoutServiceUrl:           "https://example.com/logout/1",
			SigningCertificate:               acctest.RandString(10),
			SignatureAlgorithm:               "RSA_SHA512",
			XmlSignKeyInfoKeyNameTransformer: "KEY_ID",
			PostBindingAuthnRequest:          keycloak.KeycloakBoolQuoted(firstPostBindingRequest),
			PostBindingResponse:              keycloak.KeycloakBoolQuoted(firstPostBindingResponse),
			PostBindingLogout:                keycloak.KeycloakBoolQuoted(firstPostBindingLogout),
			ForceAuthn:                       keycloak.KeycloakBoolQuoted(firstForceAuthn),
			WantAssertionsSigned:             keycloak.KeycloakBoolQuoted(firstAssertionsSigned),
			WantAssertionsEncrypted:          keycloak.KeycloakBoolQuoted(firstAssertionsEncrypted),
			GuiOrder:                         strconv.Itoa(acctest.RandIntRange(1, 3)),
			SyncMode:                         randomStringInSlice(syncModes),
		},
	}

	secondSaml := &keycloak.IdentityProvider{
		Realm:   realmName,
		Alias:   acctest.RandString(10),
		Enabled: !firstEnabled,
		Config: &keycloak.IdentityProviderConfig{
			EntityId:                         "https://example.com/entity_id/2",
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
			GuiOrder:                         strconv.Itoa(acctest.RandIntRange(1, 3)),
			SyncMode:                         randomStringInSlice(syncModes),
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

func testAccCheckKeycloakSamlIdentityProviderHasCustomConfigValue(resourceName, customConfigValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedSaml, err := getKeycloakSamlIdentityProviderFromState(s, resourceName)
		if err != nil {
			return err
		}

		if fetchedSaml.Config.ExtraConfig["dummyConfig"].(string) != customConfigValue {
			return fmt.Errorf("expected custom saml provider to have config with a custom key 'dummyConfig' with a value %s, but value was %s", customConfigValue, fetchedSaml.Config.ExtraConfig["dummyConfig"].(string))
		}

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

func testKeycloakSamlIdentityProvider_basic(saml string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm             			= data.keycloak_realm.realm.id
	alias             			= "%s"
	entity_id					= "https://example.com/entity_id"
	single_sign_on_service_url  = "https://example.com/auth"
}
	`, testAccRealm.Realm, saml)
}

func testKeycloakSamlIdentityProvider_extra_config(alias, configKey, configValue string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm             			= data.keycloak_realm.realm.id
	alias             			= "%s"
	entity_id					= "https://example.com/entity_id"
	single_sign_on_service_url  = "https://example.com/auth"
	extra_config                = {
		%s = "%s"
	}
}
	`, testAccRealm.Realm, alias, configKey, configValue)
}

func testKeycloakSamlIdentityProvider_basicFromInterface(saml *keycloak.IdentityProvider) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm             			= data.keycloak_realm.realm.id
	alias             			= "%s"
	enabled           			= %t
	entity_id					= "%s"
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
	gui_order                  = %s
	sync_mode                  = "%s"
}
	`, testAccRealm.Realm, saml.Alias, saml.Enabled, saml.Config.EntityId, saml.Config.SingleSignOnServiceUrl, bool(saml.Config.BackchannelSupported), bool(saml.Config.ValidateSignature), bool(saml.Config.HideOnLoginPage), saml.Config.NameIDPolicyFormat, saml.Config.SingleLogoutServiceUrl, saml.Config.SigningCertificate, saml.Config.SignatureAlgorithm, saml.Config.XmlSignKeyInfoKeyNameTransformer, bool(saml.Config.PostBindingAuthnRequest), bool(saml.Config.PostBindingResponse), bool(saml.Config.PostBindingLogout), bool(saml.Config.ForceAuthn), bool(saml.Config.WantAssertionsSigned), bool(saml.Config.WantAssertionsEncrypted), saml.Config.GuiOrder, saml.Config.SyncMode)
}
