package keycloak

import (
	"fmt"
	"net/url"
)

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

func (keycloakClient *KeycloakClient) getUserByUsername(realmId, username string) (*User, error) {
	var users []*User

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users?username=%s", realmId, url.QueryEscape(username)), &users)
	if err != nil {
		return nil, err
	}

	// more than one user could be returned so we need to search through all of the results and return the correct one
	// ex: foo and foo-user could both exist, but searching for "foo" will return both
	for _, user := range users {
		if user.Username == username {
			user.RealmId = realmId

			return user, nil
		}
	}

	// the requested user does not exist
	// I don't think raising an error is appropriate here - consumers should check if the user is nil
	return nil, nil
}

func (keycloakClient *KeycloakClient) addUserToGroup(user *User, groupId string) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/users/%s/groups/%s", user.RealmId, user.Id, groupId), nil)
}

func (keycloakClient *KeycloakClient) AddUsersToGroup(realmId, groupId string, users []interface{}) error {
	for _, username := range users {
		user, err := keycloakClient.getUserByUsername(realmId, username.(string)) // we need the user's id in order to add them to a group
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user with username %s does not exist", username.(string))
		}

		err = keycloakClient.addUserToGroup(user, groupId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) RemoveUserFromGroup(user *User, groupId string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s/groups/%s", user.RealmId, user.Id, groupId))
}

func (keycloakClient *KeycloakClient) RemoveUsersFromGroup(realmId, groupId string, usernames []interface{}) error {
	for _, username := range usernames {
		user, err := keycloakClient.getUserByUsername(realmId, username.(string)) // we need the user's id in order to remove them from a group
		if err != nil {
			return err
		}

		err = keycloakClient.RemoveUserFromGroup(user, groupId)
		if err != nil {
			return err
		}
	}

	return nil
}
