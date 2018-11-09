package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakIdentityProviderMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	identityProviderMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProviderMapper_basic(realmName, identityProviderMapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperExists("keycloak_identity_provider_mapper.saml_mapper"),
			},
			{
				ResourceName:      "keycloak_identity_provider_mapper.saml_mapper",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericMapperImportId("keycloak_identity_provider_mapper.saml_mapper"),
			},
		},
	})
}

func TestAccKeycloakIdentityProviderMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	identityProviderMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProviderMapper_basic(realmName, identityProviderMapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperFetch("keycloak_identity_provider_mapper.saml_mapper", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakIdentityProviderMapper_basic(realmName, identityProviderMapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperFetch("keycloak_identity_provider_mapper.saml_mapper", mapper),
			},
		},
	})
}

func TestAccKeycloakIdentityProviderMapper_updateUserFederation(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProviderMapper_updateUserFederationBefore(realmOne, realmTwo, mapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperExists("keycloak_identity_provider_mapper.saml_mapper"),
			},
			{
				Config: testKeycloakIdentityProviderMapper_updateUserFederationAfter(realmOne, realmTwo, mapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperExists("keycloak_identity_provider_mapper.saml_mapper"),
			},
		},
	})
}

func testAccCheckKeycloakIdentityProviderMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakIdentityProviderMapperFetch(resourceName string, mapper *keycloak.IdentityProviderMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.Realm = fetchedMapper.Realm

		return nil
	}
}

func testAccCheckKeycloakIdentityProviderMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_identity_provider_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm"]
			alias := rs.Primary.Attributes["identity_provider_alias"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			IdentityProviderMapper, _ := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
			if IdentityProviderMapper != nil {
				return fmt.Errorf("identity provider mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getIdentityProviderMapperFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm"]
	alias := rs.Primary.Attributes["identity_provider_alias"]

	IdentityProviderMapper, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return nil, fmt.Errorf("error getting identity provider mapper with id %s: %s", id, err)
	}

	return IdentityProviderMapper, nil
}

func getGenericMapperImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		Realm := rs.Primary.Attributes["realm_id"]
		UserFederationId := rs.Primary.Attributes["_user_federation_id"]

		return fmt.Sprintf("%s/%s/%s", Realm, UserFederationId, id), nil
	}
}

func testKeycloakIdentityProviderMapper_basic(realm, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak__user_federation" "open" {
   name                    = "open"
   realm_id                = "${keycloak_realm.realm.id}"

   enabled                 = true

   username__attribute = "cn"
   rdn__attribute      = "cn"
   uuid__attribute     = "entryDN"
   user_object_classes     = [
      "simpleSecurityObject",
      "organizationalRole"
   ]
   connection_url          = "://open"
   users_dn                = "dc=example,dc=org"
   bind_dn                 = "cn=admin,dc=example,dc=org"
   bind_credential         = "admin"
}

resource "keycloak__full_name_mapper" "full_name_mapper" {
   name                     = "%s"
   realm_id                 = "${keycloak_realm.realm.id}"
   _user_federation_id  = "${keycloak__user_federation.open.id}"

   _full_name_attribute = "cn"
}
   `, realm, mapperName)
}

func testKeycloakIdentityProviderMapper_basicFromInterface(realm string, mapper *keycloak.IdentityProviderMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource keycloak_identity_provider saml {
  alias   = "saml"
  realm   = "${keycloak_realm.realm.name}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider_mapper saml {
  realm   = "${keycloak_realm.realm.name}"
  name = "%s"
  identity_provider_alias = "saml"
  identity_provider_mapper = "user-attribute-mapper"
  social {
    template = "asdasdasdadsdad"
  }
}
   `, realm, mapper.Name)
}

func testKeycloakIdentityProviderMapper_updateUserFederationBefore(realmOne, realmTwo, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_one" {
   realm = "%s"
}

resource "keycloak_realm" "realm_two" {
   realm = "%s"
}

resource "keycloak_identity_provider" "saml_one" {
  alias   = "saml"
  realm   = "${keycloak_realm.realm_one.name}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource "keycloak_identity_provider" "saml_two" {
  alias   = "saml"
  realm   = "${keycloak_realm.realm_two.name}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider_mapper saml_mapper {
  realm   = "${keycloak_realm.realm_one.id}"
  name = "%s"
  identity_provider_alias = "saml"
  identity_provider_mapper = "user-attribute-mapper"
  social {
    template = "asdasdasdadsdad"
  }
}
   `, realmOne, realmTwo, mapperName)
}

func testKeycloakIdentityProviderMapper_updateUserFederationAfter(realmOne, realmTwo, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_one" {
   realm = "%s"
}

resource "keycloak_realm" "realm_two" {
   realm = "%s"
}

resource "keycloak__user_federation" "open_one" {
   name                    = "open"
   realm_id                = "${keycloak_realm.realm_one.id}"

   enabled                 = true

   username__attribute = "cn"
   rdn__attribute      = "cn"
   uuid__attribute     = "entryDN"
   user_object_classes     = [
      "simpleSecurityObject",
      "organizationalRole"
   ]
   connection_url          = "://open"
   users_dn                = "dc=example,dc=org"
   bind_dn                 = "cn=admin,dc=example,dc=org"
   bind_credential         = "admin"
}

resource "keycloak__user_federation" "open_two" {
   name                    = "open"
   realm_id                = "${keycloak_realm.realm_two.id}"

   enabled                 = true

   username__attribute = "cn"
   rdn__attribute      = "cn"
   uuid__attribute     = "entryDN"
   user_object_classes     = [
      "simpleSecurityObject",
      "organizationalRole"
   ]
   connection_url          = "://open"
   users_dn                = "dc=example,dc=org"
   bind_dn                 = "cn=admin,dc=example,dc=org"
   bind_credential         = "admin"
}

resource "keycloak__full_name_mapper" "full_name_mapper" {
   name                     = "%s"
   realm_id                 = "${keycloak_realm.realm_two.id}"
   _user_federation_id  = "${keycloak__user_federation.open_two.id}"

   _full_name_attribute = "cn"
}
   `, realmOne, realmTwo, mapperName)
}

func testKeycloakIdentityProviderMapper_writableInvalid(realm, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak__user_federation" "open" {
   name                    = "open"
   realm_id                = "${keycloak_realm.realm.id}"

   enabled                 = true
   edit_mode               = "READ_ONLY"

   username__attribute = "cn"
   rdn__attribute      = "cn"
   uuid__attribute     = "entryDN"
   user_object_classes     = [
      "simpleSecurityObject",
      "organizationalRole"
   ]
   connection_url          = "://open"
   users_dn                = "dc=example,dc=org"
   bind_dn                 = "cn=admin,dc=example,dc=org"
   bind_credential         = "admin"
}

resource "keycloak__full_name_mapper" "full_name_mapper" {
   name                     = "%s"
   realm_id                 = "${keycloak_realm.realm.id}"
   _user_federation_id  = "${keycloak__user_federation.open.id}"

   _full_name_attribute = "cn"
   write_only               = true
}
   `, realm, mapperName)
}

func testKeycloakIdentityProviderMapper_writableValid(realm, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
   realm = "%s"
}

resource "keycloak__user_federation" "open" {
   name                    = "open"
   realm_id                = "${keycloak_realm.realm.id}"

   enabled                 = true
   edit_mode               = "WRITABLE"

   username__attribute = "cn"
   rdn__attribute      = "cn"
   uuid__attribute     = "entryDN"
   user_object_classes     = [
      "simpleSecurityObject",
      "organizationalRole"
   ]
   connection_url          = "://open"
   users_dn                = "dc=example,dc=org"
   bind_dn                 = "cn=admin,dc=example,dc=org"
   bind_credential         = "admin"
}

resource "keycloak__full_name_mapper" "full_name_mapper" {
   name                     = "%s"
   realm_id                 = "${keycloak_realm.realm.id}"
   _user_federation_id  = "${keycloak__user_federation.open.id}"

   _full_name_attribute = "cn"
   write_only               = true
}
   `, realm, mapperName)
}
