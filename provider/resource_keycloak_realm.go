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
				Computed: true,
			},
			"registration_email_as_username": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"edit_username_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"reset_password_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"remember_me": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"verify_email": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"login_with_email_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"duplicate_emails_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			//Smtp server

			"smtp_server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starttls": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"reply_to": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reply_to_display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"from_display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"envelope_from": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auth": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
										DiffSuppressFunc: func(_, smtpServerPassword, _ string, _ *schema.ResourceData) bool {
											return smtpServerPassword == "**********"
										},
									},
								},
							},
						},
					},
				},
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
				Computed: true,
			},
			"sso_session_idle_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"sso_session_max_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"offline_session_idle_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"offline_session_max_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_token_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_token_lifespan_for_implicit_flow": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_code_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_code_lifespan_login": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"access_code_lifespan_user_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"action_token_generated_by_user_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"action_token_generated_by_admin_lifespan": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDurationStringDiff,
			},

			//internationalization
			"internationalization": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"supported_locales": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							Required: true,
						},
						"default_locale": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			//Security Defenses
			"security_defenses": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"x_frame_options": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "SAMEORIGIN",
									},
									"content_security_policy": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "frame-src 'self'; frame-ancestors 'self'; object-src 'none';",
									},
									"content_security_policy_report_only": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"x_content_type_options": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "nosniff",
									},
									"x_robots_tag": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "none",
									},
									"x_xss_protection": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "1; mode=block",
									},
									"strict_transport_security": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "max-age=31536000; includeSubDomains",
									},
								},
							},
						},
					},
				},
			},
			"password_policy": {
				Type:        schema.TypeString,
				Description: "String that represents the passwordPolicies that are in place. Each policy is separated with \" and \". Supported policies can be found in the server-info providers page. example: \"upperCase(1) and length(8) and forceExpiredPasswordChange(365) and notUsername(undefined)\"",
				Optional:    true,
			},

			//flow bindings
			"browser_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for BrowserFlow",
				Optional:    true,
				Default:     "browser",
			},
			"registration_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for RegistrationFlow",
				Optional:    true,
				Default:     "registration",
			},
			"direct_grant_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for DirectGrantFlow",
				Optional:    true,
				Default:     "direct grant",
			},
			"reset_credentials_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for ResetCredentialsFlow",
				Optional:    true,
				Default:     "registration",
			},
			"client_authentication_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for ClientAuthenticationFlow",
				Optional:    true,
				Default:     "clients",
			},
			"docker_authentication_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for DockerAuthenticationFlow",
				Optional:    true,
				Default:     "docker auth",
			},
		},
	}
}

func getRealmSMTPPasswordFromData(data *schema.ResourceData) (string, bool) {
	if v, ok := data.GetOk("smtp_server"); ok {
		smtpSettings := v.([]interface{})[0].(map[string]interface{})
		authConfig := smtpSettings["auth"].([]interface{})

		if len(authConfig) == 1 {
			return authConfig[0].(map[string]interface{})["password"].(string), true
		}

		return "", false
	}

	return "", false
}

func getRealmFromData(data *schema.ResourceData) (*keycloak.Realm, error) {
	internationalizationEnabled := false
	supportLocales := make([]string, 0)
	defaultLocale := ""
	if v, ok := data.GetOk("internationalization"); ok {
		internationalizationEnabled = true
		internationalizationSettings := v.([]interface{})[0].(map[string]interface{})
		if v, ok := internationalizationSettings["supported_locales"]; ok {
			for _, supportLocale := range v.(*schema.Set).List() {
				supportLocales = append(supportLocales, supportLocale.(string))
			}
		}
		defaultLocale = internationalizationSettings["default_locale"].(string)
	}

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

		//internationalization
		InternationalizationEnabled: internationalizationEnabled,
		SupportLocales:              supportLocales,
		DefaultLocale:               defaultLocale,
	}

	//smtp
	if v, ok := data.GetOk("smtp_server"); ok {
		smtpSettings := v.([]interface{})[0].(map[string]interface{})

		smtpServer := keycloak.SmtpServer{
			StartTls:           keycloak.KeycloakBoolQuoted(smtpSettings["starttls"].(bool)),
			Port:               smtpSettings["port"].(string),
			Host:               smtpSettings["host"].(string),
			ReplyTo:            smtpSettings["reply_to"].(string),
			ReplyToDisplayName: smtpSettings["reply_to_display_name"].(string),
			From:               smtpSettings["from"].(string),
			FromDisplayName:    smtpSettings["from_display_name"].(string),
			EnvelopeFrom:       smtpSettings["envelope_from"].(string),
			Ssl:                keycloak.KeycloakBoolQuoted(smtpSettings["ssl"].(bool)),
		}

		authConfig := smtpSettings["auth"].([]interface{})
		if len(authConfig) == 1 {
			auth := authConfig[0].(map[string]interface{})

			smtpServer.Auth = true
			smtpServer.User = auth["username"].(string)
			smtpServer.Password = auth["password"].(string)
		} else {
			smtpServer.Auth = false
		}

		realm.SmtpServer = smtpServer
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

	if ssoSessionIdleTimeout := data.Get("sso_session_idle_timeout").(string); ssoSessionIdleTimeout != "" {
		ssoSessionIdleTimeoutDurationString, err := getSecondsFromDurationString(ssoSessionIdleTimeout)
		if err != nil {
			return nil, err
		}
		realm.SsoSessionIdleTimeout = ssoSessionIdleTimeoutDurationString
	}

	if ssoSessionMaxLifespan := data.Get("sso_session_max_lifespan").(string); ssoSessionMaxLifespan != "" {
		ssoSessionMaxLifespanDurationString, err := getSecondsFromDurationString(ssoSessionMaxLifespan)
		if err != nil {
			return nil, err
		}
		realm.SsoSessionMaxLifespan = ssoSessionMaxLifespanDurationString
	}

	if offlineSessionIdleTimeout := data.Get("offline_session_idle_timeout").(string); offlineSessionIdleTimeout != "" {
		offlineSessionIdleTimeoutDurationString, err := getSecondsFromDurationString(offlineSessionIdleTimeout)
		if err != nil {
			return nil, err
		}
		realm.OfflineSessionIdleTimeout = offlineSessionIdleTimeoutDurationString
	}

	if offlineSessionMaxLifespan := data.Get("offline_session_max_lifespan").(string); offlineSessionMaxLifespan != "" {
		offlineSessionMaxLifespanDurationString, err := getSecondsFromDurationString(offlineSessionMaxLifespan)
		if err != nil {
			return nil, err
		}
		realm.OfflineSessionMaxLifespan = offlineSessionMaxLifespanDurationString
	}

	if accessTokenLifespan := data.Get("access_token_lifespan").(string); accessTokenLifespan != "" {
		accessTokenLifespanDurationString, err := getSecondsFromDurationString(accessTokenLifespan)
		if err != nil {
			return nil, err
		}
		realm.AccessTokenLifespan = accessTokenLifespanDurationString
	}

	if accessTokenLifespanForImplicitFlow := data.Get("access_token_lifespan_for_implicit_flow").(string); accessTokenLifespanForImplicitFlow != "" {
		accessTokenLifespanForImplicitFlowDurationString, err := getSecondsFromDurationString(accessTokenLifespanForImplicitFlow)
		if err != nil {
			return nil, err
		}
		realm.AccessTokenLifespanForImplicitFlow = accessTokenLifespanForImplicitFlowDurationString
	}

	if accessCodeLifespan := data.Get("access_code_lifespan").(string); accessCodeLifespan != "" {
		accessCodeLifespanDurationString, err := getSecondsFromDurationString(accessCodeLifespan)
		if err != nil {
			return nil, err
		}
		realm.AccessCodeLifespan = accessCodeLifespanDurationString
	}

	if accessCodeLifespanLogin := data.Get("access_code_lifespan_login").(string); accessCodeLifespanLogin != "" {
		accessCodeLifespanLoginDurationString, err := getSecondsFromDurationString(accessCodeLifespanLogin)
		if err != nil {
			return nil, err
		}
		realm.AccessCodeLifespanLogin = accessCodeLifespanLoginDurationString
	}

	if accessCodeLifespanUserAction := data.Get("access_code_lifespan_user_action").(string); accessCodeLifespanUserAction != "" {
		accessCodeLifespanUserActionDurationString, err := getSecondsFromDurationString(accessCodeLifespanUserAction)
		if err != nil {
			return nil, err
		}
		realm.AccessCodeLifespanUserAction = accessCodeLifespanUserActionDurationString
	}

	if actionTokenGeneratedByUserLifespan := data.Get("action_token_generated_by_user_lifespan").(string); actionTokenGeneratedByUserLifespan != "" {
		actionTokenGeneratedByUserLifespanDurationString, err := getSecondsFromDurationString(actionTokenGeneratedByUserLifespan)
		if err != nil {
			return nil, err
		}
		realm.ActionTokenGeneratedByUserLifespan = actionTokenGeneratedByUserLifespanDurationString
	}

	if actionTokenGeneratedByAdminLifespan := data.Get("action_token_generated_by_admin_lifespan").(string); actionTokenGeneratedByAdminLifespan != "" {
		actionTokenGeneratedByAdminLifespanDurationString, err := getSecondsFromDurationString(actionTokenGeneratedByAdminLifespan)
		if err != nil {
			return nil, err
		}
		realm.ActionTokenGeneratedByAdminLifespan = actionTokenGeneratedByAdminLifespanDurationString
	}

	//security defenses
	if v, ok := data.GetOk("security_defenses"); ok {
		securityDefensesSettings := v.([]interface{})[0].(map[string]interface{})

		headersConfig := securityDefensesSettings["headers"].([]interface{})
		if len(headersConfig) == 1 {
			headerSettings := headersConfig[0].(map[string]interface{})

			realm.Attributes = keycloak.Attributes{
				BrowserHeaderContentSecurityPolicy:           headerSettings["content_security_policy"].(string),
				BrowserHeaderContentSecurityPolicyReportOnly: headerSettings["content_security_policy_report_only"].(string),
				BrowserHeaderStrictTransportSecurity:         headerSettings["strict_transport_security"].(string),
				BrowserHeaderXContentTypeOptions:             headerSettings["x_content_type_options"].(string),
				BrowserHeaderXFrameOptions:                   headerSettings["x_frame_options"].(string),
				BrowserHeaderXRobotsTag:                      headerSettings["x_robots_tag"].(string),
				BrowserHeaderXXSSProtection:                  headerSettings["x_xss_protection"].(string),
			}
		} else {
			setDefaultSecuritySettings(realm)
		}
	} else {
		setDefaultSecuritySettings(realm)
	}

	if passwordPolicy, ok := data.GetOk("password_policy"); ok {
		realm.PasswordPolicy = passwordPolicy.(string)
	}

	//Flow Bindings
	if flow, ok := data.GetOk("browser_flow"); ok {
		realm.BrowserFlow = flow.(string)
	}

	if flow, ok := data.GetOk("registration_flow"); ok {
		realm.RegistrationFlow = flow.(string)
	}

	if flow, ok := data.GetOk("direct_grant_flow"); ok {
		realm.DirectGrantFlow = flow.(string)
	}

	if flow, ok := data.GetOk("reset_credentials_flow"); ok {
		realm.ResetCredentialsFlow = flow.(string)
	}

	if flow, ok := data.GetOk("client_authentication_flow"); ok {
		realm.ClientAuthenticationFlow = flow.(string)
	}

	if flow, ok := data.GetOk("docker_authentication_flow"); ok {
		realm.DockerAuthenticationFlow = flow.(string)
	}

	return realm, nil
}

func setDefaultSecuritySettings(realm *keycloak.Realm) {
	realm.Attributes = keycloak.Attributes{
		BrowserHeaderContentSecurityPolicy:           "frame-src 'self'; frame-ancestors 'self'; object-src 'none';",
		BrowserHeaderContentSecurityPolicyReportOnly: "",
		BrowserHeaderStrictTransportSecurity:         "max-age=31536000; includeSubDomains",
		BrowserHeaderXContentTypeOptions:             "nosniff",
		BrowserHeaderXFrameOptions:                   "SAMEORIGIN",
		BrowserHeaderXRobotsTag:                      "none",
		BrowserHeaderXXSSProtection:                  "1; mode=block",
	}
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

	// Smtp Config

	if (keycloak.SmtpServer{}) == realm.SmtpServer {
		data.Set("smtp_server", nil)
	} else {
		smtpSettings := make(map[string]interface{})

		smtpSettings["starttls"] = realm.SmtpServer.StartTls
		smtpSettings["port"] = realm.SmtpServer.Port
		smtpSettings["host"] = realm.SmtpServer.Host
		smtpSettings["reply_to"] = realm.SmtpServer.ReplyTo
		smtpSettings["reply_to_display_name"] = realm.SmtpServer.ReplyToDisplayName
		smtpSettings["from"] = realm.SmtpServer.From
		smtpSettings["from_display_name"] = realm.SmtpServer.FromDisplayName
		smtpSettings["envelope_from"] = realm.SmtpServer.EnvelopeFrom
		smtpSettings["ssl"] = realm.SmtpServer.Ssl

		if realm.SmtpServer.Auth {
			auth := make(map[string]interface{})

			auth["username"] = realm.SmtpServer.User
			auth["password"] = realm.SmtpServer.Password

			smtpSettings["auth"] = []interface{}{auth}
		}

		data.Set("smtp_server", []interface{}{smtpSettings})
	}

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

	//internationalization
	if realm.InternationalizationEnabled {
		internationalizationSettings := make(map[string]interface{})
		internationalizationSettings["supported_locales"] = realm.SupportLocales
		internationalizationSettings["default_locale"] = realm.DefaultLocale
		data.Set("internationalization", []interface{}{internationalizationSettings})
	} else {
		data.Set("internationalization", nil)
	}

	if _, ok := data.GetOk("security_defenses"); ok {

		if (keycloak.Attributes{}) == realm.Attributes {
			data.Set("security_defenses", nil)
		} else {
			securityDefensesSettings := make(map[string]interface{})

			headersSettings := make(map[string]interface{})

			headersSettings["content_security_policy"] = realm.Attributes.BrowserHeaderContentSecurityPolicy
			headersSettings["content_security_policy_report_only"] = realm.Attributes.BrowserHeaderContentSecurityPolicyReportOnly
			headersSettings["strict_transport_security"] = realm.Attributes.BrowserHeaderStrictTransportSecurity
			headersSettings["x_content_type_options"] = realm.Attributes.BrowserHeaderXContentTypeOptions
			headersSettings["x_frame_options"] = realm.Attributes.BrowserHeaderXFrameOptions
			headersSettings["x_robots_tag"] = realm.Attributes.BrowserHeaderXRobotsTag
			headersSettings["x_xss_protection"] = realm.Attributes.BrowserHeaderXXSSProtection

			securityDefensesSettings["headers"] = []interface{}{headersSettings}

			data.Set("security_defenses", []interface{}{securityDefensesSettings})
		}
	}

	data.Set("password_policy", realm.PasswordPolicy)

	//Flow Bindings
	data.Set("browser_flow", realm.BrowserFlow)
	data.Set("registration_flow", realm.RegistrationFlow)
	data.Set("direct_grant_flow", realm.DirectGrantFlow)
	data.Set("reset_credentials_flow", realm.ResetCredentialsFlow)
	data.Set("client_authentication_flow", realm.ClientAuthenticationFlow)
	data.Set("docker_authentication_flow", realm.DockerAuthenticationFlow)
}

func resourceKeycloakRealmCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := getRealmFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateRealm(realm)
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
		return handleNotFoundError(err, data)
	}

	// we can't trust the API to set this field correctly since it just responds with "**********" this implies a 'password only' change will not detected
	if smtpPassword, ok := getRealmSMTPPasswordFromData(data); ok {
		realm.SmtpServer.Password = smtpPassword
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

	err = keycloakClient.ValidateRealm(realm)
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
