package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAuthenticationBindings() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationBindingsCreate,
		Read:   resourceKeycloakAuthenticationBindingsRead,
		Delete: resourceKeycloakAuthenticationBindingsDelete,
		Update: resourceKeycloakAuthenticationBindingsUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
				Default:     "reset credentials",
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

func getAuthenticationBindingsFromData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData) (*keycloak.Realm, error) {
	realm, err := keycloakClient.GetRealm(data.Get("realm_id").(string))
	if err != nil {
		return nil, err
	}

	realm.BrowserFlow = data.Get("browser_flow").(string)
	realm.RegistrationFlow = data.Get("registration_flow").(string)
	realm.DirectGrantFlow = data.Get("direct_grant_flow").(string)
	realm.ResetCredentialsFlow = data.Get("reset_credentials_flow").(string)
	realm.ClientAuthenticationFlow = data.Get("client_authentication_flow").(string)
	realm.DockerAuthenticationFlow = data.Get("docker_authentication_flow").(string)

	return realm, nil
}

func setAuthenticationBindingsData(data *schema.ResourceData, realm *keycloak.Realm) {
	data.SetId(realm.Realm)
	data.Set("browser_flow", realm.BrowserFlow)
	data.Set("registration_flow", realm.RegistrationFlow)
	data.Set("direct_grant_flow", realm.DirectGrantFlow)
	data.Set("reset_credentials_flow", realm.ResetCredentialsFlow)
	data.Set("client_authentication_flow", realm.ClientAuthenticationFlow)
	data.Set("docker_authentication_flow", realm.DockerAuthenticationFlow)
}

func resetAuthenticationBindingsForRealm(realm *keycloak.Realm) {
	realm.BrowserFlow = "browser"
	realm.RegistrationFlow = "registration"
	realm.DirectGrantFlow = "direct grant"
	realm.ResetCredentialsFlow = "reset credentials"
	realm.ClientAuthenticationFlow = "clients"
	realm.DockerAuthenticationFlow = "docker auth"
}

func resourceKeycloakAuthenticationBindingsCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := keycloakClient.GetRealm(data.Get("realm_id").(string))
	if err != nil {
		return err
	}

	realm, err = getAuthenticationBindingsFromData(keycloakClient, data)
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

	realm, err = keycloakClient.GetRealm(realm.Id)
	if err != nil {
		return err
	}

	setAuthenticationBindingsData(data, realm)

	return nil
}

func resourceKeycloakAuthenticationBindingsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := keycloakClient.GetRealm(data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setAuthenticationBindingsData(data, realm)

	return nil
}

func resourceKeycloakAuthenticationBindingsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := keycloakClient.GetRealm(data.Id())
	if err != nil {
		return nil
	}

	resetAuthenticationBindingsForRealm(realm)

	err = keycloakClient.UpdateRealm(realm)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakAuthenticationBindingsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := getAuthenticationBindingsFromData(keycloakClient, data)
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

	setAuthenticationBindingsData(data, realm)

	return nil
}
