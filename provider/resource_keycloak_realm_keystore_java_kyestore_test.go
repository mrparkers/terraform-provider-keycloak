package provider

//import (
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
//	"regexp"
//	"strconv"
//	"testing"
//)
//
//func TestAccKeycloakRealmKeystoreJavaKeystoreGenerated_basic(t *testing.T) {
//	t.Parallel()
//
//	javaKeystoreName := acctest.RandomWithPrefix("tf-acc")
//
//	resource.Test(t, resource.TestCase{
//		ProviderFactories: testAccProviderFactories,
//		PreCheck:          func() { testAccPreCheck(t) },
//		CheckDestroy:      testAccCheckRealmKeystoreJavaKeystoreGeneratedDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testKeycloakRealmKeystoreJavaKeystoreGenerated_basic(javaKeystoreName),
//				Check:  testAccCheckRealmKeystoreJavaKeystoreGeneratedExists("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore"),
//			},
//			{
//				ResourceName:      "keycloak_realm_key_javaKeystore_generated.realm_javaKeystore",
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: getRealmKeystoreJavaKeystoreGeneratedImportId("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore"),
//			},
//		},
//	})
//}
//
//func TestAccKeycloakRealmKeystoreJavaKeystoreGenerated_createAfterManualDestroy(t *testing.T) {
//	t.Parallel()
//
//	var javaKeystore = &keycloak.RealmKeystoreJavaKeystore{}
//
//	fullNameMapperName := acctest.RandomWithPrefix("tf-acc")
//
//	resource.Test(t, resource.TestCase{
//		ProviderFactories: testAccProviderFactories,
//		PreCheck:          func() { testAccPreCheck(t) },
//		CheckDestroy:      testAccCheckRealmKeystoreJavaKeystoreGeneratedDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testKeycloakRealmKeystoreJavaKeystoreGenerated_basic(fullNameMapperName),
//				Check:  testAccCheckRealmKeystoreJavaKeystoreGeneratedFetch("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore", javaKeystore),
//			},
//			{
//				PreConfig: func() {
//					err := keycloakClient.DeleteRealmKeystoreJavaKeystore(javaKeystore.RealmId, javaKeystore.Id)
//					if err != nil {
//						t.Fatal(err)
//					}
//				},
//				Config: testKeycloakRealmKeystoreJavaKeystoreGenerated_basic(fullNameMapperName),
//				Check:  testAccCheckRealmKeystoreJavaKeystoreGeneratedFetch("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore", javaKeystore),
//			},
//		},
//	})
//}
//
//func TestAccKeycloakRealmKeystoreJavaKeystoreGenerated_ellipticCurveValidation(t *testing.T) {
//	t.Parallel()
//
//	javaKeystoreName := acctest.RandomWithPrefix("tf-acc")
//	ellipticCurve := randomStringInSlice(keycloakRealmKeystoreJavaKeystoreGeneratedEllipticCurve)
//
//	resource.Test(t, resource.TestCase{
//		ProviderFactories: testAccProviderFactories,
//		PreCheck:          func() { testAccPreCheck(t) },
//		CheckDestroy:      testAccCheckRealmKeystoreJavaKeystoreGeneratedDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config:      testKeycloakRealmKeystoreJavaKeystoreGenerated_basicWithAttrValidation(javaKeystoreName, "elliptic_curve_key", acctest.RandString(10)),
//				ExpectError: regexp.MustCompile("expected mode to be one of .+ got .+"),
//			},
//			{
//				Config: testKeycloakRealmKeystoreJavaKeystoreGenerated_basicWithAttrValidation(javaKeystoreName, "elliptic_curve_key", ellipticCurve),
//				Check:  testAccCheckRealmKeystoreJavaKeystoreGeneratedExists("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore"),
//			},
//		},
//	})
//}
//
//func TestAccKeycloakRealmKeystoreJavaKeystoreGenerated_updateLdapUserFederationInPlace(t *testing.T) {
//	t.Parallel()
//
//	enabled := randomBool()
//	active := randomBool()
//
//	groupMapperOne := &keycloak.RealmKeystoreJavaKeystore{
//		Name:     acctest.RandString(10),
//		RealmId:  testAccRealmUserFederation.Realm,
//		Enabled:  enabled,
//		Active:   active,
//		Priority: acctest.RandInt(),
//	}
//
//	groupMapperTwo := &keycloak.RealmKeystoreJavaKeystore{
//		Name:     acctest.RandString(10),
//		RealmId:  testAccRealmUserFederation.Realm,
//		Enabled:  enabled,
//		Active:   active,
//		Priority: acctest.RandInt(),
//	}
//
//	resource.Test(t, resource.TestCase{
//		ProviderFactories: testAccProviderFactories,
//		PreCheck:          func() { testAccPreCheck(t) },
//		CheckDestroy:      testAccCheckRealmKeystoreJavaKeystoreGeneratedDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testKeycloakRealmKeystoreJavaKeystoreGenerated_basicFromInterface(groupMapperOne),
//				Check:  testAccCheckRealmKeystoreJavaKeystoreGeneratedExists("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore"),
//			},
//			{
//				Config: testKeycloakRealmKeystoreJavaKeystoreGenerated_basicFromInterface(groupMapperTwo),
//				Check:  testAccCheckRealmKeystoreJavaKeystoreGeneratedExists("keycloak_realm_key_javaKeystore_generated.realm_javaKeystore"),
//			},
//		},
//	})
//}
//
//func testAccCheckRealmKeystoreJavaKeystoreGeneratedExists(resourceName string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		_, err := getKeycloakRealmKeystoreJavaKeystoreGeneratedFromState(s, resourceName)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
//}
//
//func testAccCheckRealmKeystoreJavaKeystoreGeneratedFetch(resourceName string, mapper *keycloak.RealmKeystoreJavaKeystoreGenerated) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		fetchedMapper, err := getKeycloakRealmKeystoreJavaKeystoreGeneratedFromState(s, resourceName)
//		if err != nil {
//			return err
//		}
//
//		mapper.Id = fetchedMapper.Id
//		mapper.RealmId = fetchedMapper.RealmId
//
//		return nil
//	}
//}
//
//func testAccCheckRealmKeystoreJavaKeystoreGeneratedDestroy() resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		for _, rs := range s.RootModule().Resources {
//			if rs.Type != "keycloak_realm_key_javaKeystore_generated" {
//				continue
//			}
//
//			id := rs.Primary.ID
//			realm := rs.Primary.Attributes["realm_id"]
//
//			ldapGroupMapper, _ := keycloakClient.GetRealmKeystoreJavaKeystoreGenerated(realm, id)
//			if ldapGroupMapper != nil {
//				return fmt.Errorf("javaKeystore keystore with id %s still exists", id)
//			}
//		}
//
//		return nil
//	}
//}
//
//func getKeycloakRealmKeystoreJavaKeystoreGeneratedFromState(s *terraform.State,
//	resourceName string) (*keycloak.RealmKeystoreJavaKeystoreGenerated,
//	error) {
//	rs, ok := s.RootModule().Resources[resourceName]
//	if !ok {
//		return nil, fmt.Errorf("resource not found: %s", resourceName)
//	}
//
//	id := rs.Primary.ID
//	realm := rs.Primary.Attributes["realm_id"]
//
//	realmKeystore, err := keycloakClient.GetRealmKeystoreJavaKeystoreGenerated(realm, id)
//	if err != nil {
//		return nil, fmt.Errorf("error getting javaKeystore keystore with id %s: %s", id, err)
//	}
//
//	return realmKeystore, nil
//}
//
//func getRealmKeystoreJavaKeystoreGeneratedImportId(resourceName string) resource.ImportStateIdFunc {
//	return func(s *terraform.State) (string, error) {
//		rs, ok := s.RootModule().Resources[resourceName]
//		if !ok {
//			return "", fmt.Errorf("resource not found: %s", resourceName)
//		}
//
//		id := rs.Primary.ID
//		realmId := rs.Primary.Attributes["realm_id"]
//		providerId := "java-keystore"
//
//		return fmt.Sprintf("%s/%s/%s", realmId, providerId, id), nil
//	}
//}
//
//func testKeycloakRealmKeystoreJavaKeystoreGenerated_basic(javaKeystoreName string) string {
//	return fmt.Sprintf(`
//data "keycloak_realm" "realm" {
//	realm = "%s"
//}
//
//resource "keycloak_realm_key_javaKeystore_generated" "realm_java_keystore" {
//	name      = "%s"
//	realm_id  = data.keycloak_realm.realm.id
//	parent_id = data.keycloak_realm.realm.id
//
//   priority           = 100
//   elliptic_curve_key = "P-384"
//}
//	`, testAccRealmUserFederation.Realm, javaKeystoreName)
//}
//
//func testKeycloakRealmKeystoreJavaKeystoreGenerated_basicWithAttrValidation(javaKeystoreName, attr, val string) string {
//	return fmt.Sprintf(`
//data "keycloak_realm" "realm" {
//	realm = "%s"
//}
//
//resource "keycloak_ldap_group_mapper" "realm_java_keystore" {
//	name      = "%s"
//	realm_id  = data.keycloak_realm.realm.id
//	parent_id = data.keycloak_realm.realm.id
//
//	%s         = "%s"
//}
//	`, testAccRealmUserFederation.Realm, javaKeystoreName, attr, val)
//}
//
//func testKeycloakRealmKeystoreJavaKeystoreGenerated_basicFromInterface(mapper *keycloak.RealmKeystoreJavaKeystore) string {
//	return fmt.Sprintf(`
//data "keycloak_realm" "realm" {
//	realm = "%s"
//}
//
//resource "keycloak_ldap_group_mapper" "realm_java_keystore" {
//	name      = "%s"
//	realm_id  = data.keycloak_realm.realm.id
//	parent_id = data.keycloak_realm.realm.id
//
//   priority           = "%s"
//   elliptic_curve_key = "%s"
//}
//	`, testAccRealmUserFederation.Realm, mapper.Name, strconv.Itoa(mapper.Priority), mapper.EllipticCurve)
//}
