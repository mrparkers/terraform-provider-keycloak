package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakRealm() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmCreate,
		Read:   resourceKeycloakRealmRead,
		Delete: resourceKeycloakRealmDelete,
		Update: resourceKeycloakRealmUpdate,
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getRealm(data *schema.ResourceData) *keycloak.Realm {
	return &keycloak.Realm{
		Id:          data.Get("realm").(string),
		Realm:       data.Get("realm").(string),
		Enabled:     data.Get("enabled").(bool),
		DisplayName: data.Get("display_name").(string),
	}
}

func setData(data *schema.ResourceData, realm *keycloak.Realm) {
	data.SetId(realm.Realm)
	data.Set("realm", realm.Realm)
	data.Set("enabled", realm.Enabled)
	data.Set("display_name", realm.DisplayName)
}

func resourceKeycloakRealmCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := getRealm(data)

	err := keycloakClient.NewRealm(realm)
	if err != nil {
		return err
	}

	setData(data, realm)

	return resourceKeycloakRealmRead(data, meta)
}

func resourceKeycloakRealmRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := keycloakClient.GetRealm(data.Id())
	if err != nil {
		return err
	}

	setData(data, realm)

	return nil
}

func resourceKeycloakRealmUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := getRealm(data)

	err := keycloakClient.UpdateRealm(realm)
	if err != nil {
		return err
	}

	setData(data, realm)

	return nil
}

func resourceKeycloakRealmDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	return keycloakClient.DeleteRealm(data.Id())
}
