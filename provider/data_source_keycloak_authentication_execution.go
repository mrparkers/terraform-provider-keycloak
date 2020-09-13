package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakAuthenticationExecution() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakAuthenticationExecutionRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceKeycloakAuthenticationExecutionRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmID := data.Get("realm_id").(string)
	parentFlowAlias := data.Get("parent_flow_alias").(string)
	providerID := data.Get("provider_id").(string)

	authenticationExecutionInfo, err := keycloakClient.GetAuthenticationExecutionInfoFromProviderId(realmID, parentFlowAlias, providerID)
	if err != nil {
		return err
	}

	mapFromAuthenticationExecutionInfoToData(data, authenticationExecutionInfo)

	return nil
}
