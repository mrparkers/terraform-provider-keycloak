package keycloak

import (
	"fmt"
	"strconv"
)

type LdapFullNameMapper struct {
	Id                   string
	Name                 string
	RealmId              string
	LdapUserFederationId string

	LdapFullNameAttribute string
	ReadOnly              bool
	WriteOnly             bool
}

func convertFromLdapFullNameMapperToComponent(ldapFullNameMapper *LdapFullNameMapper) *component {
	return &component{
		Id:           ldapFullNameMapper.Id,
		Name:         ldapFullNameMapper.Name,
		ProviderId:   "full-name-ldap-mapper",
		ProviderType: "org.keycloak.storage.ldap.mappers.LDAPStorageMapper",
		ParentId:     ldapFullNameMapper.LdapUserFederationId,
		Config: map[string][]string{
			"ldap.full.name.attribute": {
				ldapFullNameMapper.LdapFullNameAttribute,
			},
			"read.only": {
				strconv.FormatBool(ldapFullNameMapper.ReadOnly),
			},
			"write.only": {
				strconv.FormatBool(ldapFullNameMapper.WriteOnly),
			},
		},
	}
}

func convertFromComponentToLdapFullNameMapper(component *component, realmId string) (*LdapFullNameMapper, error) {
	readOnly, err := strconv.ParseBool(component.getConfig("read.only"))
	if err != nil {
		return nil, err
	}

	writeOnly, err := strconv.ParseBool(component.getConfig("write.only"))
	if err != nil {
		return nil, err
	}

	return &LdapFullNameMapper{
		Id:                   component.Id,
		Name:                 component.Name,
		RealmId:              realmId,
		LdapUserFederationId: component.ParentId,

		LdapFullNameAttribute: component.getConfig("ldap.full.name.attribute"),
		ReadOnly:              readOnly,
		WriteOnly:             writeOnly,
	}, nil
}

func (keycloakClient *KeycloakClient) NewLdapFullNameMapper(ldapFullNameMapper *LdapFullNameMapper) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", ldapFullNameMapper.RealmId), convertFromLdapFullNameMapperToComponent(ldapFullNameMapper))
	if err != nil {
		return err
	}

	ldapFullNameMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetLdapFullNameMapper(realmId, id string) (*LdapFullNameMapper, error) {
	var component *component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToLdapFullNameMapper(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateLdapFullNameMapper(ldapFullNameMapper *LdapFullNameMapper) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", ldapFullNameMapper.RealmId, ldapFullNameMapper.Id), convertFromLdapFullNameMapperToComponent(ldapFullNameMapper))
}

func (keycloakClient *KeycloakClient) DeleteLdapFullNameMapper(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id))
}
