package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

const MAX_ATTRIBUTE_VALUE_LEN = 255

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
			"email_verified": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"federated_identity": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_provider": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func mapFromDataToUser(data *schema.ResourceData) *keycloak.User {
	attributes := map[string][]string{}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = splitLen(value.(string), MAX_ATTRIBUTE_VALUE_LEN)
		}
	}

	federatedIdentities := &keycloak.FederatedIdentities{}

	if v, ok := data.GetOk("federated_identity"); ok {
		federatedIdentities = getUserFederatedIdentitiesFromData(v.(*schema.Set).List())
	}

	return &keycloak.User{
		Id:                  data.Id(),
		RealmId:             data.Get("realm_id").(string),
		Username:            data.Get("username").(string),
		Email:               data.Get("email").(string),
		EmailVerified:       data.Get("email_verified").(bool),
		FirstName:           data.Get("first_name").(string),
		LastName:            data.Get("last_name").(string),
		Enabled:             data.Get("enabled").(bool),
		Attributes:          attributes,
		FederatedIdentities: *federatedIdentities,
	}
}

func getUserFederatedIdentitiesFromData(data []interface{}) *keycloak.FederatedIdentities {
	var federatedIdentities keycloak.FederatedIdentities
	for _, d := range data {
		federatedIdentitiesData := d.(map[string]interface{})
		federatedIdentity := &keycloak.FederatedIdentity{
			IdentityProvider: federatedIdentitiesData["identity_provider"].(string),
			UserId:           federatedIdentitiesData["user_id"].(string),
			UserName:         federatedIdentitiesData["user_name"].(string),
		}
		federatedIdentities = append(federatedIdentities, federatedIdentity)
	}
	return &federatedIdentities
}

func mapFromUserToData(data *schema.ResourceData, user *keycloak.User) {
	federatedIdentities := []interface{}{}
	for _, federatedIdentity := range user.FederatedIdentities {
		identity := map[string]interface{}{
			"identity_provider": federatedIdentity.IdentityProvider,
			"user_id":           federatedIdentity.UserId,
			"user_name":         federatedIdentity.UserName,
		}
		federatedIdentities = append(federatedIdentities, identity)
	}
	attributes := map[string]string{}
	for k, v := range user.Attributes {
		attributes[k] = strings.Join(v, "")
	}
	data.SetId(user.Id)
	data.Set("realm_id", user.RealmId)
	data.Set("username", user.Username)
	data.Set("email", user.Email)
	data.Set("email_verified", user.EmailVerified)
	data.Set("first_name", user.FirstName)
	data.Set("last_name", user.LastName)
	data.Set("enabled", user.Enabled)
	data.Set("attributes", attributes)
	data.Set("federated_identity", federatedIdentities)
}

func resourceKeycloakUserCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	user := mapFromDataToUser(data)

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
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
