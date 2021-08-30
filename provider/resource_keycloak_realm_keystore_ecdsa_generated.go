package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeystoreEcdsaGeneratedEllipticCurve = []string{"P-256", "P-384", "P-521"}
)

func resourceKeycloakRealmKeystoreEcdsaGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeystoreEcdsaGeneratedCreate,
		Read:   resourceKeycloakRealmKeystoreEcdsaGeneratedRead,
		Update: resourceKeycloakRealmKeystoreEcdsaGeneratedUpdate,
		Delete: resourceKeycloakRealmKeystoreEcdsaGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeystoreEcdsaGeneratedImport,
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
			"elliptic_curve_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakRealmKeystoreEcdsaGeneratedEllipticCurve, false),
				Default:      "P-256",
				Description:  "Elliptic Curve used in ECDSA",
			},
		},
	}
}

func getRealmKeystoreEcdsaGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeystoreEcdsaGenerated, error) {
	mapper := &keycloak.RealmKeystoreEcdsaGenerated{
		Id:       data.Id(),
		Name:     data.Get("name").(string),
		RealmId:  data.Get("realm_id").(string),
		ParentId: data.Get("parent_id").(string),

		Active:        data.Get("active").(bool),
		Enabled:       data.Get("enabled").(bool),
		Priority:      data.Get("priority").(int),
		EllipticCurve: data.Get("elliptic_curve_key").(string),
	}

	return mapper, nil
}

func setRealmKeystoreEcdsaGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeystoreEcdsaGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("elliptic_curve_key", realmKey.EllipticCurve)

	return nil
}

func resourceKeycloakRealmKeystoreEcdsaGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreEcdsaGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeystoreEcdsaGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreEcdsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeystoreEcdsaGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeystoreEcdsaGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeystoreEcdsaGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeystoreEcdsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreEcdsaGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreEcdsaGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeystoreEcdsaGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreEcdsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeystoreEcdsaGenerated(realmKey)
}

func resourceKeycloakRealmKeystoreEcdsaGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeystoreEcdsaGenerated(realmId, id)
}

func resourceKeycloakRealmKeystoreEcdsaGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{keystoreId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
