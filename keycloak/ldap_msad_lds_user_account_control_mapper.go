package keycloak

import (
	"context"
	"fmt"
)

type LdapMsadLdsUserAccountControlMapper struct {
	Id                   string
	Name                 string
	RealmId              string
	LdapUserFederationId string
}

func convertFromLdapMsadLdsUserAccountControlMapperToComponent(ldapMsadLdsUserAccountControlMapper *LdapMsadLdsUserAccountControlMapper) *component {
	return &component{
		Id:           ldapMsadLdsUserAccountControlMapper.Id,
		Name:         ldapMsadLdsUserAccountControlMapper.Name,
		ProviderId:   "msad-lds-user-account-control-mapper",
		ProviderType: "org.keycloak.storage.ldap.mappers.LDAPStorageMapper",
		ParentId:     ldapMsadLdsUserAccountControlMapper.LdapUserFederationId,
	}
}

func convertFromComponentToLdapMsadLdsUserAccountControlMapper(component *component, realmId string) (*LdapMsadLdsUserAccountControlMapper, error) {
	return &LdapMsadLdsUserAccountControlMapper{
		Id:                   component.Id,
		Name:                 component.Name,
		RealmId:              realmId,
		LdapUserFederationId: component.ParentId,
	}, nil
}

func (keycloakClient *KeycloakClient) NewLdapMsadLdsUserAccountControlMapper(ctx context.Context, ldapMsadLdsUserAccountControlMapper *LdapMsadLdsUserAccountControlMapper) error {
	_, location, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/components", ldapMsadLdsUserAccountControlMapper.RealmId), convertFromLdapMsadLdsUserAccountControlMapperToComponent(ldapMsadLdsUserAccountControlMapper))
	if err != nil {
		return err
	}

	ldapMsadLdsUserAccountControlMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetLdapMsadLdsUserAccountControlMapper(ctx context.Context, realmId, id string) (*LdapMsadLdsUserAccountControlMapper, error) {
	var component *component

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToLdapMsadLdsUserAccountControlMapper(component, realmId)
}

func (keycloakClient *KeycloakClient) UpdateLdapMsadLdsUserAccountControlMapper(ctx context.Context, ldapMsadLdsUserAccountControlMapper *LdapMsadLdsUserAccountControlMapper) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/components/%s", ldapMsadLdsUserAccountControlMapper.RealmId, ldapMsadLdsUserAccountControlMapper.Id), convertFromLdapMsadLdsUserAccountControlMapperToComponent(ldapMsadLdsUserAccountControlMapper))
}

func (keycloakClient *KeycloakClient) DeleteLdapMsadLdsUserAccountControlMapper(ctx context.Context, realmId, id string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
