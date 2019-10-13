package keycloak

import (
	"fmt"
	"log"
)

type IdentityProviderMapperConfig struct {
	UserAttribute                  string `json:"user.attribute,omitempty"`
	Claim                          string `json:"claim,omitempty"`
	ClaimValue                     string `json:"claim.value,omitempty"`
	HardcodedAttribute             string `json:"attribute,omitempty"`
	Attribute                      string `json:"attribute.name,omitempty"`
	AttributeValue                 string `json:"attribute.value,omitempty"`
	AttributeFriendlyName          string `json:"attribute.friendly.name,omitempty"`
	Template                       string `json:"template,omitempty"`
	Role                           string `json:"role,omitempty"`
	ProxyIdUserAttributeName       string `json:"proxy_id_user_attribute_name,omitempty"`
	HealthCloudIdUserAttributeName string `json:"health_cloud_id_user_attribute_name,omitempty"`
	HcAppIdUserAttributeName       string `json:"hc_app_id_user_attribute_name,omitempty"`
	ProxyIdApiUrl                  string `json:"proxy_id_api_url,omitempty"`
	HealthCloudIdApiUrl            string `json:"health_cloud_id_api_url,omitempty"`
	Tenant                         string `json:"tenant,omitempty"`
	SourceCode                     string `json:"source_code,omitempty"`
	SourceUidName                  string `json:"source_uid_name,omitempty"`
	IdpIdKeyName                   string `json:"idp_id_key_name,omitempty"`
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
