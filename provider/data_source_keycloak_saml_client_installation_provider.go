package provider

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakSamlClientInstallationProvider() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakSamlClientInstallationProviderRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakSamlClientInstallationProviderRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	cliendId := data.Get("client_id").(string)
	providerId := data.Get("provider_id").(string)

	value, err := keycloakClient.GetSamlClientInstallationProvider(realmId, cliendId, providerId)
	if err != nil {
		return err
	}

	h := sha1.New()
	h.Write(value)
	id := base64.URLEncoding.EncodeToString(h.Sum(nil))

	data.SetId(id)
	data.Set("realm_id", realmId)
	data.Set("client_id", cliendId)
	data.Set("provider_id", providerId)
	data.Set("value", string(value))

	return nil
}
