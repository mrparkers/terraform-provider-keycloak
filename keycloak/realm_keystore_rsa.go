package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeystoreRsa struct {
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

func convertFromRealmKeystoreRsaToComponent(realmKey *RealmKeystoreRsa) *component {
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

func convertFromComponentToRealmKeystoreRsa(component *component, realmId string) (*RealmKeystoreRsa, error) {
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

	realmKey := &RealmKeystoreRsa{
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

func (keycloakClient *KeycloakClient) NewRealmKeystoreRsa(realmKey *RealmKeystoreRsa) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeystoreRsaToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeystoreRsa(realmId, id string) (*RealmKeystoreRsa, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeystoreRsa(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeystoreRsa(realmKey *RealmKeystoreRsa) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeystoreRsaToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeystoreRsa(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
