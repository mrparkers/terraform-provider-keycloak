package keycloak

import "fmt"

// https://www.keycloak.org/docs-api/4.2/rest-api/index.html#_protocolmapperrepresentation
type protocolMapper struct {
	Id             string            `json:"id,omitempty"`
	Name           string            `json:"name"`
	Protocol       string            `json:"protocol"`
	ProtocolMapper string            `json:"protocolMapper"`
	Config         map[string]string `json:"config"`
}

var (
	accessTokenClaimField   = "access.token.claim"
	addToAccessTokenField   = "access.token.claim"
	addToIdTokenField       = "id.token.claim"
	addToUserInfoField      = "userinfo.token.claim"
	claimNameField          = "claim.name"
	claimValueTypeField     = "jsonType.label"
	fullPathField           = "full.path"
	idTokenClaimField       = "id.token.claim"
	multivaluedField        = "multivalued"
	userAttributeField      = "user.attribute"
	userinfoTokenClaimField = "userinfo.token.claim"
)

func protocolMapperPath(realmId, clientId, clientScopeId string) string {
	parentResourceId := clientId
	parentResourcePath := "clients"

	if clientScopeId != "" {
		parentResourceId = clientScopeId
		parentResourcePath = "client-scopes"
	}

	return fmt.Sprintf("/realms/%s/%s/%s/protocol-mappers/models", realmId, parentResourcePath, parentResourceId)
}

func individualProtocolMapperPath(realmId, clientId, clientScopeId, mapperId string) string {
	return fmt.Sprintf("%s/%s", protocolMapperPath(realmId, clientId, clientScopeId), mapperId)
}
