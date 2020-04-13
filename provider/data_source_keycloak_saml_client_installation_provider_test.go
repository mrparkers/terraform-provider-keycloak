package provider

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKeycloakSamlClientInstallationProvider_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resourceName := "keycloak_saml_client.saml_client"
	dataSourceName := "data.keycloak_saml_client_installation_provider.saml_sp_descriptor"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakSamlClientDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakSamlClientInstallationProvider_basic(realmName, clientId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_id", resourceName, "realm_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "client_id", resourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_id", "saml-sp-descriptor"),
					testAccCheckDataKeycloakSamlClientInstallationProvider(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakSamlClientInstallationProvider(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		value := rs.Primary.Attributes["value"]

		err := xml.Unmarshal([]byte(value), new(interface{}))
		if err != nil {
			return fmt.Errorf("invalid XML: %s\n%s", err, value)
		}

		return nil
	}
}

func testDataSourceKeycloakSamlClientInstallationProvider_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

data "keycloak_saml_client_installation_provider" "saml_sp_descriptor" {
  realm_id    = "${keycloak_realm.realm.id}"
  client_id   = "${keycloak_saml_client.saml_client.id}"
  provider_id = "saml-sp-descriptor"
}
	`, realm, clientId)
}
