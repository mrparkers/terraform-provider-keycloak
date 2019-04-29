package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenidClientServiceAccountRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientServiceAccountRoleCreate,
		Read:   resourceKeycloakOpenidClientServiceAccountRoleRead,
		Delete: resourceKeycloakOpenidClientServiceAccountRoleDelete,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientServiceAccountRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"source_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getOpenidClientServiceAccountRoleFromData(data *schema.ResourceData) *keycloak.OpenidClientServiceAccountRole {
	return &keycloak.OpenidClientServiceAccountRole{
		Id:          data.Id(),
		ContainerId: data.Get("source_client_id").(string),
		Name:        data.Get("role").(string),
		RealmId:     data.Get("realm_id").(string),
		ClientId:    data.Get("target_client_id").(string),
	}
}

func setOpenidClientServiceAccountRoleData(data *schema.ResourceData, serviceAccountRole *keycloak.OpenidClientServiceAccountRole) {
	data.SetId(serviceAccountRole.Id)
	data.Set("realm_id", serviceAccountRole.RealmId)
	data.Set("source_client_id", serviceAccountRole.ContainerId)
	data.Set("target_client_id", serviceAccountRole.ClientId)
	data.Set("role", serviceAccountRole.Name)
}

func resourceKeycloakOpenidClientServiceAccountRoleCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	serviceAccountRole := getOpenidClientServiceAccountRoleFromData(data)
	err := keycloakClient.NewOpenidClientServiceAccountRole(serviceAccountRole)
	if err != nil {
		return err
	}
	setOpenidClientServiceAccountRoleData(data, serviceAccountRole)
	return resourceKeycloakOpenidClientServiceAccountRoleRead(data, meta)
}

func resourceKeycloakOpenidClientServiceAccountRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	sourceClientId := data.Get("source_client_id").(string)
	targetClientId := data.Get("target_client_id").(string)
	id := data.Id()

	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountRole(realmId, sourceClientId, targetClientId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientServiceAccountRoleData(data, serviceAccountRole)

	return nil
}

func resourceKeycloakOpenidClientServiceAccountRoleDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	sourceClientId := data.Get("source_client_id").(string)
	targetClientId := data.Get("target_client_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientServiceAccountRole(realmId, sourceClientId, targetClientId, id)
}

func resourceKeycloakOpenidClientServiceAccountRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{sourceClientId}}/{{targetClientId}}/{{role}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("source_client_id", parts[1])
	d.Set("target_client_id", parts[2])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
