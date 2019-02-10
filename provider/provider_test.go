package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

var requiredEnvironmentVariables = []string{
	"KEYCLOAK_CLIENT_ID",
	"KEYCLOAK_URL",
}

func init() {
	testAccProvider = KeycloakProvider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"keycloak": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			t.Fatalf("%s must be set before running acceptance tests.", requiredEnvironmentVariable)
		}
	}
	if v := os.Getenv("KEYCLOAK_CLIENT_SECRET"); v == "" {
		if v := os.Getenv("KEYCLOAK_USER"); v == "" {
			t.Fatal("KEYCLOAK_USER must be set for acceptance tests")
		}
		if v := os.Getenv("KEYCLOAK_PASSWORD"); v == "" {
			t.Fatal("KEYCLOAK_PASSWORD must be set for acceptance tests")
		}
	}
}
