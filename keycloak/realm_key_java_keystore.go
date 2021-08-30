package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeyJavaKeystore struct {
	Id       string
	Name     string
	RealmId  string
	ParentId string

	Active    bool
	Enabled   bool
	Priority  int
	Algorithm string

	Keystore         string
	KeystorePassword string
	KeyAlias         string
	KeyPassword      string
}

func convertFromRealmKeyJavaKeystoreToComponent(realmKey *RealmKeyJavaKeystore) *component {
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
		"keystore": {
			realmKey.Keystore,
		},
		"keystorePassword": {
			realmKey.KeystorePassword,
		},
		"keyAlias": {
			realmKey.KeyAlias,
		},
		"keyPassword": {
			realmKey.KeyPassword,
		},
	}

	return &component{
		Id:           realmKey.Id,
		Name:         realmKey.Name,
		ParentId:     realmKey.ParentId,
		ProviderId:   "java-keystore",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeyJavaKeystore(component *component, realmId string) (*RealmKeyJavaKeystore, error) {
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

	realmKey := &RealmKeyJavaKeystore{
		Id:       component.Id,
		Name:     component.Name,
		ParentId: component.ParentId,
		RealmId:  realmId,

		Active:           active,
		Enabled:          enabled,
		Priority:         priority,
		Algorithm:        component.getConfig("algorithm"),
		Keystore:         component.getConfig("keystore"),
		KeystorePassword: component.getConfig("keystorePassword"),
		KeyAlias:         component.getConfig("keyAlias"),
		KeyPassword:      component.getConfig("keyPassword"),
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeyJavaKeystore(realmKey *RealmKeyJavaKeystore) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeyJavaKeystoreToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeyJavaKeystore(realmId, id string) (*RealmKeyJavaKeystore, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeyJavaKeystore(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeyJavaKeystore(realmKey *RealmKeyJavaKeystore) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeyJavaKeystoreToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeyJavaKeystore(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
