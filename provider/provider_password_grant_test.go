package provider

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"keycloak": provider,
		},
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClient_basic(realmName, clientId),
			},
		},
	})
}
