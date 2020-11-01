package keycloak

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
	NameIdFormat                    string                 `json:"saml_name_id_format"`
	SigningCertificate              *string                `json:"saml.signing.certificate,omitempty"`
	SigningPrivateKey               *string                `json:"saml.signing.private.key"`
	IDPInitiatedSSOURLName          string                 `json:"saml_idp_initiated_sso_url_name"`
	IDPInitiatedSSORelayState       string                 `json:"saml_idp_initiated_sso_relay_state"`
	AssertionConsumerPostURL        string                 `json:"saml_assertion_consumer_url_post"`
	AssertionConsumerRedirectURL    string                 `json:"saml_assertion_consumer_url_redirect"`
	LogoutServicePostBindingURL     string                 `json:"saml_single_logout_service_url_post"`
	LogoutServiceRedirectBindingURL string                 `json:"saml_single_logout_service_url_redirect"`
	OtherAttributes                 map[string]interface{} `json:"-"`
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

func (f *SamlClientAttributes) UnmarshalJSON(data []byte) error {
	f.OtherAttributes = map[string]interface{}{}
	err := json.Unmarshal(data, &f.OtherAttributes)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(f).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		jsonKey := strings.Split(structField.Tag.Get("json"), ",")[0]
		if jsonKey != "-" {
			value, ok := f.OtherAttributes[jsonKey]
			if ok {
				field := v.FieldByName(structField.Name)
				if field.IsValid() && field.CanSet() {
					if field.Kind() == reflect.Ptr {
						temp := value.(string)
						if temp != "" {
							field.Set(reflect.ValueOf(&temp))
						}
					} else if field.Kind() == reflect.String {
						field.SetString(value.(string))
					} else if field.Kind() == reflect.Bool {
						boolVal, err := strconv.ParseBool(value.(string))
						if err == nil {
							field.Set(reflect.ValueOf(KeycloakBoolQuoted(boolVal)))
						}
					}
					delete(f.OtherAttributes, jsonKey)
				}
			}
		}
	}
	return nil
}

func (f *SamlClientAttributes) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{}

	for k, v := range f.OtherAttributes {
		out[k] = v
	}
	v := reflect.ValueOf(f).Elem()
	for i := 0; i < v.NumField(); i++ {
		jsonKey := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0]
		if jsonKey != "-" {
			field := v.Field(i)
			if field.IsValid() && field.CanSet() {
				if field.Kind() == reflect.Ptr {
					out[jsonKey] = field.String()
				} else if field.Kind() == reflect.String {
					out[jsonKey] = field.String()
				} else if field.Kind() == reflect.Bool {
					out[jsonKey] = KeycloakBoolQuoted(field.Bool())
				}
			}
		}
	}
	return json.Marshal(out)
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
