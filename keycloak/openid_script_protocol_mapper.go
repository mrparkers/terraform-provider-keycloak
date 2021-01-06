package keycloak

import (
	"fmt"
	"strconv"
)

type OpenIdScriptProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	AddToIdToken     bool
	AddToAccessToken bool
	AddToUserInfo    bool

	Script         string
	ClaimName      string
	ClaimValueType string

	Multivalued bool // indicates whether is this an array of attributes or a single attribute
}

func (mapper *OpenIdScriptProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "openid-connect",
		ProtocolMapper: "oidc-script-based-protocol-mapper",
		Config: map[string]string{
			addToIdTokenField:     strconv.FormatBool(mapper.AddToIdToken),
			addToAccessTokenField: strconv.FormatBool(mapper.AddToAccessToken),
			addToUserInfoField:    strconv.FormatBool(mapper.AddToUserInfo),
			scriptField:           mapper.Script,
			claimNameField:        mapper.ClaimName,
			claimValueTypeField:   mapper.ClaimValueType,
			multivaluedField:      strconv.FormatBool(mapper.Multivalued),
		},
	}
}

func (protocolMapper *protocolMapper) convertToOpenIdScriptProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdScriptProtocolMapper, error) {
	addToIdToken, err := strconv.ParseBool(protocolMapper.Config[addToIdTokenField])
	if err != nil {
		return nil, err
	}

	addToAccessToken, err := strconv.ParseBool(protocolMapper.Config[addToAccessTokenField])
	if err != nil {
		return nil, err
	}

	addToUserInfo, err := strconv.ParseBool(protocolMapper.Config[addToUserInfoField])
	if err != nil {
		return nil, err
	}

	// multivalued's default is "", this is an issue when importing an existing mapper
	multivalued, err := parseBoolAndTreatEmptyStringAsFalse(protocolMapper.Config[multivaluedField])
	if err != nil {
		return nil, err
	}

	return &OpenIdScriptProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		AddToIdToken:     addToIdToken,
		AddToAccessToken: addToAccessToken,
		AddToUserInfo:    addToUserInfo,

		Script:         protocolMapper.Config[scriptField],
		ClaimName:      protocolMapper.Config[claimNameField],
		ClaimValueType: protocolMapper.Config[claimValueTypeField],
		Multivalued:    multivalued,
	}, nil
}

func (keycloakClient *KeycloakClient) GetOpenIdScriptProtocolMapper(realmId, clientId, clientScopeId, mapperId string) (*OpenIdScriptProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper, nil)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdScriptProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteOpenIdScriptProtocolMapper(realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (keycloakClient *KeycloakClient) NewOpenIdScriptProtocolMapper(mapper *OpenIdScriptProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	_, location, err := keycloakClient.post(path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenIdScriptProtocolMapper(mapper *OpenIdScriptProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateOpenIdScriptProtocolMapper(mapper *OpenIdScriptProtocolMapper) error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}

	protocolMappers, err := keycloakClient.listGenericProtocolMappers(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
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
