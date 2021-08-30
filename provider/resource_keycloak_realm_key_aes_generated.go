package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeyAesGeneratedSize = []int{16, 24, 32}
)

func resourceKeycloakRealmKeyAesGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeyAesGeneratedCreate,
		Read:   resourceKeycloakRealmKeyAesGeneratedRead,
		Update: resourceKeycloakRealmKeyAesGeneratedUpdate,
		Delete: resourceKeycloakRealmKeyAesGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeyAesGeneratedImport,
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
			"secret_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(keycloakRealmKeyAesGeneratedSize),
				Default:      16,
				Description:  "Size in bytes for the generated AES Key. Size 16 is for AES-128, Size 24 for AES-192 and Size 32 for AES-256. WARN: Bigger keys then 128 bits are not allowed on some JDK implementations",
			},
		},
	}
}

func getRealmKeyAesGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeyAesGenerated, error) {
	mapper := &keycloak.RealmKeyAesGenerated{
		Id:       data.Id(),
		Name:     data.Get("name").(string),
		RealmId:  data.Get("realm_id").(string),
		ParentId: data.Get("parent_id").(string),

		Active:     data.Get("active").(bool),
		Enabled:    data.Get("enabled").(bool),
		Priority:   data.Get("priority").(int),
		SecretSize: data.Get("secret_size").(int),
	}

	return mapper, nil
}

func setRealmKeyAesGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeyAesGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("secretSize", realmKey.SecretSize)

	return nil
}

func resourceKeycloakRealmKeyAesGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyAesGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeyAesGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyAesGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeyAesGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeyAesGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeyAesGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeyAesGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeyAesGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyAesGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeyAesGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyAesGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeyAesGenerated(realmKey)
}

func resourceKeycloakRealmKeyAesGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeyAesGenerated(realmId, id)
}

func resourceKeycloakRealmKeyAesGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
