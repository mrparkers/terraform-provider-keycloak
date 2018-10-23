package keycloak

import (
	"fmt"
	"strconv"
)

type OpenIdUserAttributeProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	AddToIdToken     bool
	AddToAccessToken bool
	AddToUserInfo    bool

	UserAttribute  string
	ClaimName      string
	ClaimValueType string

	Multivalued bool // indicates whether is this an array of attributes or a single attribute
}

var (
	addToIdTokenField     = "id.token.claim"
	addToAccessTokenField = "access.token.claim"
	addToUserInfoField    = "userinfo.token.claim"
	userAttributeField    = "user.attribute"
	claimNameField        = "claim.name"
	claimValueTypeField   = "jsonType.label"
	multivaluedField      = "multivalued"
)

func (mapper *OpenIdUserAttributeProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "openid-connect",
		ProtocolMapper: "oidc-usermodel-attribute-mapper",
		Config: map[string]string{
			addToIdTokenField:     strconv.FormatBool(mapper.AddToIdToken),
			addToAccessTokenField: strconv.FormatBool(mapper.AddToAccessToken),
			addToUserInfoField:    strconv.FormatBool(mapper.AddToUserInfo),
			userAttributeField:    mapper.UserAttribute,
			claimNameField:        mapper.ClaimName,
			claimValueTypeField:   mapper.ClaimValueType,
			multivaluedField:      strconv.FormatBool(mapper.Multivalued),
		},
	}
}

func parseBoolAndTreatEmptyStringAsFalse(b string) (bool, error) {
	if b == "" {
		return false, nil
	}

	return strconv.ParseBool(b)
}

func (protocolMapper *protocolMapper) convertToOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdUserAttributeProtocolMapper, error) {
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

	return &OpenIdUserAttributeProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		AddToIdToken:     addToIdToken,
		AddToAccessToken: addToAccessToken,
		AddToUserInfo:    addToUserInfo,

		UserAttribute:  protocolMapper.Config[userAttributeField],
		ClaimName:      protocolMapper.Config[claimNameField],
		ClaimValueType: protocolMapper.Config[claimValueTypeField],
		Multivalued:    multivalued,
	}, nil
}

func (keycloakClient *KeycloakClient) GetOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId, mapperId string) (*OpenIdUserAttributeProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper)

	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId))
}

func (keycloakClient *KeycloakClient) NewOpenIdUserAttributeProtocolMapper(mapper *OpenIdUserAttributeProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
	location, err := keycloakClient.post(path, mapper.convertToGenericProtocolMapper())

	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenIdUserAttributeProtocolMapper(mapper *OpenIdUserAttributeProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (mapper *OpenIdUserAttributeProtocolMapper) Validate() error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}

	return nil
}
