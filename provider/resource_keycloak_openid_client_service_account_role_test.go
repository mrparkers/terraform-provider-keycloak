package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
	"testing"
)

func TestAccKeycloakOpenidClientServiceAccountRole_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	resourceName := "keycloak_openid_client_service_account_role.test"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientServiceAccountRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountRoleExists(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getKeycloakOpenidClientServiceAccountRoleImportId(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenidClientServiceAccountRole_createAfterManualDestroy(t *testing.T) {
	var serviceAccountRole = &keycloak.OpenidClientServiceAccountRole{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientServiceAccountRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountRoleFetch("keycloak_openid_client_service_account_role.test", serviceAccountRole),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClientServiceAccountRole(serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ContainerId, serviceAccountRole.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientServiceAccountRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountRoleExists("keycloak_openid_client_service_account_role.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientServiceAccountRole_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientServiceAccountRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountRole_basic(firstRealm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientServiceAccountRoleExists("keycloak_openid_client_service_account_role.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_service_account_role.test", "realm_id", firstRealm),
				),
			},
			{
				Config: testKeycloakOpenidClientServiceAccountRole_basic(secondRealm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientServiceAccountRoleExists("keycloak_openid_client_service_account_role.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_service_account_role.test", "realm_id", secondRealm),
				),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientServiceAccountRoleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakOpenidClientServiceAccountRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientServiceAccountRoleFetch(resourceName string, serviceAccountRole *keycloak.OpenidClientServiceAccountRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedServiceAccountRole, err := getKeycloakOpenidClientServiceAccountRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		serviceAccountRole.ServiceAccountUserId = fetchedServiceAccountRole.ServiceAccountUserId
		serviceAccountRole.RealmId = fetchedServiceAccountRole.RealmId
		serviceAccountRole.ClientRole = fetchedServiceAccountRole.ClientRole
		serviceAccountRole.ContainerId = fetchedServiceAccountRole.ContainerId
		serviceAccountRole.Id = fetchedServiceAccountRole.Id

		return nil
	}
}

func testAccCheckKeycloakOpenidClientServiceAccountRoleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_service_account_role" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
			clientId := rs.Primary.Attributes["client_id"]
			id := strings.Split(rs.Primary.ID, "/")[1]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			serviceAccountRole, _ := keycloakClient.GetOpenidClientServiceAccountRole(realmId, serviceAccountUserId, clientId, id)
			if serviceAccountRole != nil {
				return fmt.Errorf("service account role exists")
			}
		}

		return nil
	}
}

func getKeycloakOpenidClientServiceAccountRoleFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientServiceAccountRole, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
	clientId := rs.Primary.Attributes["client_id"]
	id := strings.Split(rs.Primary.ID, "/")[1]

	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountRole(realmId, serviceAccountUserId, clientId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting service account role mapping: %s", err)
	}

	return serviceAccountRole, nil
}

func getKeycloakOpenidClientServiceAccountRoleImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		realmId := rs.Primary.Attributes["realm_id"]
		serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
		clientId := rs.Primary.Attributes["client_id"]
		roleId := strings.Split(rs.Primary.ID, "/")[1]

		return fmt.Sprintf("%s/%s/%s/%s", realmId, serviceAccountUserId, clientId, roleId), nil
	}
}

func testKeycloakOpenidClientServiceAccountRole_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource keycloak_realm test {
	realm = "%s"
}

resource keycloak_openid_client test {
	client_id                = "%s"
	realm_id                 = "${keycloak_realm.test.id}"
	access_type              = "CONFIDENTIAL"
	service_accounts_enabled = true
}

data keycloak_openid_client broker {
  realm_id  = "${keycloak_realm.test.id}"
  client_id = "broker"
}

resource keycloak_openid_client_service_account_role test {
	service_account_user_id = "${keycloak_openid_client.test.service_account_user_id}"
	realm_id 					= "${keycloak_realm.test.id}"
	client_id 					= "${data.keycloak_openid_client.broker.id}"
	role 							= "read-token"
}
	`, realm, clientId)
}
