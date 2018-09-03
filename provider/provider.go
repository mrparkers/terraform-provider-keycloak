package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func KeycloakProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"keycloak_realm": resourceKeycloakRealm(),
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"client_secret": {
				Required: true,
				Type:     schema.TypeString,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL of the Keycloak instance, before `/auth`",
			},
		},
		ConfigureFunc: configureKeycloakProvider,
	}
}

func configureKeycloakProvider(data *schema.ResourceData) (interface{}, error) {
	url := data.Get("url").(string)
	clientId := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)

	return keycloak.NewKeycloakClient(url, clientId, clientSecret)
}
