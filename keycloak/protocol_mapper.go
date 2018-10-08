package keycloak

// https://www.keycloak.org/docs-api/4.2/rest-api/index.html#_protocolmapperrepresentation
type protocolMapper struct {
	Id             string            `json:"id,omitempty"`
	Name           string            `json:"name"`
	Protocol       string            `json:"protocol"`
	ProtocolMapper string            `json:"protocolMapper"`
	Config         map[string]string `json:"config"`
}
