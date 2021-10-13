package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeystoreAesGeneratedSize = []int{16, 24, 32}
)

func resourceKeycloakRealmKeystoreAesGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeystoreAesGeneratedCreate,
		Read:   resourceKeycloakRealmKeystoreAesGeneratedRead,
		Update: resourceKeycloakRealmKeystoreAesGeneratedUpdate,
		Delete: resourceKeycloakRealmKeystoreAesGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeystoreAesGeneratedImport,
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
			"secret_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(keycloakRealmKeystoreAesGeneratedSize),
				Default:      16,
				Description:  "Size in bytes for the generated AES Key. Size 16 is for AES-128, Size 24 for AES-192 and Size 32 for AES-256. WARN: Bigger keys then 128 bits are not allowed on some JDK implementations",
			},
		},
	}
}

func getRealmKeystoreAesGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeystoreAesGenerated, error) {
	keystore := &keycloak.RealmKeystoreAesGenerated{
		Id:      data.Id(),
		Name:    data.Get("name").(string),
		RealmId: data.Get("realm_id").(string),

		Active:     data.Get("active").(bool),
		Enabled:    data.Get("enabled").(bool),
		Priority:   data.Get("priority").(int),
		SecretSize: data.Get("secret_size").(int),
	}

	return keystore, nil
}

func setRealmKeystoreAesGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeystoreAesGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("secret_size", realmKey.SecretSize)

	return nil
}

func resourceKeycloakRealmKeystoreAesGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreAesGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeystoreAesGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreAesGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeystoreAesGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeystoreAesGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeystoreAesGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeystoreAesGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreAesGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreAesGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeystoreAesGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreAesGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreAesGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeystoreAesGenerated(realmId, id)
}

func resourceKeycloakRealmKeystoreAesGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{keystoreId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
