package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeyEcdsaGeneratedEllipticCurve = []string{"P-256", "P-384", "P-521"}
)

func resourceKeycloakRealmKeyEcdsaGenerated() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeyEcdsaGeneratedCreate,
		Read:   resourceKeycloakRealmKeyEcdsaGeneratedRead,
		Update: resourceKeycloakRealmKeyEcdsaGeneratedUpdate,
		Delete: resourceKeycloakRealmKeyEcdsaGeneratedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeyEcdsaGeneratedImport,
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
				ValidateFunc: validation.StringInSlice(keycloakRealmKeyEcdsaGeneratedEllipticCurve, false),
				Default:      "P-256",
				Description:  "Elliptic Curve used in ECDSA",
			},
		},
	}
}

func getRealmKeyEcdsaGeneratedFromData(data *schema.ResourceData) (*keycloak.RealmKeyEcdsaGenerated, error) {
	mapper := &keycloak.RealmKeyEcdsaGenerated{
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

func setRealmKeyEcdsaGeneratedData(data *schema.ResourceData, realmKey *keycloak.RealmKeyEcdsaGenerated) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("ecdsaEllipticCurveKey", realmKey.EllipticCurve)

	return nil
}

func resourceKeycloakRealmKeyEcdsaGeneratedCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyEcdsaGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeyEcdsaGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyEcdsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeyEcdsaGeneratedRead(data, meta)
}

func resourceKeycloakRealmKeyEcdsaGeneratedRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeyEcdsaGenerated(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeyEcdsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeyEcdsaGeneratedUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyEcdsaGeneratedFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeyEcdsaGenerated(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyEcdsaGeneratedData(data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeyEcdsaGenerated(realmKey)
}

func resourceKeycloakRealmKeyEcdsaGeneratedDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeyEcdsaGenerated(realmId, id)
}

func resourceKeycloakRealmKeyEcdsaGeneratedImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
