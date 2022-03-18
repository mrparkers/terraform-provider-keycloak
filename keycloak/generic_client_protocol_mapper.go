package keycloak

import (
	"context"
	"fmt"
)

type GenericClientProtocolMapper struct {
	ClientId       string            `json:"-"`
	ClientScopeId  string            `json:"-"`
	Config         map[string]string `json:"config"`
	Id             string            `json:"id,omitempty"`
	Name           string            `json:"name"`
	Protocol       string            `json:"protocol"`
	ProtocolMapper string            `json:"protocolMapper"`
	RealmId        string            `json:"-"`
}

type OpenidClientWithGenericClientProtocolMappers struct {
	OpenidClient
	ProtocolMappers []*GenericClientProtocolMapper
}

func (keycloakClient *KeycloakClient) NewGenericClientProtocolMapper(ctx context.Context, genericClientProtocolMapper *GenericClientProtocolMapper) error {
	path := protocolMapperPath(genericClientProtocolMapper.RealmId, genericClientProtocolMapper.ClientId, genericClientProtocolMapper.ClientScopeId)

	_, location, err := keycloakClient.post(ctx, path, genericClientProtocolMapper)
	if err != nil {
		return err
	}

	genericClientProtocolMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetGenericClientProtocolMappers(ctx context.Context, realmId string, clientId string) (*OpenidClientWithGenericClientProtocolMappers, error) {
	var openidClientWithGenericClientProtocolMappers OpenidClientWithGenericClientProtocolMappers

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/clients/%s", realmId, clientId), &openidClientWithGenericClientProtocolMappers, nil)
	if err != nil {
		return nil, err
	}

	openidClientWithGenericClientProtocolMappers.RealmId = realmId
	openidClientWithGenericClientProtocolMappers.ClientId = clientId

	for _, protocolMapper := range openidClientWithGenericClientProtocolMappers.ProtocolMappers {
		protocolMapper.RealmId = realmId
		protocolMapper.ClientId = clientId
	}

	return &openidClientWithGenericClientProtocolMappers, nil

}

func (keycloakClient *KeycloakClient) GetGenericClientProtocolMapper(ctx context.Context, realmId string, clientId string, clientScopeId string, mapperId string) (*GenericClientProtocolMapper, error) {
	var genericClientProtocolMapper GenericClientProtocolMapper

	err := keycloakClient.get(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &genericClientProtocolMapper, nil)
	if err != nil {
		return nil, err
	}

	// these values are not provided by the keycloak API
	genericClientProtocolMapper.ClientId = clientId
	genericClientProtocolMapper.ClientScopeId = clientScopeId
	genericClientProtocolMapper.RealmId = realmId

	return &genericClientProtocolMapper, nil
}

func (keycloakClient *KeycloakClient) UpdateGenericClientProtocolMapper(ctx context.Context, genericClientProtocolMapper *GenericClientProtocolMapper) error {
	path := individualProtocolMapperPath(genericClientProtocolMapper.RealmId, genericClientProtocolMapper.ClientId, genericClientProtocolMapper.ClientScopeId, genericClientProtocolMapper.Id)

	return keycloakClient.put(ctx, path, genericClientProtocolMapper)
}

func (keycloakClient *KeycloakClient) DeleteGenericClientProtocolMapper(ctx context.Context, realmId string, clientId string, clientScopeId string, mapperId string) error {
	return keycloakClient.delete(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (mapper *GenericClientProtocolMapper) Validate(ctx context.Context, keycloakClient *KeycloakClient) error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}
	if mapper.ClientId != "" && mapper.ClientScopeId != "" {
		return fmt.Errorf("validation error: only one of ClientId or ClientScopeId must be set")
	}

	protocolMappers, err := keycloakClient.listGenericProtocolMappers(ctx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
	if err != nil {
		return err
	}

	for _, protocolMapper := range protocolMappers {
		if protocolMapper.Name == mapper.Name {
			return fmt.Errorf("validation error: a protocol mapper with name %s already exists for this client", mapper.Name)
		}
	}

	return nil
}
