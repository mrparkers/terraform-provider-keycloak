package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeyJavaKeystoreAlgorithm = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512"}
)

func resourceKeycloakRealmKeyJavaKeystore() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeyJavaKeystoreCreate,
		Read:   resourceKeycloakRealmKeyJavaKeystoreRead,
		Update: resourceKeycloakRealmKeyJavaKeystoreUpdate,
		Delete: resourceKeycloakRealmKeyJavaKeystoreDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeyJavaKeystoreImport,
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
				ValidateFunc: validation.StringInSlice(keycloakRealmKeyJavaKeystoreAlgorithm, false),
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

func getRealmKeyJavaKeystoreFromData(data *schema.ResourceData) (*keycloak.RealmKeyJavaKeystore, error) {
	mapper := &keycloak.RealmKeyJavaKeystore{
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

func setRealmKeyJavaKeystoreData(data *schema.ResourceData, realmKey *keycloak.RealmKeyJavaKeystore) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("keystore", realmKey.Keystore)
	data.Set("keystorePassword", realmKey.KeystorePassword)
	data.Set("keyAlias", realmKey.KeyAlias)
	data.Set("keyPassword", realmKey.KeyPassword)

	return nil
}

func resourceKeycloakRealmKeyJavaKeystoreCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyJavaKeystoreFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeyJavaKeystore(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyJavaKeystoreData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeyJavaKeystoreRead(data, meta)
}

func resourceKeycloakRealmKeyJavaKeystoreRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeyJavaKeystore(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeyJavaKeystoreData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeyJavaKeystoreUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyJavaKeystoreFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeyJavaKeystore(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyJavaKeystoreData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeyJavaKeystore(realmKey)
}

func resourceKeycloakRealmKeyJavaKeystoreDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeyJavaKeystore(realmId, id)
}

func resourceKeycloakRealmKeyJavaKeystoreImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
