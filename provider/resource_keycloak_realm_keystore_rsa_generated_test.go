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

func TestAccKeycloakRealmKeystoreRsaGenerated_basic(t *testing.T) {
	t.Parallel()

	rsaName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basic(rsaName),
				Check:  testAccCheckRealmKeystoreRsaGeneratedExists("keycloak_realm_key_rsa_generated.realm_rsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaGenerated_createAfterManualDestroy(t *testing.T) {
	t.Parallel()

	var rsa = &keycloak.RealmKeystoreRsaGenerated{}

	fullNameMapperName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basic(fullNameMapperName),
				Check:  testAccCheckRealmKeystoreRsaGeneratedFetch("keycloak_realm_key_rsa_generated.realm_rsa", rsa),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteRealmKeystoreRsaGenerated(rsa.RealmId, rsa.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakRealmKeystoreRsaGenerated_basic(fullNameMapperName),
				Check:  testAccCheckRealmKeystoreRsaGeneratedFetch("keycloak_realm_key_rsa_generated.realm_rsa", rsa),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaGenerated_keySizeValidation(t *testing.T) {
	t.Parallel()

	rsaName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basicWithAttrValidation(rsaName, "key_size",
					strconv.Itoa(acctest.RandIntRange(0, 10000))),
				ExpectError: regexp.MustCompile("expected key_size to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basicWithAttrValidation(rsaName, "key_size", "2048"),
				Check:  testAccCheckRealmKeystoreRsaGeneratedExists("keycloak_realm_key_rsa_generated.realm_rsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaGenerated_algorithmValidation(t *testing.T) {
	t.Parallel()

	algorithm := randomStringInSlice(keycloakRealmKeystoreRsaAlgorithm)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basicWithAttrValidation(algorithm, "algorithm",
					acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected algorithm to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basicWithAttrValidation(algorithm, "algorithm", algorithm),
				Check:  testAccCheckRealmKeystoreRsaGeneratedExists("keycloak_realm_key_rsa_generated.realm_rsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsaGenerated_updateRsaKeystoreGenerated(t *testing.T) {
	t.Parallel()

	enabled := randomBool()
	active := randomBool()

	groupMapperOne := &keycloak.RealmKeystoreRsaGenerated{
		Name:      acctest.RandString(10),
		RealmId:   testAccRealmUserFederation.Realm,
		Enabled:   enabled,
		Active:    active,
		Priority:  acctest.RandIntRange(0, 100),
		KeySize:   1024,
		Algorithm: randomStringInSlice(keycloakRealmKeystoreRsaAlgorithm),
	}

	groupMapperTwo := &keycloak.RealmKeystoreRsaGenerated{
		Name:      acctest.RandString(10),
		RealmId:   testAccRealmUserFederation.Realm,
		Enabled:   enabled,
		Active:    active,
		Priority:  acctest.RandIntRange(0, 100),
		KeySize:   2048,
		Algorithm: randomStringInSlice(keycloakRealmKeystoreRsaAlgorithm),
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaGeneratedDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basicFromInterface(groupMapperOne),
				Check:  testAccCheckRealmKeystoreRsaGeneratedExists("keycloak_realm_key_rsa_generated.realm_rsa"),
			},
			{
				Config: testKeycloakRealmKeystoreRsaGenerated_basicFromInterface(groupMapperTwo),
				Check:  testAccCheckRealmKeystoreRsaGeneratedExists("keycloak_realm_key_rsa_generated.realm_rsa"),
			},
		},
	})
}

func testAccCheckRealmKeystoreRsaGeneratedExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakRealmKeystoreRsaGeneratedFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckRealmKeystoreRsaGeneratedFetch(resourceName string, mapper *keycloak.RealmKeystoreRsaGenerated) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakRealmKeystoreRsaGeneratedFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckRealmKeystoreRsaGeneratedDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm_key_rsa_generated" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			ldapGroupMapper, _ := keycloakClient.GetRealmKeystoreRsaGenerated(realm, id)
			if ldapGroupMapper != nil {
				return fmt.Errorf("rsa keystore with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakRealmKeystoreRsaGeneratedFromState(s *terraform.State,
	resourceName string) (*keycloak.RealmKeystoreRsaGenerated,
	error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	realmKeystore, err := keycloakClient.GetRealmKeystoreRsaGenerated(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting rsa keystore with id %s: %s", id, err)
	}

	return realmKeystore, nil
}

func getRealmKeystoreRsaGeneratedImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		providerId := "rsa-generated"

		return fmt.Sprintf("%s/%s/%s", realmId, providerId, id), nil
	}
}

func testKeycloakRealmKeystoreRsaGenerated_basic(rsaName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_rsa_generated" "realm_rsa" {

	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

    priority  = 100
    algorithm = "RS384"
}
	`, testAccRealmUserFederation.Realm, rsaName)
}

func testKeycloakRealmKeystoreRsaGenerated_basicWithAttrValidation(rsaName, attr, val string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_rsa_generated" "realm_rsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

	%s        = "%s"
}
	`, testAccRealmUserFederation.Realm, rsaName, attr, val)
}

func testKeycloakRealmKeystoreRsaGenerated_basicFromInterface(mapper *keycloak.RealmKeystoreRsaGenerated) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_rsa_generated" "realm_rsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

    priority  = %s
    algorithm = "%s"
    key_size  = %s
}
	`, testAccRealmUserFederation.Realm, mapper.Name, strconv.Itoa(mapper.Priority), mapper.Algorithm,
		strconv.Itoa(mapper.KeySize))
}
