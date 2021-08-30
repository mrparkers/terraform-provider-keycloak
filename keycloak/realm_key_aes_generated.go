package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeyAesGenerated struct {
	Id       string
	Name     string
	RealmId  string
	ParentId string

	Active     bool
	Enabled    bool
	Priority   int
	SecretSize int
}

func convertFromRealmKeyAesGeneratedToComponent(realmKey *RealmKeyAesGenerated) *component {
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
		ParentId:     realmKey.ParentId,
		ProviderId:   "aes-generated",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeyAesGenerated(component *component, realmId string) (*RealmKeyAesGenerated, error) {
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

	realmKey := &RealmKeyAesGenerated{
		Id:       component.Id,
		Name:     component.Name,
		ParentId: component.ParentId,
		RealmId:  realmId,

		Active:     active,
		Enabled:    enabled,
		Priority:   priority,
		SecretSize: secretSize,
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeyAesGenerated(realmKey *RealmKeyAesGenerated) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeyAesGeneratedToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeyAesGenerated(realmId, id string) (*RealmKeyAesGenerated, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeyAesGenerated(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeyAesGenerated(realmKey *RealmKeyAesGenerated) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeyAesGeneratedToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeyAesGenerated(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
