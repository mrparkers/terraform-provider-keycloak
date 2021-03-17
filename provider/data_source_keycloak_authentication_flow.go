package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakAuthenticationFlow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakAuthenticationFlowRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceKeycloakAuthenticationFlowRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmID := data.Get("realm_id").(string)
	alias := data.Get("alias").(string)

	authenticationFlowInfo, err := keycloakClient.GetAuthenticationFlowFromAlias(realmID, alias)
	if err != nil {
		return err
	}

	mapFromAuthenticationFlowInfoToData(data, authenticationFlowInfo)

	return nil
}
