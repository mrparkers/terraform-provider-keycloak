package keycloak

import (
	"fmt"
)

type OpenidClientScope struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Protocol    string `json:"protocol"`
	Attributes  struct {
		DisplayOnConsentScreen string `json:"display.on.consent.screen"` // boolean in string form
		ConsentScreenText      string `json:"consent.screen.text"`
	} `json:"attributes"`
}

type OpenidClientScopeFilterFunc func(OpenidClientScope) (*OpenidClientScope, bool)

func (keycloakClient *KeycloakClient) NewOpenidClientScope(clientScope *OpenidClientScope) error {
	clientScope.Protocol = "openid-connect"

	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/client-scopes", clientScope.RealmId), clientScope)
	if err != nil {
		return err
	}

	clientScope.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientScope(realmId, id string) (*OpenidClientScope, error) {
	var clientScope OpenidClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/client-scopes/%s", realmId, id), &clientScope)
	if err != nil {
		return nil, err
	}

	clientScope.RealmId = realmId

	return &clientScope, nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientScope(clientScope *OpenidClientScope) error {
	clientScope.Protocol = "openid-connect"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/client-scopes/%s", clientScope.RealmId, clientScope.Id), clientScope)
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientScope(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/client-scopes/%s", realmId, id))
}

func (keycloakClient *KeycloakClient) listOpenidClientScopesWithFilter(realmId string, filter OpenidClientScopeFilterFunc) ([]*OpenidClientScope, error) {
	var clientScopes []OpenidClientScope
	var openidClientScopes []*OpenidClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/client-scopes", realmId), &clientScopes)
	if err != nil {
		return nil, err
	}

	for _, clientScope := range clientScopes {
		if filteredClientScope, ok := filter(clientScope); ok && clientScope.Protocol == "openid-connect" {
			openidClientScopes = append(openidClientScopes, filteredClientScope)
		}
	}

	return openidClientScopes, nil
}

func includeOpenidClientScopesMatchingNames(scopeNames []string) OpenidClientScopeFilterFunc {
	return func(scope OpenidClientScope) (*OpenidClientScope, bool) {
		for _, scopeName := range scopeNames {
			if scopeName == scope.Name {
				return &scope, true
			}
		}

		return &scope, false
	}
}
