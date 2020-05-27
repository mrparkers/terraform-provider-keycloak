package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				ResourceName:            "keycloak_openid_client.client",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     realmName + "/",
				ImportStateVerifyIgnore: []string{"exclude_session_state_from_auth_response"},
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

func TestAccKeycloakOpenidClient_adminUrl(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	adminUrl := "https://www.example.com/admin"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_adminUrl(realmName, clientId, adminUrl),
				Check:  testAccCheckKeycloakOpenidClientAdminUrl("keycloak_openid_client.client", adminUrl),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_baseUrl(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	baseUrl := "https://www.example.com"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_baseUrl(realmName, clientId, baseUrl),
				Check:  testAccCheckKeycloakOpenidClientBaseUrl("keycloak_openid_client.client", baseUrl),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_rootUrl(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	rootUrl := "https://www.example.com"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_rootUrl(realmName, clientId, rootUrl),
				Check:  testAccCheckKeycloakOpenidClientRootUrl("keycloak_openid_client.client", rootUrl),
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

	if !standardFlowEnabled {
		implicitFlowEnabled = !standardFlowEnabled
	}

	rootUrlBefore := "http://localhost:2222/" + acctest.RandString(20)
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
		AdminUrl:                  acctest.RandString(20),
		BaseUrl:                   "http://localhost:2222/" + acctest.RandString(20),
		RootUrl:                   &rootUrlBefore,
	}

	standardFlowEnabled, implicitFlowEnabled = implicitFlowEnabled, standardFlowEnabled

	rootUrlAfter := "http://localhost:2222/" + acctest.RandString(20)
	openidClientAfter := &keycloak.OpenidClient{
		RealmId:                   realm,
		ClientId:                  clientId,
		Name:                      acctest.RandString(10),
		Enabled:                   !enabled,
		Description:               acctest.RandString(50),
		ClientSecret:              acctest.RandString(10),
		StandardFlowEnabled:       standardFlowEnabled,
		ImplicitFlowEnabled:       implicitFlowEnabled,
		DirectAccessGrantsEnabled: !directAccessGrantsEnabled,
		ServiceAccountsEnabled:    !serviceAccountsEnabled,
		ValidRedirectUris:         []string{acctest.RandString(10), acctest.RandString(10)},
		WebOrigins:                []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
		AdminUrl:                  acctest.RandString(20),
		BaseUrl:                   "http://localhost:2222/" + acctest.RandString(20),
		RootUrl:                   &rootUrlAfter,
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

func TestAccKeycloakOpenidClient_AccessToken_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	accessTokenLifespan := "1800"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_AccessToken_basic(realmName, clientId, accessTokenLifespan),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectLifespan("keycloak_openid_client.client", accessTokenLifespan),
			},
			{
				ResourceName:            "keycloak_openid_client.client",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     realmName + "/",
				ImportStateVerifyIgnore: []string{"exclude_session_state_from_auth_response"},
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

func TestAccKeycloakOpenidClient_pkceCodeChallengeMethod(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_pkceChallengeMethod(realmName, clientId, "invalidMethod"),
				ExpectError: regexp.MustCompile(`config is invalid: expected pkce_code_challenge_method to be one of \[\ plain S256\], got invalidMethod`),
			},
			{
				Config: testKeycloakOpenidClient_omitPkceChallengeMethod(realmName, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
			{
				Config: testKeycloakOpenidClient_pkceChallengeMethod(realmName, clientId, "plain"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", "plain"),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
			{
				Config: testKeycloakOpenidClient_pkceChallengeMethod(realmName, clientId, "S256"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", "S256"),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
			{
				Config: testKeycloakOpenidClient_pkceChallengeMethod(realmName, clientId, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_excludeSessionStateFromAuthResponse(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_omitExcludeSessionStateFromAuthResponse(realmName, clientId, "plain"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", "plain"),
				),
			},
			{
				Config: testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(realmName, clientId, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
				),
			},
			{
				Config: testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(realmName, clientId, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", true),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
				),
			},
			{
				Config: testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(realmName, clientId, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_authenticationFlowBindingOverrides(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_authenticationFlowBindingOverrides(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientAuthenticationFlowBindingOverrides("keycloak_openid_client.client", "keycloak_authentication_flow.another_flow"),
			},
			{
				Config: testKeycloakOpenidClient_withoutAuthenticationFlowBindingOverrides(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientAuthenticationFlowBindingOverrides("keycloak_openid_client.client", ""),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_loginTheme(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	loginThemeKeycloak := "keycloak"
	loginThemeBase := "base"
	loginThemeRandom := "theme-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_loginTheme(realmName, clientId, loginThemeKeycloak),
				Check:  testAccCheckKeycloakOpenidClientLoginTheme("keycloak_openid_client.client", loginThemeKeycloak),
			},
			{
				Config: testKeycloakOpenidClient_loginTheme(realmName, clientId, loginThemeBase),
				Check:  testAccCheckKeycloakOpenidClientLoginTheme("keycloak_openid_client.client", loginThemeBase),
			},
			{
				Config:      testKeycloakOpenidClient_loginTheme(realmName, clientId, loginThemeRandom),
				ExpectError: regexp.MustCompile("validation error: theme \".+\" does not exist on the server"),
			},
			{
				Config: testKeycloakOpenidClient_loginTheme(realmName, clientId, loginThemeKeycloak),
				Check:  testAccCheckKeycloakOpenidClientLoginTheme("keycloak_openid_client.client", loginThemeKeycloak),
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

func testAccCheckKeycloakOpenidClientExistsWithCorrectLifespan(resourceName string, accessTokenLifespan string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Attributes.AccessTokenLifespan != accessTokenLifespan {
			return fmt.Errorf("expected openid client to have access token lifespan set to %s, but got %s", accessTokenLifespan, client.Attributes.AccessTokenLifespan)
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

func testAccCheckKeycloakOpenidClientBaseUrl(resourceName string, baseUrl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.BaseUrl != baseUrl {
			return fmt.Errorf("expected openid client to have baseUrl set to %s, but got %s", baseUrl, client.BaseUrl)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientRootUrl(resourceName string, rootUrl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if *client.RootUrl != rootUrl {
			return fmt.Errorf("expected openid client to have rootUrl set to %s, but got %s", rootUrl, *client.RootUrl)
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

func testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod(resourceName, pkceCodeChallengeMethod string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Attributes.PkceCodeChallengeMethod != pkceCodeChallengeMethod {
			return fmt.Errorf("expected openid client %s to have pkce code challenge method value of %s, but got %s", client.ClientId, pkceCodeChallengeMethod, client.Attributes.PkceCodeChallengeMethod)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse(resourceName string, excludeSessionStateFromAuthResponse keycloak.KeycloakBoolQuoted) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Attributes.ExcludeSessionStateFromAuthResponse != excludeSessionStateFromAuthResponse {
			return fmt.Errorf("expected openid client %s to have exclude_session_state_from_auth_response value of %t, but got %t", client.ClientId, excludeSessionStateFromAuthResponse, client.Attributes.ExcludeSessionStateFromAuthResponse)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAdminUrl(resourceName string, adminUrl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.AdminUrl != adminUrl {
			return fmt.Errorf("expected openid client to have adminUrl set to %s, but got %s", adminUrl, client.AdminUrl)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientAuthenticationFlowBindingOverrides(resourceName, flowResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
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

func testAccCheckKeycloakOpenidClientLoginTheme(resourceName string, loginTheme string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Attributes.LoginTheme != loginTheme {
			return fmt.Errorf("expected openid client to have login theme set to %s, but got %s", loginTheme, client.Attributes.LoginTheme)
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

func testKeycloakOpenidClient_AccessToken_basic(realm, clientId, accessTokenLifespan string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   		  = "%s"
	realm_id    		  = "${keycloak_realm.realm.id}"
	access_type 		  = "CONFIDENTIAL"
	access_token_lifespan = "%s"
}
	`, realm, clientId, accessTokenLifespan)
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

func testKeycloakOpenidClient_adminUrl(realm, clientId, adminUrl string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	admin_url = "%s"
	access_type = "PUBLIC"
}
	`, realm, clientId, adminUrl)
}

func testKeycloakOpenidClient_baseUrl(realm, clientId, baseUrl string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	base_url = "%s"
	access_type = "PUBLIC"
}
	`, realm, clientId, baseUrl)
}

func testKeycloakOpenidClient_rootUrl(realm, clientId, rootUrl string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id             = "%s"
	realm_id			  = "${keycloak_realm.realm.id}"
	root_url			  = "%s"
	valid_redirect_uris   = ["http://example.com"]
	web_origins			  = ["http://example.com"]
	admin_url			  = "http://example.com"
	access_type           = "CONFIDENTIAL"
	standard_flow_enabled = true
}
	`, realm, clientId, rootUrl)
}

func testKeycloakOpenidClient_pkceChallengeMethod(realm, clientId, pkceChallengeMethod string) string {

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
	pkce_code_challenge_method = "%s"
}
	`, realm, clientId, pkceChallengeMethod)
}

func testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(realm, clientId string, excludeSessionStateFromAuthResponse bool) string {

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
	exclude_session_state_from_auth_response = %t
}
	`, realm, clientId, excludeSessionStateFromAuthResponse)
}

func testKeycloakOpenidClient_omitPkceChallengeMethod(realm, clientId string) string {

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

func testKeycloakOpenidClient_omitExcludeSessionStateFromAuthResponse(realm, clientId, pkceChallengeMethod string) string {

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
    pkce_code_challenge_method = "%s"
}
	`, realm, clientId, pkceChallengeMethod)
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
	admin_url					 = "%s"
	base_url                     = "%s"
	root_url                     = "%s"
}
	`, openidClient.RealmId, openidClient.ClientId, openidClient.Name, openidClient.Enabled, openidClient.Description, openidClient.ClientSecret, openidClient.StandardFlowEnabled, openidClient.ImplicitFlowEnabled, openidClient.DirectAccessGrantsEnabled, openidClient.ServiceAccountsEnabled, arrayOfStringsForTerraformResource(openidClient.ValidRedirectUris), arrayOfStringsForTerraformResource(openidClient.WebOrigins), openidClient.AdminUrl, openidClient.BaseUrl, *openidClient.RootUrl)
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

func testKeycloakOpenidClient_authenticationFlowBindingOverrides(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "another_flow" {
  alias    = "anotherFlow"
  realm_id = "${keycloak_realm.realm.id}"
  description = "this is another flow"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
	authentication_flow_binding_overrides {
		browser_id = "${keycloak_authentication_flow.another_flow.id}"
		direct_grant_id = "${keycloak_authentication_flow.another_flow.id}"
	}
}
	`, realm, clientId)
}

func testKeycloakOpenidClient_withoutAuthenticationFlowBindingOverrides(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "another_flow" {
  alias    = "anotherFlow"
  realm_id = "${keycloak_realm.realm.id}"
  description = "this is another flow"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
}
	`, realm, clientId)
}

func testKeycloakOpenidClient_loginTheme(realm, clientId, loginTheme string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
	login_theme = "%s"
}
	`, realm, clientId, loginTheme)
}
