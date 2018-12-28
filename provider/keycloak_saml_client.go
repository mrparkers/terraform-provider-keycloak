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
				Computed: true,
			},
			"sign_documents": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sign_assertions": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"client_signature_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"force_post_binding": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

	samlAttributes := &keycloak.SamlClientAttributes{
		NameIdFormat:       data.Get("name_id_format").(string),
		SigningCertificate: data.Get("signing_certificate").(string),
		SigningPrivateKey:  data.Get("signing_private_key").(string),
	}

	if includeAuthnStatement, ok := data.GetOk("include_authn_statement"); ok {
		includeAuthnStatementString := strconv.FormatBool(includeAuthnStatement.(bool))
		samlAttributes.IncludeAuthnStatement = &includeAuthnStatementString
	}

	if signDocuments, ok := data.GetOk("sign_documents"); ok {
		signDocumentsString := strconv.FormatBool(signDocuments.(bool))
		samlAttributes.SignDocuments = &signDocumentsString
	}

	if signAssertions, ok := data.GetOk("sign_assertions"); ok {
		signAssertionsString := strconv.FormatBool(signAssertions.(bool))
		samlAttributes.SignAssertions = &signAssertionsString
	}

	if clientSignatureRequired, ok := data.GetOk("client_signature_required"); ok {
		clientSignatureRequiredString := strconv.FormatBool(clientSignatureRequired.(bool))
		samlAttributes.ClientSignatureRequired = &clientSignatureRequiredString
	}

	if forcePostBinding, ok := data.GetOk("force_post_binding"); ok {
		forcePostBindingString := strconv.FormatBool(forcePostBinding.(bool))
		samlAttributes.ForcePostBinding = &forcePostBindingString
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
		Attributes:              samlAttributes,
	}

	return samlClient
}

func mapToDataFromSamlClient(data *schema.ResourceData, client *keycloak.SamlClient) error {
	data.SetId(client.Id)

	if client.Attributes.IncludeAuthnStatement != nil {
		includeAuthnStatement, err := strconv.ParseBool(*client.Attributes.IncludeAuthnStatement)
		if err != nil {
			return err
		}

		data.Set("include_authn_statement", includeAuthnStatement)
	}

	if client.Attributes.SignDocuments != nil {
		signDocuments, err := strconv.ParseBool(*client.Attributes.SignDocuments)
		if err != nil {
			return err
		}

		data.Set("sign_documents", signDocuments)
	}

	if client.Attributes.SignAssertions != nil {
		signAssertions, err := strconv.ParseBool(*client.Attributes.SignAssertions)
		if err != nil {
			return err
		}

		data.Set("sign_assertions", signAssertions)
	}

	if client.Attributes.ClientSignatureRequired != nil {
		clientSignatureRequired, err := strconv.ParseBool(*client.Attributes.ClientSignatureRequired)
		if err != nil {
			return err
		}

		data.Set("client_signature_required", clientSignatureRequired)
	}

	if client.Attributes.ForcePostBinding != nil {
		forcePostBinding, err := strconv.ParseBool(*client.Attributes.ForcePostBinding)
		if err != nil {
			return err
		}

		data.Set("force_post_binding", forcePostBinding)
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

	data.SetId(client.Id)

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
