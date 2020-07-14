package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var keycloakClient *keycloak.KeycloakClient

var requiredEnvironmentVariables = []string{
	"KEYCLOAK_CLIENT_ID",
	"KEYCLOAK_CLIENT_SECRET",
	"KEYCLOAK_REALM",
	"KEYCLOAK_URL",
}

func init() {
	testAccProvider = KeycloakProvider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"keycloak": testAccProvider,
	}

	keycloakClient, _ = keycloak.NewKeycloakClient(os.Getenv("KEYCLOAK_URL"), "/auth", os.Getenv("KEYCLOAK_CLIENT_ID"), os.Getenv("KEYCLOAK_CLIENT_SECRET"), os.Getenv("KEYCLOAK_REALM"), "", "", true, 5, "", false, httpclient.TerraformUserAgent(testAccProvider.TerraformVersion))
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
}
