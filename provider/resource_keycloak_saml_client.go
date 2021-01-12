package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var (
	keycloakSamlClientNameIdFormats       = []string{"username", "email", "transient", "persistent"}
	keycloakSamlClientSignatureAlgorithms = []string{"RSA_SHA1", "RSA_SHA256", "RSA_SHA512", "DSA_SHA1"}
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
			"encrypt_assertions": {
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
				Computed: true,
			},
			"force_name_id_format": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"signature_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakSamlClientSignatureAlgorithms, false),
			},
			"name_id_format": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
			"encryption_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return old == formatCertificate(new)
				},
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return old == formatCertificate(new)
				},
			},
			"signing_private_key": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return old == formatSigningPrivateKey(new)
				},
			},
			"idp_initiated_sso_url_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_initiated_sso_relay_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assertion_consumer_post_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assertion_consumer_redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logout_service_post_binding_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logout_service_redirect_binding_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_scope_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"authentication_flow_binding_overrides": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"browser_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"direct_grant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func formatCertificate(signingCertificate string) string {
	r := strings.NewReplacer(
		"-----BEGIN CERTIFICATE-----", "",
		"-----END CERTIFICATE-----", "",
		"\n", "",
	)

	return r.Replace(signingCertificate)
}

func formatSigningPrivateKey(signingPrivateKey string) string {
	r := strings.NewReplacer(
		"-----BEGIN PRIVATE KEY-----", "",
		"-----END PRIVATE KEY-----", "",
		"\n", "",
	)

	return r.Replace(signingPrivateKey)
}

func mapToSamlClientFromData(data *schema.ResourceData) *keycloak.SamlClient {
	var validRedirectUris []string

	if v, ok := data.GetOk("valid_redirect_uris"); ok {
		for _, validRedirectUri := range v.(*schema.Set).List() {
			validRedirectUris = append(validRedirectUris, validRedirectUri.(string))
		}
	}

	samlAttributes := &keycloak.SamlClientAttributes{
		SignatureAlgorithm:              data.Get("signature_algorithm").(string),
		NameIdFormat:                    data.Get("name_id_format").(string),
		IDPInitiatedSSOURLName:          data.Get("idp_initiated_sso_url_name").(string),
		IDPInitiatedSSORelayState:       data.Get("idp_initiated_sso_relay_state").(string),
		AssertionConsumerPostURL:        data.Get("assertion_consumer_post_url").(string),
		AssertionConsumerRedirectURL:    data.Get("assertion_consumer_redirect_url").(string),
		LogoutServicePostBindingURL:     data.Get("logout_service_post_binding_url").(string),
		LogoutServiceRedirectBindingURL: data.Get("logout_service_redirect_binding_url").(string),
	}

	if encryptionCertificate, ok := data.GetOkExists("encryption_certificate"); ok {
		encryptionCertificateString := formatCertificate(encryptionCertificate.(string))
		samlAttributes.EncryptionCertificate = &encryptionCertificateString
	}

	if signingCertificate, ok := data.GetOkExists("signing_certificate"); ok {
		signingCertificateString := formatCertificate(signingCertificate.(string))
		samlAttributes.SigningCertificate = &signingCertificateString
	}

	if signingPrivateKey, ok := data.GetOkExists("signing_private_key"); ok {
		signingPrivateKeyString := formatSigningPrivateKey(signingPrivateKey.(string))
		samlAttributes.SigningPrivateKey = &signingPrivateKeyString
	}

	if includeAuthnStatement, ok := data.GetOkExists("include_authn_statement"); ok {
		includeAuthnStatementString := strconv.FormatBool(includeAuthnStatement.(bool))
		samlAttributes.IncludeAuthnStatement = &includeAuthnStatementString
	}

	if forceNameIdFormat, ok := data.GetOkExists("force_name_id_format"); ok {
		forceNameIdFormatString := strconv.FormatBool(forceNameIdFormat.(bool))
		samlAttributes.ForceNameIdFormat = &forceNameIdFormatString
	}

	if signDocuments, ok := data.GetOkExists("sign_documents"); ok {
		signDocumentsString := strconv.FormatBool(signDocuments.(bool))
		samlAttributes.SignDocuments = &signDocumentsString
	}

	if signAssertions, ok := data.GetOkExists("sign_assertions"); ok {
		signAssertionsString := strconv.FormatBool(signAssertions.(bool))
		samlAttributes.SignAssertions = &signAssertionsString
	}

	if encryptAssertions, ok := data.GetOkExists("encrypt_assertions"); ok {
		encryptAssertionsString := strconv.FormatBool(encryptAssertions.(bool))
		samlAttributes.EncryptAssertions = &encryptAssertionsString
	}

	if clientSignatureRequired, ok := data.GetOkExists("client_signature_required"); ok {
		clientSignatureRequiredString := strconv.FormatBool(clientSignatureRequired.(bool))
		samlAttributes.ClientSignatureRequired = &clientSignatureRequiredString
	}

	if forcePostBinding, ok := data.GetOkExists("force_post_binding"); ok {
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
		FullScopeAllowed:        data.Get("full_scope_allowed").(bool),
		Attributes:              samlAttributes,
	}

	if v, ok := data.GetOk("authentication_flow_binding_overrides"); ok {
		authenticationFlowBindingOverridesData := v.(*schema.Set).List()[0]
		authenticationFlowBindingOverrides := authenticationFlowBindingOverridesData.(map[string]interface{})
		samlClient.AuthenticationFlowBindingOverrides = keycloak.SamlAuthenticationFlowBindingOverrides{
			BrowserId:     authenticationFlowBindingOverrides["browser_id"].(string),
			DirectGrantId: authenticationFlowBindingOverrides["direct_grant_id"].(string),
		}
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
	if client.Attributes.ForceNameIdFormat != nil {
		forceNameIdFormat, err := strconv.ParseBool(*client.Attributes.ForceNameIdFormat)
		if err != nil {
			return err
		}

		data.Set("force_name_id_format", forceNameIdFormat)
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

	if client.Attributes.EncryptAssertions != nil {
		encryptAssertions, err := strconv.ParseBool(*client.Attributes.EncryptAssertions)
		if err != nil {
			return err
		}

		data.Set("encrypt_assertions", encryptAssertions)
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

	if _, exists := data.GetOkExists("encryption_certificate"); client.Attributes.EncryptionCertificate != nil && exists {
		data.Set("encryption_certificate", client.Attributes.EncryptionCertificate)
	}

	if _, exists := data.GetOkExists("signing_certificate"); client.Attributes.SigningCertificate != nil && exists {
		data.Set("signing_certificate", *client.Attributes.SigningCertificate)
	}

	if _, exists := data.GetOkExists("signing_certificate"); client.Attributes.SigningPrivateKey != nil && exists {
		data.Set("signing_private_key", *client.Attributes.SigningPrivateKey)
	}

	if (keycloak.SamlAuthenticationFlowBindingOverrides{}) == client.AuthenticationFlowBindingOverrides {
		data.Set("authentication_flow_binding_overrides", nil)
	} else {
		authenticationFlowBindingOverridesSettings := make(map[string]interface{})
		authenticationFlowBindingOverridesSettings["browser_id"] = client.AuthenticationFlowBindingOverrides.BrowserId
		authenticationFlowBindingOverridesSettings["direct_grant_id"] = client.AuthenticationFlowBindingOverrides.DirectGrantId
		data.Set("authentication_flow_binding_overrides", []interface{}{authenticationFlowBindingOverridesSettings})
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
	data.Set("signature_algorithm", client.Attributes.SignatureAlgorithm)
	data.Set("name_id_format", client.Attributes.NameIdFormat)
	data.Set("idp_initiated_sso_url_name", client.Attributes.IDPInitiatedSSOURLName)
	data.Set("idp_initiated_sso_relay_state", client.Attributes.IDPInitiatedSSORelayState)
	data.Set("assertion_consumer_post_url", client.Attributes.AssertionConsumerPostURL)
	data.Set("assertion_consumer_redirect_url", client.Attributes.AssertionConsumerRedirectURL)
	data.Set("logout_service_post_binding_url", client.Attributes.LogoutServicePostBindingURL)
	data.Set("logout_service_redirect_binding_url", client.Attributes.LogoutServiceRedirectBindingURL)
	data.Set("full_scope_allowed", client.FullScopeAllowed)

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

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{samlClientId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
