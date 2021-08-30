package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeyRsaAlgorithm = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512"}
	keycloakRealmKeyRsaSize      = []int{1024, 2048, 4096}
)

func resourceKeycloakRealmKeyRsa() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeyRsaCreate,
		Read:   resourceKeycloakRealmKeyRsaRead,
		Update: resourceKeycloakRealmKeyRsaUpdate,
		Delete: resourceKeycloakRealmKeyRsaDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeyRsaImport,
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
				ValidateFunc: validation.StringInSlice(keycloakRealmKeyRsaAlgorithm, false),
				Default:      "RS256",
				Description:  "Intended algorithm for the key",
			},
			"key_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(keycloakRealmKeyRsaSize),
				Default:      2048,
				Description:  "Size for the generated keys",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private RSA Key encoded in PEM format",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "X509 Certificate encoded in PEM format",
			},
		},
	}
}

func getRealmKeyRsaFromData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData) (*keycloak.RealmKeyRsa, error) {
	mapper := &keycloak.RealmKeyRsa{
		Id:       data.Id(),
		Name:     data.Get("name").(string),
		RealmId:  data.Get("realm_id").(string),
		ParentId: data.Get("parent_id").(string),

		Active:      data.Get("active").(bool),
		Enabled:     data.Get("enabled").(bool),
		Priority:    data.Get("priority").(int),
		KeySize:     data.Get("keySize").(int),
		Algorithm:   data.Get("algorithm").(string),
		PrivateKey:  data.Get("privateKey").(string),
		Certificate: data.Get("certificate").(string),
	}
	_, err := keycloakClient.VersionIsGreaterThanOrEqualTo(keycloak.Version_11)
	if err != nil {
		return nil, err
	}

	return mapper, nil
}

func setRealmKeyRsaData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData, realmKey *keycloak.RealmKeyRsa) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)
	data.Set("parent_id", realmKey.ParentId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("keySize", realmKey.KeySize)
	data.Set("algorithm", realmKey.Algorithm)
	data.Set("privateKey", realmKey.PrivateKey)
	data.Set("certificate", realmKey.Certificate)

	_, err := keycloakClient.VersionIsGreaterThanOrEqualTo(keycloak.Version_11)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeyRsaCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyRsaFromData(keycloakClient, data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeyRsa(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyRsaData(keycloakClient, data, realmKey)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeyRsaRead(data, meta)
}

func resourceKeycloakRealmKeyRsaRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeyRsa(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeyRsaData(keycloakClient, data, realmKey)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeyRsaUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeyRsaFromData(keycloakClient, data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeyRsa(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeyRsaData(keycloakClient, data, realmKey)
	if err != nil {
		return err
	}

	return keycloakClient.UpdateRealmKeyRsa(realmKey)
}

func resourceKeycloakRealmKeyRsaDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeyRsa(realmId, id)
}

func resourceKeycloakRealmKeyRsaImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
