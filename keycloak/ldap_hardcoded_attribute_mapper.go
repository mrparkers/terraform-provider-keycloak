package keycloak

import "fmt"

type LdapHardcodedAttributeMapper struct {
	Id                   string
	Name                 string
	RealmId              string
	LdapUserFederationId string
	AttributeName        string
	AttributeValue       string
}

func convertFromLdapHardcodedAttributeMapperToComponent(ldapMapper *LdapHardcodedAttributeMapper) *component {
	return &component{
		Id:           ldapMapper.Id,
		Name:         ldapMapper.Name,
		ProviderId:   "hardcoded-ldap-attribute-mapper",
		ProviderType: "org.keycloak.storage.ldap.mappers.LDAPStorageMapper",
		ParentId:     ldapMapper.LdapUserFederationId,

		Config: map[string][]string{
			"ldap.attribute.name": {
				ldapMapper.AttributeName,
			},
			"ldap.attribute.value": {
				ldapMapper.AttributeValue,
			},
		},
	}
}

func convertFromComponentToLdapHardcodedAttributeMapper(component *component, realmId string) *LdapHardcodedAttributeMapper {
	return &LdapHardcodedAttributeMapper{
		Id:                   component.Id,
		Name:                 component.Name,
		RealmId:              realmId,
		LdapUserFederationId: component.ParentId,

		AttributeName: component.getConfig("ldap.attribute.name"),
		AttributeValue: component.getConfig("ldap.attribute.value"),
	}
}

func (keycloakClient *KeycloakClient) ValidateLdapHardcodedAttributeMapper(ldapMapper *LdapHardcodedAttributeMapper) error {
	if len(ldapMapper.AttributeName) == 0 {
		return fmt.Errorf("validation error: hardcoded attribute name must not be empty")
	}
	if len(ldapMapper.AttributeValue) == 0 {
		return fmt.Errorf("validation error: hardcoded attribute value must not be empty")
	}
	return nil
}

func (keycloakClient *KeycloakClient) NewLdapHardcodedAttributeMapper(ldapMapper *LdapHardcodedAttributeMapper) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", ldapMapper.RealmId), convertFromLdapHardcodedAttributeMapperToComponent(ldapMapper))
	if err != nil {
		return err
	}

	ldapMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetLdapHardcodedAttributeMapper(realmId, id string) (*LdapHardcodedAttributeMapper, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToLdapHardcodedAttributeMapper(component, realmId), nil
}

func (keycloakClient *KeycloakClient) UpdateLdapHardcodedAttributeMapper(ldapMapper *LdapHardcodedAttributeMapper) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", ldapMapper.RealmId, ldapMapper.Id), convertFromLdapHardcodedAttributeMapperToComponent(ldapMapper))
}

func (keycloakClient *KeycloakClient) DeleteLdapHardcodedAttributeMapper(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
