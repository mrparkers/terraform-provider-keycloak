package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var (
	keycloakClientScopeProtocols = []string{"openid-connect", "saml"}
)

func resourceKeycloakClientScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakClientScopeCreate,
		Read:   resourceKeycloakClientScopeRead,
		Delete: resourceKeycloakClientScopeDelete,
		Update: resourceKeycloakClientScopeUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "openid-connect",
				ValidateFunc: validation.StringInSlice(keycloakClientScopeProtocols, false),
			},
			"consent_screen_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getClientScopeFromData(data *schema.ResourceData) *keycloak.ClientScope {
	clientScope := &keycloak.ClientScope{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Protocol:    data.Get("protocol").(string),
	}

	if consentScreenText, ok := data.GetOk("consent_screen_text"); ok {
		clientScope.Attributes.ConsentScreenText = consentScreenText.(string)
		clientScope.Attributes.DisplayOnConsentScreen = "true"
	} else {
		clientScope.Attributes.DisplayOnConsentScreen = "false"
	}

	return clientScope
}

func setClientScopeData(data *schema.ResourceData, clientScope *keycloak.ClientScope) {
	data.SetId(clientScope.Id)

	data.Set("realm_id", clientScope.RealmId)
	data.Set("name", clientScope.Name)
	data.Set("description", clientScope.Description)
	data.Set("protocol", clientScope.Protocol)

	if clientScope.Attributes.DisplayOnConsentScreen == "true" {
		data.Set("consent_screen_text", clientScope.Attributes.ConsentScreenText)
	}
}

func resourceKeycloakClientScopeCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getClientScopeFromData(data)

	err := keycloakClient.NewClientScope(clientScope)
	if err != nil {
		return err
	}

	setClientScopeData(data, clientScope)

	return resourceKeycloakClientScopeRead(data, meta)
}

func resourceKeycloakClientScopeRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	clientScope, err := keycloakClient.GetClientScope(realmId, id)
	if err != nil {
		return err
	}

	setClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakClientScopeUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getClientScopeFromData(data)

	err := keycloakClient.UpdateClientScope(clientScope)
	if err != nil {
		return err
	}

	setClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakClientScopeDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteClientScope(realmId, id)
}
