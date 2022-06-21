package keycloak

import (
	"context"
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

type RealmUserProfileValidationLength struct {
	Min          int  `json:"min,omitempty"`
	Max          int  `json:"max,omitempty"`
	TrimDisabled bool `json:"trim-disabled,omitempty"`
}

type RealmUserProfileValidationInteger struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type RealmUserProfileValidationDouble struct {
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
}

type RealmUserProfileValidationPattern struct {
	Pattern      string `json:"pattern,omitempty"`
	ErrorMessage string `json:"error-message,omitempty"`
}

type RealmUserProfileValidationProhibited struct {
	ErrorMessage string `json:"error-message,omitempty"`
}

type RealmUserProfileValidationOptions struct {
	Options []string `json:"options,omitempty"`
}

type RealmUserProfileValidationConfig struct {
	Length                    *RealmUserProfileValidationLength     `json:"length,omitempty"`
	Integer                   *RealmUserProfileValidationInteger    `json:"integer,omitempty"`
	Double                    *RealmUserProfileValidationDouble     `json:"double,omitempty"`
	URI                       *map[string]interface{}               `json:"uri,omitempty"`
	Pattern                   *RealmUserProfileValidationPattern    `json:"pattern,omitempty"`
	Email                     *map[string]interface{}               `json:"email,omitempty"`
	LocalDate                 *map[string]interface{}               `json:"local-date,omitempty"`
	PersonNameProhibitedChars *RealmUserProfileValidationProhibited `json:"person-name-prohibited-characters,omitempty"`
	UsernameProhibitedChars   *RealmUserProfileValidationProhibited `json:"username-prohibited-characters,omitempty"`
	Options                   *RealmUserProfileValidationOptions    `json:"options,omitempty"`
}

type RealmUserProfileAttribute struct {
	Annotations map[string]string                 `json:"annotations,omitempty"`
	DisplayName string                            `json:"displayName,omitempty"`
	Group       string                            `json:"group,omitempty"`
	Name        string                            `json:"name"`
	Permissions *RealmUserProfilePermissions      `json:"permissions,omitempty"`
	Required    *RealmUserProfileRequired         `json:"required,omitempty"`
	Selector    *RealmUserProfileSelector         `json:"selector,omitempty"`
	Validations *RealmUserProfileValidationConfig `json:"validations,omitempty"`
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

func (keycloakClient *KeycloakClient) UpdateRealmUserProfile(ctx context.Context, realmId string, realmUserProfile *RealmUserProfile) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/users/profile", realmId), realmUserProfile)
}

func (keycloakClient *KeycloakClient) GetRealmUserProfile(ctx context.Context, realmId string) (*RealmUserProfile, error) {
	var realmUserProfile RealmUserProfile
	body, err := keycloakClient.getRaw(ctx, fmt.Sprintf("/realms/%s/users/profile", realmId), nil)
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
