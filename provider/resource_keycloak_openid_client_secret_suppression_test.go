package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeycloakOpenidClient_suppressSecretDuringAccessTypeDrift(t *testing.T) {
	clientID := acctest.RandomWithPrefix("tf-acc")
	config := fmt.Sprintf(`resource "keycloak_openid_client" "client" {
 realm_id = "%s"
 client_id = "%s"
 access_type = "PUBLIC"
 standard_flow_enabled = false
 implicit_flow_enabled = false
 direct_access_grants_enabled = false
 service_accounts_enabled = false
 suppress_client_secret = true
}`, testAccRealm.Realm, clientID)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("keycloak_openid_client.client", "access_type", "PUBLIC"),
		resource.TestCheckResourceAttr("keycloak_openid_client.client", "suppress_client_secret", "true"),
		resource.TestCheckResourceAttr("keycloak_openid_client.client", "client_secret", ""),
	)
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckKeycloakOpenidClientDestroy(),
		Steps: []resource.TestStep{
			{Config: config, Check: check},
			{RefreshState: true, ExpectNonEmptyPlan: true, PreConfig: func() {
				client, err := keycloakClient.GetOpenidClientByClientId(testCtx, testAccRealm.Realm, clientID)
				if err != nil {
					t.Fatal(err)
				}
				client.PublicClient = false
				if err := keycloakClient.UpdateOpenidClient(testCtx, client); err != nil {
					t.Fatal(err)
				}
				if _, err := keycloakClient.RegenerateOpenIdClientSecret(testCtx, client); err != nil {
					t.Fatal(err)
				}
			}, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("keycloak_openid_client.client", "access_type", "CONFIDENTIAL"),
				resource.TestCheckResourceAttr("keycloak_openid_client.client", "client_secret", ""),
			)},
			{Config: config, Check: check},
		},
	})
}
