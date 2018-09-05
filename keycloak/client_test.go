package keycloak

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"
)

func TestClientLocationParse(t *testing.T) {
	expectedRealm := "terraform-" + acctest.RandString(10)
	expectedId := "aed9f702-c0d3-42cd-91b6-a894da769324"
	locationHeader := fmt.Sprintf("http://localhost:8080/auth/admin/realms/%s/clients/%s", expectedRealm, expectedId)

	actualId := parseClientLocation(locationHeader)

	if expectedId != actualId {
		t.Fatalf("parsed keycloak client location header did not return correct ID")
	}
}
