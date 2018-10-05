package keycloak

import (
	"fmt"
	"strconv"
)

type CustomUserFederation struct {
	Id         string
	Name       string
	RealmId    string
	ProviderId string

	Enabled  bool
	Priority int

	CachePolicy string
}

var (
	userStorageProviderType = "org.keycloak.storage.UserStorageProvider"
)

func convertFromCustomUserFederationToComponent(custom *CustomUserFederation) *component {
	componentConfig := map[string][]string{
		"cachePolicy": {
			custom.CachePolicy,
		},
		"enabled": {
			strconv.FormatBool(custom.Enabled),
		},
		"priority": {
			strconv.Itoa(custom.Priority),
		},
	}

	return &component{
		Id:           custom.Id,
		Name:         custom.Name,
		ProviderId:   custom.ProviderId,
		ProviderType: userStorageProviderType,
		ParentId:     custom.RealmId,
		Config:       componentConfig,
	}
}

func convertFromComponentToCustomUserFederation(component *component) (*CustomUserFederation, error) {
	enabled, err := strconv.ParseBool(component.getConfig("enabled"))
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(component.getConfig("priority"))
	if err != nil {
		return nil, err
	}

	custom := &CustomUserFederation{
		Id:         component.Id,
		Name:       component.Name,
		RealmId:    component.ParentId,
		ProviderId: component.ProviderId,

		Enabled:  enabled,
		Priority: priority,

		CachePolicy: component.getConfig("cachePolicy"),
	}

	return custom, nil
}

func (custom *CustomUserFederation) Validate(keycloakClient *KeycloakClient) error {
	// validate if the given custom user storage provider exists on the server.
	serverInfo, err := keycloakClient.GetServerInfo()
	if err != nil {
		return err
	}

	if !serverInfo.ComponentTypeIsInstalled(userStorageProviderType, custom.ProviderId) {
		return fmt.Errorf("custom user federation provider with id %s is not installed on the server", custom.ProviderId)
	}

	return nil
}

func (keycloakClient *KeycloakClient) NewCustomUserFederation(customUserFederation *CustomUserFederation) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", customUserFederation.RealmId), convertFromCustomUserFederationToComponent(customUserFederation))
	if err != nil {
		return err
	}

	customUserFederation.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetCustomUserFederation(realmId, id string) (*CustomUserFederation, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToCustomUserFederation(component)
}

func (keycloakClient *KeycloakClient) UpdateCustomUserFederation(customUserFederation *CustomUserFederation) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", customUserFederation.RealmId, customUserFederation.Id), convertFromCustomUserFederationToComponent(customUserFederation))
}

func (keycloakClient *KeycloakClient) DeleteCustomUserFederation(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id))
}
