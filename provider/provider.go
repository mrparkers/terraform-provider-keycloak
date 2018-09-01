package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func KeycloakProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
	}
}
