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

func protocolMapperPath(realmId, clientId, clientScopeId string) string {
	parentResourceId := clientId
	parentResourcePath := "clients"

	if clientScopeId != "" {
		parentResourceId = clientScopeId
		parentResourcePath = "client-scopes"
	}

	return fmt.Sprintf("/realms/%s/%s/%s/protocol-mappers/models", realmId, parentResourceId, parentResourcePath)
}

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

	multivalued, err := strconv.ParseBool(protocolMapper.Config[multivaluedField])

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

func (keycloakClient *KeycloakClient) getOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdUserAttributeProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(protocolMapperPath(realmId, clientId, clientScopeId), &protocolMapper)

	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) deleteOpenIdUserAttributeProtocolMapper(realmId, clientId, clientScopeId string) error {
	path := protocolMapperPath(realmId, clientId, clientScopeId)

	return keycloakClient.delete(path)
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

func (keycloakClient *KeycloakClient) GetOpenIdUserAttributeProtocolMapperForClient(realmId, clientId string) (*OpenIdUserAttributeProtocolMapper, error) {
	return keycloakClient.getOpenIdUserAttributeProtocolMapper(realmId, clientId, "")
}

func (keycloakClient *KeycloakClient) GetOpenIdUserAttributeProtocolMapperForClientScope(realmId, clientScopeId string) (*OpenIdUserAttributeProtocolMapper, error) {
	return keycloakClient.getOpenIdUserAttributeProtocolMapper(realmId, "", clientScopeId)
}

func (keycloakClient *KeycloakClient) UpdateOpenIdUserAttributeProtocolMapper(mapper *OpenIdUserAttributeProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) DeleteOpenIdUserAttributeProtocolMapperForClient(realmId, clientId string) error {
	return keycloakClient.deleteOpenIdUserAttributeProtocolMapper(realmId, clientId, "")
}

func (keycloakClient *KeycloakClient) DeleteOpenIdUserAttributeProtocolMapperForClientScope(realmId, clientScopeId string) error {
	return keycloakClient.deleteOpenIdUserAttributeProtocolMapper(realmId, "", clientScopeId)
}
