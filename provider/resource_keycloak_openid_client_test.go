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

func TestAccKeycloakOpenidClient_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
			},
			{
				ResourceName:        "keycloak_openid_client.client",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakOpenidClient_createAfterManualDestroy(t *testing.T) {
	var client = &keycloak.OpenidClient{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientFetch("keycloak_openid_client.client", client),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClient(client.RealmId, client.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_updateRealmBefore(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientBelongsToRealm("keycloak_openid_client.client", realmOne),
				),
			},
			{
				Config: testKeycloakOpenidClient_updateRealmAfter(realmOne, realmTwo, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientBelongsToRealm("keycloak_openid_client.client", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_accessType(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_accessType(realmName, clientId, "CONFIDENTIAL"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", false, false),
			},
			{
				Config: testKeycloakOpenidClient_accessType(realmName, clientId, "PUBLIC"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", true, false),
			},
			{
				Config: testKeycloakOpenidClient_accessType(realmName, clientId, "BEARER-ONLY"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", false, true),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_updateInPlace(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	enabled := randomBool()
	standardFlowEnabled := randomBool()
	implicitFlowEnabled := randomBool()
	directAccessGrantsEnabled := randomBool()
	serviceAccountsEnabled := randomBool()

	openidClientBefore := &keycloak.OpenidClient{
		RealmId:                   realm,
		ClientId:                  clientId,
		Name:                      acctest.RandString(10),
		Enabled:                   enabled,
		Description:               acctest.RandString(50),
		ClientSecret:              acctest.RandString(10),
		StandardFlowEnabled:       standardFlowEnabled,
		ImplicitFlowEnabled:       implicitFlowEnabled,
		DirectAccessGrantsEnabled: directAccessGrantsEnabled,
		ServiceAccountsEnabled:    serviceAccountsEnabled,
		ValidRedirectUris:         []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
		WebOrigins:                []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
	}

	openidClientAfter := &keycloak.OpenidClient{
		RealmId:                   realm,
		ClientId:                  clientId,
		Name:                      acctest.RandString(10),
		Enabled:                   !enabled,
		Description:               acctest.RandString(50),
		ClientSecret:              acctest.RandString(10),
		StandardFlowEnabled:       !standardFlowEnabled,
		ImplicitFlowEnabled:       !implicitFlowEnabled,
		DirectAccessGrantsEnabled: !directAccessGrantsEnabled,
		ServiceAccountsEnabled:    !serviceAccountsEnabled,
		ValidRedirectUris:         []string{acctest.RandString(10), acctest.RandString(10)},
		WebOrigins:                []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_fromInterface(openidClientBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
				),
			},
			{
				Config: testKeycloakOpenidClient_fromInterface(openidClientAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
				),
			},
			{
				Config: testKeycloakOpenidClient_basic(realm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_secret(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	clientSecret := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientHasNonEmptyClientSecret("keycloak_openid_client.client"),
				),
			},
			{
				Config: testKeycloakOpenidClient_secret(realmName, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientHasClientSecret("keycloak_openid_client.client", clientSecret),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_redirectUrisValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	accessType := randomStringInSlice([]string{"PUBLIC", "CONFIDENTIAL"})

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_invalidRedirectUris(realmName, clientId, accessType, true, false),
				ExpectError: regexp.MustCompile("validation error: standard \\(authorization code\\) and implicit flows require at least one valid redirect uri"),
			},
			{
				Config:      testKeycloakOpenidClient_invalidRedirectUris(realmName, clientId, accessType, false, true),
				ExpectError: regexp.MustCompile("validation error: standard \\(authorization code\\) and implicit flows require at least one valid redirect uri"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_publicClientCredentialsValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_invalidPublicClientWithClientCredentials(realmName, clientId),
				ExpectError: regexp.MustCompile("validation error: service accounts \\(client credentials flow\\) cannot be enabled on public clients"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_bearerClientNoGrantsValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(realmName, clientId, true, false, false, false),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(realmName, clientId, false, true, false, false),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(realmName, clientId, false, false, true, false),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(realmName, clientId, false, false, false, true),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Protocol != "openid-connect" {
			return fmt.Errorf("expected openid client to have openid-connect protocol, but got %s", client.Protocol)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientFetch(resourceName string, client *keycloak.OpenidClient) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClient, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		client.Id = fetchedClient.Id
		client.RealmId = fetchedClient.RealmId

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAccessType(resourceName string, public, bearer bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.PublicClient != public {
			return fmt.Errorf("expected openid client to have public set to %t, but got %t", public, client.PublicClient)
		}

		if client.BearerOnly != bearer {
			return fmt.Errorf("expected openid client to have bearer set to %t, but got %t", bearer, client.BearerOnly)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.RealmId != realm {
			return fmt.Errorf("expected openid client %s to have realm_id of %s, but got %s", client.ClientId, realm, client.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientHasClientSecret(resourceName, secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.ClientSecret != secret {
			return fmt.Errorf("expected openid client %s to have secret value of %s, but got %s", client.ClientId, secret, client.ClientSecret)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientHasNonEmptyClientSecret(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.ClientSecret == "" {
			return fmt.Errorf("expected openid client %s to have non empty secret value", client.ClientId)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			client, _ := keycloakClient.GetOpenidClient(realm, id)
			if client != nil {
				return fmt.Errorf("openid client %s still exists", id)
			}
		}

		return nil
	}
}

func getOpenidClientFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClient, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetOpenidClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client %s: %s", id, err)
	}

	return client, nil
}

func testKeycloakOpenidClient_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}
	`, realm, clientId)
}

func testKeycloakOpenidClient_accessType(realm, clientId, accessType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "%s"
}
	`, realm, clientId, accessType)
}

func testKeycloakOpenidClient_updateRealmBefore(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm_1.id}"
	access_type = "BEARER-ONLY"
}
	`, realmOne, realmTwo, clientId)
}

func testKeycloakOpenidClient_updateRealmAfter(realmOne, realmTwo, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm_2.id}"
	access_type = "BEARER-ONLY"
}
	`, realmOne, realmTwo, clientId)
}

func testKeycloakOpenidClient_fromInterface(openidClient *keycloak.OpenidClient) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id                    = "%s"
	realm_id                     = "${keycloak_realm.realm.id}"
	name                         = "%s"
	enabled                      = %t
	description                  = "%s"

	access_type                  = "CONFIDENTIAL"
	client_secret                = "%s"

	standard_flow_enabled        = %t
	implicit_flow_enabled        = %t
	direct_access_grants_enabled = %t
	service_accounts_enabled     = %t

	valid_redirect_uris          = %s
	web_origins                  = %s
}
	`, openidClient.RealmId, openidClient.ClientId, openidClient.Name, openidClient.Enabled, openidClient.Description, openidClient.ClientSecret, openidClient.StandardFlowEnabled, openidClient.ImplicitFlowEnabled, openidClient.ServiceAccountsEnabled, openidClient.DirectAccessGrantsEnabled, arrayOfStringsForTerraformResource(openidClient.ValidRedirectUris), arrayOfStringsForTerraformResource(openidClient.WebOrigins))
}

func testKeycloakOpenidClient_secret(realm, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id     = "%s"
	realm_id      = "${keycloak_realm.realm.id}"
	access_type   = "CONFIDENTIAL"
	client_secret = "%s"
}
	`, realm, clientId, clientSecret)
}

func testKeycloakOpenidClient_invalidRedirectUris(realm, clientId, accessType string, standardFlowEnabled, implicitFlowEnabled bool) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id             = "%s"
	realm_id              = "${keycloak_realm.realm.id}"
	access_type           = "%s"

	standard_flow_enabled = %t
	implicit_flow_enabled = %t
}
	`, realm, clientId, accessType, standardFlowEnabled, implicitFlowEnabled)
}

func testKeycloakOpenidClient_invalidPublicClientWithClientCredentials(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id                = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	access_type              = "PUBLIC"

	service_accounts_enabled = true
}
	`, realm, clientId)
}

func testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(realm, clientId string, standardFlowEnabled, implicitFlowEnabled, directAccessGrantsEnabled, serviceAccountsEnabled bool) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id                    = "%s"
	realm_id                     = "${keycloak_realm.realm.id}"
	access_type                  = "BEARER-ONLY"

	standard_flow_enabled        = %t
	implicit_flow_enabled        = %t
	direct_access_grants_enabled = %t
	service_accounts_enabled     = %t
}
	`, realm, clientId, standardFlowEnabled, implicitFlowEnabled, directAccessGrantsEnabled, serviceAccountsEnabled)
}
