package keycloak

import (
	"encoding/json"
	"fmt"
)

type RealmUserProfilePermissions struct {
	Edit []string `json:"edit"`
	View []string `json:"view"`
}

type RealmUserProfileRequired struct {
	Roles  []string `json:"roles,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

type RealmUserProfileSelector struct {
	Scopes []string `json:"scopes,omitempty"`
}

type RealmUserProfileValidationConfig map[string]interface{}

type RealmUserProfileAttribute struct {
	Annotations map[string]string                           `json:"annotations,omitempty"`
	DisplayName string                                      `json:"displayName,omitempty"`
	Group       string                                      `json:"group,omitempty"`
	Name        string                                      `json:"name"`
	Permissions *RealmUserProfilePermissions                `json:"permissions,omitempty"`
	Required    *RealmUserProfileRequired                   `json:"required,omitempty"`
	Selector    *RealmUserProfileSelector                   `json:"selector,omitempty"`
	Validations map[string]RealmUserProfileValidationConfig `json:"validations,omitempty"`
}

type RealmUserProfileGroup struct {
	Annotations        map[string]string `json:"annotations,omitempty"`
	DisplayDescription string            `json:"displayDescription,omitempty"`
	DisplayHeader      string            `json:"displayHeader,omitempty"`
	Name               string            `json:"name"`
}

type RealmUserProfile struct {
	Attributes []*RealmUserProfileAttribute `json:"attributes"`
	Groups     []*RealmUserProfileGroup     `json:"groups,omitempty"`
}

func (keycloakClient *KeycloakClient) UpdateRealmUserProfile(realmId string, realmUserProfile *RealmUserProfile) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s/users/profile", realmId), realmUserProfile)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmUserProfile(realmId string) (*RealmUserProfile, error) {
	var realmUserProfile RealmUserProfile
	body, err := keycloakClient.getRaw(fmt.Sprintf("/realms/%s/users/profile", realmId), nil)
	if err != nil {
		return nil, err
	}

	if string(body) == "" {
		return nil, fmt.Errorf("User Profile is disabled for the %s realm", realmId)
	}

	err = json.Unmarshal(body, &realmUserProfile)
	if err != nil {
		return nil, err
	}

	return &realmUserProfile, nil
}
