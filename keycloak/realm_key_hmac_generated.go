package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeyHmacGenerated struct {
	Id       string
	Name     string
	RealmId  string
	ParentId string

	Active     bool
	Enabled    bool
	Priority   int
	SecretSize int
	Algorithm  string
}

func convertFromRealmKeyHmacGeneratedToComponent(realmKey *RealmKeyHmacGenerated) *component {
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
		"algorithm": {
			realmKey.Algorithm,
		},
	}

	return &component{
		Id:           realmKey.Id,
		Name:         realmKey.Name,
		ParentId:     realmKey.ParentId,
		ProviderId:   "hmac-generated",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeyHmacGenerated(component *component, realmId string) (*RealmKeyHmacGenerated, error) {
	active, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("active"))
	if err != nil {
		return nil, err
	}

	enabled, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("enabled"))
	if err != nil {
		return nil, err
	}

	priority := 0 // Default priority
	if component.getConfig("priority") != "" {
		priority, err = strconv.Atoi(component.getConfig("priority"))
		if err != nil {
			return nil, err
		}
	}

	secretSize := 16 // Default key size for aes key
	if component.getConfig("secretSize") != "" {
		secretSize, err = strconv.Atoi(component.getConfig("secretSize"))
		if err != nil {
			return nil, err
		}
	}

	realmKey := &RealmKeyHmacGenerated{
		Id:       component.Id,
		Name:     component.Name,
		ParentId: component.ParentId,
		RealmId:  realmId,

		Active:     active,
		Enabled:    enabled,
		Priority:   priority,
		Algorithm:  component.getConfig("algorithm"),
		SecretSize: secretSize,
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeyHmacGenerated(realmKey *RealmKeyHmacGenerated) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeyHmacGeneratedToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeyHmacGenerated(realmId, id string) (*RealmKeyHmacGenerated, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeyHmacGenerated(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeyHmacGenerated(realmKey *RealmKeyHmacGenerated) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeyHmacGeneratedToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeyHmacGenerated(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
