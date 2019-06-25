package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakRequiredAction() *schema.Resource {

	return &schema.Resource{
		Create: resourceKeycloakRequiredActionsCreate,
		Read:   resourceKeycloakRequiredActionsRead,
		Delete: resourceKeycloakRequiredActionsDelete,
		Update: resourceKeycloakRequiredActionsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_action": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func getRequiredActionFromData(data *schema.ResourceData) (*keycloak.RequiredAction, error) {
	action := &keycloak.RequiredAction{
		Id:            fmt.Sprintf("%s/%s", data.Get("realm_id").(string), data.Get("alias").(string)),
		RealmId:       data.Get("realm_id").(string),
		Alias:         data.Get("alias").(string),
		Name:          data.Get("name").(string),
		Enabled:       data.Get("enabled").(bool),
		DefaultAction: data.Get("default_action").(bool),
		Priority:      data.Get("priority").(int),
		Config:        make(map[string][]string),
	}

	return action, nil
}

func setRequiredActionData(data *schema.ResourceData, action *keycloak.RequiredAction) {
	data.SetId(fmt.Sprintf("%s/%s", action.RealmId, action.Alias))
	data.Set("realm_id", action.RealmId)
	data.Set("alias", action.Alias)
	data.Set("name", action.Name)
	data.Set("enabled", action.Enabled)
	data.Set("default_action", action.DefaultAction)
	data.Set("priority", action.Priority)
}

func resourceKeycloakRequiredActionsCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	action, err := getRequiredActionFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.CreateRequiredAction(action)
	if err != nil {
		return err
	}

	setRequiredActionData(data, action)

	return resourceKeycloakRequiredActionsRead(data, meta)
}

func resourceKeycloakRequiredActionsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	action, err := keycloakClient.GetRequiredAction(data.Get("realm_id").(string), data.Get("alias").(string))
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setRequiredActionData(data, action)

	return nil
}

func resourceKeycloakRequiredActionsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	action, err := getRequiredActionFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRequiredAction(action)
	if err != nil {
		return err
	}

	setRequiredActionData(data, action)

	return nil
}

func resourceKeycloakRequiredActionsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmName := data.Get("realm_id").(string)
	alias := data.Get("alias").(string)

	return keycloakClient.DeleteRequiredAction(realmName, alias)
}
