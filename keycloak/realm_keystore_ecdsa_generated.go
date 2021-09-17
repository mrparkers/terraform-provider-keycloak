package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeystoreEcdsaGenerated struct {
	Id      string
	Name    string
	RealmId string

	Active        bool
	Enabled       bool
	Priority      int
	EllipticCurve string
}

func convertFromRealmKeystoreEcdsaGeneratedToComponent(realmKey *RealmKeystoreEcdsaGenerated) *component {
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
		"ecdsaEllipticCurveKey": {
			realmKey.EllipticCurve,
		},
	}

	return &component{
		Id:           realmKey.Id,
		Name:         realmKey.Name,
		ParentId:     realmKey.RealmId,
		ProviderId:   "ecdsa-generated",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeystoreEcdsaGenerated(component *component, realmId string) (*RealmKeystoreEcdsaGenerated, error) {
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

	realmKey := &RealmKeystoreEcdsaGenerated{
		Id:      component.Id,
		Name:    component.Name,
		RealmId: realmId,

		Active:        active,
		Enabled:       enabled,
		Priority:      priority,
		EllipticCurve: component.getConfig("ecdsaEllipticCurveKey"),
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeystoreEcdsaGenerated(realmKey *RealmKeystoreEcdsaGenerated) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeystoreEcdsaGeneratedToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeystoreEcdsaGenerated(realmId, id string) (*RealmKeystoreEcdsaGenerated, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeystoreEcdsaGenerated(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeystoreEcdsaGenerated(realmKey *RealmKeystoreEcdsaGenerated) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeystoreEcdsaGeneratedToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeystoreEcdsaGenerated(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
