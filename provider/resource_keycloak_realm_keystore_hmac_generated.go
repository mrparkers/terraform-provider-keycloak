package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeystoreHmacGeneratedSize      = []int{16, 24, 32, 64, 128, 256, 512}
	keycloakRealmKeystoreHmacGeneratedAlgorithm = []string{"HS256", "HS384", "HS512"}
)

func resourceKeycloakRealmKeystoreHmacGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeystoreHmacGeneratedCreate,
		Read:   resourceKeycloakRealmKeystoreHmacGeneratedRead,
		Update: resourceKeycloakRealmKeystoreHmacGeneratedUpdate,
		Delete: resourceKeycloakRealmKeystoreHmacGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeystoreHmacGeneratedImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of provider when linked in admin console.",
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice(keycloakRealmKeystoreHmacGeneratedAlgorithm, false),
				Default:      "HS256",
				Description:  "Intended algorithm for the key",
			},
			"secret_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(keycloakRealmKeystoreHmacGeneratedSize),
				Default:      64,
				Description:  "Size in bytes for the generated secret",
			},
		},
	}
}

func getRealmKeystoreHmacGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeystoreHmacGenerated, error) {
	keystore := &keycloak.RealmKeystoreHmacGenerated{
		Id:      data.Id(),
		Name:    data.Get("name").(string),
		RealmId: data.Get("realm_id").(string),

		Active:     data.Get("active").(bool),
		Enabled:    data.Get("enabled").(bool),
		Priority:   data.Get("priority").(int),
		SecretSize: data.Get("secret_size").(int),
		Algorithm:  data.Get("algorithm").(string),
	}

	return keystore, nil
}

func setRealmKeystoreHmacGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeystoreHmacGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("secret_size", realmKey.SecretSize)
	data.Set("algorithm", realmKey.Algorithm)

	return nil
}

func resourceKeycloakRealmKeystoreHmacGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreHmacGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeystoreHmacGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreHmacGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeystoreHmacGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeystoreHmacGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeystoreHmacGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeystoreHmacGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreHmacGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreHmacGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeystoreHmacGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreHmacGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreHmacGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeystoreHmacGenerated(realmId, id)
}

func resourceKeycloakRealmKeystoreHmacGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{keystoreId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
