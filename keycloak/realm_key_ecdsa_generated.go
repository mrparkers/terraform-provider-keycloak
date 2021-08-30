package keycloak

import (
	"fmt"
	"strconv"
)

type RealmKeyEcdsaGenerated struct {
	Id       string
	Name     string
	RealmId  string
	ParentId string

	Active        bool
	Enabled       bool
	Priority      int
	EllipticCurve string
}

func convertFromRealmKeyEcdsaGeneratedToComponent(realmKey *RealmKeyEcdsaGenerated) *component {
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
		ParentId:     realmKey.ParentId,
		ProviderId:   "ecdsa-generated",
		ProviderType: "org.keycloak.keys.KeyProvider",
		Config:       componentConfig,
	}
}

func convertFromComponentToRealmKeyEcdsaGenerated(component *component, realmId string) (*RealmKeyEcdsaGenerated, error) {
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

	realmKey := &RealmKeyEcdsaGenerated{
		Id:       component.Id,
		Name:     component.Name,
		ParentId: component.ParentId,
		RealmId:  realmId,

		Active:        active,
		Enabled:       enabled,
		Priority:      priority,
		EllipticCurve: component.getConfig("ecdsaEllipticCurveKey"),
	}

	return realmKey, nil
}

func (keycloakClient *KeycloakClient) NewRealmKeyEcdsaGenerated(realmKey *RealmKeyEcdsaGenerated) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realmKey.RealmId), convertFromRealmKeyEcdsaGeneratedToComponent(realmKey))
	if err != nil {
		return err
	}

	realmKey.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetRealmKeyEcdsaGenerated(realmId, id string) (*RealmKeyEcdsaGenerated, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToRealmKeyEcdsaGenerated(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateRealmKeyEcdsaGenerated(realmKey *RealmKeyEcdsaGenerated) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", realmKey.RealmId, realmKey.Id), convertFromRealmKeyEcdsaGeneratedToComponent(realmKey))
}

func (keycloakClient *KeycloakClient) DeleteRealmKeyEcdsaGenerated(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
