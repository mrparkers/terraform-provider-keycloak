package keycloak

import (
	"fmt"
)

const AudienceResolveMapperName = "audience resolve"

type OpenIdAudienceResolveProtocolMapper struct {
	Id            string
	RealmId       string
	ClientId      string
	ClientScopeId string
}

func (mapper *OpenIdAudienceResolveProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           AudienceResolveMapperName,
		Protocol:       "openid-connect",
		ProtocolMapper: "oidc-audience-resolve-mapper",
		Config:         map[string]string{},
	}
}

func (protocolMapper *protocolMapper) convertToOpenIdAudienceResolveProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdAudienceResolveProtocolMapper, error) {
	return &OpenIdAudienceResolveProtocolMapper{
		Id:            protocolMapper.Id,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,
	}, nil
}

func (keycloakClient *KeycloakClient) GetOpenIdAudienceResolveProtocolMapper(realmId, clientId, clientScopeId, mapperId string) (*OpenIdAudienceResolveProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper, nil)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdAudienceResolveProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteOpenIdAudienceResolveProtocolMapper(realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (keycloakClient *KeycloakClient) NewOpenIdAudienceResolveProtocolMapper(mapper *OpenIdAudienceResolveProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	_, location, err := keycloakClient.post(path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenIdAudienceResolveProtocolMapper(mapper *OpenIdAudienceResolveProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateOpenIdAudienceResolveProtocolMapper(mapper *OpenIdAudienceResolveProtocolMapper) error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}

	if mapper.ClientId != "" && mapper.ClientScopeId != "" {
		return fmt.Errorf("validation error: ClientId and ClientScopeId cannot both be set")
	}

	protocolMappers, err := keycloakClient.listGenericProtocolMappers(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
	if err != nil {
		return err
	}

	for _, protocolMapper := range protocolMappers {
		if protocolMapper.Name == AudienceResolveMapperName && protocolMapper.Id != mapper.Id {
			return fmt.Errorf("validation error: a protocol mapper with name %s already exists for this client", AudienceResolveMapperName)
		}
	}

	return nil
}
