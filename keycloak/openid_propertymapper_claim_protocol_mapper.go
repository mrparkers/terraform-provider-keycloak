package keycloak

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"
)

type OpenIdPropertyMapperClaimProtocolMapper struct {
	Id             string
	Name           string
	Protocol       string
	ProtocolMapper string
	RealmId        string
	ClientId       string
	ClientScopeId  string

	AddToIdToken            bool
	AddToAccessToken        bool
	AddToUserInfo           bool
	AddToIntrospectionToken bool
	AddToLightweightClaim   bool

	ClaimName string
	JsonType  string

	AdditionalConfig map[string]string
}

func (mapper *OpenIdPropertyMapperClaimProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {

	config := map[string]string{
		addToIdTokenField:            strconv.FormatBool(mapper.AddToIdToken),
		addToAccessTokenField:        strconv.FormatBool(mapper.AddToAccessToken),
		addToUserInfoField:           strconv.FormatBool(mapper.AddToUserInfo),
		addToIntrospectionTokenField: strconv.FormatBool(mapper.AddToIntrospectionToken),
		addToLightweightClaimField:   strconv.FormatBool(mapper.AddToLightweightClaim),
		jsonTypeField:                mapper.JsonType,
		claimNameField:               mapper.ClaimName,
	}

	maps.Copy(config, mapper.AdditionalConfig)
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       mapper.Protocol,
		ProtocolMapper: mapper.ProtocolMapper,
		Config:         config,
	}
}

func (protocolMapper *protocolMapper) convertToOpenIdPropertyMapperClaimProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdPropertyMapperClaimProtocolMapper, error) {
	addToIdToken, err := parseBoolAndTreatEmptyStringAsFalse(protocolMapper.Config[addToIdTokenField])
	if err != nil {
		return nil, err
	}

	addToAccessToken, err := parseBoolAndTreatEmptyStringAsFalse(protocolMapper.Config[addToAccessTokenField])
	if err != nil {
		return nil, err
	}

	addToUserInfo, err := parseBoolAndTreatEmptyStringAsFalse(protocolMapper.Config[addToUserInfoField])
	if err != nil {
		return nil, err
	}

	addToIntrospectionTokenField, err := parseBoolAndTreatEmptyStringAsFalse(protocolMapper.Config[addToIntrospectionTokenField])
	if err != nil {
		return nil, err
	}

	addToLightweightClaimField, err := parseBoolAndTreatEmptyStringAsFalse(protocolMapper.Config[addToLightweightClaimField])
	if err != nil {
		return nil, err
	}

	additionalConfig := map[string]string{}
	for k, v := range protocolMapper.Config {
		if !slices.Contains(protocolMapperIgnore, k) {
			additionalConfig[k] = v
		}
	}

	return &OpenIdPropertyMapperClaimProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		AddToIdToken:            addToIdToken,
		AddToAccessToken:        addToAccessToken,
		AddToUserInfo:           addToUserInfo,
		AddToIntrospectionToken: addToIntrospectionTokenField,
		AddToLightweightClaim:   addToLightweightClaimField,

		Protocol:         protocolMapper.Protocol,
		ProtocolMapper:   protocolMapper.ProtocolMapper,
		ClaimName:        protocolMapper.Config[claimNameField],
		JsonType:         protocolMapper.Config[jsonTypeField],
		AdditionalConfig: additionalConfig,
	}, nil
}

func (keycloakClient *KeycloakClient) GetOpenIdPropertyMapperClaimProtocolMapper(ctx context.Context, realmId, clientId, clientScopeId, mapperId string) (*OpenIdPropertyMapperClaimProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper, nil)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdPropertyMapperClaimProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteOpenIdPropertyMapperClaimProtocolMapper(ctx context.Context, realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(ctx, individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (keycloakClient *KeycloakClient) NewOpenIdPropertyMapperClaimProtocolMapper(ctx context.Context, mapper *OpenIdPropertyMapperClaimProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	_, location, err := keycloakClient.post(ctx, path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenIdPropertyMapperClaimProtocolMapper(ctx context.Context, mapper *OpenIdPropertyMapperClaimProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(ctx, path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateOpenIdPropertyMapperClaimProtocolMapper(ctx context.Context, mapper *OpenIdPropertyMapperClaimProtocolMapper) error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}

	protocolMappers, err := keycloakClient.listGenericProtocolMappers(ctx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
	if err != nil {
		return err
	}

	for _, protocolMapper := range protocolMappers {
		if protocolMapper.Name == mapper.Name && protocolMapper.Id != mapper.Id {
			return fmt.Errorf("validation error: a protocol mapper with name %s already exists for this client", mapper.Name)
		}
	}

	return nil
}
