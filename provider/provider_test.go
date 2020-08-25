package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"os"
	"testing"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var keycloakClient *keycloak.KeycloakClient

var requiredEnvironmentVariables = []string{
	"KEYCLOAK_CLIENT_ID",
	"KEYCLOAK_CLIENT_SECRET",
	"KEYCLOAK_REALM",
	"KEYCLOAK_URL",
}

func init() {
	ctx := context.Background()
	userAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", schema.Provider{}.TerraformVersion, meta.SDKVersionString())
	keycloakClient, _ = keycloak.NewKeycloakClient(ctx, os.Getenv("KEYCLOAK_URL"), "/auth", os.Getenv("KEYCLOAK_CLIENT_ID"), os.Getenv("KEYCLOAK_CLIENT_SECRET"), os.Getenv("KEYCLOAK_REALM"), "", "", true, 5, "", false, userAgent)

	testAccProvider = KeycloakProvider()
	providerConfigureFunc := testAccProvider.ConfigureContextFunc
	// override the default context so a failing test doesn't cancel other tests
	testAccProvider.ConfigureContextFunc = func(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigureFunc(ctx, data)
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"keycloak": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
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
}
