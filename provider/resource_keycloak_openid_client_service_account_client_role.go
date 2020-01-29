package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenidClientServiceAccountClientRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientServiceAccountClientRoleCreate,
		Read:   resourceKeycloakOpenidClientServiceAccountClientRoleRead,
		Delete: resourceKeycloakOpenidClientServiceAccountClientRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientServiceAccountClientRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"service_account_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getOpenidClientServiceAccountClientRoleFromData(data *schema.ResourceData, keycloakClient *keycloak.KeycloakClient) (*keycloak.OpenidClientServiceAccountClientRole, error) {
	roleName := data.Get("role").(string)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	serviceAccountRoleId := data.Get("service_account_user_id").(string)

	role, err := keycloakClient.GetRoleByName(realmId, clientId, roleName)
	if err != nil {
		if keycloak.ErrorIs404(err) {
			role = &keycloak.Role{Id: ""}
		} else {
			return nil, err
		}
	}

	return &keycloak.OpenidClientServiceAccountClientRole{
		Id:                   role.Id,
		Name:                 roleName,
		RealmId:              realmId,
		ClientId:             clientId,
		ServiceAccountUserId: serviceAccountRoleId,
	}, nil
}

func setOpenidClientServiceAccountClientRoleData(data *schema.ResourceData, serviceAccountRole *keycloak.OpenidClientServiceAccountClientRole) {
	data.SetId(fmt.Sprintf("%s/%s", serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id))
	data.Set("realm_id", serviceAccountRole.RealmId)
	data.Set("client_id", serviceAccountRole.ClientId)
	data.Set("service_account_user_id", serviceAccountRole.ServiceAccountUserId)
	data.Set("role", serviceAccountRole.Name)
}

func resourceKeycloakOpenidClientServiceAccountClientRoleCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	serviceAccountRole, err := getOpenidClientServiceAccountClientRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClientServiceAccountClientRole(serviceAccountRole)
	if err != nil {
		return err
	}
	setOpenidClientServiceAccountClientRoleData(data, serviceAccountRole)
	return resourceKeycloakOpenidClientServiceAccountClientRoleRead(data, meta)
}

func resourceKeycloakOpenidClientServiceAccountClientRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	serviceAccountRole, err := getOpenidClientServiceAccountClientRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	serviceAccountRole, err = keycloakClient.GetOpenidClientServiceAccountClientRole(serviceAccountRole.RealmId, serviceAccountRole.ClientId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientServiceAccountClientRoleData(data, serviceAccountRole)

	return nil
}

func resourceKeycloakOpenidClientServiceAccountClientRoleDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	serviceAccountRole, err := getOpenidClientServiceAccountClientRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.DeleteOpenidClientServiceAccountClientRole(serviceAccountRole.RealmId, serviceAccountRole.ClientId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	return nil
}

func resourceKeycloakOpenidClientServiceAccountClientRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{serviceAccountUserId}}/{{roleId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("client_id", parts[1])
	d.Set("service_account_user_id", parts[2])
	d.SetId(fmt.Sprintf("%s/%s/%s", parts[1], parts[2], parts[3]))

	return []*schema.ResourceData{d}, nil
}
