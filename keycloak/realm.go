package keycloak

import (
	"fmt"
)

type Realm struct {
	Id          string `json:"id"`
	Realm       string `json:"realm"`
	Enabled     bool   `json:"enabled"`
	DisplayName string `json:"displayName"`

	// Login
	RegistrationAllowed bool `json:"registrationAllowed"`
	EmailAsUsername     bool `json:"registrationEmailAsUsername"`
	EditUsername        bool `json:"editUsernameAllowed"`
	ForgotPassword      bool `json:"resetPasswordAllowed"`
	RememberMe          bool `json:"rememberMe"`
	VerifyEmail         bool `json:"verifyEmail"`
	LoginWithEmail      bool `json:"loginWithEmailAllowed"`
}

func (keycloakClient *KeycloakClient) NewRealm(realm *Realm) error {
	return keycloakClient.post("/realms/", realm)
}

func (keycloakClient *KeycloakClient) GetRealm(id string) (*Realm, error) {
	var realm Realm

	err := keycloakClient.get(fmt.Sprintf("/realms/%s", id), &realm)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

func (keycloakClient *KeycloakClient) UpdateRealm(realm *Realm) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s", realm.Id), realm)
}

func (keycloakClient *KeycloakClient) DeleteRealm(id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s", id))
}
