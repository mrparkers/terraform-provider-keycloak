package keycloak

import (
	"context"
	"fmt"
	"strconv"
)

type RealmKeystoreAesGenerated struct {
	Id      string
	Name    string
	RealmId string

	Active     bool
	Enabled    bool
	Priority   int
	SecretSize int
}

func convertFromRealmKeystoreAesGeneratedToComponent(realmKey *RealmKeystoreAesGenerated) *component {
	componentConfig := map[string][]string{
		"active": {
			strconv.FormatBool(realmKey.Active),
		},
		"enabled": {
			strconv.FormatBool(realmKey.Enabled),
		},
		"priority": {
			strconv.Itoa(realmKey.Priority),
		},
		"secretSize": {
			strconv.Itoa(realmKey.SecretSize),
		},
	}

	return &component{
		Id:           realmKey.Id,
		Name:         realmKey.Name,
		ParentId:     realmKey.RealmId,
		ProviderId:   "aes-generated",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeystoreAesGenerated(component *component, realmId string) (*RealmKeystoreAesGenerated, error) {
	active, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfigFirstOrDefault("active"))
	if err != nil {
		return nil, err
	}

	enabled, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfigFirstOrDefault("enabled"))
	if err != nil {
		return nil, err
	}

	priority := 0 // Default priority
	if component.getConfigFirstOrDefault("priority") != "" {
		priority, err = strconv.Atoi(component.getConfigFirstOrDefault("priority"))
		if err != nil {
			return nil, err
		}
	}

	secretSize := 16 // Default key size for aes key
	if component.getConfigFirstOrDefault("secretSize") != "" {
		secretSize, err = strconv.Atoi(component.getConfigFirstOrDefault("secretSize"))
		if err != nil {
			return nil, err
		}
	}

	realmKey := &RealmKeystoreAesGenerated{
		Id:      component.Id,
		Name:    component.Name,
		RealmId: realmId,

		Active:     active,
		Enabled:    enabled,
		Priority:   priority,
		SecretSize: secretSize,
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeystoreAesGenerated(ctx context.Context, realmKey *RealmKeystoreAesGenerated) error {
	_, location, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeystoreAesGeneratedToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeystoreAesGenerated(ctx context.Context, realmId, id string) (*RealmKeystoreAesGenerated, error) {
	var component *component

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeystoreAesGenerated(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeystoreAesGenerated(ctx context.Context, realmKey *RealmKeystoreAesGenerated) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeystoreAesGeneratedToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeystoreAesGenerated(ctx context.Context, realmId, id string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
