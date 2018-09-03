package provider_test

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

func TestAccKeycloakRealm_loginConfigBasic(t *testing.T) {
	realm := &keycloak.Realm{
		Realm:                       "terraform-" + acctest.RandString(10),
		RegistrationAllowed:         true,
		RegistrationEmailAsUsername: true,
		EditUsernameAllowed:         RandomBool(),
		ResetPasswordAllowed:        RandomBool(),
		RememberMe:                  RandomBool(),
		VerifyEmail:                 RandomBool(),
		LoginWithEmailAllowed:       RandomBool(),
		DuplicateEmailsAllowed:      false,
	}

	updatedRealm := &keycloak.Realm{
		Realm:                       realm.Realm,
		RegistrationAllowed:         true,
		RegistrationEmailAsUsername: false,
		EditUsernameAllowed:         RandomBool(),
		ResetPasswordAllowed:        RandomBool(),
		RememberMe:                  RandomBool(),
		VerifyEmail:                 RandomBool(),
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
