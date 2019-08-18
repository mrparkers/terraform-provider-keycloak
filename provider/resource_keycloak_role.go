package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"strings"
)

func resourceKeycloakRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRoleCreate,
		Read:   resourceKeycloakRoleRead,
		Delete: resourceKeycloakRoleDelete,
		Update: resourceKeycloakRoleUpdate,
		// This resource can be imported using {{realm}}/{{role_name}} or {{realm}}/{{client_id}}/{{role_name}}.
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func mapFromDataToRole(data *schema.ResourceData) *keycloak.Role {
	role := &keycloak.Role{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		ClientId:    data.Get("client_id").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}

	return role
}

func mapFromRoleToData(data *schema.ResourceData, role *keycloak.Role) {
	data.SetId(role.Id)

	data.Set("realm_id", role.RealmId)
	data.Set("client_id", role.ClientId)
	data.Set("name", role.Name)
	data.Set("description", role.Description)
}

func resourceKeycloakRoleCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToRole(data)

	err := keycloakClient.CreateRole(role)
	if err != nil {
		return err
	}

	mapFromRoleToData(data, role)

	return resourceKeycloakRoleRead(data, meta)
}

func resourceKeycloakRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	var role *keycloak.Role
	var err error

	if ok, _ := regexp.MatchString(`.+/.+(/.+)?`, id); ok { // the ID is set to {{realm}}/{{role}} or {{realm}}/{{client}}/{{role}} during import
		clientId := data.Get("client_id").(string)
		name := data.Get("name").(string)

		role, err = keycloakClient.GetRoleByName(realmId, clientId, name)
	} else {
		role, err = keycloakClient.GetRole(realmId, id)
	}

	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromRoleToData(data, role)

	return nil
}

func resourceKeycloakRoleUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToRole(data)

	err := keycloakClient.UpdateRole(role)
	if err != nil {
		return err
	}

	mapFromRoleToData(data, role)

	return nil
}

func resourceKeycloakRoleDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRole(realmId, id)
}

func resourceKeycloakRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) == 2 {
		d.Set("realm_id", parts[0])
		d.Set("name", parts[1])
	} else if len(parts) == 3 {
		d.Set("realm_id", parts[0])
		d.Set("client_id", parts[1])
		d.Set("name", parts[2])
	} else {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realm}}/{{role_name}} or {{realm}}/{{client_id}}/{{role_name}}")
	}

	return []*schema.ResourceData{d}, nil
}
