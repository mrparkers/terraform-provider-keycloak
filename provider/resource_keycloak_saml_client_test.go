package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakSamlClient_basic(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_basic(clientId),
				Check:  testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
			},
			{
				ResourceName:        "keycloak_saml_client.saml_client",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: testAccRealm.Realm + "/",
			},
		},
	})
}

func TestAccKeycloakSamlClient_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	var client = &keycloak.SamlClient{}

	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_basic(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientFetch("keycloak_saml_client.saml_client", client),
				),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteSamlClient(client.RealmId, client.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakSamlClient_basic(clientId),
				Check:  testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
			},
		},
	})
}

func TestAccKeycloakSamlClient_updateRealm(t *testing.T) {
	t.Parallel()

	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_updateRealmBefore(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					resource.TestCheckResourceAttr("keycloak_saml_client.saml_client", "realm_id", testAccRealm.Realm),
				),
			},
			{
				Config: testKeycloakSamlClient_updateRealmAfter(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					resource.TestCheckResourceAttr("keycloak_saml_client.saml_client", "realm_id", testAccRealmTwo.Realm),
				),
			},
		},
	})
}

// Keycloak typically sets some values as default if they aren't provided
// This test asserts that these default values are present if none are provided
func TestAccKeycloakSamlClient_keycloakDefaults(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_basic(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientHasDefaultBooleanAttributes("keycloak_saml_client.saml_client"),
					TestCheckResourceAttrNot("keycloak_saml_client.saml_client", "signing_certificate", ""),
					TestCheckResourceAttrNot("keycloak_saml_client.saml_client", "signing_private_key", ""),
				),
			},
		},
	})
}

func TestAccKeycloakSamlClient_updateInPlace(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	enabled := randomBool()
	frontChannelLogout := randomBool()
	clientSignatureRequired := "true"

	encryptionCertificateBefore := acctest.RandomWithPrefix("tf-acc")
	encryptionCertificateAfter := acctest.RandomWithPrefix("tf-acc")
	signingCertificateBefore := acctest.RandomWithPrefix("tf-acc")
	signingCertificateAfter := acctest.RandomWithPrefix("tf-acc")
	signingPrivateKeyBefore := acctest.RandomWithPrefix("tf-acc")
	signingPrivateKeyAfter := acctest.RandomWithPrefix("tf-acc")

	samlClientBefore := &keycloak.SamlClient{
		RealmId:  testAccRealm.Realm,
		ClientId: clientId,
		Name:     acctest.RandString(10),

		Enabled:     enabled,
		Description: acctest.RandString(50),

		FrontChannelLogout: frontChannelLogout,

		RootUrl: "http://localhost:2222/" + acctest.RandString(20),
		ValidRedirectUris: []string{
			acctest.RandString(20),
			acctest.RandString(20),
			acctest.RandString(20),
		},
		BaseUrl:                 "http://localhost:2222/" + acctest.RandString(20),
		MasterSamlProcessingUrl: acctest.RandString(20),

		Attributes: &keycloak.SamlClientAttributes{
			IncludeAuthnStatement:           randomBoolAsStringPointer(),
			SignDocuments:                   randomBoolAsStringPointer(),
			SignAssertions:                  randomBoolAsStringPointer(),
			EncryptAssertions:               randomBoolAsStringPointer(),
			ClientSignatureRequired:         &clientSignatureRequired,
			ForcePostBinding:                randomBoolAsStringPointer(),
			ForceNameIdFormat:               randomBoolAsStringPointer(),
			SignatureAlgorithm:              randomStringInSlice(keycloakSamlClientSignatureAlgorithms),
			NameIdFormat:                    randomStringInSlice(keycloakSamlClientNameIdFormats),
			EncryptionCertificate:           &encryptionCertificateBefore,
			SigningCertificate:              &signingCertificateBefore,
			SigningPrivateKey:               &signingPrivateKeyBefore,
			IDPInitiatedSSOURLName:          acctest.RandString(20),
			IDPInitiatedSSORelayState:       acctest.RandString(20),
			AssertionConsumerPostURL:        acctest.RandString(20),
			AssertionConsumerRedirectURL:    acctest.RandString(20),
			LogoutServicePostBindingURL:     acctest.RandString(20),
			LogoutServiceRedirectBindingURL: acctest.RandString(20),
		},
	}

	samlClientAfter := &keycloak.SamlClient{
		RealmId:  testAccRealm.Realm,
		ClientId: clientId,
		Name:     acctest.RandString(10),

		Enabled:     !enabled,
		Description: acctest.RandString(50),

		FrontChannelLogout: !frontChannelLogout,

		RootUrl: "http://localhost:2222/" + acctest.RandString(20),
		ValidRedirectUris: []string{
			acctest.RandString(20),
		},
		BaseUrl:                 "http://localhost:2222/" + acctest.RandString(20),
		MasterSamlProcessingUrl: acctest.RandString(20),

		Attributes: &keycloak.SamlClientAttributes{
			IncludeAuthnStatement:           randomBoolAsStringPointer(),
			SignDocuments:                   randomBoolAsStringPointer(),
			SignAssertions:                  randomBoolAsStringPointer(),
			EncryptAssertions:               randomBoolAsStringPointer(),
			ClientSignatureRequired:         &clientSignatureRequired,
			ForcePostBinding:                randomBoolAsStringPointer(),
			ForceNameIdFormat:               randomBoolAsStringPointer(),
			SignatureAlgorithm:              randomStringInSlice(keycloakSamlClientSignatureAlgorithms),
			NameIdFormat:                    randomStringInSlice(keycloakSamlClientNameIdFormats),
			EncryptionCertificate:           &encryptionCertificateAfter,
			SigningCertificate:              &signingCertificateAfter,
			SigningPrivateKey:               &signingPrivateKeyAfter,
			IDPInitiatedSSOURLName:          acctest.RandString(20),
			IDPInitiatedSSORelayState:       acctest.RandString(20),
			AssertionConsumerPostURL:        acctest.RandString(20),
			AssertionConsumerRedirectURL:    acctest.RandString(20),
			LogoutServicePostBindingURL:     acctest.RandString(20),
			LogoutServiceRedirectBindingURL: acctest.RandString(20),
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_fromInterface(samlClientBefore),
				Check:  testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
			},
			{
				Config: testKeycloakSamlClient_fromInterface(samlClientAfter),
				Check:  testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
			},
		},
	})
}

func TestAccKeycloakSamlClient_certificateAndKey(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_signingCertificateAndKey(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientHasSigningCertificate("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientHasPrivateKey("keycloak_saml_client.saml_client"),
				),
			},
			{
				Config: testKeycloakSamlClient_signingCertificateNoKey(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientHasSigningCertificate("keycloak_saml_client.saml_client"),
					resource.TestCheckResourceAttr("keycloak_saml_client.saml_client", "signing_private_key", ""),
				),
			},
		},
	})
}

func TestAccKeycloakSamlClient_encryptionCertificate(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_encryptionCertificate(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					testAccCheckKeycloakSamlClientHasEncryptionCertificate("keycloak_saml_client.saml_client"),
				),
			},
			{
				Config: testKeycloakSamlClient_NoEncryptionCertificate(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakSamlClientExistsWithCorrectProtocol("keycloak_saml_client.saml_client"),
					resource.TestCheckResourceAttr("keycloak_saml_client.saml_client", "encryption_certificate", ""),
				),
			},
		},
	})
}

func TestAccCheckKeycloakSamlClient_authenticationFlowBindingOverrides(t *testing.T) {
	t.Parallel()

	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakSamlClient_authenticationFlowBindingOverrides(clientId),
				Check:  testAccCheckKeycloakSamlClientAuthenticationFlowBindingOverrides("keycloak_saml_client.client", "keycloak_authentication_flow.another_flow"),
			},
			{
				Config: testKeycloakSamlClient_withoutFlowBindingOverrides(clientId),
				Check:  testAccCheckKeycloakSamlClientAuthenticationFlowBindingOverrides("keycloak_saml_client.client", ""),
			},
		},
	})
}

func testAccCheckKeycloakSamlClientExistsWithCorrectProtocol(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Protocol != "saml" {
			return fmt.Errorf("expected saml client to have saml protocol, but got %s", client.Protocol)
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientHasEncryptionCertificate(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if *client.Attributes.EncryptionCertificate == "" {
			return fmt.Errorf("expected saml client to have a encryption certificate")
		}

		if strings.Contains(*client.Attributes.EncryptionCertificate, "-----BEGIN CERTIFICATE-----") || strings.Contains(*client.Attributes.EncryptionCertificate, "-----END CERTIFICATE-----") {
			return fmt.Errorf("expected saml client encryption certificate to not contain headers")
		}

		if strings.ContainsAny(*client.Attributes.EncryptionCertificate, "\n\r ") {
			return fmt.Errorf("expected saml client encryption certificate to not contain whitespace")
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientHasSigningCertificate(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if *client.Attributes.SigningCertificate == "" {
			return fmt.Errorf("expected saml client to have a signing certificate")
		}

		if strings.Contains(*client.Attributes.SigningCertificate, "-----BEGIN CERTIFICATE-----") || strings.Contains(*client.Attributes.SigningCertificate, "-----END CERTIFICATE-----") {
			return fmt.Errorf("expected saml client signing certificate to not contain headers")
		}

		if strings.ContainsAny(*client.Attributes.SigningCertificate, "\n\r ") {
			return fmt.Errorf("expected saml client signing certificate to not contain whitespace")
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientHasPrivateKey(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if *client.Attributes.SigningPrivateKey == "" {
			return fmt.Errorf("expected saml client to have a signing private key")
		}

		if strings.Contains(*client.Attributes.SigningPrivateKey, "-----BEGIN PRIVATE KEY-----") || strings.Contains(*client.Attributes.SigningPrivateKey, "-----END PRIVATE KEY-----") {
			return fmt.Errorf("expected saml client signing private key to not contain headers")
		}

		if strings.ContainsAny(*client.Attributes.SigningPrivateKey, "\n\r ") {
			return fmt.Errorf("expected saml client signing private key to not contain whitespace")
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientFetch(resourceName string, client *keycloak.SamlClient) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClient, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		client.Id = fetchedClient.Id
		client.RealmId = fetchedClient.RealmId

		return nil
	}
}

func testAccCheckKeycloakSamlClientDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_saml_client" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			client, _ := keycloakClient.GetSamlClient(realm, id)
			if client != nil {
				return fmt.Errorf("saml client %s still exists", id)
			}
		}

		return nil
	}
}

func getSamlClientFromState(s *terraform.State, resourceName string) (*keycloak.SamlClient, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetSamlClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting saml client %s: %s", id, err)
	}

	return client, nil
}

func testAccCheckKeycloakSamlClientHasDefaultBooleanAttributes(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		includeAuthnStatement, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["include_authn_statement"])
		if err != nil {
			return err
		}

		signDocuments, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["sign_documents"])
		if err != nil {
			return err
		}

		signAssertions, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["sign_assertions"])
		if err != nil {
			return err
		}

		encryptAssertions, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["encrypt_assertions"])
		if err != nil {
			return err
		}

		clientSignatureRequired, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["client_signature_required"])
		if err != nil {
			return err
		}

		forcePostBinding, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["force_post_binding"])
		if err != nil {
			return err
		}

		forceNameIdFormat, err := parseBoolAndTreatEmptyStringAsFalse(rs.Primary.Attributes["force_name_id_format"])
		if err != nil {
			return err
		}

		if !includeAuthnStatement && !signDocuments && !signAssertions && !encryptAssertions && !clientSignatureRequired && !forcePostBinding && !forceNameIdFormat {
			return fmt.Errorf("expected saml client with id %s to have some defaults set by Keycloak", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckKeycloakSamlClientAuthenticationFlowBindingOverrides(resourceName, flowResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getSamlClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if flowResourceName == "" {
			if client.AuthenticationFlowBindingOverrides.BrowserId != "" {
				return fmt.Errorf("expected openid client to have browserId set to empty, but got %s", client.AuthenticationFlowBindingOverrides.BrowserId)
			}

			if client.AuthenticationFlowBindingOverrides.DirectGrantId != "" {
				return fmt.Errorf("expected openid client to have directGrantId set to empty, but got %s", client.AuthenticationFlowBindingOverrides.DirectGrantId)
			}

		} else {
			flow, err := getAuthenticationFlowFromState(s, flowResourceName)
			if err != nil {
				return err
			}

			if client.AuthenticationFlowBindingOverrides.BrowserId != flow.Id {
				return fmt.Errorf("expected openid client to have browserId set to %s, but got %s", flow.Id, client.AuthenticationFlowBindingOverrides.BrowserId)
			}

			if client.AuthenticationFlowBindingOverrides.DirectGrantId != flow.Id {
				return fmt.Errorf("expected openid client to have directGrantId set to %s, but got %s", flow.Id, client.AuthenticationFlowBindingOverrides.DirectGrantId)
			}
		}

		return nil
	}
}

func parseBoolAndTreatEmptyStringAsFalse(b string) (bool, error) {
	if b == "" {
		return false, nil
	}

	return strconv.ParseBool(b)
}

func randomBoolAsStringPointer() *string {
	s := strconv.FormatBool(randomBool())

	return &s
}

func testKeycloakSamlClient_basic(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = data.keycloak_realm.realm.id
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakSamlClient_updateRealmBefore(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm_1" {
	realm = "%s"
}

data "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = data.keycloak_realm.realm_1.id
}
	`, testAccRealm.Realm, testAccRealmTwo.Realm, clientId)
}

func testKeycloakSamlClient_updateRealmAfter(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm_1" {
	realm = "%s"
}

data "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = data.keycloak_realm.realm_2.id
}
	`, testAccRealm.Realm, testAccRealmTwo.Realm, clientId)
}

func testKeycloakSamlClient_fromInterface(client *keycloak.SamlClient) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	realm_id    = data.keycloak_realm.realm.id
	client_id   = "%s"
	name        = "%s"
	description = "%s"
	enabled     = %t

	# below attributes are bools, but the model (and API) uses strings
	include_authn_statement    = %s
	sign_documents             = %s
	sign_assertions            = %s
	encrypt_assertions         = %s
	client_signature_required  = %s
	force_post_binding         = %s
	force_name_id_format       = %s

	front_channel_logout       = %t
	signature_algorithm        = "%s"
	name_id_format             = "%s"
	root_url                   = "%s"
	valid_redirect_uris        = %s
	base_url                   = "%s"
	master_saml_processing_url = "%s"

	encryption_certificate     = "%s"
	signing_certificate        = "%s"
	signing_private_key        = "%s"

	idp_initiated_sso_url_name    = "%s"
	idp_initiated_sso_relay_state = "%s"

	assertion_consumer_post_url         = "%s"
	assertion_consumer_redirect_url     = "%s"
	logout_service_post_binding_url     = "%s"
	logout_service_redirect_binding_url = "%s"
}
	`, client.RealmId,
		client.ClientId,
		client.Name,
		client.Description,
		client.Enabled,
		*client.Attributes.IncludeAuthnStatement,
		*client.Attributes.SignDocuments,
		*client.Attributes.SignAssertions,
		*client.Attributes.EncryptAssertions,
		*client.Attributes.ClientSignatureRequired,
		*client.Attributes.ForcePostBinding,
		*client.Attributes.ForceNameIdFormat,
		client.FrontChannelLogout,
		client.Attributes.SignatureAlgorithm,
		client.Attributes.NameIdFormat,
		client.RootUrl,
		arrayOfStringsForTerraformResource(client.ValidRedirectUris),
		client.BaseUrl, client.MasterSamlProcessingUrl,
		*client.Attributes.EncryptionCertificate,
		*client.Attributes.SigningCertificate,
		*client.Attributes.SigningPrivateKey,
		client.Attributes.IDPInitiatedSSOURLName,
		client.Attributes.IDPInitiatedSSORelayState,
		client.Attributes.AssertionConsumerPostURL,
		client.Attributes.AssertionConsumerRedirectURL,
		client.Attributes.LogoutServicePostBindingURL,
		client.Attributes.LogoutServiceRedirectBindingURL,
	)
}

func testKeycloakSamlClient_signingCertificateAndKey(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id               = "%s"
	realm_id                = data.keycloak_realm.realm.id
	name                    = "test-saml-client"

	sign_documents          = false
	sign_assertions         = true
	encrypt_assertions      = false
	include_authn_statement = true

	signing_certificate     = file("misc/saml-cert.pem")
	signing_private_key     = file("misc/saml-key.pem")
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakSamlClient_signingCertificateNoKey(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id               = "%s"
	realm_id                = data.keycloak_realm.realm.id
	name                    = "test-saml-client"

	sign_documents          = false
	sign_assertions         = true
	encrypt_assertions      = false
	include_authn_statement = true

	signing_certificate     = file("misc/saml-cert.pem")
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakSamlClient_encryptionCertificate(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id               = "%s"
	realm_id                = data.keycloak_realm.realm.id
	name                    = "test-saml-client"

	encrypt_assertions      = true
	include_authn_statement = true

	encryption_certificate     = file("misc/saml-cert.pem")
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakSamlClient_NoEncryptionCertificate(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id               = "%s"
	realm_id                = data.keycloak_realm.realm.id
	name                    = "test-saml-client"

	encrypt_assertions      = true
	include_authn_statement = true
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakSamlClient_authenticationFlowBindingOverrides(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "another_flow" {
  alias       = "anotherFlow"
  realm_id    = data.keycloak_realm.realm.id
  description = "this is another flow"
}

resource "keycloak_saml_client" "client" {
	client_id = "%s"
	realm_id  = data.keycloak_realm.realm.id
	name      = "test-saml-client"

	authentication_flow_binding_overrides {
		browser_id      = keycloak_authentication_flow.another_flow.id
		direct_grant_id = keycloak_authentication_flow.another_flow.id
	}
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakSamlClient_withoutFlowBindingOverrides(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "another_flow" {
  alias       = "anotherFlow"
  realm_id    = data.keycloak_realm.realm.id
  description = "this is another flow"
}

resource "keycloak_saml_client" "client" {
	client_id = "%s"
	realm_id  = data.keycloak_realm.realm.id
	name      = "test-saml-client"
}
	`, testAccRealm.Realm, clientId)
}
