package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakCustomUserFederation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakCustomUserFederationCreate,
		Read:   resourceKeycloakCustomUserFederationRead,
		Update: resourceKeycloakCustomUserFederationUpdate,
		Delete: resourceKeycloakCustomUserFederationDelete,
		// This resource can be imported using {{realm}}/{{provider_id}}. The Provider ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakCustomUserFederationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the provider when displayed in the console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The realm this provider will provide user federation for.",
			},
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the custom provider, specified in the `getId` implementation for the UserStorageProviderFactory interface",
			},

			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "When false, this provider will not be used when performing queries for users.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Priority of this provider when looking up users. Lower values are first.",
			},

			"cache_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice(keycloakUserFederationCachePolicies, false),
			},
		},
	}
}

func getCustomUserFederationFromData(data *schema.ResourceData, client *keycloak.KeycloakClient) *keycloak.CustomUserFederation {

	var realmId string
	if v, ok := data.GetOk("realm_id"); ok {
		realmId = v.(string)
	} else {
		realmId = client.RealmId
	}
	return &keycloak.CustomUserFederation{
		Id:         data.Id(),
		Name:       data.Get("name").(string),
		RealmId:    realmId,
		ProviderId: data.Get("provider_id").(string),

		Enabled:  data.Get("enabled").(bool),
		Priority: data.Get("priority").(int),

		CachePolicy: data.Get("cache_policy").(string),
	}
}

func setCustomUserFederationData(data *schema.ResourceData, custom *keycloak.CustomUserFederation) {
	data.SetId(custom.Id)

	data.Set("name", custom.Name)
	data.Set("realm_id", custom.RealmId)
	data.Set("provider_id", custom.ProviderId)

	data.Set("enabled", custom.Enabled)
	data.Set("priority", custom.Priority)

	data.Set("cache_policy", custom.CachePolicy)
}

func resourceKeycloakCustomUserFederationCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	custom := getCustomUserFederationFromData(data, keycloakClient)

	err := keycloakClient.ValidateCustomUserFederation(custom)
	if err != nil {
		return err
	}

	err = keycloakClient.NewCustomUserFederation(custom)
	if err != nil {
		return err
	}

	setCustomUserFederationData(data, custom)

	return resourceKeycloakCustomUserFederationRead(data, meta)
}

func resourceKeycloakCustomUserFederationRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	custom, err := keycloakClient.GetCustomUserFederation(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setCustomUserFederationData(data, custom)

	return nil
}

func resourceKeycloakCustomUserFederationUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	custom := getCustomUserFederationFromData(data, keycloakClient)

	err := keycloakClient.ValidateCustomUserFederation(custom)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateCustomUserFederation(custom)
	if err != nil {
		return err
	}

	setCustomUserFederationData(data, custom)

	return nil
}

func resourceKeycloakCustomUserFederationDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteCustomUserFederation(realmId, id)
}

func resourceKeycloakCustomUserFederationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	keycloakClient := meta.(*keycloak.KeycloakClient)

	var realmId, id string
	switch len(parts) {
	case 1:
		realmId = keycloakClient.RealmId
		id = parts[0]
	case 2:
		realmId = parts[0]
		id = parts[1]
	default:
		return nil, fmt.Errorf("Resouce %s cannot be imported", d.Id())
	}

	d.Set("realm_id", realmId)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
