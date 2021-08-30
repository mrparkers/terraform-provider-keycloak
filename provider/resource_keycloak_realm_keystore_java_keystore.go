package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeystoreJavaKeystoreAlgorithm = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512"}
)

func resourceKeycloakRealmKeystoreJavaKeystore() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeystoreJavaKeystoreCreate,
		Read:   resourceKeycloakRealmKeystoreJavaKeystoreRead,
		Update: resourceKeycloakRealmKeystoreJavaKeystoreUpdate,
		Delete: resourceKeycloakRealmKeystoreJavaKeystoreDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeystoreJavaKeystoreImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of provider when linked in admin console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm in which the ldap user federation provider exists.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm in which the ldap user federation provider exists.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set if the keys can be used for signing",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set if the keys are enabled",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Priority for the provider",
			},
			"algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakRealmKeystoreJavaKeystoreAlgorithm, false),
				Default:      "RS256",
				Description:  "Intended algorithm for the key",
			},
			"keystore": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Intended algorithm for the key",
			},
			"keystore_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size for the generated keys",
			},
			"key_alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Intended algorithm for the key",
			},
			"key_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size for the generated keys",
			},
		},
	}
}

func getRealmKeystoreJavaKeystoreFromData(data *schema.ResourceData) (*keycloak.RealmKeystoreJavaKeystore, error) {
	mapper := &keycloak.RealmKeystoreJavaKeystore{
		Id:       data.Id(),
		Name:     data.Get("name").(string),
		RealmId:  data.Get("realm_id").(string),
		ParentId: data.Get("parent_id").(string),

		Active:           data.Get("active").(bool),
		Enabled:          data.Get("enabled").(bool),
		Priority:         data.Get("priority").(int),
		Keystore:         data.Get("keystore").(string),
		KeystorePassword: data.Get("keystore_password").(string),
		KeyAlias:         data.Get("key_alias").(string),
		KeyPassword:      data.Get("key_password").(string),
	}

	return mapper, nil
}

func setRealmKeystoreJavaKeystoreData(data *schema.ResourceData, realmKey *keycloak.RealmKeystoreJavaKeystore) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("keystore", realmKey.Keystore)
	data.Set("keystore_password", realmKey.KeystorePassword)
	data.Set("key_alias", realmKey.KeyAlias)
	data.Set("key_password", realmKey.KeyPassword)

	return nil
}

func resourceKeycloakRealmKeystoreJavaKeystoreCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreJavaKeystoreFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeystoreJavaKeystore(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreJavaKeystoreData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeystoreJavaKeystoreRead(data, meta)
}

func resourceKeycloakRealmKeystoreJavaKeystoreRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeystoreJavaKeystore(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeystoreJavaKeystoreData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreJavaKeystoreUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreJavaKeystoreFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeystoreJavaKeystore(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreJavaKeystoreData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeystoreJavaKeystore(realmKey)
}

func resourceKeycloakRealmKeystoreJavaKeystoreDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeystoreJavaKeystore(realmId, id)
}

func resourceKeycloakRealmKeystoreJavaKeystoreImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
