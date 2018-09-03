package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakClientCreate,
		Read:   resourceKeycloakClientRead,
		Delete: resourceKeycloakClientDelete,
		Update: resourceKeycloakClientUpdate,
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func getClientFromData(data *schema.ResourceData) *keycloak.Client {
	return &keycloak.Client{
		Id:       data.Id(),
		ClientId: data.Get("client_id").(string),
		RealmId:  data.Get("realm_id").(string),
	}
}

func setClientData(data *schema.ResourceData, client *keycloak.Client) {
	data.SetId(client.Id)

	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
}

func resourceKeycloakClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := getClientFromData(data)

	err := keycloakClient.NewClient(client)
	if err != nil {
		return err
	}

	setClientData(data, client)

	return resourceKeycloakClientRead(data, meta)
}

func resourceKeycloakClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetClient(realmId, id)
	if err != nil {
		return err
	}

	setClientData(data, client)

	return nil
}

func resourceKeycloakClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := getClientFromData(data)

	err := keycloakClient.UpdateClient(client)
	if err != nil {
		return err
	}

	setClientData(data, client)

	return nil
}

func resourceKeycloakClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteClient(realmId, id)
}
