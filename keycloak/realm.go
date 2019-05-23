package keycloak

import (
	"fmt"
)

type Realm struct {
	Id          string `json:"id"`
	Realm       string `json:"realm"`
	Enabled     bool   `json:"enabled"`
	DisplayName string `json:"displayName"`

	Hostname string `json:"hostname,omitempty"`

	// Login Config
	RegistrationAllowed         bool `json:"registrationAllowed"`
	RegistrationEmailAsUsername bool `json:"registrationEmailAsUsername"`
	EditUsernameAllowed         bool `json:"editUsernameAllowed"`
	ResetPasswordAllowed        bool `json:"resetPasswordAllowed"`
	RememberMe                  bool `json:"rememberMe"`
	VerifyEmail                 bool `json:"verifyEmail"`
	LoginWithEmailAllowed       bool `json:"loginWithEmailAllowed"`
	DuplicateEmailsAllowed      bool `json:"duplicateEmailsAllowed"`

	// Themes
	LoginTheme   string `json:"loginTheme,omitempty"`
	AccountTheme string `json:"accountTheme,omitempty"`
	AdminTheme   string `json:"adminTheme,omitempty"`
	EmailTheme   string `json:"emailTheme,omitempty"`

	// Tokens
	RevokeRefreshToken                  bool `json:"revokeRefreshToken,omitempty"`
	RefreshTokenMaxReuse                int  `json:"refreshTokenMaxReuse,omitempty"`
	SsoSessionIdleTimeout               int  `json:"ssoSessionIdleTimeout,omitempty"`
	SsoSessionMaxLifespan               int  `json:"ssoSessionMaxLifespan,omitempty"`
	OfflineSessionIdleTimeout           int  `json:"offlineSessionIdleTimeout,omitempty"`
	OfflineSessionMaxLifespan           int  `json:"offlineSessionMaxLifespan,omitempty"`
	AccessTokenLifespan                 int  `json:"accessTokenLifespan,omitempty"`
	AccessTokenLifespanForImplicitFlow  int  `json:"accessTokenLifespanForImplicitFlow,omitempty"`
	AccessCodeLifespan                  int  `json:"accessCodeLifespan,omitempty"`
	AccessCodeLifespanLogin             int  `json:"accessCodeLifespanLogin,omitempty"`
	AccessCodeLifespanUserAction        int  `json:"accessCodeLifespanUserAction,omitempty"`
	ActionTokenGeneratedByUserLifespan  int  `json:"actionTokenGeneratedByUserLifespan,omitempty"`
	ActionTokenGeneratedByAdminLifespan int  `json:"actionTokenGeneratedByAdminLifespan,omitempty"`
}

func (keycloakClient *KeycloakClient) NewRealm(realm *Realm) error {
	_, _, err := keycloakClient.post("/realms", realm)

	return err
}

func (keycloakClient *KeycloakClient) GetRealm(id string) (*Realm, error) {
	var realm Realm

	err := keycloakClient.get(fmt.Sprintf("/realms/%s", id), &realm, nil)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

func (keycloakClient *KeycloakClient) UpdateRealm(realm *Realm) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s", realm.Id), realm)
}

func (keycloakClient *KeycloakClient) DeleteRealm(id string) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s", id), nil)
	if err != nil {
		// For whatever reason, this fails sometimes with a 500 during acceptance tests. try again
		return keycloakClient.delete(fmt.Sprintf("/realms/%s", id), nil)
	}

	return nil
}

func (keycloakClient *KeycloakClient) ValidateRealm(realm *Realm) error {
	if realm.RegistrationAllowed == false && realm.RegistrationEmailAsUsername == true {
		return fmt.Errorf("validation error: RegistrationEmailAsUsername cannot be true if RegistrationAllowed is false")
	}

	if realm.DuplicateEmailsAllowed == true && realm.RegistrationEmailAsUsername == true {
		return fmt.Errorf("validation error: DuplicateEmailsAllowed cannot be true if RegistrationEmailAsUsername is true")
	}

	if realm.DuplicateEmailsAllowed == true && realm.LoginWithEmailAllowed == true {
		return fmt.Errorf("validation error: DuplicateEmailsAllowed cannot be true if LoginWithEmailAllowed is true")
	}

	// validate if the given theme exists on the server. the keycloak API allows you to use any random string for a theme
	serverInfo, err := keycloakClient.GetServerInfo()
	if err != nil {
		return err
	}

	if realm.LoginTheme != "" && !serverInfo.ThemeIsInstalled("login", realm.LoginTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.LoginTheme)
	}

	if realm.AccountTheme != "" && !serverInfo.ThemeIsInstalled("account", realm.AccountTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.AccountTheme)
	}

	if realm.AdminTheme != "" && !serverInfo.ThemeIsInstalled("admin", realm.AdminTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.AdminTheme)
	}

	if realm.EmailTheme != "" && !serverInfo.ThemeIsInstalled("email", realm.EmailTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.EmailTheme)
	}

	return nil
}
