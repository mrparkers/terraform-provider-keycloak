package keycloak

import "fmt"

type AuthenticationFlow struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	ProviderId  string `json:"providerId"` // always "basic-flow"
	TopLevel    bool   `json:"topLevel"`   // should only be false if this is a subflow
	BuiltIn     bool   `json:"builtIn"`    // this controls whether or not this flow can be edited from the console. it can be updated, but this provider will only set it to `true`
}

func (keycloakClient *KeycloakClient) NewAuthenticationFlow(authenticationFlow *AuthenticationFlow) error {
	authenticationFlow.BuiltIn = false
	authenticationFlow.TopLevel = true
	authenticationFlow.ProviderId = "basic-flow"

	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/authentication/flows", authenticationFlow.RealmId), authenticationFlow)
	if err != nil {
		return err
	}

	authenticationFlow.Id = getIdFromLocationHeader(location)

	return err
}

func (keycloakClient *KeycloakClient) GetAuthenticationFlow(realmId, id string) (*AuthenticationFlow, error) {
	var authenticationFlow AuthenticationFlow

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), &authenticationFlow)
	if err != nil {
		return nil, err
	}

	authenticationFlow.RealmId = realmId

	return &authenticationFlow, nil
}

func (keycloakClient *KeycloakClient) UpdateAuthenticationFlow(authenticationFlow *AuthenticationFlow) error {
	authenticationFlow.BuiltIn = false
	authenticationFlow.TopLevel = true
	authenticationFlow.ProviderId = "basic-flow"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/authentication/flows/%s", authenticationFlow.RealmId, authenticationFlow.Id), authenticationFlow)
}

func (keycloakClient *KeycloakClient) DeleteAuthenticationFlow(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id))
}
