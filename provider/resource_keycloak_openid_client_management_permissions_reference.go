package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientManagementPermissionsReference() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenIdClientManagementPermissionsReferenceCreate,
		Read:   resourceKeycloakOpenIdClientManagementPermissionsReferenceRead,
		Delete: resourceKeycloakOpenIdClientManagementPermissionsReferenceDelete,
		Update: resourceKeycloakOpenIdClientManagementPermissionsReferenceUpdate,
		// This resource can be imported using {{realm}}/clients/{{client_id}}/management/permissions. The Client Id is displayed in the URL when editing it from the GUI.
		Importer: &schema.ResourceImporter{
			State: genericManagementPermissionsReferenceImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scopePermissions": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func mapFromDataToOpenIdClientManagementPermissionsReference(data *schema.ResourceData) *keycloak.OpenIdClientManagementPermissionsReference {
	return &keycloak.OpenIdClientManagementPermissionsReference{
		RealmId:          data.Get("realm_id").(string),
		ClientId:         data.Get("client_id").(string),
		Enabled:          data.Get("enabled").(bool),
		Resource:         data.Get("resource").(string),
		ScopePermissions: data.Get("scopePermissions").(map[string]string),
	}
}

func mapFromOpenIdClientManagementPermissionsReferenceToData(reference *keycloak.OpenIdClientManagementPermissionsReference, data *schema.ResourceData) {
	data.Set("realm_id", reference.RealmId)
	data.Set("client_id", reference.ClientId)
	data.Set("enabled", reference.Enabled)
	data.Set("resource", reference.Resource)
	data.Set("scopePermissions", reference.ScopePermissions)
}

func resourceKeycloakOpenIdClientManagementPermissionsReferenceCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	err := keycloakClient.CreateOpenIdClientManagementPermissionsReference(realmId, clientId)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdClientManagementPermissionsReferenceRead(data, meta)
}

func resourceKeycloakOpenIdClientManagementPermissionsReferenceRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	reference, err := keycloakClient.GetClientManagementPermissionsReference(realmId, clientId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdClientManagementPermissionsReferenceToData(reference, data)

	return nil
}

func resourceKeycloakOpenIdClientManagementPermissionsReferenceUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientManagementPermissionsReference := mapFromDataToOpenIdClientManagementPermissionsReference(data)

	err := keycloakClient.UpdateOpenIdClientManagementPermissionsReference(clientManagementPermissionsReference)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenIdClientManagementPermissionsReferenceRead(data, meta)
}

func resourceKeycloakOpenIdClientManagementPermissionsReferenceDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	return keycloakClient.DeleteOpenIdClientManagementPermissionsReference(realmId, clientId)
}
