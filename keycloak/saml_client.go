package keycloak

import (
	"fmt"
)

type SamlClientAttributes struct {
	IncludeAuthnStatement   *string `json:"saml.authnstatement"`
	SignDocuments           *string `json:"saml.server.signature"`
	SignAssertions          *string `json:"saml.assertion.signature"`
	EncryptAssertions       *string `json:"saml.encrypt"`
	ClientSignatureRequired *string `json:"saml.client.signature"`
	ForcePostBinding        *string `json:"saml.force.post.binding"`
	ForceNameIdFormat       *string `json:"saml_force_name_id_format"`
	// attributes above are actually booleans, but the Keycloak API expects strings
	SignatureAlgorithm              string  `json:"saml.signature.algorithm"`
	NameIdFormat                    string  `json:"saml_name_id_format"`
	SigningCertificate              *string `json:"saml.signing.certificate,omitempty"`
	SigningPrivateKey               *string `json:"saml.signing.private.key"`
	EncryptionCertificate           *string `json:"saml.encryption.certificate"`
	IDPInitiatedSSOURLName          string  `json:"saml_idp_initiated_sso_url_name"`
	IDPInitiatedSSORelayState       string  `json:"saml_idp_initiated_sso_relay_state"`
	AssertionConsumerPostURL        string  `json:"saml_assertion_consumer_url_post"`
	AssertionConsumerRedirectURL    string  `json:"saml_assertion_consumer_url_redirect"`
	LogoutServicePostBindingURL     string  `json:"saml_single_logout_service_url_post"`
	LogoutServiceRedirectBindingURL string  `json:"saml_single_logout_service_url_redirect"`
}

type SamlAuthenticationFlowBindingOverrides struct {
	BrowserId     string `json:"browser"`
	DirectGrantId string `json:"direct_grant"`
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

	AuthenticationFlowBindingOverrides SamlAuthenticationFlowBindingOverrides `json:"authenticationFlowBindingOverrides,omitempty"`
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

func (keycloakClient *KeycloakClient) GetSamlClientByClientId(realmId, clientId string) (*SamlClient, error) {
	var clients []SamlClient

	params := map[string]string{
		"clientId": clientId,
	}

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients", realmId), &clients, params)
	if err != nil {
		return nil, err
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("saml client with name %s does not exist", clientId)
	}

	client := clients[0]

	client.RealmId = realmId

	return &client, nil
}

func (keycloakClient *KeycloakClient) UpdateSamlClient(client *SamlClient) error {
	client.Protocol = "saml"
	client.ClientAuthenticatorType = "client-secret"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s", client.RealmId, client.Id), client)
}

func (keycloakClient *KeycloakClient) DeleteSamlClient(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), nil)
}

func (keycloakClient *KeycloakClient) getSamlClientScopes(realmId, clientId, t string) ([]*SamlClientScope, error) {
	var scopes []*SamlClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/%s-client-scopes", realmId, clientId, t), &scopes, nil)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (keycloakClient *KeycloakClient) GetSamlClientDefaultScopes(realmId, clientId string) ([]*SamlClientScope, error) {
	return keycloakClient.getSamlClientScopes(realmId, clientId, "default")
}

func (keycloakClient *KeycloakClient) attachSamlClientScopes(realmId, clientId, t string, scopeNames []string) error {
	_, err := keycloakClient.GetSamlClient(realmId, clientId)
	if err != nil && ErrorIs404(err) {
		return fmt.Errorf("validation error: client with id %s does not exist", clientId)
	} else if err != nil {
		return err
	}

	allSamlClientScopes, err := keycloakClient.ListSamlClientScopesWithFilter(realmId, includeSamlClientScopesMatchingNames(scopeNames))
	if err != nil {
		return err
	}

	for _, samlClientScope := range allSamlClientScopes {
		err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/%s-client-scopes/%s", realmId, clientId, t, samlClientScope.Id), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) AttachSamlClientDefaultScopes(realmId, clientId string, scopeNames []string) error {
	return keycloakClient.attachSamlClientScopes(realmId, clientId, "default", scopeNames)
}

func (keycloakClient *KeycloakClient) detachSamlClientScopes(realmId, clientId, t string, scopeNames []string) error {
	allSamlClientScopes, err := keycloakClient.ListSamlClientScopesWithFilter(realmId, includeSamlClientScopesMatchingNames(scopeNames))
	if err != nil {
		return err
	}

	for _, samlClientScope := range allSamlClientScopes {
		err := keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/%s-client-scopes/%s", realmId, clientId, t, samlClientScope.Id), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) DetachSamlClientDefaultScopes(realmId, clientId string, scopeNames []string) error {
	return keycloakClient.detachSamlClientScopes(realmId, clientId, "default", scopeNames)
}
