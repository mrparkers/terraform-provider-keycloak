package keycloak

import (
	"fmt"
	"strconv"
)

type CustomUserFederation struct {
	Id         string
	Name       string
	RealmId    string
	ParentId   string
	ProviderId string

	Enabled  bool
	Priority int

	CachePolicy string

	Config map[string][]string
}

var (
	userStorageProviderType = "org.keycloak.storage.UserStorageProvider"
)

func convertFromCustomUserFederationToComponent(custom *CustomUserFederation) *component {
	componentConfig := make(map[string][]string)

	if custom.Config != nil {
		for k, j := range custom.Config {
			componentConfig[k] = append(componentConfig[k], j[0])
		}
	}
	componentConfig["cachePolicy"] = append(componentConfig["cachePolicy"], custom.CachePolicy)
	componentConfig["enabled"] = append(componentConfig["enabled"], strconv.FormatBool(custom.Enabled))
	componentConfig["priority"] = append(componentConfig["priority"], strconv.Itoa(custom.Priority))
	parentId := custom.RealmId
	if custom.ParentId != "" {
		parentId = custom.ParentId
	}
	return &component{
		Id:           custom.Id,
		Name:         custom.Name,
		ProviderId:   custom.ProviderId,
		ProviderType: userStorageProviderType,
		ParentId:     parentId,
		Config:       componentConfig,
	}
}

func convertFromComponentToCustomUserFederation(component *component, realmName string) (*CustomUserFederation, error) {
	enabled, err := strconv.ParseBool(component.getConfig("enabled"))
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(component.getConfig("priority"))
	if err != nil {
		return nil, err
	}

	config := make(map[string][]string)
	for k := range component.Config {
		if k != "enabled" && k != "priority" && k != "cachePolicy" {
			config[k] = append(config[k], component.getConfig(k))
		}
	}

	custom := &CustomUserFederation{
		Id:         component.Id,
		Name:       component.Name,
		ParentId:   component.ParentId,
		RealmId:    realmName,
		ProviderId: component.ProviderId,

		Enabled:  enabled,
		Priority: priority,

		CachePolicy: component.getConfig("cachePolicy"),

		Config: config,
	}

	return custom, nil
}

func (keycloakClient *KeycloakClient) ValidateCustomUserFederation(custom *CustomUserFederation) error {
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
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", customUserFederation.RealmId), convertFromCustomUserFederationToComponent(customUserFederation))
	if err != nil {
		return err
	}

	customUserFederation.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetCustomUserFederation(realmName, id string) (*CustomUserFederation, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmName, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToCustomUserFederation(component, realmName)
}

func (keycloakClient *KeycloakClient) GetCustomUserFederations(realmName, realmId string) (*[]CustomUserFederation, error) {
	var components []*component
	var customUserFederations []CustomUserFederation
	var customUserFederation *CustomUserFederation

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components?parent=%s&type=%s", realmName, realmId, userStorageProviderType), &components, nil)
	if err != nil {
		return nil, err
	}

	for _, component := range components {
		customUserFederation, err = convertFromComponentToCustomUserFederation(component, realmName)
		if err != nil {
			return nil, err
		}
		customUserFederations = append(customUserFederations, *customUserFederation)
	}
	return &customUserFederations, nil
}

func (keycloakClient *KeycloakClient) UpdateCustomUserFederation(customUserFederation *CustomUserFederation) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", customUserFederation.RealmId, customUserFederation.Id), convertFromCustomUserFederationToComponent(customUserFederation))
}

func (keycloakClient *KeycloakClient) DeleteCustomUserFederation(realmName, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmName, id), nil)
}
