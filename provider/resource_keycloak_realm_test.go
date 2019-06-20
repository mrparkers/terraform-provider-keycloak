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

func TestAccKeycloakRealm_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmDisplayName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				Config: testKeycloakRealm_notEnabled(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmEnabled("keycloak_realm.realm", false),
			},
			{
				Config: testKeycloakRealm_basic(realmName, fmt.Sprintf("%s-changed", realmDisplayName)),
				Check:  testAccCheckKeycloakRealmDisplayName("keycloak_realm.realm", fmt.Sprintf("%s-changed", realmDisplayName)),
			},
		},
	})
}

func TestAccKeycloakRealm_createAfterManualDestroy(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmDisplayName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteRealm(realmName)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakRealm_basic(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
		},
	})
}

func TestAccKeycloakRealm_import(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmDisplayName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName, realmDisplayName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				ResourceName:      "keycloak_realm.realm",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeycloakRealm_SmtpServer(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_WithSmtpServer(realm, "myhost.com", "My Host", "user"),
				Check:  testAccCheckKeycloakRealmSmtp("keycloak_realm.realm", "myhost.com", "My Host", "user"),
			},
			{
				Config: testKeycloakRealm_basic(realm, realm),
				Check:  testAccCheckKeycloakRealmSmtp("keycloak_realm.realm", "", "", ""),
			},
		},
	})
}

func TestAccKeycloakRealm_SmtpServerUpdate(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_WithSmtpServer(realm, "myhost.com", "My Host", "user"),
				Check:  testAccCheckKeycloakRealmSmtp("keycloak_realm.realm", "myhost.com", "My Host", "user"),
			},
			{
				Config: testKeycloakRealm_WithSmtpServer(realm, "myhost2.com", "My Host2", "user2"),
				Check:  testAccCheckKeycloakRealmSmtp("keycloak_realm.realm", "myhost2.com", "My Host2", "user2"),
			},
		},
	})
}

func TestAccKeycloakRealm_SmtpServerInvalid(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRealm_WithSmtpServerWithoutHost(realm, "My Host"),
				ExpectError: regexp.MustCompile("config is invalid: Missing required argument: The argument \"host\" is required, but no definition was found."),
			},
			{
				Config:      testKeycloakRealm_WithSmtpServerWithoutFrom(realm, "myhost.com"),
				ExpectError: regexp.MustCompile("config is invalid: Missing required argument: The argument \"from\" is required, but no definition was found."),
			},
		},
	})
}

func TestAccKeycloakRealm_themes(t *testing.T) {
	realmOne := &keycloak.Realm{
		Realm:        "terraform-" + acctest.RandString(10),
		DisplayName:  "terraform-" + acctest.RandString(10),
		LoginTheme:   randomStringInSlice([]string{"base", "keycloak"}),
		AccountTheme: randomStringInSlice([]string{"base", "keycloak"}),
		AdminTheme:   randomStringInSlice([]string{"base", "keycloak"}),
		EmailTheme:   randomStringInSlice([]string{"base", "keycloak"}),
	}

	realmTwo := &keycloak.Realm{
		Realm:        realmOne.Realm,
		DisplayName:  realmOne.DisplayName,
		LoginTheme:   randomStringInSlice([]string{"base", "keycloak"}),
		AccountTheme: randomStringInSlice([]string{"base", "keycloak"}),
		AdminTheme:   randomStringInSlice([]string{"base", "keycloak"}),
		EmailTheme:   randomStringInSlice([]string{"base", "keycloak"}),
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_themes(realmOne),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				Config: testKeycloakRealm_themes(realmTwo),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
		},
	})
}

func TestAccKeycloakRealm_themesValidation(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRealm_themesValidation(realm, "login", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("validation error: theme \".+\" does not exist on the server"),
			},
			{
				Config:      testKeycloakRealm_themesValidation(realm, "account", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("validation error: theme \".+\" does not exist on the server"),
			},
			{
				Config:      testKeycloakRealm_themesValidation(realm, "admin", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("validation error: theme \".+\" does not exist on the server"),
			},
			{
				Config:      testKeycloakRealm_themesValidation(realm, "email", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("validation error: theme \".+\" does not exist on the server"),
			},
		},
	})
}

func TestAccKeycloakRealm_InternationalizationValidation(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRealm_internationalizationValidationWithoutSupportedLocales(realm, "en"),
				ExpectError: regexp.MustCompile("config is invalid: Missing required argument: The argument \"supported_locales\" is required, but no definition was found."),
			},
			{
				Config:      testKeycloakRealm_internationalizationValidation(realm, "en", "de"),
				ExpectError: regexp.MustCompile("validation error: DefaultLocale should be in the SupportLocales"),
			},
		},
	})
}

func TestAccKeycloakRealm_Internationalization(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_internationalizationValidation(realm, "en", "en"),
				Check:  testAccCheckKeycloakRealmInternationalizationIsEnabled("keycloak_realm.realm", "en"),
			},
			{
				Config: testKeycloakRealm_internationalizationValidation(realm, "es", "es"),
				Check:  testAccCheckKeycloakRealmInternationalizationIsEnabled("keycloak_realm.realm", "es"),
			},
			{
				Config: testKeycloakRealm_basic(realm, realm),
				Check:  testAccCheckKeycloakRealmInternationalizationIsDisabled("keycloak_realm.realm"),
			},
		},
	})
}

func TestAccKeycloakRealm_InternationalizationDisabled(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realm, realm),
				Check:  testAccCheckKeycloakRealmInternationalizationIsDisabled("keycloak_realm.realm"),
			},
		},
	})
}

func TestAccKeycloakRealm_loginConfigBasic(t *testing.T) {
	realm := &keycloak.Realm{
		Realm:                       "terraform-" + acctest.RandString(10),
		RegistrationAllowed:         true,
		RegistrationEmailAsUsername: true,
		EditUsernameAllowed:         randomBool(),
		ResetPasswordAllowed:        randomBool(),
		RememberMe:                  randomBool(),
		VerifyEmail:                 randomBool(),
		LoginWithEmailAllowed:       randomBool(),
		DuplicateEmailsAllowed:      false,
	}

	updatedRealm := &keycloak.Realm{
		Realm:                       realm.Realm,
		RegistrationAllowed:         true,
		RegistrationEmailAsUsername: false,
		EditUsernameAllowed:         randomBool(),
		ResetPasswordAllowed:        randomBool(),
		RememberMe:                  randomBool(),
		VerifyEmail:                 randomBool(),
		LoginWithEmailAllowed:       false,
		DuplicateEmailsAllowed:      true,
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_loginConfigBasic(realm),
				Check:  testKeycloakRealmLoginInfo("keycloak_realm.realm", realm),
			},
			{
				Config: testKeycloakRealm_loginConfigBasic(updatedRealm),
				Check:  testKeycloakRealmLoginInfo("keycloak_realm.realm", updatedRealm),
			},
		},
	})
}

func TestAccKeycloakRealm_loginConfigValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRealm_invalidRegistrationEmailAsUsernameWithoutRegistrationAllowed(realmName),
				ExpectError: regexp.MustCompile("validation error: RegistrationEmailAsUsername cannot be true if RegistrationAllowed is false"),
			},
			{
				Config:      testKeycloakRealm_invalidRegistrationEmailAsUsernameAndDuplicateEmailsAllowed(realmName),
				ExpectError: regexp.MustCompile("validation error: DuplicateEmailsAllowed cannot be true if RegistrationEmailAsUsername is true"),
			},
			{
				Config:      testKeycloakRealm_invalidLoginWithEmailAllowedAndDuplicateEmailsAllowed(realmName),
				ExpectError: regexp.MustCompile("validation error: DuplicateEmailsAllowed cannot be true if LoginWithEmailAllowed is true"),
			},
		},
	})
}

func TestAccKeycloakRealm_tokenSettings(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName, realmName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				Config: testKeycloakRealm_tokenSettings(realmName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			// This is duplicated so another set of random value is used, effectively an update test
			{
				Config: testKeycloakRealm_tokenSettings(realmName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
		},
	})
}

func TestAccKeycloakRealm_computedTokenSettings(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmDisplayName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName, realmDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRealmExists("keycloak_realm.realm"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "sso_session_idle_timeout"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "sso_session_idle_timeout", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "sso_session_max_lifespan"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "sso_session_max_lifespan", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "offline_session_idle_timeout"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "offline_session_idle_timeout", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "offline_session_max_lifespan"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "offline_session_max_lifespan", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "access_token_lifespan"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "access_token_lifespan", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "access_token_lifespan_for_implicit_flow"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "access_token_lifespan_for_implicit_flow", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "access_code_lifespan"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "access_code_lifespan", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "access_code_lifespan_login"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "access_code_lifespan_login", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "access_code_lifespan_user_action"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "access_code_lifespan_user_action", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "action_token_generated_by_user_lifespan"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "action_token_generated_by_user_lifespan", "0s"),

					resource.TestCheckResourceAttrSet("keycloak_realm.realm", "action_token_generated_by_admin_lifespan"),
					TestCheckResourceAttrNot("keycloak_realm.realm", "action_token_generated_by_admin_lifespan", "0s"),
				),
			},
		},
	})
}

func testKeycloakRealmLoginInfo(resourceName string, realm *keycloak.Realm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realmFromState, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realmFromState.Realm != realm.Realm {
			return fmt.Errorf("expected realm in state to have name %s, but was %s", realmFromState.Realm, realm.Realm)
		}

		if realmFromState.RegistrationAllowed != realm.RegistrationAllowed {
			return fmt.Errorf("expected realm %s to have registration_allowed set to %t, but was %t", realm.Realm, realm.RegistrationAllowed, realmFromState.RegistrationAllowed)
		}

		if realmFromState.RegistrationEmailAsUsername != realm.RegistrationEmailAsUsername {
			return fmt.Errorf("expected realm %s to have registration_email_as_username set to %t, but was %t", realm.Realm, realm.RegistrationEmailAsUsername, realmFromState.RegistrationEmailAsUsername)
		}

		if realmFromState.EditUsernameAllowed != realm.EditUsernameAllowed {
			return fmt.Errorf("expected realm %s to have edit_username_allowed set to %t, but was %t", realm.Realm, realm.EditUsernameAllowed, realmFromState.EditUsernameAllowed)
		}

		if realmFromState.ResetPasswordAllowed != realm.ResetPasswordAllowed {
			return fmt.Errorf("expected realm %s to have reset_password_allowed set to %t, but was %t", realm.Realm, realm.ResetPasswordAllowed, realmFromState.ResetPasswordAllowed)
		}

		if realmFromState.RememberMe != realm.RememberMe {
			return fmt.Errorf("expected realm %s to have remember_me set to %t, but was %t", realm.Realm, realm.RememberMe, realmFromState.RememberMe)
		}

		if realmFromState.VerifyEmail != realm.VerifyEmail {
			return fmt.Errorf("expected realm %s to have verify_email set to %t, but was %t", realm.Realm, realm.VerifyEmail, realmFromState.VerifyEmail)
		}

		if realmFromState.LoginWithEmailAllowed != realm.LoginWithEmailAllowed {
			return fmt.Errorf("expected realm %s to have login_with_email_allowed set to %t, but was %t", realm.Realm, realm.LoginWithEmailAllowed, realmFromState.LoginWithEmailAllowed)
		}

		if realmFromState.DuplicateEmailsAllowed != realm.DuplicateEmailsAllowed {
			return fmt.Errorf("expected realm %s to have duplicate_emails_allowed set to %t, but was %t", realm.Realm, realm.DuplicateEmailsAllowed, realmFromState.DuplicateEmailsAllowed)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakRealmEnabled(resourceName string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realm.Enabled != enabled {
			return fmt.Errorf("expected realm %s to have enabled set to %t, but was %t", realm.Realm, enabled, realm.Enabled)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmDisplayName(resourceName string, displayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realm.DisplayName != displayName {
			return fmt.Errorf("expected realm %s to have display name set to %s, but was %s", realm.Realm, displayName, realm.DisplayName)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmSmtp(resourceName, host, from, user string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realm.SmtpServer.Host != host {
			return fmt.Errorf("expected realm %s to have smtp host set to %s, but was %s", realm.Realm, host, realm.SmtpServer.Host)
		}

		if realm.SmtpServer.From != from {
			return fmt.Errorf("expected realm %s to have smtp from set to %s, but was %s", realm.Realm, from, realm.SmtpServer.From)
		}

		if realm.SmtpServer.User != user {
			return fmt.Errorf("expected realm %s to have smtp user set to %s, but was %s", realm.Realm, user, realm.SmtpServer.User)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmInternationalizationIsEnabled(resourceName string, defaultLocale string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if !realm.InternationalizationEnabled {
			return fmt.Errorf("expected realm %s to have internationalization enabled but was disabled", realm.Realm)
		}

		if realm.DefaultLocale != defaultLocale {
			return fmt.Errorf("expected realm %s to have defaultLocale set to %s, but was %s", realm.Realm, defaultLocale, realm.DefaultLocale)
		}

		if !contains(realm.SupportLocales, defaultLocale) {
			return fmt.Errorf("expected realm %s to contain defaultLocale %s, but was %s", realm.Realm, defaultLocale, realm.SupportLocales)
		}
		return nil
	}
}

func testAccCheckKeycloakRealmInternationalizationIsDisabled(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realm, err := getRealmFromState(s, resourceName)
		if err != nil {
			return err
		}

		if realm.InternationalizationEnabled {
			return fmt.Errorf("expected realm %s to have internationalization disabled but was enabled", realm.Realm)
		}
		return nil
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func testAccCheckKeycloakRealmDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm" {
				continue
			}

			realmName := rs.Primary.ID
			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			realm, _ := keycloakClient.GetRealm(realmName)
			if realm != nil {
				return fmt.Errorf("realm %s still exists", realmName)
			}
		}

		return nil
	}
}

func getRealmFromState(s *terraform.State, resourceName string) (*keycloak.Realm, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmName := rs.Primary.Attributes["realm"]

	realm, err := keycloakClient.GetRealm(realmName)
	if err != nil {
		return nil, fmt.Errorf("error getting realm %s: %s", realmName, err)
	}

	return realm, nil
}

func testKeycloakRealm_basic(realm, realmDisplayName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = true
	display_name = "%s"
}
	`, realm, realmDisplayName)
}

func testKeycloakRealm_WithSmtpServer(realm, host, from, user string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
	enabled = true
	display_name = "%s"
	smtp_server {
		host = "%s"
		port = 25
		from_display_name = "Tom"
		from = "%s"
		reply_to_display_name = "Tom"
		reply_to = "tom@myhost.com"
		ssl = true
		starttls = true
		envelope_from = "nottom@myhost.com"
		auth {
			username = "%s"
			password = "tom"
		}
	}
}
	`, realm, realm, host, from, user)
}

func testKeycloakRealm_WithSmtpServerWithoutHost(realm, from string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
	enabled = true
	display_name = "%s"
	smtp_server {
		port = 25
		from_display_name = "Tom"
		from = "%s"
		reply_to_display_name = "Tom"
		reply_to = "tom@myhost.com"
		ssl = true
		starttls = true
		envelope_from = "nottom@myhost.com"
		auth {
			username = "tom"
			password = "tom"
		}
	}
}
	`, realm, realm, from)
}

func testKeycloakRealm_WithSmtpServerWithoutFrom(realm, host string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
	enabled = true
	display_name = "%s"
	smtp_server {
		host = "%s"
		port = 25
		from_display_name = "Tom"
		reply_to_display_name = "Tom"
		reply_to = "tom@myhost.com"
		ssl = true
		starttls = true
		envelope_from = "nottom@myhost.com"
		auth {
			username = "tom"
			password = "tom"
		}
	}
}
	`, realm, realm, host)
}

func testKeycloakRealm_themes(realm *keycloak.Realm) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = true
	display_name = "%s"

	login_theme   = "%s"
	account_theme = "%s"
	admin_theme   = "%s"
	email_theme   = "%s"
}
	`, realm.Realm, realm.DisplayName, realm.LoginTheme, realm.AccountTheme, realm.AdminTheme, realm.EmailTheme)
}

func testKeycloakRealm_themesValidation(realm, theme, value string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = true
	display_name = "%s"

	%s_theme     = "%s"
}
	`, realm, realm, theme, value)
}

func testKeycloakRealm_internationalizationValidation(realm, supportedLocale, defaultLocale string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = true
	display_name = "%s"
	internationalization {
		supported_locales	= ["nl", "%s", "fr"]
		default_locale		= "%s"
	}
}
	`, realm, realm, supportedLocale, defaultLocale)
}

func testKeycloakRealm_internationalizationValidationWithoutSupportedLocales(realm, defaultLocale string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = true
	display_name = "%s"
	internationalization {
		default_locale		= "%s"
	}

}
	`, realm, realm, defaultLocale)
}

func testKeycloakRealm_notEnabled(realm, realmDisplayName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"
	enabled      = false
	display_name = "%s"
}
	`, realm, realmDisplayName)
}

func testKeycloakRealm_loginConfigBasic(realm *keycloak.Realm) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm        = "%s"

	registration_allowed           = "%t"
	registration_email_as_username = "%t"
	edit_username_allowed          = "%t"
	reset_password_allowed         = "%t"
	remember_me                    = "%t"
	verify_email                   = "%t"
	login_with_email_allowed       = "%t"
	duplicate_emails_allowed       = "%t"
}
	`, realm.Realm, realm.RegistrationAllowed, realm.RegistrationEmailAsUsername, realm.EditUsernameAllowed, realm.ResetPasswordAllowed, realm.RememberMe, realm.VerifyEmail, realm.LoginWithEmailAllowed, realm.DuplicateEmailsAllowed)
}

func testKeycloakRealm_invalidRegistrationEmailAsUsernameWithoutRegistrationAllowed(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm                          = "%s"

	registration_allowed           = false
	registration_email_as_username = true
}
	`, realm)
}

func testKeycloakRealm_invalidRegistrationEmailAsUsernameAndDuplicateEmailsAllowed(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm                          = "%s"

	registration_allowed           = true
	registration_email_as_username = true
	duplicate_emails_allowed       = true
}
	`, realm)
}

func testKeycloakRealm_invalidLoginWithEmailAllowedAndDuplicateEmailsAllowed(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm                          = "%s"

	login_with_email_allowed       = true
	duplicate_emails_allowed       = true
}
	`, realm)
}

func testKeycloakRealm_tokenSettings(realm string) string {
	ssoSessionIdleTimeout := randomDurationString()
	ssoSessionMaxLifespan := randomDurationString()
	offlineSessionIdleTimeout := randomDurationString()
	offlineSessionMaxLifespan := randomDurationString()
	accessTokenLifespan := randomDurationString()
	accessTokenLifespanForImplicitFlow := randomDurationString()
	accessCodeLifespan := randomDurationString()
	accessCodeLifespanLogin := randomDurationString()
	accessCodeLifespanUserAction := randomDurationString()
	actionTokenGeneratedByUserLifespan := randomDurationString()
	actionTokenGeneratedByAdminLifespan := randomDurationString()

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm                                    = "%s"
	enabled                                  = true
	display_name                             = "%s"

	sso_session_idle_timeout                 = "%s"
	sso_session_max_lifespan                 = "%s"
	offline_session_idle_timeout             = "%s"
	offline_session_max_lifespan             = "%s"
	access_token_lifespan                    = "%s"
	access_token_lifespan_for_implicit_flow  = "%s"
	access_code_lifespan                     = "%s"
	access_code_lifespan_login               = "%s"
	access_code_lifespan_user_action         = "%s"
	action_token_generated_by_user_lifespan  = "%s"
	action_token_generated_by_admin_lifespan = "%s"
}
	`, realm, realm, ssoSessionIdleTimeout, ssoSessionMaxLifespan, offlineSessionIdleTimeout, offlineSessionMaxLifespan, accessTokenLifespan, accessTokenLifespanForImplicitFlow, accessCodeLifespan, accessCodeLifespanLogin, accessCodeLifespanUserAction, actionTokenGeneratedByUserLifespan, actionTokenGeneratedByAdminLifespan)
}
