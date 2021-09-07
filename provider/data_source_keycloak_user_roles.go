package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakUserRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakUserRolesRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_role_names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakUserRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)

	roles, err := keycloakClient.GetUserRoleMappings(realmId, userId)
	if err != nil {
		return err
	}

	var roleNames []string

	for _, realmRole := range roles.RealmMappings {
		roleNames = append(roleNames, realmRole.Name)
	}

	data.Set("realm_role_names", roleNames)
	data.SetId(realmId + "/" + userId)

	return nil
}
