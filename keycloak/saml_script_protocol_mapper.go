package keycloak

import (
	"fmt"
	"strconv"
)

type SamlScriptProtocolMapper struct {
	Id            string
	Name          string
	RealmId       string
	ClientId      string
	ClientScopeId string

	SingleValueAttribute bool

	SamlScript              string
	FriendlyName            string
	SamlAttributeName       string
	SamlAttributeNameFormat string
}

func (mapper *SamlScriptProtocolMapper) convertToGenericProtocolMapper() *protocolMapper {
	return &protocolMapper{
		Id:             mapper.Id,
		Name:           mapper.Name,
		Protocol:       "saml",
		ProtocolMapper: "saml-javascript-mapper",
		Config: map[string]string{
			attributeNameField:        mapper.SamlAttributeName,
			attributeNameFormatField:  mapper.SamlAttributeNameFormat,
			friendlyNameField:         mapper.FriendlyName,
			samlScriptField:           mapper.SamlScript,
			singleValueAttributeField: strconv.FormatBool(mapper.SingleValueAttribute),
		},
	}
}

func (protocolMapper *protocolMapper) convertToSamlScriptProtocolMapper(realmId, clientId, clientScopeId string) (*SamlScriptProtocolMapper, error) {
	singleValueAttribute, err := strconv.ParseBool(protocolMapper.Config[singleValueAttributeField])
	if err != nil {
		return nil, err
	}

	return &SamlScriptProtocolMapper{
		Id:            protocolMapper.Id,
		Name:          protocolMapper.Name,
		RealmId:       realmId,
		ClientId:      clientId,
		ClientScopeId: clientScopeId,

		SingleValueAttribute: singleValueAttribute,

		SamlScript:              protocolMapper.Config[samlScriptField],
		FriendlyName:            protocolMapper.Config[friendlyNameField],
		SamlAttributeName:       protocolMapper.Config[attributeNameField],
		SamlAttributeNameFormat: protocolMapper.Config[attributeNameFormatField],
	}, nil
}

func (keycloakClient *KeycloakClient) GetSamlScriptProtocolMapper(realmId, clientId, clientScopeId, mapperId string) (*SamlScriptProtocolMapper, error) {
	var protocolMapper *protocolMapper

	err := keycloakClient.get(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), &protocolMapper, nil)
	if err != nil {
		return nil, err
	}

	return protocolMapper.convertToSamlScriptProtocolMapper(realmId, clientId, clientScopeId)
}

func (keycloakClient *KeycloakClient) DeleteSamlScriptProtocolMapper(realmId, clientId, clientScopeId, mapperId string) error {
	return keycloakClient.delete(individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId), nil)
}

func (keycloakClient *KeycloakClient) NewSamlScriptProtocolMapper(mapper *SamlScriptProtocolMapper) error {
	path := protocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId)

	_, location, err := keycloakClient.post(path, mapper.convertToGenericProtocolMapper())
	if err != nil {
		return err
	}

	mapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) UpdateSamlScriptProtocolMapper(mapper *SamlScriptProtocolMapper) error {
	path := individualProtocolMapperPath(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)

	return keycloakClient.put(path, mapper.convertToGenericProtocolMapper())
}

func (keycloakClient *KeycloakClient) ValidateSamlScriptProtocolMapper(mapper *SamlScriptProtocolMapper) error {
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
