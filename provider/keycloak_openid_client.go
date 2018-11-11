package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakOpenidClientAccessTypes = []string{"CONFIDENTIAL", "PUBLIC", "BEARER-ONLY"}
)

func resourceKeycloakOpenidClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientCreate,
		Read:   resourceKeycloakOpenidClientRead,
		Delete: resourceKeycloakOpenidClientDelete,
		Update: resourceKeycloakOpenidClientUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientImport,
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientAccessTypes, false),
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
		},
	}
}

func getOpenidClientFromData(data *schema.ResourceData, client *keycloak.KeycloakClient) *keycloak.OpenidClient {
	var validRedirectUris []string

	if v, ok := data.GetOk("valid_redirect_uris"); ok {
		for _, validRedirectUri := range v.(*schema.Set).List() {
			validRedirectUris = append(validRedirectUris, validRedirectUri.(string))
		}
	}

	var realmId string
	if v, ok := data.GetOk("realm_id"); ok {
		realmId = v.(string)
	} else {
		realmId = client.RealmId
	}

	openidClient := &keycloak.OpenidClient{
		Id:                data.Id(),
		ClientId:          data.Get("client_id").(string),
		RealmId:           realmId,
		Name:              data.Get("name").(string),
		Enabled:           data.Get("enabled").(bool),
		Description:       data.Get("description").(string),
		ClientSecret:      data.Get("client_secret").(string),
		ValidRedirectUris: validRedirectUris,
	}

	// access type
	if accessType := data.Get("access_type").(string); accessType == "PUBLIC" {
		openidClient.PublicClient = true
	} else if accessType == "BEARER-ONLY" {
		openidClient.BearerOnly = true
	}

	return openidClient
}

func setOpenidClientData(data *schema.ResourceData, client *keycloak.OpenidClient) {
	data.SetId(client.Id)

	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
	data.Set("name", client.Name)
	data.Set("enabled", client.Enabled)
	data.Set("description", client.Description)
	data.Set("client_secret", client.ClientSecret)
	data.Set("valid_redirect_uris", client.ValidRedirectUris)

	// access type
	if client.PublicClient {
		data.Set("access_type", "PUBLIC")
	} else if client.BearerOnly {
		data.Set("access_type", "BEARER-ONLY")
	} else {
		data.Set("access_type", "CONFIDENTIAL")
	}
}

func resourceKeycloakOpenidClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := getOpenidClientFromData(data, keycloakClient)

	err := keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClient(client)
	if err != nil {
		return err
	}

	setOpenidClientData(data, client)

	return resourceKeycloakOpenidClientRead(data, meta)
}

func resourceKeycloakOpenidClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetOpenidClient(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientData(data, client)

	return nil
}

func resourceKeycloakOpenidClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := getOpenidClientFromData(data, keycloakClient)

	err := keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenidClient(client)
	if err != nil {
		return err
	}

	setOpenidClientData(data, client)

	return nil
}

func resourceKeycloakOpenidClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClient(realmId, id)
}

func resourceKeycloakOpenidClientImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	id := parts[1]

	d.Set("realm_id", realm)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
