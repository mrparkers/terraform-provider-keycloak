package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlClientScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlClientScopeCreate,
		Read:   resourceKeycloakSamlClientScopeRead,
		Delete: resourceKeycloakSamlClientScopeDelete,
		Update: resourceKeycloakSamlClientScopeUpdate,
		// This resource can be imported using {{realm}}/{{client_scope_id}}. The Client Scope ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakSamlClientScopeImport,
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
			"gui_order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func getSamlClientScopeFromData(data *schema.ResourceData) *keycloak.SamlClientScope {
	clientScope := &keycloak.SamlClientScope{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}

	if consentScreenText, ok := data.GetOk("consent_screen_text"); ok {
		clientScope.Attributes.ConsentScreenText = consentScreenText.(string)
		clientScope.Attributes.DisplayOnConsentScreen = true
	} else {
		clientScope.Attributes.DisplayOnConsentScreen = false
	}

	// Treat 0 as an empty string for the purpose of omitting the attribute to reset the order
	if guiOrder := data.Get("gui_order").(int); guiOrder != 0 {
		clientScope.Attributes.GuiOrder = strconv.Itoa(guiOrder)
	}

	return clientScope
}

func setSamlClientScopeData(data *schema.ResourceData, clientScope *keycloak.SamlClientScope) {
	data.SetId(clientScope.Id)

	data.Set("realm_id", clientScope.RealmId)
	data.Set("name", clientScope.Name)
	data.Set("description", clientScope.Description)

	if clientScope.Attributes.DisplayOnConsentScreen {
		data.Set("consent_screen_text", clientScope.Attributes.ConsentScreenText)
	}

	if guiOrder, err := strconv.Atoi(clientScope.Attributes.GuiOrder); err == nil {
		data.Set("gui_order", guiOrder)
	}
}

func resourceKeycloakSamlClientScopeCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getSamlClientScopeFromData(data)

	err := keycloakClient.NewSamlClientScope(clientScope)
	if err != nil {
		return err
	}

	setSamlClientScopeData(data, clientScope)

	return resourceKeycloakSamlClientScopeRead(data, meta)
}

func resourceKeycloakSamlClientScopeRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	clientScope, err := keycloakClient.GetSamlClientScope(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setSamlClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakSamlClientScopeUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getSamlClientScopeFromData(data)

	err := keycloakClient.UpdateSamlClientScope(clientScope)
	if err != nil {
		return err
	}

	setSamlClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakSamlClientScopeDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteSamlClientScope(realmId, id)
}

func resourceKeycloakSamlClientScopeImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{samlClientScopeId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
