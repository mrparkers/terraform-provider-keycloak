package keycloak

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var requiredEnvironmentVariables = []string{
	"KEYCLOAK_CLIENT_ID",
	"KEYCLOAK_CLIENT_SECRET",
	"KEYCLOAK_URL",
}

// Some actions, such as creating a realm, require a refresh
// before a GET can be performed on that realm
//
// This test ensures that, after creating a realm and performing
// a GET, the access token and refresh token have changed
//
// Any action that returns a 403 or a 401 could be used for this test
// Creating a realm is just the only one I'm aware of
func TestAccKeycloakApiClientRefresh(t *testing.T) {
	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			t.Fatalf("%s must be set before running acceptance tests.", requiredEnvironmentVariable)
		}
	}

	// Disable [DEBUG] logs which terraform typically handles for you. Re-enable when finished
	if tfLogLevel := os.Getenv("TF_LOG"); tfLogLevel == "" {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stdout)
	}

	keycloakClient, err := NewKeycloakClient(os.Getenv("KEYCLOAK_URL"), os.Getenv("KEYCLOAK_CLIENT_ID"), os.Getenv("KEYCLOAK_CLIENT_SECRET"))
	if err != nil {
		t.Fatalf("%s", err)
	}

	realmName := "terraform-" + acctest.RandString(10)
	realm := &Realm{
		Realm: realmName,
		Id:    realmName,
	}

	err = keycloakClient.NewRealm(realm)
	if err != nil {
		t.Fatalf("%s", err)
	}

	// A following GET for this realm will result in a 403, so we should save the current access and refresh token
	oldAccessToken := keycloakClient.clientCredentials.AccessToken
	oldRefreshToken := keycloakClient.clientCredentials.RefreshToken
	oldTokenType := keycloakClient.clientCredentials.TokenType

	_, err = keycloakClient.GetRealm(realmName) // This should not fail since it will automatically refresh and try again
	if err != nil {
		t.Fatalf("%s", err)
	}

	// Clean up - the realm doesn't need to exist in order for us to assert against the refreshed tokens
	err = keycloakClient.DeleteRealm(realmName)
	if err != nil {
		t.Fatalf("%s", err)
	}

	newAccessToken := keycloakClient.clientCredentials.AccessToken
	newRefreshToken := keycloakClient.clientCredentials.RefreshToken
	newTokenType := keycloakClient.clientCredentials.TokenType

	if oldAccessToken == newAccessToken {
		t.Fatalf("expected access token to update after refresh")
	}

	if oldRefreshToken == newRefreshToken {
		t.Fatalf("expected refresh token to update after refresh")
	}

	if oldTokenType != newTokenType {
		t.Fatalf("expected token type to remain the same after refresh")
	}
}
