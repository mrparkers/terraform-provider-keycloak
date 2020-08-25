package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestAccKeycloakProvider_passwordGrant(t *testing.T) {
	skipIfEnvNotSet(t, "KEYCLOAK_TEST_PASSWORD_GRANT")

	os.Setenv("KEYCLOAK_USER", "keycloak")
	os.Setenv("KEYCLOAK_PASSWORD", "password")

	defer func() {
		os.Unsetenv("KEYCLOAK_USER")
		os.Unsetenv("KEYCLOAK_PASSWORD")
	}()

	provider := KeycloakProvider()
	providerConfigureFunc := provider.ConfigureContextFunc
	// override the default context so a failing test doesn't cancel other tests
	provider.ConfigureContextFunc = func(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigureFunc(context.Background(), data)
	}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"keycloak": func() (*schema.Provider, error) {
				return provider, nil
			},
		},
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
			},
		},
	})
}
