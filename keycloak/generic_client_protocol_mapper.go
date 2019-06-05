package keycloak

import (
	"fmt"
)

type GenericClientProtocolMapper struct {
	ClientId       string            `json:"-"`
	Config         map[string]string `json:"config"`
	Id             string            `json:"id,omitempty"`
	Name           string            `json:"name"`
	Protocol       string            `json:"protocol"`
	ProtocolMapper string            `json:"protocolMapper"`
	RealmId        string            `json:"-"`
}

func (keycloakClient *KeycloakClient) NewGenericClientProtocolMapper(genericClientProtocolMapper *GenericClientProtocolMapper) error {
	_, location, err := keycloakClient.post(
		fmt.Sprintf("/realms/%s/clients/%s/protocol-mappers/models", genericClientProtocolMapper.RealmId, genericClientProtocolMapper.ClientId),
		genericClientProtocolMapper)
	if err != nil {
		return err
	}

	genericClientProtocolMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetGenericClientProtocolMapper(realmId string, clientId string, id string) (*GenericClientProtocolMapper, error) {
	var genericClientProtocolMapper GenericClientProtocolMapper

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/protocol-mappers/models/%s", realmId, clientId, id), &genericClientProtocolMapper, nil)
	if err != nil {
		return nil, err
	}

	// these values are not provided by the keycloak API
	genericClientProtocolMapper.ClientId = clientId
	genericClientProtocolMapper.RealmId = realmId

	return &genericClientProtocolMapper, nil
}

func (keycloakClient *KeycloakClient) UpdateGenericClientProtocolMapper(genericClientProtocolMapper *GenericClientProtocolMapper) error {
	return keycloakClient.put(
		fmt.Sprintf("/realms/%s/clients/%s/protocol-mappers/models/%s", genericClientProtocolMapper.RealmId, genericClientProtocolMapper.ClientId, genericClientProtocolMapper.Id),
		genericClientProtocolMapper)
}

func (keycloakClient *KeycloakClient) DeleteGenericClientProtocolMapper(realmId string, clientId string, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/protocol-mappers/models/%s", realmId, clientId, id), nil)
}
