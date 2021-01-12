package keycloak

import "github.com/hashicorp/go-version"

type Version string

const (
	Version_6  Version = "6.0.0"
	Version_7  Version = "7.0.0"
	Version_8  Version = "8.0.0"
	Version_9  Version = "9.0.0"
	Version_10 Version = "10.0.0"
	Version_11 Version = "11.0.0"
	Version_12 Version = "12.0.0"
)

func (keycloakClient *KeycloakClient) VersionIsGreaterThanOrEqualTo(versionString Version) bool {
	v, _ := version.NewVersion(string(versionString))

	return keycloakClient.version.GreaterThanOrEqual(v)
}
