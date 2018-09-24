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
				Default:  false,
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
		},
	}
}

func getRealmFromData(data *schema.ResourceData) *keycloak.Realm {
	loginTheme := data.Get("login_theme").(string)
	accountTheme := data.Get("account_theme").(string)
	adminTheme := data.Get("admin_theme").(string)
	emailTheme := data.Get("email_theme").(string)

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

	if loginTheme != "" {
		realm.LoginTheme = loginTheme
	}

	if accountTheme != "" {
		realm.AccountTheme = accountTheme
	}

	if adminTheme != "" {
		realm.AdminTheme = adminTheme
	}

	if emailTheme != "" {
		realm.EmailTheme = emailTheme
	}

	return realm
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

	data.Set("login_theme", realm.LoginTheme)
	data.Set("account_theme", realm.AccountTheme)
	data.Set("admin_theme", realm.AdminTheme)
	data.Set("email_theme", realm.EmailTheme)
}

func resourceKeycloakRealmCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := getRealmFromData(data)

	err := realm.Validate(keycloakClient)
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

	realm := getRealmFromData(data)

	err := realm.Validate(keycloakClient)
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
