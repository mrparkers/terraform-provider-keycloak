package keycloak

import "fmt"

type Realm struct {
	Id    string `json:"id"`
	Realm string `json:"realm"`
}

func (keycloakClient *KeycloakClient) NewRealm(realm *Realm) error {
	err := keycloakClient.post("/realms/", realm)

	return err
}

func (keycloakClient *KeycloakClient) GetRealm(id string) (*Realm, error) {
	var realm Realm

	url := fmt.Sprintf("/realms/%s", id)

	err := keycloakClient.get(url, &realm)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

func (keycloakClient *KeycloakClient) DeleteRealm(id string) error {
	url := fmt.Sprintf("/realms/%s", id)

	err := keycloakClient.delete(url)

	return err
}
