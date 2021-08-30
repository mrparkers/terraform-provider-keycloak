package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeyRsa struct {
	Id       string
	Name     string
	RealmId  string
	ParentId string

	Active    bool
	Enabled   bool
	Priority  int
	KeySize   int
	Algorithm string

	PrivateKey  string
	Certificate string
}

func convertFromRealmKeyRsaToComponent(realmKey *RealmKeyRsa) *component {
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
		"algorithm": {
			realmKey.Algorithm,
		},
		"keySize": {
			strconv.Itoa(realmKey.KeySize),
		},
		"private_key": {
			realmKey.PrivateKey,
		},
		"certificate": {
			realmKey.Certificate,
		},
	}

	return &component{
		Id:           realmKey.Id,
		Name:         realmKey.Name,
		ParentId:     realmKey.ParentId,
		ProviderId:   "rsa",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeyRsa(component *component, realmId string) (*RealmKeyRsa, error) {
	active, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("active"))
	if err != nil {
		return nil, err
	}

	enabled, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("enabled"))
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(component.getConfig("priority"))
	if err != nil {
		return nil, err
	}

	keySize, err := strconv.Atoi(component.getConfig("keySize"))
	if err != nil {
		return nil, err
	}

	realmKey := &RealmKeyRsa{
		Id:       component.Id,
		Name:     component.Name,
		RealmId:  realmId,
		ParentId: component.ParentId,

		Active:      active,
		Enabled:     enabled,
		Priority:    priority,
		Algorithm:   component.getConfig("algorithm"),
		KeySize:     keySize,
		PrivateKey:  component.getConfig("privateKey"),
		Certificate: component.getConfig("certificate"),
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeyRsa(realmKey *RealmKeyRsa) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeyRsaToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeyRsa(realmId, id string) (*RealmKeyRsa, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeyRsa(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeyRsa(realmKey *RealmKeyRsa) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeyRsaToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeyRsa(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
