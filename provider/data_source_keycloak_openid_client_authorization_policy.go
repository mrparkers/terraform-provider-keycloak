package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakOpenidClientAuthorizationPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakOpenidClientAuthorizationPolicyRead,

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
				Required: true,
			},
			"decision_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func setOpenidClientAuthorizationPolicyData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationPolicy) {
	data.SetId(policy.Id)

	data.Set("client_id", policy.ClientId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("owner", policy.Owner)
	data.Set("logic", policy.Logic)
	data.Set("policies", policy.Policies)
	data.Set("resources", policy.Resources)
	data.Set("scopes", policy.Scopes)
	data.Set("type", policy.Type)
}

func dataSourceKeycloakOpenidClientAuthorizationPolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	name := data.Get("name").(string)

	client, err := keycloakClient.GetClientAuthorizationPolicyByName(realmId, clientId, name)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationPolicyData(data, client)

	return nil
}
