package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func KeycloakProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
		Schema: map[string]*schema.Schema{
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

	return keycloak.NewKeycloakClient(url), nil
}
