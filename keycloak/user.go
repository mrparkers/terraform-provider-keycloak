package keycloak

import "fmt"

type User struct {
	Id      string `json:"id,omitempty"`
	RealmId string `json:"-"`

	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Enabled   bool   `json:"enabled"`
}

func (keycloakClient *KeycloakClient) NewUser(user *User) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/users", user.RealmId), user)
	if err != nil {
		return err
	}

	user.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetUser(realmId, id string) (*User, error) {
	var user User

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users/%s", realmId, id), &user)
	if err != nil {
		return nil, err
	}

	user.RealmId = realmId

	return &user, nil
}

func (keycloakClient *KeycloakClient) UpdateUser(user *User) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/users/%s", user.RealmId, user.Id), user)
}

func (keycloakClient *KeycloakClient) DeleteUser(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s", realmId, id))
}
