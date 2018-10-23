package keycloak

import (
	"fmt"
	"strconv"
)

type OpenIdGroupMembershipProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	ClaimName          string
	FullPath           bool
	IdTokenClaim       bool
	AccessTokenClaim   bool
	UserinfoTokenClaim bool
}

var (
	fullPathField           = "full.path"
	idTokenClaimField       = "id.token.claim"
	accessTokenClaimField   = "access.token.claim"
	userinfoTokenClaimField = "userinfo.token.claim"
)

func (mapper *OpenIdGroupMembershipProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "openid-connect",
		ProtocolMapper: "oidc-group-membership-mapper",
		Config: map[string]string{
			fullPathField:           strconv.FormatBool(mapper.FullPath),
			idTokenClaimField:       strconv.FormatBool(mapper.IdTokenClaim),
			accessTokenClaimField:   strconv.FormatBool(mapper.AccessTokenClaim),
			userinfoTokenClaimField: strconv.FormatBool(mapper.UserinfoTokenClaim),
			claimNameField:          mapper.ClaimName,
		},
	}
}

func (protocolMapper *protocolMapper) convertToOpenIdGroupMembershipProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdGroupMembershipProtocolMapper, error) {
	fullPath, err := strconv.ParseBool(protocolMapper.Config[fullPathField])
	if err != nil {
		return nil, err
	}
	idTokenClaim, err := strconv.ParseBool(protocolMapper.Config[idTokenClaimField])
	if err != nil {
		return nil, err
	}
	accessTokenClaim, err := strconv.ParseBool(protocolMapper.Config[accessTokenClaimField])
	if err != nil {
		return nil, err
	}
	userinfoTokenClaim, err := strconv.ParseBool(protocolMapper.Config[userinfoTokenClaimField])
	if err != nil {
		return nil, err
	}

	return &OpenIdGroupMembershipProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		ClaimName:          protocolMapper.Config[claimNameField],
		FullPath:           fullPath,
		IdTokenClaim:       idTokenClaim,
		AccessTokenClaim:   accessTokenClaim,
		UserinfoTokenClaim: userinfoTokenClaim,
	}, nil
}

func (keycloakClient *KeycloakClient) GetOpenIdGroupMembershipProtocolMapper(realmId, clientId, clientScopeId, mapperId string) (*OpenIdGroupMembershipProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper)

	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdGroupMembershipProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteOpenIdGroupMembershipProtocolMapper(realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId))
}

func (keycloakClient *KeycloakClient) NewOpenIdGroupMembershipProtocolMapper(mapper *OpenIdGroupMembershipProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)
	location, err := keycloakClient.post(path, mapper.convertToGenericProtocolMapper())

	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenIdGroupMembershipProtocolMapper(mapper *OpenIdGroupMembershipProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (mapper *OpenIdGroupMembershipProtocolMapper) Validate() error {
	if mapper.ClientId == "" && mapper.ClientScopeId == "" {
		return fmt.Errorf("validation error: one of ClientId or ClientScopeId must be set")
	}

	return nil
}
