package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakGenericClientRoleMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakGenericClientRoleMapperCreate,
		Read:   resourceKeycloakGenericClientRoleMapperRead,
		Delete: resourceKeycloakGenericClientRoleMapperDelete,

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
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKeycloakGenericClientRoleMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	roleId := data.Get("role_id").(string)

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return err
	}

	err = keycloakClient.CreateRoleScopeMapping(realmId, clientId, role)
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/client/%s/scope-mappings/%s/%s", realmId, clientId, role.ClientId, role.Id))

	return resourceKeycloakGenericClientRoleMapperRead(data, meta)
}

func resourceKeycloakGenericClientRoleMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	roleId := data.Get("role_id").(string)

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return err
	}

	mappedRole, err := keycloakClient.GetRoleScopeMapping(realmId, clientId, role)

	if mappedRole == nil {
		data.SetId("")
	}

	return nil
}

func resourceKeycloakGenericClientRoleMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	roleId := data.Get("role_id").(string)

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return err
	}

	return keycloakClient.DeleteRoleScopeMapping(realmId, clientId, role)
}


