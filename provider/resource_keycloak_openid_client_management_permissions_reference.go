package provider

import (
	"fmt"

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
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope_permissions": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func openIdClientManagementReferenceId(realmId, clientId string) string {
	return fmt.Sprintf("%s/%s", realmId, clientId)
}

func mapFromOpenIdClientManagementPermissionsReferenceToData(reference *keycloak.OpenIdClientManagementPermissionsReference, data *schema.ResourceData) {
	data.Set("realm_id", reference.RealmId)
	data.Set("client_id", reference.ClientId)
	data.Set("resource", reference.Resource)
	data.Set("scope_permissions", reference.ScopePermissions)
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

	data.SetId(openIdClientManagementReferenceId(realmId, clientId))

	reference, err := keycloakClient.GetOpenIdClientManagementPermissionsReference(realmId, clientId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromOpenIdClientManagementPermissionsReferenceToData(reference, data)

	return nil
}

func resourceKeycloakOpenIdClientManagementPermissionsReferenceUpdate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakOpenIdClientManagementPermissionsReferenceRead(data, meta)
}

func resourceKeycloakOpenIdClientManagementPermissionsReferenceDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	return keycloakClient.DeleteOpenIdClientManagementPermissionsReference(realmId, clientId)
}
