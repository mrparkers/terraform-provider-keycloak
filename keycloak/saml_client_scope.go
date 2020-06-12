package keycloak

import (
	"fmt"
)

type SamlClientScope struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Protocol    string `json:"protocol"`
	Attributes  struct {
		DisplayOnConsentScreen KeycloakBoolQuoted `json:"display.on.consent.screen"` // boolean in string form
		ConsentScreenText      string             `json:"consent.screen.text"`
		GuiOrder               string             `json:"gui.order"`
	} `json:"attributes"`
}

type SamlClientScopeFilterFunc func(*SamlClientScope) bool

func (keycloakClient *KeycloakClient) NewSamlClientScope(clientScope *SamlClientScope) error {
	clientScope.Protocol = "saml"

	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/client-scopes", clientScope.RealmId), clientScope)
	if err != nil {
		return err
	}

	clientScope.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetSamlClientScope(realmId, id string) (*SamlClientScope, error) {
	var clientScope SamlClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/client-scopes/%s", realmId, id), &clientScope, nil)
	if err != nil {
		return nil, err
	}

	clientScope.RealmId = realmId

	return &clientScope, nil
}

func (keycloakClient *KeycloakClient) GetSamlDefaultClientScopes(realmId, clientId string) (*[]SamlClientScope, error) {
	var clientScopes []SamlClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/default-client-scopes", realmId, clientId), &clientScopes, nil)
	if err != nil {
		return nil, err
	}

	for _, clientScope := range clientScopes {
		clientScope.RealmId = realmId
	}

	return &clientScopes, nil
}

func (keycloakClient *KeycloakClient) UpdateSamlClientScope(clientScope *SamlClientScope) error {
	clientScope.Protocol = "saml"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/client-scopes/%s", clientScope.RealmId, clientScope.Id), clientScope)
}

func (keycloakClient *KeycloakClient) DeleteSamlClientScope(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/client-scopes/%s", realmId, id), nil)
}

func (keycloakClient *KeycloakClient) ListSamlClientScopesWithFilter(realmId string, filter SamlClientScopeFilterFunc) ([]*SamlClientScope, error) {
	var clientScopes []SamlClientScope
	var samlClientScopes []*SamlClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/client-scopes", realmId), &clientScopes, nil)
	if err != nil {
		return nil, err
	}

	for _, clientScope := range clientScopes {
		if clientScope.Protocol == "saml" && filter(&clientScope) {
			scope := new(SamlClientScope)
			*scope = clientScope

			samlClientScopes = append(samlClientScopes, scope)
		}
	}

	return samlClientScopes, nil
}

func includeSamlClientScopesMatchingNames(scopeNames []string) SamlClientScopeFilterFunc {
	return func(scope *SamlClientScope) bool {
		for _, scopeName := range scopeNames {
			if scopeName == scope.Name {
				return true
			}
		}

		return false
	}
}
