package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakRealm() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmCreate,
		Read:   resourceKeycloakRealmRead,
		Delete: resourceKeycloakRealmDelete,
		Update: resourceKeycloakRealmUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Login Config

			"registration_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"registration_email_as_username": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"edit_username_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"reset_password_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"remember_me": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"verify_email": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"login_with_email_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"duplicate_emails_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Themes

			"login_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Tokens

			"refresh_token_max_reuse": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"sso_session_idle_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "30m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"sso_session_max_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "10h",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"offline_session_idle_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "10h",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"offline_session_max_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "0",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_token_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "5m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_token_lifespan_for_implicit_flow": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "15m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_code_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "1m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_code_lifespan_login": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "30m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_code_lifespan_user_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "5m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"action_token_generated_by_user_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "5m",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"action_token_generated_by_admin_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "12h",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
		},
	}
}

func getRealmFromData(data *schema.ResourceData) (*keycloak.Realm, error) {
	realm := &keycloak.Realm{
		Id:          data.Get("realm").(string),
		Realm:       data.Get("realm").(string),
		Enabled:     data.Get("enabled").(bool),
		DisplayName: data.Get("display_name").(string),

		// Login Config
		RegistrationAllowed:         data.Get("registration_allowed").(bool),
		RegistrationEmailAsUsername: data.Get("registration_email_as_username").(bool),
		EditUsernameAllowed:         data.Get("edit_username_allowed").(bool),
		ResetPasswordAllowed:        data.Get("reset_password_allowed").(bool),
		RememberMe:                  data.Get("remember_me").(bool),
		VerifyEmail:                 data.Get("verify_email").(bool),
		LoginWithEmailAllowed:       data.Get("login_with_email_allowed").(bool),
		DuplicateEmailsAllowed:      data.Get("duplicate_emails_allowed").(bool),
	}

	// Themes

	if loginTheme, ok := data.GetOk("login_theme"); ok {
		realm.LoginTheme = loginTheme.(string)
	}

	if accountTheme, ok := data.GetOk("account_theme"); ok {
		realm.AccountTheme = accountTheme.(string)
	}

	if adminTheme, ok := data.GetOk("admin_theme"); ok {
		realm.AdminTheme = adminTheme.(string)
	}

	if emailTheme, ok := data.GetOk("email_theme"); ok {
		realm.EmailTheme = emailTheme.(string)
	}

	// Tokens

	if refreshTokenMaxReuse := data.Get("refresh_token_max_reuse").(int); refreshTokenMaxReuse > 0 {
		realm.RevokeRefreshToken = true
		realm.RefreshTokenMaxReuse = refreshTokenMaxReuse
	} else {
		realm.RevokeRefreshToken = false
	}

	ssoSessionIdleTimeout, err := getSecondsFromDurationString(data.Get("sso_session_idle_timeout").(string))
	if err != nil {
		return nil, err
	}
	realm.SsoSessionIdleTimeout = ssoSessionIdleTimeout

	ssoSessionMaxLifespan, err := getSecondsFromDurationString(data.Get("sso_session_max_lifespan").(string))
	if err != nil {
		return nil, err
	}
	realm.SsoSessionMaxLifespan = ssoSessionMaxLifespan

	offlineSessionIdleTimeout, err := getSecondsFromDurationString(data.Get("offline_session_idle_timeout").(string))
	if err != nil {
		return nil, err
	}
	realm.OfflineSessionIdleTimeout = offlineSessionIdleTimeout

	offlineSessionMaxLifespan, err := getSecondsFromDurationString(data.Get("offline_session_max_lifespan").(string))
	if err != nil {
		return nil, err
	}
	realm.OfflineSessionMaxLifespan = offlineSessionMaxLifespan

	accessTokenLifespan, err := getSecondsFromDurationString(data.Get("access_token_lifespan").(string))
	if err != nil {
		return nil, err
	}
	realm.AccessTokenLifespan = accessTokenLifespan

	accessTokenLifespanForImplicitFlow, err := getSecondsFromDurationString(data.Get("access_token_lifespan_for_implicit_flow").(string))
	if err != nil {
		return nil, err
	}
	realm.AccessTokenLifespanForImplicitFlow = accessTokenLifespanForImplicitFlow

	accessCodeLifespan, err := getSecondsFromDurationString(data.Get("access_code_lifespan").(string))
	if err != nil {
		return nil, err
	}
	realm.AccessCodeLifespan = accessCodeLifespan

	accessCodeLifespanLogin, err := getSecondsFromDurationString(data.Get("access_code_lifespan_login").(string))
	if err != nil {
		return nil, err
	}
	realm.AccessCodeLifespanLogin = accessCodeLifespanLogin

	accessCodeLifespanUserAction, err := getSecondsFromDurationString(data.Get("access_code_lifespan_user_action").(string))
	if err != nil {
		return nil, err
	}
	realm.AccessCodeLifespanUserAction = accessCodeLifespanUserAction

	actionTokenGeneratedByUserLifespan, err := getSecondsFromDurationString(data.Get("action_token_generated_by_user_lifespan").(string))
	if err != nil {
		return nil, err
	}
	realm.ActionTokenGeneratedByUserLifespan = actionTokenGeneratedByUserLifespan

	actionTokenGeneratedByAdminLifespan, err := getSecondsFromDurationString(data.Get("action_token_generated_by_admin_lifespan").(string))
	if err != nil {
		return nil, err
	}
	realm.ActionTokenGeneratedByAdminLifespan = actionTokenGeneratedByAdminLifespan

	return realm, nil
}

func setRealmData(data *schema.ResourceData, realm *keycloak.Realm) {
	data.SetId(realm.Realm)

	data.Set("realm", realm.Realm)
	data.Set("enabled", realm.Enabled)
	data.Set("display_name", realm.DisplayName)

	// Login Config
	data.Set("registration_allowed", realm.RegistrationAllowed)
	data.Set("registration_email_as_username", realm.RegistrationEmailAsUsername)
	data.Set("edit_username_allowed", realm.EditUsernameAllowed)
	data.Set("reset_password_allowed", realm.ResetPasswordAllowed)
	data.Set("remember_me", realm.RememberMe)
	data.Set("verify_email", realm.VerifyEmail)
	data.Set("login_with_email_allowed", realm.LoginWithEmailAllowed)
	data.Set("duplicate_emails_allowed", realm.DuplicateEmailsAllowed)

	// Themes
	data.Set("login_theme", realm.LoginTheme)
	data.Set("account_theme", realm.AccountTheme)
	data.Set("admin_theme", realm.AdminTheme)
	data.Set("email_theme", realm.EmailTheme)

	// Tokens

	data.Set("refresh_token_max_reuse", realm.RefreshTokenMaxReuse)
	data.Set("sso_session_idle_timeout", getDurationStringFromSeconds(realm.SsoSessionIdleTimeout))
	data.Set("sso_session_max_lifespan", getDurationStringFromSeconds(realm.SsoSessionMaxLifespan))
	data.Set("offline_session_idle_timeout", getDurationStringFromSeconds(realm.OfflineSessionIdleTimeout))
	data.Set("offline_session_max_lifespan", getDurationStringFromSeconds(realm.OfflineSessionMaxLifespan))
	data.Set("access_token_lifespan", getDurationStringFromSeconds(realm.AccessTokenLifespan))
	data.Set("access_token_lifespan_for_implicit_flow", getDurationStringFromSeconds(realm.AccessTokenLifespanForImplicitFlow))
	data.Set("access_code_lifespan", getDurationStringFromSeconds(realm.AccessCodeLifespan))
	data.Set("access_code_lifespan_login", getDurationStringFromSeconds(realm.AccessCodeLifespanLogin))
	data.Set("access_code_lifespan_user_action", getDurationStringFromSeconds(realm.AccessCodeLifespanUserAction))
	data.Set("action_token_generated_by_user_lifespan", getDurationStringFromSeconds(realm.ActionTokenGeneratedByUserLifespan))
	data.Set("action_token_generated_by_admin_lifespan", getDurationStringFromSeconds(realm.ActionTokenGeneratedByAdminLifespan))
}

func resourceKeycloakRealmCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := getRealmFromData(data)
	if err != nil {
		return err
	}

	err = realm.Validate(keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealm(realm)
	if err != nil {
		return err
	}

	setRealmData(data, realm)

	return resourceKeycloakRealmRead(data, meta)
}

func resourceKeycloakRealmRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := keycloakClient.GetRealm(data.Id())
	if err != nil {
		return err
	}

	setRealmData(data, realm)

	return nil
}

func resourceKeycloakRealmUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := getRealmFromData(data)
	if err != nil {
		return err
	}

	err = realm.Validate(keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealm(realm)
	if err != nil {
		return err
	}

	setRealmData(data, realm)

	return nil
}

func resourceKeycloakRealmDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	return keycloakClient.DeleteRealm(data.Id())
}
