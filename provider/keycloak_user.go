package provider

import (
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
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func mapFromDataToUser(data *schema.ResourceData) *keycloak.User {
	return &keycloak.User{
		Id:        data.Id(),
		RealmId:   data.Get("realm_id").(string),
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

	user := mapFromDataToUser(data)

	err := keycloakClient.NewUser(user)
	if err != nil {
		return err
	}

	initialPassword, isPasswordSet := data.GetOk("initial_password")
	if isPasswordSet {
		err := keycloakClient.ResetUserPassword(data.Get("realm_id").(string), user.Id, initialPassword.(string))
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

	user := mapFromDataToUser(data)

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

func resourceKeycloakUserImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	id := parts[1]

	d.Set("realm_id", realm)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
