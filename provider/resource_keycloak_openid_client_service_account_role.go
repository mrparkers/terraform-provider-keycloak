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
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientServiceAccountRoleImport,
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
				Optional: true,
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

func getOpenidClientServiceAccountRoleFromData(data *schema.ResourceData, keycloakClient *keycloak.KeycloakClient) (*keycloak.OpenidClientServiceAccountRole, error) {
	containerId := data.Get("client_id").(string)
	roleName := data.Get("role").(string)
	realmId := data.Get("realm_id").(string)
	serviceAccountRoleId := data.Get("service_account_user_id").(string)

	role, err := keycloakClient.GetRoleByName(realmId, containerId, roleName)
	if err != nil {
		if keycloak.ErrorIs404(err) {
			role = &keycloak.Role{Id: ""}
		} else {
			return nil, err
		}
	}

	return &keycloak.OpenidClientServiceAccountRole{
		Id:                   role.Id,
		ContainerId:          containerId,
		Name:                 roleName,
		RealmId:              realmId,
		ServiceAccountUserId: serviceAccountRoleId,
	}, nil
}

func setOpenidClientServiceAccountRoleData(data *schema.ResourceData, serviceAccountRole *keycloak.OpenidClientServiceAccountRole) {
	data.SetId(fmt.Sprintf("%s/%s", serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id))
	data.Set("realm_id", serviceAccountRole.RealmId)
	data.Set("client_id", serviceAccountRole.ContainerId)
	data.Set("service_account_user_id", serviceAccountRole.ServiceAccountUserId)
	data.Set("role", serviceAccountRole.Name)
}

func resourceKeycloakOpenidClientServiceAccountRoleCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	serviceAccountRole, err := getOpenidClientServiceAccountRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClientServiceAccountRole(serviceAccountRole)
	if err != nil {
		return err
	}
	setOpenidClientServiceAccountRoleData(data, serviceAccountRole)
	return resourceKeycloakOpenidClientServiceAccountRoleRead(data, meta)
}

func resourceKeycloakOpenidClientServiceAccountRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	serviceAccountRole, err := getOpenidClientServiceAccountRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	serviceAccountRole, err = keycloakClient.GetOpenidClientServiceAccountRole(serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ContainerId, serviceAccountRole.Id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientServiceAccountRoleData(data, serviceAccountRole)

	return nil
}

func resourceKeycloakOpenidClientServiceAccountRoleDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	serviceAccountRole, err := getOpenidClientServiceAccountRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.DeleteOpenidClientServiceAccountRole(serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ContainerId, serviceAccountRole.Id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	return nil
}

func resourceKeycloakOpenidClientServiceAccountRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{serviceAccountUserId}}/{{clientId}}/{{roleId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("service_account_user_id", parts[1])
	d.Set("client_id", parts[2])
	d.SetId(fmt.Sprintf("%s/%s", parts[1], parts[3]))

	return []*schema.ResourceData{d}, nil
}
