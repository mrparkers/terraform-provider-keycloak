package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakSamlClient() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakSamlClientRead,

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
			"include_authn_statement": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sign_documents": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sign_assertions": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encrypt_assertions": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"client_signature_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force_post_binding": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"front_channel_logout": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force_name_id_format": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_id_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"base_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_saml_processing_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_initiated_sso_url_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_initiated_sso_relay_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assertion_consumer_post_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assertion_consumer_redirect_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logout_service_post_binding_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logout_service_redirect_binding_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_scope_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"authentication_flow_binding_overrides": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"browser_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direct_grant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKeycloakSamlClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	client, err := keycloakClient.GetSamlClientByClientId(realmId, clientId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = mapToDataFromSamlClient(data, client)
	if err != nil {
		return err
	}

	return nil
}
