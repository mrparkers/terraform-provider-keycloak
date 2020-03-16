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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm id where the associated client or client scope exists.",
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The destination client of the client role. Cannot be used at the same time as client_scope_id.",
				ConflictsWith: []string{"client_scope_id"},
			},
			"client_scope_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The destination client scope of the client role. Cannot be used at the same time as client_id.",
				ConflictsWith: []string{"client_id"},
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the role to assign",
			},
		},
	}
}

func resourceKeycloakGenericClientRoleMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)
	roleId := data.Get("role_id").(string)

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return err
	}

	err = keycloakClient.CreateRoleScopeMapping(realmId, clientId, clientScopeId, role)
	if err != nil {
		return err
	}

	if clientId != "" {
		data.SetId(fmt.Sprintf("%s/client/%s/scope-mappings/%s/%s", realmId, clientId, role.ClientId, role.Id))
	} else {
		data.SetId(fmt.Sprintf("%s/client-scope/%s/scope-mappings/%s/%s", realmId, clientScopeId, role.ClientId, role.Id))
	}

	return resourceKeycloakGenericClientRoleMapperRead(data, meta)
}

func resourceKeycloakGenericClientRoleMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)
	roleId := data.Get("role_id").(string)

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return err
	}

	mappedRole, err := keycloakClient.GetRoleScopeMapping(realmId, clientId, clientScopeId, role)

	if mappedRole == nil {
		data.SetId("")
	}

	return nil
}

func resourceKeycloakGenericClientRoleMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)
	roleId := data.Get("role_id").(string)

	role, err := keycloakClient.GetRole(realmId, roleId)
	if err != nil {
		return err
	}

	return keycloakClient.DeleteRoleScopeMapping(realmId, clientId, clientScopeId, role)
}
