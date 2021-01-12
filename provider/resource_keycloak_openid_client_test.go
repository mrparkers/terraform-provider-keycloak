package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenidClient_basic(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(clientId),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
			},
			{
				ResourceName:            "keycloak_openid_client.client",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     testAccRealm.Realm + "/",
				ImportStateVerifyIgnore: []string{"exclude_session_state_from_auth_response"},
			},
		},
	})
}

func TestAccKeycloakOpenidClient_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	var client = &keycloak.OpenidClient{}

	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientFetch("keycloak_openid_client.client", client),
				),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteOpenidClient(client.RealmId, client.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClient_basic(clientId),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_updateRealm(t *testing.T) {
	t.Parallel()

	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_updateRealmBefore(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientBelongsToRealm("keycloak_openid_client.client", testAccRealm.Realm),
				),
			},
			{
				Config: testKeycloakOpenidClient_updateRealmAfter(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientBelongsToRealm("keycloak_openid_client.client", testAccRealmTwo.Realm),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_accessType(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_accessType(clientId, "CONFIDENTIAL"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", false, false),
			},
			{
				Config: testKeycloakOpenidClient_accessType(clientId, "PUBLIC"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", true, false),
			},
			{
				Config: testKeycloakOpenidClient_accessType(clientId, "BEARER-ONLY"),
				Check:  testAccCheckKeycloakOpenidClientAccessType("keycloak_openid_client.client", false, true),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_adminUrl(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	adminUrl := "https://www.example.com/admin"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_adminUrl(clientId, adminUrl),
				Check:  testAccCheckKeycloakOpenidClientAdminUrl("keycloak_openid_client.client", adminUrl),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_baseUrl(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	baseUrl := "https://www.example.com"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_baseUrl(clientId, baseUrl),
				Check:  testAccCheckKeycloakOpenidClientBaseUrl("keycloak_openid_client.client", baseUrl),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_rootUrl(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	rootUrl := "https://www.example.com"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_rootUrl(clientId, rootUrl),
				Check:  testAccCheckKeycloakOpenidClientRootUrl("keycloak_openid_client.client", rootUrl),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_updateInPlace(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
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
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
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
				Config: testKeycloakOpenidClient_basic(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_AccessToken_basic(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	accessTokenLifespan := "1800"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_AccessToken_basic(clientId, accessTokenLifespan),
				Check:  testAccCheckKeycloakOpenidClientExistsWithCorrectLifespan("keycloak_openid_client.client", accessTokenLifespan),
			},
			{
				ResourceName:            "keycloak_openid_client.client",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     testAccRealm.Realm + "/",
				ImportStateVerifyIgnore: []string{"exclude_session_state_from_auth_response"},
			},
		},
	})
}

func TestAccKeycloakOpenidClient_ClientTimeouts_basic(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	offlineSessionIdleTimeout := "1800"
	offlineSessionMaxLifespan := "1900"
	sessionIdleTimeout := "2000"
	sessionMaxLifespan := "2100"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_ClientTimeouts(clientId,
					offlineSessionIdleTimeout, offlineSessionMaxLifespan, sessionIdleTimeout, sessionMaxLifespan),
				Check: testAccCheckKeycloakOpenidClientExistsWithCorrectClientTimeouts("keycloak_openid_client.client",
					offlineSessionIdleTimeout, offlineSessionMaxLifespan, sessionIdleTimeout, sessionMaxLifespan,
				),
			},
			{
				ResourceName:            "keycloak_openid_client.client",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     testAccRealm.Realm + "/",
				ImportStateVerifyIgnore: []string{"exclude_session_state_from_auth_response"},
			},
		},
	})
}

func TestAccKeycloakOpenidClient_secret(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	clientSecret := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientHasNonEmptyClientSecret("keycloak_openid_client.client"),
				),
			},
			{
				Config: testKeycloakOpenidClient_secret(clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientExistsWithCorrectProtocol("keycloak_openid_client.client"),
					testAccCheckKeycloakOpenidClientHasClientSecret("keycloak_openid_client.client", clientSecret),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_redirectUrisValidation(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	accessType := randomStringInSlice([]string{"PUBLIC", "CONFIDENTIAL"})

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_invalidRedirectUris(clientId, accessType, true, false),
				ExpectError: regexp.MustCompile("validation error: standard \\(authorization code\\) and implicit flows require at least one valid redirect uri"),
			},
			{
				Config:      testKeycloakOpenidClient_invalidRedirectUris(clientId, accessType, false, true),
				ExpectError: regexp.MustCompile("validation error: standard \\(authorization code\\) and implicit flows require at least one valid redirect uri"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_publicClientCredentialsValidation(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_invalidPublicClientWithClientCredentials(clientId),
				ExpectError: regexp.MustCompile("validation error: service accounts \\(client credentials flow\\) cannot be enabled on public clients"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_bearerClientNoGrantsValidation(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(clientId, true, false, false, false),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(clientId, false, true, false, false),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(clientId, false, false, true, false),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
			{
				Config:      testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(clientId, false, false, false, true),
				ExpectError: regexp.MustCompile("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client"),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_pkceCodeChallengeMethod(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClient_pkceChallengeMethod(clientId, "invalidMethod"),
				ExpectError: regexp.MustCompile(`expected pkce_code_challenge_method to be one of \[\ plain S256\], got invalidMethod`),
			},
			{
				Config: testKeycloakOpenidClient_omitPkceChallengeMethod(clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
			{
				Config: testKeycloakOpenidClient_pkceChallengeMethod(clientId, "plain"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", "plain"),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
			{
				Config: testKeycloakOpenidClient_pkceChallengeMethod(clientId, "S256"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", "S256"),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
			{
				Config: testKeycloakOpenidClient_pkceChallengeMethod(clientId, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_excludeSessionStateFromAuthResponse(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_omitExcludeSessionStateFromAuthResponse(clientId, "plain"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", "plain"),
				),
			},
			{
				Config: testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(clientId, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
				),
			},
			{
				Config: testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(clientId, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", true),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
				),
			},
			{
				Config: testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(clientId, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasExcludeSessionStateFromAuthResponse("keycloak_openid_client.client", false),
					testAccCheckKeycloakOpenidClientHasPkceCodeChallengeMethod("keycloak_openid_client.client", ""),
				),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_authenticationFlowBindingOverrides(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_authenticationFlowBindingOverrides(clientId),
				Check:  testAccCheckKeycloakOpenidClientAuthenticationFlowBindingOverrides("keycloak_openid_client.client", "keycloak_authentication_flow.another_flow"),
			},
			{
				Config: testKeycloakOpenidClient_withoutAuthenticationFlowBindingOverrides(clientId),
				Check:  testAccCheckKeycloakOpenidClientAuthenticationFlowBindingOverrides("keycloak_openid_client.client", ""),
			},
		},
	})
}

func TestAccKeycloakOpenidClient_loginTheme(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	loginThemeKeycloak := "keycloak"
	loginThemeBase := "base"
	loginThemeRandom := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_loginTheme(clientId, loginThemeKeycloak),
				Check:  testAccCheckKeycloakOpenidClientLoginTheme("keycloak_openid_client.client", loginThemeKeycloak),
			},
			{
				Config: testKeycloakOpenidClient_loginTheme(clientId, loginThemeBase),
				Check:  testAccCheckKeycloakOpenidClientLoginTheme("keycloak_openid_client.client", loginThemeBase),
			},
			{
				Config:      testKeycloakOpenidClient_loginTheme(clientId, loginThemeRandom),
				ExpectError: regexp.MustCompile("validation error: theme \".+\" does not exist on the server"),
			},
			{
				Config: testKeycloakOpenidClient_loginTheme(clientId, loginThemeKeycloak),
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

func testAccCheckKeycloakOpenidClientExistsWithCorrectClientTimeouts(resourceName string,
	offlineSessionIdleTimeout string, offlineSessionMaxLifespan string,
	sessionIdleTimeout string, sessionMaxLifespan string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getOpenidClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		if client.Attributes.ClientOfflineSessionIdleTimeout != offlineSessionIdleTimeout {
			return fmt.Errorf("expected openid client to have client offline session idle timeout set to %s, but got %s", offlineSessionIdleTimeout, client.Attributes.ClientOfflineSessionIdleTimeout)
		}

		if client.Attributes.ClientOfflineSessionMaxLifespan != offlineSessionMaxLifespan {
			return fmt.Errorf("expected openid client to have client offline session max lifespan set to %s, but got %s", offlineSessionMaxLifespan, client.Attributes.ClientOfflineSessionMaxLifespan)
		}

		if client.Attributes.ClientSessionIdleTimeout != sessionIdleTimeout {
			return fmt.Errorf("expected openid client to have client session idle timeout set to %s, but got %s", sessionIdleTimeout, client.Attributes.ClientSessionIdleTimeout)
		}

		if client.Attributes.ClientSessionMaxLifespan != sessionMaxLifespan {
			return fmt.Errorf("expected openid client to have client session max lifespan set to %s, but got %s", sessionMaxLifespan, client.Attributes.ClientSessionMaxLifespan)
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

func testKeycloakOpenidClient_basic(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "CONFIDENTIAL"
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakOpenidClient_AccessToken_basic(clientId, accessTokenLifespan string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   		  = "%s"
	realm_id    		  = data.keycloak_realm.realm.id
	access_type 		  = "CONFIDENTIAL"
	access_token_lifespan = "%s"
}
	`, testAccRealm.Realm, clientId, accessTokenLifespan)
}

func testKeycloakOpenidClient_ClientTimeouts(clientId,
	offlineSessionIdleTimeout string, offlineSessionMaxLifespan string,
	sessionIdleTimeout string, sessionMaxLifespan string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   		  = "%s"
	realm_id    		  = data.keycloak_realm.realm.id
	access_type 		  = "CONFIDENTIAL"

	client_offline_session_idle_timeout = "%s"
	client_offline_session_max_lifespan = "%s"
	client_session_idle_timeout         = "%s"
	client_session_max_lifespan         = "%s"
}
	`, testAccRealm.Realm, clientId, offlineSessionIdleTimeout, offlineSessionMaxLifespan, sessionIdleTimeout, sessionMaxLifespan)
}

func testKeycloakOpenidClient_accessType(clientId, accessType string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "%s"
}
	`, testAccRealm.Realm, clientId, accessType)
}

func testKeycloakOpenidClient_adminUrl(clientId, adminUrl string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	admin_url = "%s"
	access_type = "PUBLIC"
}
	`, testAccRealm.Realm, clientId, adminUrl)
}

func testKeycloakOpenidClient_baseUrl(clientId, baseUrl string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	base_url = "%s"
	access_type = "PUBLIC"
}
	`, testAccRealm.Realm, clientId, baseUrl)
}

func testKeycloakOpenidClient_rootUrl(clientId, rootUrl string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id             = "%s"
	realm_id			  = data.keycloak_realm.realm.id
	root_url			  = "%s"
	valid_redirect_uris   = ["http://example.com"]
	web_origins			  = ["http://example.com"]
	admin_url			  = "http://example.com"
	access_type           = "CONFIDENTIAL"
	standard_flow_enabled = true
}
	`, testAccRealm.Realm, clientId, rootUrl)
}

func testKeycloakOpenidClient_pkceChallengeMethod(clientId, pkceChallengeMethod string) string {

	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "CONFIDENTIAL"
	pkce_code_challenge_method = "%s"
}
	`, testAccRealm.Realm, clientId, pkceChallengeMethod)
}

func testKeycloakOpenidClient_excludeSessionStateFromAuthResponse(clientId string, excludeSessionStateFromAuthResponse bool) string {

	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "CONFIDENTIAL"
	exclude_session_state_from_auth_response = %t
}
	`, testAccRealm.Realm, clientId, excludeSessionStateFromAuthResponse)
}

func testKeycloakOpenidClient_omitPkceChallengeMethod(clientId string) string {

	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "CONFIDENTIAL"
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakOpenidClient_omitExcludeSessionStateFromAuthResponse(clientId, pkceChallengeMethod string) string {

	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "CONFIDENTIAL"
    pkce_code_challenge_method = "%s"
}
	`, testAccRealm.Realm, clientId, pkceChallengeMethod)
}

func testKeycloakOpenidClient_updateRealmBefore(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm_1" {
	realm = "%s"
}

data "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm_1.id
	access_type = "BEARER-ONLY"
}
	`, testAccRealm.Realm, testAccRealmTwo.Realm, clientId)
}

func testKeycloakOpenidClient_updateRealmAfter(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm_1" {
	realm = "%s"
}

data "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm_2.id
	access_type = "BEARER-ONLY"
}
	`, testAccRealm.Realm, testAccRealmTwo.Realm, clientId)
}

func testKeycloakOpenidClient_fromInterface(openidClient *keycloak.OpenidClient) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id                    = "%s"
	realm_id                     = data.keycloak_realm.realm.id
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
	`, testAccRealm.Realm, openidClient.ClientId, openidClient.Name, openidClient.Enabled, openidClient.Description, openidClient.ClientSecret, openidClient.StandardFlowEnabled, openidClient.ImplicitFlowEnabled, openidClient.DirectAccessGrantsEnabled, openidClient.ServiceAccountsEnabled, arrayOfStringsForTerraformResource(openidClient.ValidRedirectUris), arrayOfStringsForTerraformResource(openidClient.WebOrigins), openidClient.AdminUrl, openidClient.BaseUrl, *openidClient.RootUrl)
}

func testKeycloakOpenidClient_secret(clientId, clientSecret string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id     = "%s"
	realm_id      = data.keycloak_realm.realm.id
	access_type   = "CONFIDENTIAL"
	client_secret = "%s"
}
	`, testAccRealm.Realm, clientId, clientSecret)
}

func testKeycloakOpenidClient_invalidRedirectUris(clientId, accessType string, standardFlowEnabled, implicitFlowEnabled bool) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id             = "%s"
	realm_id              = data.keycloak_realm.realm.id
	access_type           = "%s"

	standard_flow_enabled = %t
	implicit_flow_enabled = %t
}
	`, testAccRealm.Realm, clientId, accessType, standardFlowEnabled, implicitFlowEnabled)
}

func testKeycloakOpenidClient_invalidPublicClientWithClientCredentials(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id                = "%s"
	realm_id                 = data.keycloak_realm.realm.id
	access_type              = "PUBLIC"

	service_accounts_enabled = true
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakOpenidClient_bearerOnlyClientsCannotIssueTokens(clientId string, standardFlowEnabled, implicitFlowEnabled, directAccessGrantsEnabled, serviceAccountsEnabled bool) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id                    = "%s"
	realm_id                     = data.keycloak_realm.realm.id
	access_type                  = "BEARER-ONLY"

	standard_flow_enabled        = %t
	implicit_flow_enabled        = %t
	direct_access_grants_enabled = %t
	service_accounts_enabled     = %t
}
	`, testAccRealm.Realm, clientId, standardFlowEnabled, implicitFlowEnabled, directAccessGrantsEnabled, serviceAccountsEnabled)
}

func testKeycloakOpenidClient_authenticationFlowBindingOverrides(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "another_flow" {
  alias    = "anotherFlow"
  realm_id = data.keycloak_realm.realm.id
  description = "this is another flow"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "PUBLIC"
	authentication_flow_binding_overrides {
		browser_id = "${keycloak_authentication_flow.another_flow.id}"
		direct_grant_id = "${keycloak_authentication_flow.another_flow.id}"
	}
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakOpenidClient_withoutAuthenticationFlowBindingOverrides(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "another_flow" {
  alias    = "anotherFlow"
  realm_id = data.keycloak_realm.realm.id
  description = "this is another flow"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "PUBLIC"
}
	`, testAccRealm.Realm, clientId)
}

func testKeycloakOpenidClient_loginTheme(clientId, loginTheme string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = data.keycloak_realm.realm.id
	access_type = "PUBLIC"
	login_theme = "%s"
}
	`, testAccRealm.Realm, clientId, loginTheme)
}
