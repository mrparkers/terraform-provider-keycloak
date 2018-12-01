package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakOpenIdFullNameProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-full-name-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-full-name-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-full-name-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client"
	clientScopeResourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdFullNameProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdFullNameProtocolMapperExists(clientScopeResourceName),
				),
			},
			{
				ResourceName:      clientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(clientResourceName),
			},
			{
				ResourceName:      clientScopeResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClientScope(clientScopeResourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_update(t *testing.T) {
	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper"

	mapperOne := &keycloak.OpenIdFullNameProtocolMapper{
		Name:             acctest.RandString(10),
		RealmId:          "terraform-realm-" + acctest.RandString(10),
		ClientId:         "terraform-client-" + acctest.RandString(10),
		AddToIdToken:     randomBool(),
		AddToAccessToken: randomBool(),
		AddToUserInfo:    randomBool(),
	}

	mapperTwo := &keycloak.OpenIdFullNameProtocolMapper{
		Name:             mapperOne.Name,
		RealmId:          mapperOne.RealmId,
		ClientId:         mapperOne.ClientId,
		AddToIdToken:     randomBool(),
		AddToAccessToken: randomBool(),
		AddToUserInfo:    randomBool(),
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_fromInterface(mapperOne),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_fromInterface(mapperTwo),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdFullNameProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-full-name-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdFullNameProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdUserAttributeProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdFullNameProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_updateMapperNameForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperNameOne := acctest.RandString(10)
	mapperNameTwo := acctest.RandString(10)

	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_basic_client(realmName, clientId, mapperNameOne),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_basic_client(realmName, clientId, mapperNameTwo),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientIdOne := "terraform-client-" + acctest.RandString(10)
	clientIdTwo := "terraform-client-" + acctest.RandString(10)

	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_updateClientForceNew(realmName, clientIdOne, clientIdTwo, "openid_client_one"),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_updateClientForceNew(realmName, clientIdOne, clientIdTwo, "openid_client_two"),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdFullNameProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeOne := "terraform-client-scope-" + acctest.RandString(10)
	clientScopeTwo := "terraform-client-scope-" + acctest.RandString(10)

	resourceName := "keycloak_openid_full_name_protocol_mapper.full_name_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_updateClientScopeForceNew(realmName, clientScopeOne, clientScopeTwo, "client_scope_one"),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdFullNameProtocolMapper_updateClientScopeForceNew(realmName, clientScopeOne, clientScopeTwo, "client_scope_two"),
				Check:  testKeycloakOpenIdFullNameProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdFullNameProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_full_name_protocol_mapper" {
				continue
			}

			mapper, _ := getFullNameMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdFullNameProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getFullNameMapperUsingState(state, resourceName)

		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdFullNameProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdFullNameProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getFullNameMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.ClientId = fetchedMapper.ClientId
		mapper.ClientScopeId = fetchedMapper.ClientScopeId
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func getFullNameMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdFullNameProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdFullNameProtocolMapper(realm, clientId, clientScopeId, id)
}

func getGenericProtocolMapperIdForClient(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		clientId := rs.Primary.Attributes["client_id"]

		return fmt.Sprintf("%s/client/%s/%s", realmId, clientId, id), nil
	}
}

func getGenericProtocolMapperIdForClientScope(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		clientScopeId := rs.Primary.Attributes["client_scope_id"]

		return fmt.Sprintf("%s/client-scope/%s/%s", realmId, clientScopeId, id), nil
	}
}

func testKeycloakOpenIdFullNameProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client" {
	name       = "%s"
	realm_id   = "${keycloak_realm.realm.id}"
	client_id  = "${keycloak_openid_client.openid_client.id}"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdFullNameProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdFullNameProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client" {
	name       = "%s"
	realm_id   = "${keycloak_realm.realm.id}"
	client_id  = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client_scope" {
	name            = "%s"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdFullNameProtocolMapper_fromInterface(mapper *keycloak.OpenIdFullNameProtocolMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper" {
	name                = "%s"
	realm_id            = "${keycloak_realm.realm.id}"
	client_id           = "${keycloak_openid_client.openid_client.id}"

	add_to_id_token     = %t
	add_to_access_token = %t
	add_to_userinfo     = %t
}`, mapper.RealmId, mapper.ClientId, mapper.Name, mapper.AddToIdToken, mapper.AddToAccessToken, mapper.AddToUserInfo)
}

func testKeycloakOpenIdFullNameProtocolMapper_updateClientForceNew(realmId, clientIdOne, clientIdTwo, currentClient string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client_one" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_client" "openid_client_two" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client" {
	name       = "group-mapper"
	realm_id   = "${keycloak_realm.realm.id}"
	client_id  = "${keycloak_openid_client.%s.id}"
}`, realmId, clientIdOne, clientIdTwo, currentClient)
}

func testKeycloakOpenIdFullNameProtocolMapper_updateClientScopeForceNew(realmId, clientScopeIdOne, clientScopeIdTwo, currentClientScope string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope_one" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_client_scope" "client_scope_two" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_full_name_protocol_mapper" "full_name_mapper_client_scope" {
	name            = "group-mapper"
	realm_id        = "${keycloak_realm.realm.id}"
	client_scope_id = "${keycloak_openid_client_scope.%s.id}"
}`, realmId, clientScopeIdOne, clientScopeIdTwo, currentClientScope)
}
