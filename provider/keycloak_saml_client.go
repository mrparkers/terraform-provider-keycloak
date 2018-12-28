package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strconv"
	"strings"
)

var (
	keycloakSamlClientNameIdFormats = []string{"username", "email", "transient", "persistent"}
)

func resourceKeycloakSamlClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlClientCreate,
		Read:   resourceKeycloakSamlClientRead,
		Delete: resourceKeycloakSamlClientDelete,
		Update: resourceKeycloakSamlClientUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakSamlClientImport,
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_authn_statement": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sign_documents": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sign_assertions": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"client_signature_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_post_binding": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"front_channel_logout": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name_id_format": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "username",
				ValidateFunc: validation.StringInSlice(keycloakSamlClientNameIdFormats, false),
			},
			"root_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"master_saml_processing_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signing_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func mapToSamlClientFromData(data *schema.ResourceData) *keycloak.SamlClient {
	var validRedirectUris []string

	if v, ok := data.GetOk("valid_redirect_uris"); ok {
		for _, validRedirectUri := range v.(*schema.Set).List() {
			validRedirectUris = append(validRedirectUris, validRedirectUri.(string))
		}
	}

	samlClient := &keycloak.SamlClient{
		Id:                      data.Id(),
		ClientId:                data.Get("client_id").(string),
		RealmId:                 data.Get("realm_id").(string),
		Name:                    data.Get("name").(string),
		Enabled:                 data.Get("enabled").(bool),
		Description:             data.Get("description").(string),
		FrontChannelLogout:      data.Get("front_channel_logout").(bool),
		RootUrl:                 data.Get("root_url").(string),
		ValidRedirectUris:       validRedirectUris,
		BaseUrl:                 data.Get("base_url").(string),
		MasterSamlProcessingUrl: data.Get("master_saml_processing_url").(string),
		Attributes: &keycloak.SamlClientAttributes{
			IncludeAuthnStatement:   strconv.FormatBool(data.Get("include_authn_statement").(bool)),
			SignDocuments:           strconv.FormatBool(data.Get("sign_documents").(bool)),
			SignAssertions:          strconv.FormatBool(data.Get("sign_assertions").(bool)),
			ClientSignatureRequired: strconv.FormatBool(data.Get("client_signature_required").(bool)),
			ForcePostBinding:        strconv.FormatBool(data.Get("force_post_binding").(bool)),
			NameIdFormat:            data.Get("name_id_format").(string),
			SigningCertificate:      data.Get("signing_certificate").(string),
			SigningPrivateKey:       data.Get("signing_private_key").(string),
		},
	}

	return samlClient
}

func mapToDataFromSamlClient(data *schema.ResourceData, client *keycloak.SamlClient) error {
	data.SetId(client.Id)

	includeAuthnStatement, err := strconv.ParseBool(client.Attributes.IncludeAuthnStatement)
	if err != nil {
		return err
	}

	signDocuments, err := strconv.ParseBool(client.Attributes.SignDocuments)
	if err != nil {
		return err
	}

	signAssertions, err := strconv.ParseBool(client.Attributes.SignAssertions)
	if err != nil {
		return err
	}

	clientSignatureRequired, err := strconv.ParseBool(client.Attributes.ClientSignatureRequired)
	if err != nil {
		return err
	}

	forcePostBinding, err := strconv.ParseBool(client.Attributes.ForcePostBinding)
	if err != nil {
		return err
	}

	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
	data.Set("name", client.Name)
	data.Set("enabled", client.Enabled)
	data.Set("description", client.Description)
	data.Set("front_channel_logout", client.FrontChannelLogout)
	data.Set("root_url", client.RootUrl)
	data.Set("valid_redirect_uris", client.ValidRedirectUris)
	data.Set("base_url", client.BaseUrl)
	data.Set("master_saml_processing_url", client.MasterSamlProcessingUrl)
	data.Set("include_authn_statement", includeAuthnStatement)
	data.Set("sign_documents", signDocuments)
	data.Set("sign_assertions", signAssertions)
	data.Set("client_signature_required", clientSignatureRequired)
	data.Set("force_post_binding", forcePostBinding)
	data.Set("name_id_format", client.Attributes.NameIdFormat)
	data.Set("signing_certificate", client.Attributes.SigningCertificate)
	data.Set("signing_private_key", client.Attributes.SigningPrivateKey)

	return nil
}

func resourceKeycloakSamlClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := mapToSamlClientFromData(data)

	err := keycloakClient.NewSamlClient(client)
	if err != nil {
		return err
	}

	err = mapToDataFromSamlClient(data, client)
	if err != nil {
		return err
	}

	return resourceKeycloakSamlClientRead(data, meta)
}

func resourceKeycloakSamlClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetSamlClient(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = mapToDataFromSamlClient(data, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakSamlClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := mapToSamlClientFromData(data)

	err := keycloakClient.UpdateSamlClient(client)
	if err != nil {
		return err
	}

	err = mapToDataFromSamlClient(data, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakSamlClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteSamlClient(realmId, id)
}

func resourceKeycloakSamlClientImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	realm := parts[0]
	id := parts[1]

	d.Set("realm_id", realm)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
