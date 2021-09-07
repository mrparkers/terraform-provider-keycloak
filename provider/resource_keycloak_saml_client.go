package provider

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

var (
	keycloakSamlClientNameIdFormats           = []string{"username", "email", "transient", "persistent"}
	keycloakSamlClientSignatureAlgorithms     = []string{"RSA_SHA1", "RSA_SHA256", "RSA_SHA512", "DSA_SHA1"}
	keycloakSamlClientSignatureKeyNames       = []string{"NONE", "KEY_ID", "CERT_SUBJECT"}
	keycloakSamlClientCanonicalizationMethods = map[string]string{
		"EXCLUSIVE":               "http://www.w3.org/2001/10/xml-exc-c14n#",
		"EXCLUSIVE_WITH_COMMENTS": "http://www.w3.org/2001/10/xml-exc-c14n#WithComments",
		"INCLUSIVE":               "http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
		"INCLUSIVE_WITH_COMMENTS": "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments",
	}
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
			"signature_key_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "KEY_ID",
				ValidateFunc: validation.StringInSlice(keycloakSamlClientSignatureKeyNames, false),
			},
			"canonicalization_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "EXCLUSIVE",
				ValidateFunc: validation.StringInSlice(keys(keycloakSamlClientCanonicalizationMethods), false),
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
			"extra_config": {
				Type:             schema.TypeMap,
				Optional:         true,
				ValidateDiagFunc: validateExtraConfig(reflect.ValueOf(&keycloak.SamlClientAttributes{}).Elem()),
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
		IncludeAuthnStatement:           keycloak.KeycloakBoolQuoted(data.Get("include_authn_statement").(bool)),
		ForceNameIdFormat:               keycloak.KeycloakBoolQuoted(data.Get("force_name_id_format").(bool)),
		SignDocuments:                   keycloak.KeycloakBoolQuoted(data.Get("sign_documents").(bool)),
		SignAssertions:                  keycloak.KeycloakBoolQuoted(data.Get("sign_assertions").(bool)),
		EncryptAssertions:               keycloak.KeycloakBoolQuoted(data.Get("encrypt_assertions").(bool)),
		ClientSignatureRequired:         keycloak.KeycloakBoolQuoted(data.Get("client_signature_required").(bool)),
		ForcePostBinding:                keycloak.KeycloakBoolQuoted(data.Get("force_post_binding").(bool)),
		SignatureAlgorithm:              data.Get("signature_algorithm").(string),
		SignatureKeyName:                data.Get("signature_key_name").(string),
		CanonicalizationMethod:          keycloakSamlClientCanonicalizationMethods[data.Get("canonicalization_method").(string)],
		NameIdFormat:                    data.Get("name_id_format").(string),
		IDPInitiatedSSOURLName:          data.Get("idp_initiated_sso_url_name").(string),
		IDPInitiatedSSORelayState:       data.Get("idp_initiated_sso_relay_state").(string),
		AssertionConsumerPostURL:        data.Get("assertion_consumer_post_url").(string),
		AssertionConsumerRedirectURL:    data.Get("assertion_consumer_redirect_url").(string),
		LogoutServicePostBindingURL:     data.Get("logout_service_post_binding_url").(string),
		LogoutServiceRedirectBindingURL: data.Get("logout_service_redirect_binding_url").(string),
		ExtraConfig:                     getExtraConfigFromData(data),
	}

	if encryptionCertificate, ok := data.GetOkExists("encryption_certificate"); ok {
		samlAttributes.EncryptionCertificate = formatCertificate(encryptionCertificate.(string))
	}

	if signingCertificate, ok := data.GetOkExists("signing_certificate"); ok {
		samlAttributes.SigningCertificate = formatCertificate(signingCertificate.(string))
	}

	if signingPrivateKey, ok := data.GetOkExists("signing_private_key"); ok {
		samlAttributes.SigningPrivateKey = formatSigningPrivateKey(signingPrivateKey.(string))
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

	data.Set("include_authn_statement", client.Attributes.IncludeAuthnStatement)
	data.Set("force_name_id_format", client.Attributes.ForceNameIdFormat)
	data.Set("sign_documents", client.Attributes.SignDocuments)
	data.Set("sign_assertions", client.Attributes.SignAssertions)
	data.Set("encrypt_assertions", client.Attributes.EncryptAssertions)
	data.Set("client_signature_required", client.Attributes.ClientSignatureRequired)
	data.Set("force_post_binding", client.Attributes.ForcePostBinding)

	if _, exists := data.GetOkExists("encryption_certificate"); exists {
		data.Set("encryption_certificate", client.Attributes.EncryptionCertificate)
	}

	if _, exists := data.GetOkExists("signing_certificate"); exists {
		data.Set("signing_certificate", client.Attributes.SigningCertificate)
	}

	if _, exists := data.GetOkExists("signing_certificate"); exists {
		data.Set("signing_private_key", client.Attributes.SigningPrivateKey)
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
	data.Set("signature_key_name", client.Attributes.SignatureKeyName)
	data.Set("name_id_format", client.Attributes.NameIdFormat)
	data.Set("idp_initiated_sso_url_name", client.Attributes.IDPInitiatedSSOURLName)
	data.Set("idp_initiated_sso_relay_state", client.Attributes.IDPInitiatedSSORelayState)
	data.Set("assertion_consumer_post_url", client.Attributes.AssertionConsumerPostURL)
	data.Set("assertion_consumer_redirect_url", client.Attributes.AssertionConsumerRedirectURL)
	data.Set("logout_service_post_binding_url", client.Attributes.LogoutServicePostBindingURL)
	data.Set("logout_service_redirect_binding_url", client.Attributes.LogoutServiceRedirectBindingURL)
	data.Set("full_scope_allowed", client.FullScopeAllowed)

	if canonicalizationMethod, ok := mapKeyFromValue(keycloakSamlClientCanonicalizationMethods, client.Attributes.CanonicalizationMethod); ok {
		data.Set("canonicalization_method", canonicalizationMethod)
	}

	setExtraConfigData(data, client.Attributes.ExtraConfig)

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
