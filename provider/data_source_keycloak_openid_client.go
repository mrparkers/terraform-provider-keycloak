package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakOpenidClient() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakOpenidClientRead,

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"standard_flow_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"implicit_flow_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"direct_access_grants_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"web_origins": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"service_account_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_accounts_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"resource_server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization": {
				Type:     schema.TypeSet,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_enforcement_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_remote_resource_management": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKeycloakOpenidClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	client, err := keycloakClient.GetOpenidClientByClientId(realmId, clientId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return nil
}
