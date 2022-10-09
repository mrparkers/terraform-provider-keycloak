package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeycloakGenericClientRoleMapper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakGenericRoleMapperCreate,
		ReadContext:   resourceKeycloakGenericRoleMapperRead,
		DeleteContext: resourceKeycloakGenericRoleMapperDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakGenericRoleMapperImport,
		},
		DeprecationMessage: "please use keycloak_generic_role_mapper instead",
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
				Description:   "The destination client of the role. Cannot be used at the same time as client_scope_id.",
				ConflictsWith: []string{"client_scope_id"},
			},
			"client_scope_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The destination client scope of the role. Cannot be used at the same time as client_id.",
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
