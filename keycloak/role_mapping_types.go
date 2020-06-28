package keycloak

// struct for the MappingRepresentation
// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_mappingsrepresentation
type RoleMapping struct {
	ClientMappings map[string]*ClientRoleMapping `json:"clientMappings"`
	RealmMappings  []*Role                       `json:"realmMappings"`
}

// struct for the ClientMappingRepresentation
// https://www.keycloak.org/docs-api/8.0/rest-api/index.html#_clientmappingsrepresentation
type ClientRoleMapping struct {
	Client   string  `json:"client"`
	Id       string  `json:"id"`
	Mappings []*Role `json:"mappings"`
}
