package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
	"testing"
)

func TestAccKeycloakOpenidClientServiceAccountClientRole_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientServiceAccountClientRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountClientRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountClientRoleExists("keycloak_openid_client_service_account_realm_role.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientServiceAccountClientRole_createAfterManualDestroy(t *testing.T) {
	var serviceAccountRole = &keycloak.OpenidClientServiceAccountClientRole{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientServiceAccountClientRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountClientRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountClientRoleFetch("keycloak_openid_client_service_account_realm_role.test", serviceAccountRole),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClientServiceAccountClientRole(serviceAccountRole.RealmId,serviceAccountRole.ClientId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientServiceAccountClientRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountClientRoleExists("keycloak_openid_client_service_account_realm_role.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientServiceAccountClientRole_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakOpenidClientServiceAccountClientRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountClientRole_basic(firstRealm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientServiceAccountClientRoleExists("keycloak_openid_client_service_account_realm_role.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_service_account_realm_role.test", "realm_id", firstRealm),
				),
			},
			{
				Config: testKeycloakOpenidClientServiceAccountClientRole_basic(secondRealm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientServiceAccountClientRoleExists("keycloak_openid_client_service_account_realm_role.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_service_account_realm_role.test", "realm_id", secondRealm),
				),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientServiceAccountClientRoleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakOpenidClientServiceAccountClientRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientServiceAccountClientRoleFetch(resourceName string, serviceAccountRole *keycloak.OpenidClientServiceAccountClientRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedServiceAccountRole, err := getKeycloakOpenidClientServiceAccountClientRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		serviceAccountRole.ServiceAccountUserId = fetchedServiceAccountRole.ServiceAccountUserId
		serviceAccountRole.RealmId = fetchedServiceAccountRole.RealmId
		serviceAccountRole.ClientId = fetchedServiceAccountRole.ClientId
		serviceAccountRole.Id = fetchedServiceAccountRole.Id

		return nil
	}
}

func testAccCheckKeycloakOpenidClientServiceAccountClientRoleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_service_account_realm_role" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			clientId := rs.Primary.Attributes["client_id"]
			serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
			id := strings.Split(rs.Primary.ID, "/")[1]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			serviceAccountRole, _ := keycloakClient.GetOpenidClientServiceAccountClientRole(realmId, clientId, serviceAccountUserId, id)
			if serviceAccountRole != nil {
				return fmt.Errorf("service account role exists")
			}
		}

		return nil
	}
}

func getKeycloakOpenidClientServiceAccountClientRoleFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientServiceAccountClientRole, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
	id := strings.Split(rs.Primary.ID, "/")[1]

	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountClientRole(realmId, clientId, serviceAccountUserId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting service account role mapping: %s", err)
	}

	return serviceAccountRole, nil
}

func testKeycloakOpenidClientServiceAccountClientRole_basic(realm, clientId string) string {
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

resource keycloak_openid_client_service_account_realm_role test {
	service_account_user_id = "${keycloak_openid_client.test.service_account_user_id}"
	realm_id 					= "${keycloak_realm.test.id}"
	role 						= "offline_access"
}
	`, realm, clientId)
}
