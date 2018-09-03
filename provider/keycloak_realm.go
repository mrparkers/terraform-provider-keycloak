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
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKeycloakRealmCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := keycloak.Realm{
		Id:    data.Get("realm").(string),
		Realm: data.Get("realm").(string),
	}

	err := keycloakClient.NewRealm(&realm)
	if err != nil {
		return err
	}

	data.SetId(realm.Realm)
	data.Set("realm", realm.Realm)

	return resourceKeycloakRealmRead(data, meta)
}

func resourceKeycloakRealmRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm, err := keycloakClient.GetRealm(data.Id())
	if err != nil {
		return err
	}

	data.SetId(realm.Realm)
	data.Set("realm", realm.Realm)

	return nil
}

func resourceKeycloakRealmDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	return keycloakClient.DeleteRealm(data.Id())
}
