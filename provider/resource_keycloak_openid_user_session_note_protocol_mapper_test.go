package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper_client"
	clientScopeResourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdUserSessionNoteProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdUserSessionNoteProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_updateClaim(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	updatedClaimName := "claim-name-update-" + acctest.RandString(10)

	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(realmName, clientId, mapperName, updatedClaimName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_updateLabel(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	labelName := "session-note-label-" + acctest.RandString(10)
	updatedLabelName := "session-note-label-update-" + acctest.RandString(10)

	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_label(realmName, clientId, mapperName, labelName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_label(realmName, clientId, mapperName, updatedLabelName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdUserSessionNoteProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdUserSessionNoteProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserSessionNoteProtocolMapper_validateClaimValueType(realmName, mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(realmName, updatedClientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserSessionNoteProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-session-note-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_session_note_protocol_mapper.user_session_note_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(newRealmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdUserSessionNoteProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_user_session_note_protocol_mapper" {
				continue
			}

			mapper, _ := getUserSessionNoteMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdUserSessionNoteProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getUserSessionNoteMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdUserSessionNoteProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdUserSessionNoteProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getUserSessionNoteMapperUsingState(state, resourceName)
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

func getUserSessionNoteMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdUserSessionNoteProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdUserSessionNoteProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper_client" {
	name               = "%s"
	realm_id           = "${keycloak_realm.realm.id}"
	client_id          = "${keycloak_openid_client.openid_client.id}"
	claim_name         = "foo"
	claim_value_type   = "String"
	session_note_label = "bar"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdUserSessionNoteProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper_client_scope" {
	name               = "%s"
	realm_id           = "${keycloak_realm.realm.id}"
	client_scope_id    = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name         = "foo"
	claim_value_type   = "String"
	session_note_label = "bar"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserSessionNoteProtocolMapper_claim(realmName, clientId, mapperName, claimName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"
	claim_name       = "%s"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, claimName)
}

func testKeycloakOpenIdUserSessionNoteProtocolMapper_label(realmName, clientId, mapperName, labelName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper" {
	name               = "%s"
	realm_id           = "${keycloak_realm.realm.id}"
	client_id          = "${keycloak_openid_client.openid_client.id}"
	claim_name         = "foo"
	claim_value_type   = "String"
	session_note_label = "%s"
}`, realmName, clientId, mapperName, labelName)
}

func testKeycloakOpenIdUserSessionNoteProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper_client" {
	name               = "%s"
	realm_id           = "${keycloak_realm.realm.id}"
	client_id          = "${keycloak_openid_client.openid_client.id}"
	claim_name         = "foo"
	claim_value_type   = "String"
	session_note_label = "bar"
}
resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper_client_scope" {
	name               = "%s"
	realm_id           = "${keycloak_realm.realm.id}"
	client_scope_id    = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name         = "foo"
	claim_value_type   = "String"
	session_note_label = "bar"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserSessionNoteProtocolMapper_validateClaimValueType(realmName, mapperName, claimValueType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "openid-client"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_session_note_protocol_mapper" "user_session_note_mapper_validation" {
	name               = "%s"
	realm_id           = "${keycloak_realm.realm.id}"
	client_id          = "${keycloak_openid_client.openid_client.id}"
	claim_name         = "foo"
	claim_value_type   = "%s"
	session_note_label = "bar"
}`, realmName, mapperName, claimValueType)
}
