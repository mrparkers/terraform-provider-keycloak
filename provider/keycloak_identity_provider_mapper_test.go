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
	aliasName := "terraform-" + acctest.RandString(10)
	identityProviderMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProviderMapper_basic(realmName, aliasName, identityProviderMapperName),
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
	aliasName := "terraform-" + acctest.RandString(10)
	identityProviderMapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProviderMapper_basic(realmName, aliasName, identityProviderMapperName),
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
				Config: testKeycloakIdentityProviderMapper_basic(realmName, aliasName, identityProviderMapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperFetch("keycloak_identity_provider_mapper.saml_mapper", mapper),
			},
		},
	})
}

func TestAccKeycloakIdentityProviderMapper_updateUserFederation(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	aliasOne := "terraform-" + acctest.RandString(10)
	aliasTwo := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakIdentityProviderMapper_updateUserFederationBefore(realmOne, realmTwo, aliasOne, aliasTwo, mapperName),
				Check:  testAccCheckKeycloakIdentityProviderMapperExists("keycloak_identity_provider_mapper.saml_mapper"),
			},
			{
				Config: testKeycloakIdentityProviderMapper_updateUserFederationAfter(realmOne, realmTwo, aliasOne, aliasTwo, mapperName),
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
		mapper.IdentityProviderAlias = fetchedMapper.IdentityProviderAlias
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
		realm := rs.Primary.Attributes["realm"]
		aliasName := rs.Primary.Attributes["identity_provider_alias"]

		return fmt.Sprintf("%s/%s/%s/", realm, aliasName, id), nil
	}
}

func testKeycloakIdentityProviderMapper_basic(realm, alias, mapperName string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm {
   realm = "%s"
}

resource keycloak_identity_provider saml {
  alias   = "%s"
  realm   = "${keycloak_realm.realm.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider_mapper saml_mapper {
  realm                    = "${keycloak_realm.realm.realm}"
  name                     = "%s"
  identity_provider_alias  = "${keycloak_identity_provider.saml.alias}"
  identity_provider_mapper = "user-attribute-mapper"

  saml {
    template = "asdasdasdadsdad"
  }
}
   `, realm, alias, mapperName)
}

func testKeycloakIdentityProviderMapper_updateUserFederationBefore(realmOne, realmTwo, aliasOne, aliasTwo, mapperName string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm_one {
   realm = "%s"
}

resource keycloak_realm realm_two {
   realm = "%s"
}

resource keycloak_identity_provider saml_one {
  alias   = "%s"
  realm   = "${keycloak_realm.realm_one.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider saml_two {
  alias   = "%s"
  realm   = "${keycloak_realm.realm_two.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider_mapper saml_mapper {
  realm                    = "${keycloak_realm.realm_one.realm}"
  name                     = "%s"
  identity_provider_alias  = "${keycloak_identity_provider.saml_one.alias}"
  identity_provider_mapper = "user-attribute-mapper"

  saml {
    template = "asdasdasdadsdad"
  }
}
   `, realmOne, realmTwo, aliasOne, aliasTwo, mapperName)
}

func testKeycloakIdentityProviderMapper_updateUserFederationAfter(realmOne, realmTwo, aliasOne, aliasTwo, mapperName string) string {
	return fmt.Sprintf(`
resource keycloak_realm realm_one {
   realm = "%s"
}

resource keycloak_realm realm_two {
   realm = "%s"
}

resource keycloak_identity_provider saml_one {
  alias   = "%s"
  realm   = "${keycloak_realm.realm_one.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider saml_two {
  alias   = "%s"
  realm   = "${keycloak_realm.realm_two.realm}"
  enabled = true

  saml {
    single_sign_on_service_url = "https://example.com"
  }
}

resource keycloak_identity_provider_mapper saml_mapper {
  realm                    = "${keycloak_realm.realm_two.realm}"
  name                     = "%s"
  identity_provider_alias  = "${keycloak_identity_provider.saml_two.alias}"
  identity_provider_mapper = "user-attribute-mapper"

  saml {
    template = "asdasdasdadsdad"
  }
}
   `, realmOne, realmTwo, aliasOne, aliasTwo, mapperName)
}
