package keycloak

import (
	"context"
	"fmt"
	"reflect"
)

type CustomIdentityProviderMapperConfig struct {
	ExtraConfig map[string]interface{} `json:"-"`
}

type CustomIdentityProviderMapper struct {
	Realm                  string                              `json:"-"`
	Provider               string                              `json:"-"`
	Id                     string                              `json:"id,omitempty"`
	Name                   string                              `json:"name,omitempty"`
	IdentityProviderAlias  string                              `json:"identityProviderAlias,omitempty"`
	IdentityProviderMapper string                              `json:"identityProviderMapper,omitempty"`
	Config                 *CustomIdentityProviderMapperConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewCustomIdentityProviderMapper(ctx context.Context, customIdentityProviderMapper *CustomIdentityProviderMapper) error {
	_, location, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers", customIdentityProviderMapper.Realm, customIdentityProviderMapper.IdentityProviderAlias), customIdentityProviderMapper)
	if err != nil {
		return err
	}

	customIdentityProviderMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetCustomIdentityProviderMapper(ctx context.Context, realm, alias, id string) (*CustomIdentityProviderMapper, error) {
	var customIdentityProviderMapper CustomIdentityProviderMapper
	customIdentityProviderMapper.Realm = realm
	customIdentityProviderMapper.IdentityProviderAlias = alias

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", realm, alias, id), &customIdentityProviderMapper, nil)
	if err != nil {
		return nil, err
	}

	return &customIdentityProviderMapper, nil
}

func (keycloakClient *KeycloakClient) UpdateCustomIdentityProviderMapper(ctx context.Context, customIdentityProviderMapper *CustomIdentityProviderMapper) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", customIdentityProviderMapper.Realm, customIdentityProviderMapper.IdentityProviderAlias, customIdentityProviderMapper.Id), customIdentityProviderMapper)
}

func (keycloakClient *KeycloakClient) DeleteCustomIdentityProviderMapper(ctx context.Context, realm, alias, id string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", realm, alias, id), nil)
}

func (f *CustomIdentityProviderMapperConfig) UnmarshalJSON(data []byte) error {
	return unmarshalExtraConfig(data, reflect.ValueOf(f).Elem(), &f.ExtraConfig)
}

func (f *CustomIdentityProviderMapperConfig) MarshalJSON() ([]byte, error) {
	return marshalExtraConfig(reflect.ValueOf(f).Elem(), f.ExtraConfig)
}
