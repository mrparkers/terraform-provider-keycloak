package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUserCreate,
		Read:   resourceKeycloakUserRead,
		Delete: resourceKeycloakUserDelete,
		Update: resourceKeycloakUserUpdate,
		// This resource can be imported using {{realm}}/{{user_id}}. The User's ID is displayed in the GUI when editing
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakUserImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					username := i.(string)

					if strings.ToLower(username) != username {
						return nil, []error{fmt.Errorf("expected username %s to be all lowercase", username)}
					}

					return nil, nil
				},
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"initial_password": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: onlyDiffOnCreate,
				MaxItems:         1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"temporary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func onlyDiffOnCreate(_, _, _ string, d *schema.ResourceData) bool {
	return d.Id() != ""
}

func mapFromDataToUser(data *schema.ResourceData, client *keycloak.KeycloakClient) *keycloak.User {
	realmId := realmId(data, client)
	return &keycloak.User{
		Id:        data.Id(),
		RealmId:   realmId,
		Username:  data.Get("username").(string),
		Email:     data.Get("email").(string),
		FirstName: data.Get("first_name").(string),
		LastName:  data.Get("last_name").(string),
		Enabled:   data.Get("enabled").(bool),
	}
}

func mapFromUserToData(data *schema.ResourceData, user *keycloak.User) {
	data.SetId(user.Id)

	data.Set("realm_id", user.RealmId)
	data.Set("username", user.Username)
	data.Set("email", user.Email)
	data.Set("first_name", user.FirstName)
	data.Set("last_name", user.LastName)
	data.Set("enabled", user.Enabled)
}

func resourceKeycloakUserCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	user := mapFromDataToUser(data, keycloakClient)

	err := keycloakClient.NewUser(user)
	if err != nil {
		return err
	}

	v, isInitialPasswordSet := data.GetOk("initial_password")
	if isInitialPasswordSet {
		passwordBlock := v.([]interface{})[0].(map[string]interface{})
		passwordValue := passwordBlock["value"].(string)
		isPasswordTemporary := passwordBlock["temporary"].(bool)
		err := keycloakClient.ResetUserPassword(user.RealmId, user.Id, passwordValue, isPasswordTemporary)
		if err != nil {
			return err
		}
	}

	mapFromUserToData(data, user)

	return resourceKeycloakUserRead(data, meta)
}

func resourceKeycloakUserRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	user, err := keycloakClient.GetUser(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromUserToData(data, user)

	return nil
}

func resourceKeycloakUserUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	user := mapFromDataToUser(data, keycloakClient)

	err := keycloakClient.UpdateUser(user)
	if err != nil {
		return err
	}

	mapFromUserToData(data, user)

	return nil
}

func resourceKeycloakUserDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteUser(realmId, id)
}

func resourceKeycloakUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	keycloakClient := meta.(*keycloak.KeycloakClient)
	var realmId, id string
	switch len(parts) {
	case 1:
		realmId = keycloakClient.GetDefaultRealm()
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
