package keycloak

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type IdentityProviderMapperConfig struct {
	UserAttribute         string                 `json:"user.attribute,omitempty"`
	Claim                 string                 `json:"claim,omitempty"`
	ClaimValue            string                 `json:"claim.value,omitempty"`
	HardcodedAttribute    string                 `json:"attribute,omitempty"`
	Attribute             string                 `json:"attribute.name,omitempty"`
	AttributeValue        string                 `json:"attribute.value,omitempty"`
	AttributeFriendlyName string                 `json:"attribute.friendly.name,omitempty"`
	Template              string                 `json:"template,omitempty"`
	Role                  string                 `json:"role,omitempty"`
	ExtraConfig           map[string]interface{} `json:"-"`
}

type IdentityProviderMapper struct {
	Realm                  string                        `json:"-"`
	Provider               string                        `json:"-"`
	Id                     string                        `json:"id,omitempty"`
	Name                   string                        `json:"name,omitempty"`
	IdentityProviderAlias  string                        `json:"identityProviderAlias,omitempty"`
	IdentityProviderMapper string                        `json:"identityProviderMapper,omitempty"`
	Config                 *IdentityProviderMapperConfig `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewIdentityProviderMapper(identityProviderMapper *IdentityProviderMapper) error {
	log.Printf("[WARN] Realm: %s", identityProviderMapper.Realm)
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers", identityProviderMapper.Realm, identityProviderMapper.IdentityProviderAlias), identityProviderMapper)
	if err != nil {
		return err
	}

	identityProviderMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetIdentityProviderMapper(realm, alias, id string) (*IdentityProviderMapper, error) {
	var identityProviderMapper IdentityProviderMapper
	identityProviderMapper.Realm = realm
	identityProviderMapper.IdentityProviderAlias = alias

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", realm, alias, id), &identityProviderMapper, nil)
	if err != nil {
		return nil, err
	}

	return &identityProviderMapper, nil
}

func (keycloakClient *KeycloakClient) UpdateIdentityProviderMapper(identityProviderMapper *IdentityProviderMapper) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", identityProviderMapper.Realm, identityProviderMapper.IdentityProviderAlias, identityProviderMapper.Id), identityProviderMapper)
}

func (keycloakClient *KeycloakClient) DeleteIdentityProviderMapper(realm, alias, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", realm, alias, id), nil)
}

func (f *IdentityProviderMapperConfig) UnmarshalJSON(data []byte) error {
	f.ExtraConfig = map[string]interface{}{}
	err := json.Unmarshal(data, &f.ExtraConfig)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(f).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		jsonKey := strings.Split(structField.Tag.Get("json"), ",")[0]
		if jsonKey != "-" {
			value, ok := f.ExtraConfig[jsonKey]
			if ok {
				field := v.FieldByName(structField.Name)
				if field.IsValid() && field.CanSet() {
					if field.Kind() == reflect.String {
						field.SetString(value.(string))
					} else if field.Kind() == reflect.Bool {
						boolVal, err := strconv.ParseBool(value.(string))
						if err == nil {
							field.Set(reflect.ValueOf(KeycloakBoolQuoted(boolVal)))
						}
					}
					delete(f.ExtraConfig, jsonKey)
				}
			}
		}
	}
	return nil
}

func (f *IdentityProviderMapperConfig) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{}

	for k, v := range f.ExtraConfig {
		out[k] = v
	}
	v := reflect.ValueOf(f).Elem()
	for i := 0; i < v.NumField(); i++ {
		jsonKey := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0]
		if jsonKey != "-" {
			field := v.Field(i)
			if field.IsValid() && field.CanSet() {
				if field.Kind() == reflect.String {
					out[jsonKey] = field.String()
				} else if field.Kind() == reflect.Bool {
					out[jsonKey] = KeycloakBoolQuoted(field.Bool())
				}
			}
		}
	}
	return json.Marshal(out)
}
