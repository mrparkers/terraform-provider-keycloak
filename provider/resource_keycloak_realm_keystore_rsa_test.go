package provider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAccKeycloakRealmKeystoreRsa_basic(t *testing.T) {
	t.Parallel()

	rsaName := acctest.RandomWithPrefix("tf-acc")

	privateKey, certificate, err := generateKeyAndCert()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsa_basic(rsaName, privateKey, certificate),
				Check:  testAccCheckRealmKeystoreRsaExists("keycloak_realm_key_rsa.realm_rsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsa_createAfterManualDestroy(t *testing.T) {
	t.Parallel()

	var rsa = &keycloak.RealmKeystoreRsa{}

	fullNameMapperName := acctest.RandomWithPrefix("tf-acc")

	privateKey, certificate, err := generateKeyAndCert()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsa_basic(fullNameMapperName, privateKey, certificate),
				Check:  testAccCheckRealmKeystoreRsaFetch("keycloak_realm_key_rsa.realm_rsa", rsa),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteRealmKeystoreRsa(rsa.RealmId, rsa.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakRealmKeystoreRsa_basic(fullNameMapperName, privateKey, certificate),
				Check:  testAccCheckRealmKeystoreRsaFetch("keycloak_realm_key_rsa.realm_rsa", rsa),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsa_keySizeValidation(t *testing.T) {
	t.Parallel()

	rsaName := acctest.RandomWithPrefix("tf-acc")
	privateKey, certificate, err := generateKeyAndCert()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsa_basicWithAttrValidation(rsaName, "key_size",
					strconv.Itoa(acctest.RandIntRange(0, 10000)), privateKey, certificate),
				ExpectError: regexp.MustCompile("expected key_size to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreRsa_basicWithAttrValidation(rsaName, "key_size", "2048", privateKey,
					certificate),
				Check: testAccCheckRealmKeystoreRsaExists("keycloak_realm_key_rsa.realm_rsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsa_algorithmValidation(t *testing.T) {
	t.Parallel()

	algorithm := randomStringInSlice(keycloakRealmKeystoreRsaAlgorithm)
	privateKey, certificate, err := generateKeyAndCert()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsa_basicWithAttrValidation(algorithm, "algorithm",
					acctest.RandString(10), privateKey, certificate),
				ExpectError: regexp.MustCompile("expected algorithm to be one of .+ got .+"),
			},
			{
				Config: testKeycloakRealmKeystoreRsa_basicWithAttrValidation(algorithm, "algorithm", algorithm,
					privateKey, certificate),
				Check: testAccCheckRealmKeystoreRsaExists("keycloak_realm_key_rsa.realm_rsa"),
			},
		},
	})
}

func TestAccKeycloakRealmKeystoreRsa_updateRsaKeystoreGenerated(t *testing.T) {
	t.Parallel()

	enabled := randomBool()
	active := randomBool()
	privateKey, certificate, err := generateKeyAndCert()
	if err != nil {
		log.Fatal(err)
	}

	groupMapperOne := &keycloak.RealmKeystoreRsa{
		Name:        acctest.RandString(10),
		RealmId:     testAccRealmUserFederation.Realm,
		Enabled:     enabled,
		Active:      active,
		Priority:    acctest.RandIntRange(0, 100),
		KeySize:     1024,
		Algorithm:   randomStringInSlice(keycloakRealmKeystoreRsaAlgorithm),
		PrivateKey:  privateKey,
		Certificate: certificate,
	}

	groupMapperTwo := &keycloak.RealmKeystoreRsa{
		Name:        acctest.RandString(10),
		RealmId:     testAccRealmUserFederation.Realm,
		Enabled:     enabled,
		Active:      active,
		Priority:    acctest.RandIntRange(0, 100),
		KeySize:     2048,
		Algorithm:   randomStringInSlice(keycloakRealmKeystoreRsaAlgorithm),
		PrivateKey:  privateKey,
		Certificate: certificate,
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckRealmKeystoreRsaDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmKeystoreRsa_basicFromInterface(groupMapperOne),
				Check:  testAccCheckRealmKeystoreRsaExists("keycloak_realm_key_rsa.realm_rsa"),
			},
			{
				Config: testKeycloakRealmKeystoreRsa_basicFromInterface(groupMapperTwo),
				Check:  testAccCheckRealmKeystoreRsaExists("keycloak_realm_key_rsa.realm_rsa"),
			},
		},
	})
}

func testAccCheckRealmKeystoreRsaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakRealmKeystoreRsaFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckRealmKeystoreRsaFetch(resourceName string, mapper *keycloak.RealmKeystoreRsa) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakRealmKeystoreRsaFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckRealmKeystoreRsaDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm_key_rsa" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			ldapGroupMapper, _ := keycloakClient.GetRealmKeystoreRsa(realm, id)
			if ldapGroupMapper != nil {
				return fmt.Errorf("rsa keystore with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakRealmKeystoreRsaFromState(s *terraform.State,
	resourceName string) (*keycloak.RealmKeystoreRsa,
	error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	realmKeystore, err := keycloakClient.GetRealmKeystoreRsa(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting rsa keystore with id %s: %s", id, err)
	}

	return realmKeystore, nil
}

func getRealmKeystoreRsaImportId(resourceName string) resource.ImportStateIdFunc {
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

func generateKeyAndCert() (string, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Private key cannot be created.", err.Error())
	}

	// Generate a pem block with the private key
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	tml := x509.Certificate{
		// you can add any attr that you need
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 0, 0),
		// you have to generate a different serial number each execution
		SerialNumber: big.NewInt(123123),
		Subject: pkix.Name{
			CommonName:   "New Name",
			Organization: []string{"New Org."},
		},
		BasicConstraintsValid: true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	if err != nil {
		log.Fatal("Certificate cannot be created.", err.Error())
	}

	// Generate a pem block with the certificate
	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return strings.ReplaceAll(string(keyPem), "\n", "\\n"), strings.ReplaceAll(string(certPem), "\n", "\\n"), nil
}

func testKeycloakRealmKeystoreRsa_basic(rsaName, privateKey, certificate string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_rsa" "realm_rsa" {

	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

    priority    = 100
    algorithm   = "RS384"
    private_key = "%s"
    certificate = "%s"
}
	`, testAccRealmUserFederation.Realm, rsaName, privateKey, certificate)
}

func testKeycloakRealmKeystoreRsa_basicWithAttrValidation(rsaName, attr, val, privateKey,
	certificate string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_rsa" "realm_rsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

	%s        = "%s"

    private_key = "%s"
    certificate = "%s"
}
	`, testAccRealmUserFederation.Realm, rsaName, attr, val, privateKey, certificate)
}

func testKeycloakRealmKeystoreRsa_basicFromInterface(mapper *keycloak.RealmKeystoreRsa) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_key_rsa" "realm_rsa" {
	name      = "%s"
	realm_id  = data.keycloak_realm.realm.id
	parent_id = data.keycloak_realm.realm.id

    priority  = %s
    algorithm = "%s"
    key_size  = %s

    private_key = "%s"
    certificate = "%s"
}
	`, testAccRealmUserFederation.Realm, mapper.Name, strconv.Itoa(mapper.Priority), mapper.Algorithm,
		strconv.Itoa(mapper.KeySize), mapper.PrivateKey, mapper.Certificate)
}
