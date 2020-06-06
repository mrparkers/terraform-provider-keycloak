package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestGenericRoleMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	parentClientName := "client1-" + acctest.RandString(10)
	parentRoleName := "role-" + acctest.RandString(10)
	childClientName := "client2-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check:  testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.child-client-with-parent-client-role"),
			},
		},
	})
}

func TestGenericRoleMapper_import(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	parentClientName := "client1-" + acctest.RandString(10)
	parentRoleName := "role-" + acctest.RandString(10)
	childClientName := "client2-" + acctest.RandString(10)

	resourceName := "keycloak_generic_client_role_mapper.child-client-with-parent-client-role"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check:  testAccCheckKeycloakScopeMappingExists(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericRoleMapperId(resourceName),
			},
		},
	})
}

func TestGenericRoleMapperClientScope_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientName := "client-" + acctest.RandString(10)
	roleName := "role-" + acctest.RandString(10)
	clientScopeName := "clientscope-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMappingClientScope_basic(realmName, clientName, roleName, clientScopeName),
				Check:  testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.clientscope-with-client-role"),
			},
		},
	})
}

func TestGenericRoleMapperClientScope_import(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientName := "client-" + acctest.RandString(10)
	roleName := "role-" + acctest.RandString(10)
	clientScopeName := "clientscope-" + acctest.RandString(10)

	resourceName := "keycloak_generic_client_role_mapper.clientscope-with-client-role"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMappingClientScope_basic(realmName, clientName, roleName, clientScopeName),
				Check:  testAccCheckKeycloakScopeMappingExists(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericRoleMapperId(resourceName),
			},
		},
	})
}

func TestGenericRoleMapper_createAfterManualDestroy(t *testing.T) {
	var role = &keycloak.Role{}
	var childClient = &keycloak.GenericClient{}

	realmName := "terraform-" + acctest.RandString(10)
	parentClientName := "client1-" + acctest.RandString(10)
	parentRoleName := "role-" + acctest.RandString(10)
	childClientName := "client2-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.child-client-with-parent-client-role"),
					testAccCheckKeycloakRoleFetch("keycloak_role.parent-role", role),
					testAccCheckKeycloakGenericClientFetch("keycloak_openid_client.child-client", childClient),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteRoleScopeMapping(childClient.RealmId, childClient.Id, "", role)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check:  testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.child-client-with-parent-client-role"),
			},
		},
	})
}

func TestGenericRoleMapperClientScope_createAfterManualDestroy(t *testing.T) {
	var role = &keycloak.Role{}
	var clientScope = &keycloak.OpenidClientScope{}

	realmName := "terraform-" + acctest.RandString(10)
	clientName := "client-" + acctest.RandString(10)
	roleName := "role-" + acctest.RandString(10)
	clientScopeName := "clientscope-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMappingClientScope_basic(realmName, clientName, roleName, clientScopeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.clientscope-with-client-role"),
					testAccCheckKeycloakRoleFetch("keycloak_role.role", role),
					testAccCheckKeycloakOpenidClientScopeFetch("keycloak_openid_client_scope.clientscope", clientScope),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteRoleScopeMapping(clientScope.RealmId, "", clientScope.Id, role)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakGenericRoleMappingClientScope_basic(realmName, clientName, roleName, clientScopeName),
				Check:  testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.clientscope-with-client-role"),
			},
		},
	})
}

func testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "parent-client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "PUBLIC"
}

resource "keycloak_role" "parent-role" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "${keycloak_openid_client.parent-client.id}"
  name      = "%s"
}

resource "keycloak_openid_client" "child-client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "PUBLIC"
}

resource "keycloak_generic_client_role_mapper" "child-client-with-parent-client-role" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "${keycloak_openid_client.child-client.id}"
  role_id   = "${keycloak_role.parent-role.id}"
}
	`, realmName, parentClientName, parentRoleName, childClientName)
}

func testKeycloakGenericRoleMappingClientScope_basic(realmName, clientName, roleName, clientScopeName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "PUBLIC"
}

resource "keycloak_role" "role" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "${keycloak_openid_client.client.id}"
  name      = "%s"
}

resource "keycloak_openid_client_scope" "clientscope" {
	realm_id    = "${keycloak_realm.realm.id}"
	name        = "%s"
}

resource "keycloak_generic_client_role_mapper" "clientscope-with-client-role" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_scope_id = "${keycloak_openid_client_scope.clientscope.id}"
  role_id   = "${keycloak_role.role.id}"
}
	`, realmName, clientName, roleName, clientScopeName)
}

func testAccCheckKeycloakScopeMappingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		return nil
	}
}

func testAccCheckKeycloakGenericClientFetch(resourceName string, client *keycloak.GenericClient) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClient, err := getGenericClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		client.Id = fetchedClient.Id
		client.ClientId = fetchedClient.ClientId
		client.RealmId = fetchedClient.RealmId

		return nil
	}
}

func testAccCheckKeycloakOpenidClientScopeFetch(resourceName string, clientScope *keycloak.OpenidClientScope) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClientScope, err := getOpenidClientScopeFromState(s, resourceName)
		if err != nil {
			return err
		}

		clientScope.Id = fetchedClientScope.Id
		clientScope.RealmId = fetchedClientScope.RealmId

		return nil
	}
}

func getGenericClientFromState(s *terraform.State, resourceName string) (*keycloak.GenericClient, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetGenericClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting generic client %s: %s", id, err)
	}

	return client, nil
}

func getOpenidClientScopeFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientScope, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetOpenidClientScope(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting client scope %s: %s", id, err)
	}

	return client, nil
}

func getGenericRoleMapperId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
