package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeyHmacGeneratedSize      = []int{16, 24, 32, 64, 128, 256, 512}
	keycloakRealmKeyHmacGeneratedAlgorithm = []string{"HS256", "HS384", "HS512"}
)

func resourceKeycloakRealmKeyHmacGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeyHmacGeneratedCreate,
		Read:   resourceKeycloakRealmKeyHmacGeneratedRead,
		Update: resourceKeycloakRealmKeyHmacGeneratedUpdate,
		Delete: resourceKeycloakRealmKeyHmacGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeyHmacGeneratedImport,
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
				ValidateFunc: validation.StringInSlice(keycloakRealmKeyHmacGeneratedAlgorithm, false),
				Default:      "HS256",
				Description:  "Intended algorithm for the key",
			},
			"secret_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(keycloakRealmKeyHmacGeneratedSize),
				Default:      64,
				Description:  "Size in bytes for the generated secret",
			},
		},
	}
}

func getRealmKeyHmacGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeyHmacGenerated, error) {
	mapper := &keycloak.RealmKeyHmacGenerated{
		Id:       data.Id(),
		Name:     data.Get("name").(string),
		RealmId:  data.Get("realm_id").(string),
		ParentId: data.Get("parent_id").(string),

		Active:     data.Get("active").(bool),
		Enabled:    data.Get("enabled").(bool),
		Priority:   data.Get("priority").(int),
		SecretSize: data.Get("secret_size").(int),
		Algorithm:  data.Get("algorithm").(string),
	}

	return mapper, nil
}

func setRealmKeyHmacGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeyHmacGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("secretSize", realmKey.SecretSize)
	data.Set("algorithm", realmKey.Algorithm)

	return nil
}

func resourceKeycloakRealmKeyHmacGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyHmacGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeyHmacGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyHmacGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeyHmacGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeyHmacGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeyHmacGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeyHmacGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeyHmacGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyHmacGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeyHmacGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyHmacGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeyHmacGenerated(realmKey)
}

func resourceKeycloakRealmKeyHmacGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeyHmacGenerated(realmId, id)
}

func resourceKeycloakRealmKeyHmacGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
