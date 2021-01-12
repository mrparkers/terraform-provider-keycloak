package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenidClientServiceAccountRealmRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientServiceAccountRealmRoleCreate,
		Read:   resourceKeycloakOpenidClientServiceAccountRealmRoleRead,
		Delete: resourceKeycloakOpenidClientServiceAccountRealmRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientServiceAccountRealmRoleImport,
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
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getOpenidClientServiceAccountRealmRoleFromData(data *schema.ResourceData, keycloakClient *keycloak.KeycloakClient) (*keycloak.OpenidClientServiceAccountRealmRole, error) {
	roleName := data.Get("role").(string)
	realmId := data.Get("realm_id").(string)
	serviceAccountRoleId := data.Get("service_account_user_id").(string)

	role, err := keycloakClient.GetRoleByName(realmId, "", roleName)
	if err != nil {
		if keycloak.ErrorIs404(err) {
			role = &keycloak.Role{Id: ""}
		} else {
			return nil, err
		}
	}

	return &keycloak.OpenidClientServiceAccountRealmRole{
		Id:                   role.Id,
		Name:                 roleName,
		RealmId:              realmId,
		ServiceAccountUserId: serviceAccountRoleId,
	}, nil
}

func setOpenidClientServiceAccountRealmRoleData(data *schema.ResourceData, serviceAccountRole *keycloak.OpenidClientServiceAccountRealmRole) {
	data.SetId(fmt.Sprintf("%s/%s", serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id))
	data.Set("realm_id", serviceAccountRole.RealmId)
	data.Set("service_account_user_id", serviceAccountRole.ServiceAccountUserId)
	data.Set("role", serviceAccountRole.Name)
}

func resourceKeycloakOpenidClientServiceAccountRealmRoleCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	serviceAccountRole, err := getOpenidClientServiceAccountRealmRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClientServiceAccountRealmRole(serviceAccountRole)
	if err != nil {
		return err
	}
	setOpenidClientServiceAccountRealmRoleData(data, serviceAccountRole)
	return resourceKeycloakOpenidClientServiceAccountRealmRoleRead(data, meta)
}

func resourceKeycloakOpenidClientServiceAccountRealmRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	serviceAccountRole, err := getOpenidClientServiceAccountRealmRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	serviceAccountRole, err = keycloakClient.GetOpenidClientServiceAccountRealmRole(serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientServiceAccountRealmRoleData(data, serviceAccountRole)

	return nil
}

func resourceKeycloakOpenidClientServiceAccountRealmRoleDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	serviceAccountRole, err := getOpenidClientServiceAccountRealmRoleFromData(data, keycloakClient)
	if err != nil {
		return err
	}

	err = keycloakClient.DeleteOpenidClientServiceAccountRealmRole(serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	return nil
}

func resourceKeycloakOpenidClientServiceAccountRealmRoleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{serviceAccountUserId}}/{{roleId}}")
	}

	realmId := parts[0]
	serviceAccountUserId := parts[1]
	roleId := parts[2]

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return nil, err
	}

	d.Set("realm_id", realmId)
	d.Set("service_account_user_id", serviceAccountUserId)
	d.Set("role", role.Name)
	d.SetId(fmt.Sprintf("%s/%s", serviceAccountUserId, roleId))

	return []*schema.ResourceData{d}, nil
}
