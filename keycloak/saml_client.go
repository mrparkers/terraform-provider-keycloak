package keycloak

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SamlClientAttributes struct {
	IncludeAuthnStatement   *string `json:"saml.authnstatement"`
	SignDocuments           *string `json:"saml.server.signature"`
	SignAssertions          *string `json:"saml.assertion.signature"`
	ClientSignatureRequired *string `json:"saml.client.signature"`
	ForcePostBinding        *string `json:"saml.force.post.binding"`
	ForceNameIdFormat       *string `json:"saml_force_name_id_format"`
	// attributes above are actually booleans, but the Keycloak API expects strings
	NameIdFormat                    string  `json:"saml_name_id_format"`
	SigningCertificate              *string `json:"saml.signing.certificate,omitempty"`
	SigningPrivateKey               *string `json:"saml.signing.private.key"`
	IDPInitiatedSSOURLName          string  `json:"saml_idp_initiated_sso_url_name"`
	IDPInitiatedSSORelayState       string  `json:"saml_idp_initiated_sso_relay_state"`
	AssertionConsumerPostURL        string  `json:"saml_assertion_consumer_url_post"`
	AssertionConsumerRedirectURL    string  `json:"saml_assertion_consumer_url_redirect"`
	LogoutServicePostBindingURL     string  `json:"saml_single_logout_service_url_post"`
	LogoutServiceRedirectBindingURL string  `json:"saml_single_logout_service_url_redirect"`
	OtherAttributes                 map[string]interface{}
}

type SamlClient struct {
	Id                      string `json:"id,omitempty"`
	ClientId                string `json:"clientId"`
	RealmId                 string `json:"-"`
	Name                    string `json:"name"`
	Protocol                string `json:"protocol"`                // always saml for this resource
	ClientAuthenticatorType string `json:"clientAuthenticatorType"` // always client-secret

	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`

	FrontChannelLogout bool `json:"frontchannelLogout"`

	RootUrl                 string   `json:"rootUrl"`
	ValidRedirectUris       []string `json:"redirectUris"`
	BaseUrl                 string   `json:"baseUrl"`
	MasterSamlProcessingUrl string   `json:"adminUrl"`

	FullScopeAllowed bool `json:"fullScopeAllowed"`

	Attributes *SamlClientAttributes `json:"attributes"`
}

func (attr *SamlClientAttributes) MarshalJSON() ([]byte, error) {

	if attr.OtherAttributes == nil {
		attr.OtherAttributes = make(map[string]interface{})
	}

	attr.OtherAttributes["saml.authnstatement"] = attr.IncludeAuthnStatement
	attr.OtherAttributes["saml.server.signature"] = attr.SignDocuments
	attr.OtherAttributes["saml.assertion.signature"] = attr.SignAssertions
	attr.OtherAttributes["saml.client.signature"] = attr.ClientSignatureRequired
	attr.OtherAttributes["saml.force.post.binding"] = attr.ForcePostBinding
	attr.OtherAttributes["saml_force_name_id_format"] = attr.ForceNameIdFormat
	attr.OtherAttributes["saml_name_id_format"] = attr.NameIdFormat
	if attr.SigningCertificate != nil && *attr.SigningCertificate != "" {
		//omit empty
		attr.OtherAttributes["saml.signing.certificate"] = attr.SigningCertificate
	}
	attr.OtherAttributes["saml.signing.private.key"] = attr.SigningPrivateKey
	attr.OtherAttributes["saml_idp_initiated_sso_url_name"] = attr.IDPInitiatedSSOURLName
	attr.OtherAttributes["saml_idp_initiated_sso_relay_state"] = attr.IDPInitiatedSSORelayState
	attr.OtherAttributes["saml_assertion_consumer_url_post"] = attr.AssertionConsumerPostURL
	attr.OtherAttributes["saml_assertion_consumer_url_redirect"] = attr.AssertionConsumerRedirectURL
	attr.OtherAttributes["saml_single_logout_service_url_post"] = attr.LogoutServicePostBindingURL
	attr.OtherAttributes["saml_single_logout_service_url_redirect"] = attr.LogoutServiceRedirectBindingURL

	result, err := json.Marshal(attr.OtherAttributes)
	return result, err
}

func (attr *SamlClientAttributes) UnmarshalJSON(data []byte) error {
	var attrMap map[string]string
	if err := json.Unmarshal(data, &attrMap); err != nil {
		return err
	}

	if strings.Trim(attrMap["saml.authnstatement"], " ") != "" {
		includeAuthnStatement := attrMap["saml.authnstatement"]
		attr.IncludeAuthnStatement = &includeAuthnStatement
	}

	if strings.Trim(attrMap["saml.server.signature"], " ") != "" {
		signDocuments := attrMap["saml.server.signature"]
		attr.SignDocuments = &signDocuments
	}

	if strings.Trim(attrMap["saml.assertion.signature"], " ") != "" {
		signAssertions := attrMap["saml.assertion.signature"]
		attr.SignAssertions = &signAssertions
	}

	if strings.Trim(attrMap["saml.client.signature"], " ") != "" {
		clientSignatureRequired := attrMap["saml.client.signature"]
		attr.ClientSignatureRequired = &clientSignatureRequired
	}

	if strings.Trim(attrMap["saml.force.post.binding"], " ") != "" {
		forcePostBinding := attrMap["saml.force.post.binding"]
		attr.ForcePostBinding = &forcePostBinding
	}

	if strings.Trim(attrMap["saml_force_name_id_format"], " ") != "" {
		forceNameIDFormat := attrMap["saml_force_name_id_format"]
		attr.ForceNameIdFormat = &forceNameIDFormat
	}

	if strings.Trim(attrMap["saml.signing.certificate"], " ") != "" {
		signingCertificate := attrMap["saml.signing.certificate"]
		attr.SigningCertificate = &signingCertificate
	}

	if strings.Trim(attrMap["saml.signing.private.key"], " ") != "" {
		signingPrivateKey := attrMap["saml.signing.private.key"]
		attr.SigningPrivateKey = &signingPrivateKey
	}

	attr.NameIdFormat = attrMap["saml_name_id_format"]
	attr.IDPInitiatedSSOURLName = attrMap["saml_idp_initiated_sso_url_name"]
	attr.IDPInitiatedSSORelayState = attrMap["saml_idp_initiated_sso_relay_state"]
	attr.AssertionConsumerPostURL = attrMap["saml_assertion_consumer_url_post"]
	attr.AssertionConsumerRedirectURL = attrMap["saml_assertion_consumer_url_redirect"]
	attr.LogoutServicePostBindingURL = attrMap["saml_single_logout_service_url_post"]
	attr.LogoutServiceRedirectBindingURL = attrMap["saml_single_logout_service_url_redirect"]

	attr.OtherAttributes = make(map[string]interface{})

	reserverdKeys := map[string]bool{
		"saml.authnstatement":                     true,
		"saml.server.signature":                   true,
		"saml.assertion.signature":                true,
		"saml.client.signature":                   true,
		"saml.force.post.binding":                 true,
		"saml_force_name_id_format":               true,
		"saml_name_id_format":                     true,
		"saml.signing.certificate":                true,
		"saml.signing.private.key":                true,
		"saml_idp_initiated_sso_url_name":         true,
		"saml_idp_initiated_sso_relay_state":      true,
		"saml_assertion_consumer_url_post":        true,
		"saml_assertion_consumer_url_redirect":    true,
		"saml_single_logout_service_url_post":     true,
		"saml_single_logout_service_url_redirect": true,
	}

	for k, v := range attrMap {
		if found, _ := reserverdKeys[k]; !found {
			attr.OtherAttributes[k] = v
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) NewSamlClient(client *SamlClient) error {
	client.Protocol = "saml"
	client.ClientAuthenticatorType = "client-secret"

	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients", client.RealmId), client)
	if err != nil {
		return err
	}

	client.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetSamlClient(realmId, id string) (*SamlClient, error) {
	var client SamlClient

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), &client, nil)
	if err != nil {
		return nil, err
	}

	client.RealmId = realmId

	return &client, nil
}

func (keycloakClient *KeycloakClient) GetSamlClientInstallationProvider(realmId, id string, providerId string) ([]byte, error) {
	value, err := keycloakClient.getRaw(fmt.Sprintf("/realms/%s/clients/%s/installation/providers/%s", realmId, id, providerId), nil)
	return value, err
}

func (keycloakClient *KeycloakClient) UpdateSamlClient(client *SamlClient) error {
	client.Protocol = "saml"
	client.ClientAuthenticatorType = "client-secret"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s", client.RealmId, client.Id), client)
}

func (keycloakClient *KeycloakClient) DeleteSamlClient(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), nil)
}
