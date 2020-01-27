package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakOpenidClientAccessTypes                        = []string{"CONFIDENTIAL", "PUBLIC", "BEARER-ONLY"}
	keycloakOpenidClientAuthorizationPolicyEnforcementMode = []string{"ENFORCING", "PERMISSIVE", "DISABLED"}
	keycloakOpenidClientPkceCodeChallengeMethod            = []string{"", "plain", "S256"}
)

func resourceKeycloakOpenidClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientCreate,
		Read:   resourceKeycloakOpenidClientRead,
		Delete: resourceKeycloakOpenidClientDelete,
		Update: resourceKeycloakOpenidClientUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientImport,
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientAccessTypes, false),
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"standard_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"implicit_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"direct_access_grants_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"web_origins": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_accounts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pkce_code_challenge_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientPkceCodeChallengeMethod, false),
			},
			"exclude_session_state_from_auth_response": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_account_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_enforcement_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(keycloakOpenidClientAuthorizationPolicyEnforcementMode, false),
						},
						"allow_remote_resource_management": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"keep_defaults": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"full_scope_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func getOpenidClientFromData(data *schema.ResourceData) (*keycloak.OpenidClient, error) {
	validRedirectUris := make([]string, 0)
	webOrigins := make([]string, 0)

	if v, ok := data.GetOk("valid_redirect_uris"); ok {
		for _, validRedirectUri := range v.(*schema.Set).List() {
			validRedirectUris = append(validRedirectUris, validRedirectUri.(string))
		}
	}

	if v, ok := data.GetOk("web_origins"); ok {
		for _, webOrigin := range v.(*schema.Set).List() {
			webOrigins = append(webOrigins, webOrigin.(string))
		}
	}

	openidClient := &keycloak.OpenidClient{
		Id:                        data.Id(),
		ClientId:                  data.Get("client_id").(string),
		RealmId:                   data.Get("realm_id").(string),
		Name:                      data.Get("name").(string),
		Enabled:                   data.Get("enabled").(bool),
		Description:               data.Get("description").(string),
		ClientSecret:              data.Get("client_secret").(string),
		StandardFlowEnabled:       data.Get("standard_flow_enabled").(bool),
		ImplicitFlowEnabled:       data.Get("implicit_flow_enabled").(bool),
		DirectAccessGrantsEnabled: data.Get("direct_access_grants_enabled").(bool),
		ServiceAccountsEnabled:    data.Get("service_accounts_enabled").(bool),
		FullScopeAllowed:          data.Get("full_scope_allowed").(bool),
		Attributes: keycloak.OpenidClientAttributes{
			PkceCodeChallengeMethod:             data.Get("pkce_code_challenge_method").(string),
			ExcludeSessionStateFromAuthResponse: keycloak.KeycloakBoolQuoted(data.Get("exclude_session_state_from_auth_response").(bool)),
		},
		ValidRedirectUris: validRedirectUris,
		WebOrigins:        webOrigins,
		BaseUrl:           data.Get("base_url").(string),
	}

	if !openidClient.ImplicitFlowEnabled && !openidClient.StandardFlowEnabled {
		if _, ok := data.GetOk("valid_redirect_uris"); ok {
			return nil, errors.New("valid_redirect_uris cannot be set when standard or implicit flow is not enabled")
		}
	}

	if !openidClient.ImplicitFlowEnabled && !openidClient.StandardFlowEnabled && !openidClient.DirectAccessGrantsEnabled {
		if _, ok := data.GetOk("web_origins"); ok {
			return nil, errors.New("web_origins cannot be set when standard or implicit flow is not enabled")
		}
	}

	// access type
	if accessType := data.Get("access_type").(string); accessType == "PUBLIC" {
		openidClient.PublicClient = true
	} else if accessType == "BEARER-ONLY" {
		openidClient.BearerOnly = true
	}

	if v, ok := data.GetOk("authorization"); ok {
		openidClient.AuthorizationServicesEnabled = true
		authorizationSettingsData := v.(*schema.Set).List()[0]
		authorizationSettings := authorizationSettingsData.(map[string]interface{})
		openidClient.AuthorizationSettings = &keycloak.OpenidClientAuthorizationSettings{
			PolicyEnforcementMode:         authorizationSettings["policy_enforcement_mode"].(string),
			AllowRemoteResourceManagement: authorizationSettings["allow_remote_resource_management"].(bool),
			KeepDefaults:                  authorizationSettings["keep_defaults"].(bool),
		}
	} else {
		openidClient.AuthorizationServicesEnabled = false
	}
	return openidClient, nil
}

func setOpenidClientData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData, client *keycloak.OpenidClient) error {
	var serviceAccountUserId string
	if client.ServiceAccountsEnabled {
		serviceAccountUser, err := keycloakClient.GetOpenidClientServiceAccountUserId(client.RealmId, client.Id)
		if err != nil {
			return err
		}
		serviceAccountUserId = serviceAccountUser.Id
	}
	data.SetId(client.Id)
	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
	data.Set("name", client.Name)
	data.Set("enabled", client.Enabled)
	data.Set("description", client.Description)
	data.Set("client_secret", client.ClientSecret)
	data.Set("standard_flow_enabled", client.StandardFlowEnabled)
	data.Set("implicit_flow_enabled", client.ImplicitFlowEnabled)
	data.Set("direct_access_grants_enabled", client.DirectAccessGrantsEnabled)
	data.Set("service_accounts_enabled", client.ServiceAccountsEnabled)
	data.Set("valid_redirect_uris", client.ValidRedirectUris)
	data.Set("web_origins", client.WebOrigins)
	data.Set("base_url", client.BaseUrl)
	data.Set("authorization_services_enabled", client.AuthorizationServicesEnabled)
	data.Set("full_scope_allowed", client.FullScopeAllowed)

	if client.AuthorizationServicesEnabled {
		data.Set("resource_server_id", client.Id)
	}

	if client.ServiceAccountsEnabled {
		data.Set("service_account_user_id", serviceAccountUserId)
	} else {
		data.Set("service_account_user_id", "")
	}

	// access type
	if client.PublicClient {
		data.Set("access_type", "PUBLIC")
	} else if client.BearerOnly {
		data.Set("access_type", "BEARER-ONLY")
	} else {
		data.Set("access_type", "CONFIDENTIAL")
	}
	return nil
}

func resourceKeycloakOpenidClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client, err := getOpenidClientFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClient(client)
	if err != nil {
		return err
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenidClientRead(data, meta)
}

func resourceKeycloakOpenidClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetOpenidClient(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client, err := getOpenidClientFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenidClient(client)
	if err != nil {
		return err
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClient(realmId, id)
}

func resourceKeycloakOpenidClientImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientId}}")
	}
	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
