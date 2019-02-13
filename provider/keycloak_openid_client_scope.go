package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenidClientScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientScopeCreate,
		Read:   resourceKeycloakOpenidClientScopeRead,
		Delete: resourceKeycloakOpenidClientScopeDelete,
		Update: resourceKeycloakOpenidClientScopeUpdate,
		// This resource can be imported using {{realm}}/{{client_scope_id}}. The Client Scope ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientScopeImport,
		},
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
			"consent_screen_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getClientScopeFromData(data *schema.ResourceData) *keycloak.OpenidClientScope {
	clientScope := &keycloak.OpenidClientScope{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}

	if consentScreenText, ok := data.GetOk("consent_screen_text"); ok {
		clientScope.Attributes.ConsentScreenText = consentScreenText.(string)
		clientScope.Attributes.DisplayOnConsentScreen = "true"
	} else {
		clientScope.Attributes.DisplayOnConsentScreen = "false"
	}

	return clientScope
}

func setClientScopeData(data *schema.ResourceData, clientScope *keycloak.OpenidClientScope) {
	data.SetId(clientScope.Id)

	data.Set("realm_id", clientScope.RealmId)
	data.Set("name", clientScope.Name)
	data.Set("description", clientScope.Description)

	if clientScope.Attributes.DisplayOnConsentScreen == "true" {
		data.Set("consent_screen_text", clientScope.Attributes.ConsentScreenText)
	}
}

func resourceKeycloakOpenidClientScopeCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getClientScopeFromData(data)

	err := keycloakClient.NewOpenidClientScope(clientScope)
	if err != nil {
		return err
	}

	setClientScopeData(data, clientScope)

	return resourceKeycloakOpenidClientScopeRead(data, meta)
}

func resourceKeycloakOpenidClientScopeRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	clientScope, err := keycloakClient.GetOpenidClientScope(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakOpenidClientScopeUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getClientScopeFromData(data)

	err := keycloakClient.UpdateOpenidClientScope(clientScope)
	if err != nil {
		return err
	}

	setClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakOpenidClientScopeDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientScope(realmId, id)
}

func resourceKeycloakOpenidClientScopeImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientScopeId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
