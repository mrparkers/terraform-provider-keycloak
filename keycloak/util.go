package keycloak

import "strings"

func getIdFromLocationHeader(locationHeader string) string {
	parts := strings.Split(locationHeader, "/")

	return parts[len(parts)-1]
}
