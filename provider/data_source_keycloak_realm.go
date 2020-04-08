package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakRealm() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakRealmRead,
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name_html": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_managed_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			// Login Config

			"registration_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"registration_email_as_username": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"edit_username_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"reset_password_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remember_me": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"verify_email": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"login_with_email_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"duplicate_emails_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_required": {
				Type:     schema.TypeString,
				Computed: true,
			},

			//Smtp server

			"smtp_server": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starttls": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reply_to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reply_to_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"envelope_from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auth": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
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
				Computed: true,
			},
			"account_theme": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_theme": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_theme": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Tokens

			"refresh_token_max_reuse": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sso_session_idle_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sso_session_max_lifespan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"offline_session_idle_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"offline_session_max_lifespan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_token_lifespan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_token_lifespan_for_implicit_flow": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_code_lifespan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_code_lifespan_login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_code_lifespan_user_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action_token_generated_by_user_lifespan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action_token_generated_by_admin_lifespan": {
				Type:     schema.TypeString,
				Computed: true,
			},

			//internationalization
			"internationalization": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"supported_locales": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							Computed: true,
						},
						"default_locale": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			//Security Defenses
			"security_defenses": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"x_frame_options": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content_security_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content_security_policy_report_only": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"x_content_type_options": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"x_robots_tag": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"x_xss_protection": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"strict_transport_security": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"brute_force_detection": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"permanent_lockout": { //Permanent Lockout
										Type:     schema.TypeBool,
										Computed: true,
									},
									"max_login_failures": { //failureFactor
										Type:     schema.TypeInt,
										Computed: true,
									},
									"wait_increment_seconds": { //Wait Increment
										Type:     schema.TypeInt,
										Computed: true,
									},
									"quick_login_check_milli_seconds": { //Quick Login Check Milli Seconds
										Type:     schema.TypeInt,
										Computed: true,
									},
									"minimum_quick_login_wait_seconds": { //Minimum Quick Login Wait
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_failure_wait_seconds": { //Max Wait
										Type:     schema.TypeInt,
										Computed: true,
									},
									"failure_reset_time_seconds": { //maxDeltaTimeSeconds
										Type:     schema.TypeInt,
										Computed: true,
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
				Computed:    true,
			},

			//flow bindings
			"browser_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for BrowserFlow",
				Computed:    true,
			},
			"registration_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for RegistrationFlow",
				Computed:    true,
			},
			"direct_grant_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for DirectGrantFlow",
				Computed:    true,
			},
			"reset_credentials_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for ResetCredentialsFlow",
				Computed:    true,
			},
			"client_authentication_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for ClientAuthenticationFlow",
				Computed:    true,
			},
			"docker_authentication_flow": {
				Type:        schema.TypeString,
				Description: "Which flow should be used for DockerAuthenticationFlow",
				Computed:    true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakRealmRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmName := data.Get("realm").(string)

	realm, err := keycloakClient.GetRealm(realmName)
	if err != nil {
		return err
	}

	setRealmData(data, realm)

	return nil
}
