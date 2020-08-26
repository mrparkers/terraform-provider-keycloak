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

func TestAccKeycloakOpenidClientServiceAccountRealmRole_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientServiceAccountRealmRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountRealmRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountRealmRoleExists("keycloak_openid_client_service_account_realm_role.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientServiceAccountRealmRole_createAfterManualDestroy(t *testing.T) {
	var serviceAccountRole = &keycloak.OpenidClientServiceAccountRealmRole{}

	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientServiceAccountRealmRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountRealmRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountRealmRoleFetch("keycloak_openid_client_service_account_realm_role.test", serviceAccountRole),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenidClientServiceAccountRealmRole(serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientServiceAccountRealmRole_basic(realmName, clientId),
				Check:  testAccCheckKeycloakOpenidClientServiceAccountRealmRoleExists("keycloak_openid_client_service_account_realm_role.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientServiceAccountRealmRole_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOpenidClientServiceAccountRealmRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientServiceAccountRealmRole_basic(firstRealm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientServiceAccountRealmRoleExists("keycloak_openid_client_service_account_realm_role.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_service_account_realm_role.test", "realm_id", firstRealm),
				),
			},
			{
				Config: testKeycloakOpenidClientServiceAccountRealmRole_basic(secondRealm, clientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientServiceAccountRealmRoleExists("keycloak_openid_client_service_account_realm_role.test"),
					resource.TestCheckResourceAttr("keycloak_openid_client_service_account_realm_role.test", "realm_id", secondRealm),
				),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientServiceAccountRealmRoleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakOpenidClientServiceAccountRealmRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientServiceAccountRealmRoleFetch(resourceName string, serviceAccountRole *keycloak.OpenidClientServiceAccountRealmRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedServiceAccountRole, err := getKeycloakOpenidClientServiceAccountRealmRoleFromState(s, resourceName)
		if err != nil {
			return err
		}

		serviceAccountRole.ServiceAccountUserId = fetchedServiceAccountRole.ServiceAccountUserId
		serviceAccountRole.RealmId = fetchedServiceAccountRole.RealmId
		serviceAccountRole.Id = fetchedServiceAccountRole.Id

		return nil
	}
}

func testAccCheckKeycloakOpenidClientServiceAccountRealmRoleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_service_account_realm_role" {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
			id := strings.Split(rs.Primary.ID, "/")[1]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			serviceAccountRole, _ := keycloakClient.GetOpenidClientServiceAccountRealmRole(realmId, serviceAccountUserId, id)
			if serviceAccountRole != nil {
				return fmt.Errorf("service account role exists")
			}
		}

		return nil
	}
}

func getKeycloakOpenidClientServiceAccountRealmRoleFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientServiceAccountRealmRole, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	serviceAccountUserId := rs.Primary.Attributes["service_account_user_id"]
	id := strings.Split(rs.Primary.ID, "/")[1]

	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountRealmRole(realmId, serviceAccountUserId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting service account role mapping: %s", err)
	}

	return serviceAccountRole, nil
}

func testKeycloakOpenidClientServiceAccountRealmRole_basic(realm, clientId string) string {
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
