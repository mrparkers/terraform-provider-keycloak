package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_basicClient(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper_client"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_basicClientScope(t *testing.T) {
	t.Parallel()
	clientScopeId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_clientScope(clientScopeId, mapperName),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_import(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	clientScopeId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	clientResourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper_client"
	clientScopeResourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdFullNameProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_import(clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_update(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	claimName := acctest.RandomWithPrefix("tf-acc")
	updatedClaimName := acctest.RandomWithPrefix("tf-acc")
	claimValue := acctest.RandomWithPrefix("tf-acc")
	updatedClaimValue := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(clientId, mapperName, updatedClaimName, updatedClaimValue),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	var mapper = &keycloak.OpenIdPropertyMapperClaimProtocolMapper{}

	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper_client"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteOpenIdPropertyMapperClaimProtocolMapper(testCtx, mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_validateClaimValueType(t *testing.T) {
	t.Parallel()
	mapperName := acctest.RandomWithPrefix("tf-acc")
	invalidClaimValueType := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdPropertyMapperClaimProtocolMapper_validateClaimValueType(mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_updateClientIdForceNew(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	updatedClientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	claimName := acctest.RandomWithPrefix("tf-acc")
	claimValue := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(updatedClientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	t.Parallel()
	mapperName := acctest.RandomWithPrefix("tf-acc")
	clientScopeId := acctest.RandomWithPrefix("tf-acc")
	newClientScopeId := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_clientScope(clientScopeId, mapperName),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_clientScope(newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdPropertyMapperClaimProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")
	mapperName := acctest.RandomWithPrefix("tf-acc")

	claimName := acctest.RandomWithPrefix("tf-acc")
	claimValue := acctest.RandomWithPrefix("tf-acc")
	resourceName := "keycloak_openid_propertymapper_claim_protocol_mapper.propertymapper_claim_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(clientId, mapperName, claimName, claimValue),
				Check:  testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdPropertyMapperClaimProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_propertymapper_claim_protocol_mapper" {
				continue
			}

			mapper, _ := getPropertyMapperClaimMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getPropertyMapperClaimMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdPropertyMapperClaimProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getPropertyMapperClaimMapperUsingState(state, resourceName)
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

func getPropertyMapperClaimMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdPropertyMapperClaimProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	return keycloakClient.GetOpenIdPropertyMapperClaimProtocolMapper(testCtx, realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_client(clientId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "propertymapper_claim_mapper_client" {
	name             = "%s"
	realm_id         = data.keycloak_realm.realm.id
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, testAccRealm.Realm, clientId, mapperName)
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapper_basic_clientScope(clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "propertymapper_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = data.keycloak_realm.realm.id
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, testAccRealm.Realm, clientScopeId, mapperName)
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapper_import(clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = data.keycloak_realm.realm.id
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "propertymapper_claim_mapper_client" {
	name             = "%s"
	realm_id         = data.keycloak_realm.realm.id
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "propertymapper_claim_mapper_client_scope" {
	name             = "%s"
	realm_id         = data.keycloak_realm.realm.id
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"

	claim_name       = "foo"
	claim_value      = "bar"
	claim_value_type = "String"
}`, testAccRealm.Realm, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapper_claimNameAndValue(clientId, mapperName, claimName, claimValue string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "propertymapper_claim_mapper" {
	name        = "%s"
	realm_id    = data.keycloak_realm.realm.id
	client_id   = "${keycloak_openid_client.openid_client.id}"

	claim_name  = "%s"
	claim_value = "%s"
}`, testAccRealm.Realm, clientId, mapperName, claimName, claimValue)
}

func testKeycloakOpenIdPropertyMapperClaimProtocolMapper_validateClaimValueType(mapperName, claimValueType string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "openid-client"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "propertymapper_claim_mapper_validation" {
	name             = "%s"
	realm_id         = data.keycloak_realm.realm.id
	client_id        = "${keycloak_openid_client.openid_client.id}"

	claim_value      = "foo"
	claim_name       = "bar"
	claim_value_type = "%s"
}`, testAccRealm.Realm, mapperName, claimValueType)
}
