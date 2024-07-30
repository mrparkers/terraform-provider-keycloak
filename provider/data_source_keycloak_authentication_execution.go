package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakAuthenticationExecution() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeycloakAuthenticationExecutionRead,
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
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceKeycloakAuthenticationExecutionRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmID := data.Get("realm_id").(string)
	parentFlowAlias := data.Get("parent_flow_alias").(string)
	providerID := data.Get("provider_id").(string)
	priority := data.Get("priority").(int)

	authenticationExecution, err := keycloakClient.GetAuthenticationExecutionInfoFromProviderId(ctx, realmID, parentFlowAlias, providerID, priority)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromAuthenticationExecutionToData(data, authenticationExecution)

	return nil
}
