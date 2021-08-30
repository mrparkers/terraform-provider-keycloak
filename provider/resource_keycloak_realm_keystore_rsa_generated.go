package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeystoreRsaGeneratedSize      = []int{1024, 2048, 4096}
	keycloakRealmKeystoreRsaGeneratedAlgorithm = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512"}
)

func resourceKeycloakRealmKeystoreRsaGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeystoreRsaGeneratedCreate,
		Read:   resourceKeycloakRealmKeystoreRsaGeneratedRead,
		Update: resourceKeycloakRealmKeystoreRsaGeneratedUpdate,
		Delete: resourceKeycloakRealmKeystoreRsaGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeystoreRsaGeneratedImport,
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
				ValidateFunc: validation.StringInSlice(keycloakRealmKeystoreRsaGeneratedAlgorithm, false),
				Default:      "RS256",
				Description:  "Intended algorithm for the key",
			},
			"key_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(keycloakRealmKeystoreRsaGeneratedSize),
				Default:      2048,
				Description:  "Size for the generated keys",
			},
		},
	}
}

func getRealmKeystoreRsaGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeystoreRsaGenerated, error) {
	mapper := &keycloak.RealmKeystoreRsaGenerated{
		Id:       data.Id(),
		Name:     data.Get("name").(string),
		RealmId:  data.Get("realm_id").(string),
		ParentId: data.Get("parent_id").(string),

		Active:    data.Get("active").(bool),
		Enabled:   data.Get("enabled").(bool),
		Priority:  data.Get("priority").(int),
		KeySize:   data.Get("key_size").(int),
		Algorithm: data.Get("algorithm").(string),
	}

	return mapper, nil
}

func setRealmKeystoreRsaGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeystoreRsaGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("key_size", realmKey.KeySize)
	data.Set("algorithm", realmKey.Algorithm)

	return nil
}

func resourceKeycloakRealmKeystoreRsaGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreRsaGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeystoreRsaGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreRsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeystoreRsaGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeystoreRsaGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeystoreRsaGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeystoreRsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreRsaGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreRsaGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeystoreRsaGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreRsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeystoreRsaGenerated(realmKey)
}

func resourceKeycloakRealmKeystoreRsaGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeystoreRsaGenerated(realmId, id)
}

func resourceKeycloakRealmKeystoreRsaGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{keystoreId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
