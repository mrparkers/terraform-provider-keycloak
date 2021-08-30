package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"strconv"
	"testing"
)

func TestAccKeycloakRealmKeystoreEcdsaGenerated_basic(t *testing.T) {
	t.Parallel()

	ecdsaName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreEcdsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreEcdsaGenerated_basic(ecdsaName),
				Check:  testAccCheckRealmKeystoreEcdsaGeneratedExists("keycloak_realm_key_ecdsa_generated.realm_ecdsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreEcdsaGenerated_createAfterManualDestroy(t *testing.T) {
	t.Parallel()

	var ecdsa = &keycloak.RealmKeystoreEcdsaGenerated{}

	fullNameMapperName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreEcdsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreEcdsaGenerated_basic(fullNameMapperName),
				Check:  testAccCheckRealmKeystoreEcdsaGeneratedFetch("keycloak_realm_key_ecdsa_generated.realm_ecdsa", ecdsa),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteRealmKeystoreEcdsaGenerated(ecdsa.RealmId, ecdsa.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakRealmKeystoreEcdsaGenerated_basic(fullNameMapperName),
				Check:  testAccCheckRealmKeystoreEcdsaGeneratedFetch("keycloak_realm_key_ecdsa_generated.realm_ecdsa", ecdsa),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreEcdsaGenerated_ellipticCurveValidation(t *testing.T) {
	t.Parallel()

	ecdsaName := acctest.RandomWithPrefix("tf-acc")
	ellipticCurve := randomStringInSlice(keycloakRealmKeystoreEcdsaGeneratedEllipticCurve)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreEcdsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRealmKeystoreEcdsaGenerated_basicWithAttrValidation(ecdsaName, "elliptic_curve_key", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected elliptic_curve_key to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreEcdsaGenerated_basicWithAttrValidation(ecdsaName, "elliptic_curve_key", ellipticCurve),
				Check:  testAccCheckRealmKeystoreEcdsaGeneratedExists("keycloak_realm_key_ecdsa_generated.realm_ecdsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreEcdsaGenerated_updateRealmKeystoreEcdsaGenerated(t *testing.T) {
	t.Parallel()

	enabled := randomBool()
	active := randomBool()

	groupMapperOne := &keycloak.RealmKeystoreEcdsaGenerated{
		Name:          acctest.RandString(10),
		RealmId:       testAccRealmUserFederation.Realm,
		Enabled:       enabled,
		Active:        active,
		Priority:      acctest.RandIntRange(0, 100),
		EllipticCurve: randomStringInSlice(keycloakRealmKeystoreEcdsaGeneratedEllipticCurve),
	}

	groupMapperTwo := &keycloak.RealmKeystoreEcdsaGenerated{
		Name:          acctest.RandString(10),
		RealmId:       testAccRealmUserFederation.Realm,
		Enabled:       enabled,
		Active:        active,
		Priority:      acctest.RandIntRange(0, 100),
		EllipticCurve: randomStringInSlice(keycloakRealmKeystoreEcdsaGeneratedEllipticCurve),
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreEcdsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreEcdsaGenerated_basicFromInterface(groupMapperOne),
				Check:  testAccCheckRealmKeystoreEcdsaGeneratedExists("keycloak_realm_key_ecdsa_generated.realm_ecdsa"),
			},
			{
				Config: testKeycloakRealmKeystoreEcdsaGenerated_basicFromInterface(groupMapperTwo),
				Check:  testAccCheckRealmKeystoreEcdsaGeneratedExists("keycloak_realm_key_ecdsa_generated.realm_ecdsa"),
			},
		},
	})
}

func testAccCheckRealmKeystoreEcdsaGeneratedExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakRealmKeystoreEcdsaGeneratedFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckRealmKeystoreEcdsaGeneratedFetch(resourceName string, mapper *keycloak.RealmKeystoreEcdsaGenerated) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakRealmKeystoreEcdsaGeneratedFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckRealmKeystoreEcdsaGeneratedDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm_key_ecdsa_generated" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			ldapGroupMapper, _ := keycloakClient.GetRealmKeystoreEcdsaGenerated(realm, id)
			if ldapGroupMapper != nil {
				return fmt.Errorf("ecdsa keystore with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakRealmKeystoreEcdsaGeneratedFromState(s *terraform.State,
	resourceName string) (*keycloak.RealmKeystoreEcdsaGenerated,
	error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	realmKeystore, err := keycloakClient.GetRealmKeystoreEcdsaGenerated(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ecdsa keystore with id %s: %s", id, err)
	}

	return realmKeystore, nil
}

func getRealmKeystoreEcdsaGeneratedImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		providerId := "ecdsa-generated"

		return fmt.Sprintf("%s/%s/%s", realmId, providerId, id), nil
	}
}

func testKeycloakRealmKeystoreEcdsaGenerated_basic(ecdsaName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_ecdsa_generated" "realm_ecdsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

    priority           = 100
    elliptic_curve_key = "P-384"
}
	`, testAccRealmUserFederation.Realm, ecdsaName)
}

func testKeycloakRealmKeystoreEcdsaGenerated_basicWithAttrValidation(ecdsaName, attr, val string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_ecdsa_generated" "realm_ecdsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

	%s         = "%s"
}
	`, testAccRealmUserFederation.Realm, ecdsaName, attr, val)
}

func testKeycloakRealmKeystoreEcdsaGenerated_basicFromInterface(mapper *keycloak.RealmKeystoreEcdsaGenerated) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_ecdsa_generated" "realm_ecdsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

    priority           = "%s"
    elliptic_curve_key = "%s"
}
	`, testAccRealmUserFederation.Realm, mapper.Name, strconv.Itoa(mapper.Priority), mapper.EllipticCurve)
}
