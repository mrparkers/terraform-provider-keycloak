package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakGroupRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceKeycloakGroupRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	groupName := data.Get("name").(string)

	group, err := keycloakClient.GetGroupByName(realmId, groupName)
	if err != nil {
		return err
	}

	mapFromGroupToData(data, group)

	return nil
}
