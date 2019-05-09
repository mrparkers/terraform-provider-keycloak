package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakOpenIdAudienceProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper_client"
	clientScopeResourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdAudienceProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdAudienceProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdAudienceProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	customAudience := "terraform-audience-" + acctest.RandString(10)
	updatedCustomAudience := "terraform-audience-" + acctest.RandString(10)
	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_customAudience(realmName, clientId, mapperName, customAudience),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_customAudience(realmName, clientId, mapperName, updatedCustomAudience),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdAudienceProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdAudienceProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdAudienceProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdAudienceProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	customAudience := "terraform-audience-" + acctest.RandString(10)
	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_customAudience(realmName, clientId, mapperName, customAudience),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_customAudience(realmName, updatedClientId, mapperName, customAudience),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	customAudience := "terraform-audience-" + acctest.RandString(10)
	resourceName := "keycloak_openid_audience_protocol_mapper.audience_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_customAudience(realmName, clientId, mapperName, customAudience),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdAudienceProtocolMapper_customAudience(newRealmName, clientId, mapperName, customAudience),
				Check:  testKeycloakOpenIdAudienceProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_validateClientConflictsWithClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdAudienceProtocolMapper_validateClientConflictsWithClientScope(realmName, clientId, clientScopeId, mapperName),
				ExpectError: regexp.MustCompile(".+ conflicts with .+"),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_validateClientAudienceConflictsWithCustomAudience(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdAudienceProtocolMapper_validateClientAudienceConflictsWithCustomAudience(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile(".+ conflicts with .+"),
			},
		},
	})
}

func TestAccKeycloakOpenIdAudienceProtocolMapper_validateClientAudienceExists(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-audience-mapper-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdAudienceProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdAudienceProtocolMapper_validateClientAudienceExists(realmName, clientId, mapperName),
				ExpectError: regexp.MustCompile("validation error: client .+ does not exist"),
			},
		},
	})
}

func testAccKeycloakOpenIdAudienceProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_audience_protocol_mapper" {
				continue
			}

			mapper, _ := getAudienceMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid audience protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdAudienceProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getAudienceMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdAudienceProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdAudienceProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getAudienceMapperUsingState(state, resourceName)
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

func getAudienceMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdAudienceProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdAudienceProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdAudienceProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper_client" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_id                = "${keycloak_openid_client.openid_client.id}"

	included_custom_audience = "foo"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdAudienceProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper_client_scope" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_scope_id          = "${keycloak_openid_client_scope.client_scope.id}"

	included_custom_audience = "foo"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdAudienceProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper_client" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_id                = "${keycloak_openid_client.openid_client.id}"

	included_custom_audience = "foo"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper_client_scope" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_scope_id          = "${keycloak_openid_client_scope.client_scope.id}"

	included_custom_audience = "foo"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdAudienceProtocolMapper_customAudience(realmName, clientId, mapperName, customAudience string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "CONFIDENTIAL"

	standard_flow_enabled = true

	valid_redirect_uris = ["http://localhost:5555/callback"]
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_id                = "${keycloak_openid_client.openid_client.id}"

	included_custom_audience = "%s"
}`, realmName, clientId, mapperName, customAudience)
}

func testKeycloakOpenIdAudienceProtocolMapper_validateClientConflictsWithClientScope(realmName, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_id                = "${keycloak_openid_client.openid_client.id}"
	client_scope_id          = "${keycloak_openid_client_scope.client_scope.id}"

	included_custom_audience = "foo"
}`, realmName, clientId, clientScopeId, mapperName)
}

func testKeycloakOpenIdAudienceProtocolMapper_validateClientAudienceConflictsWithCustomAudience(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_id                = "${keycloak_openid_client.openid_client.id}"

	included_client_audience = "${keycloak_openid_client.openid_client.client_id}"
	included_custom_audience = "foo"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdAudienceProtocolMapper_validateClientAudienceExists(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "openid-client"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"
	client_id                = "${keycloak_openid_client.openid_client.id}"

	included_client_audience = "%s"

	depends_on = [ "keycloak_openid_client.openid_client" ]
}`, realmName, mapperName, clientId)
}
