package keycloak

import (
	"fmt"
	"strconv"
)

type OpenIdFullNameProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	IdTokenClaim       bool
	AccessTokenClaim   bool
	UserinfoTokenClaim bool
}

func (mapper *OpenIdFullNameProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "openid-connect",
		ProtocolMapper: "oidc-full-name-mapper",
		Config: map[string]string{
			idTokenClaimField:       strconv.FormatBool(mapper.IdTokenClaim),
			accessTokenClaimField:   strconv.FormatBool(mapper.AccessTokenClaim),
			userinfoTokenClaimField: strconv.FormatBool(mapper.UserinfoTokenClaim),
		},
	}
}

func (protocolMapper *protocolMapper) convertToOpenIdFullNameProtocolMapper(realmId, clientId, clientScopeId string) (*OpenIdFullNameProtocolMapper, error) {
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

	return &OpenIdFullNameProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		IdTokenClaim:       idTokenClaim,
		AccessTokenClaim:   accessTokenClaim,
		UserinfoTokenClaim: userinfoTokenClaim,
	}, nil
}

func (keycloakClient *KeycloakClient) GetOpenIdFullNameProtocolMapper(realmId, clientId, clientScopeId, mapperId string) (*OpenIdFullNameProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToOpenIdFullNameProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteOpenIdFullNameProtocolMapper(realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId))
}

func (keycloakClient *KeycloakClient) NewOpenIdFullNameProtocolMapper(mapper *OpenIdFullNameProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	location, err := keycloakClient.post(path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenIdFullNameProtocolMapper(mapper *OpenIdFullNameProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateOpenIdFullNameProtocolMapper(mapper *OpenIdFullNameProtocolMapper) error {
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
