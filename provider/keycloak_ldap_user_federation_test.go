package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"regexp"
	"strconv"
	"testing"
)

func TestAccKeycloakLdapUserFederation_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basic(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_import(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	bindCredentialForImport := "admin"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basic(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
			{
				ResourceName:      "keycloak_ldap_user_federation.openldap",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapUserFederationImportId("keycloak_ldap_user_federation.openldap", bindCredentialForImport),
			},
			{
				Config: testKeycloakLdapUserFederation_noAuth(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap_no_auth"),
			},
			{
				ResourceName:        "keycloak_ldap_user_federation.openldap_no_auth",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_createAfterManualDestroy(t *testing.T) {
	var ldap = &keycloak.LdapUserFederation{}

	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basic(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationFetch("keycloak_ldap_user_federation.openldap", ldap),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteLdapUserFederation(ldap.RealmId, ldap.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapUserFederation_basic(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basic(firstRealm, ldapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
					resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "realm_id", firstRealm),
				),
			},
			{
				Config: testKeycloakLdapUserFederation_basic(secondRealm, ldapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
					resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "realm_id", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	firstEnabled := randomBool()
	firstValidatePasswordPolicy := randomBool()
	firstPagination := randomBool()

	firstConnectionTimeout, _ := keycloak.GetDurationStringFromMilliseconds(strconv.Itoa(acctest.RandIntRange(1000, 3600000)))
	secondConnectionTimeout, _ := keycloak.GetDurationStringFromMilliseconds(strconv.Itoa(acctest.RandIntRange(1000, 3600000)))
	firstReadTimeout, _ := keycloak.GetDurationStringFromMilliseconds(strconv.Itoa(acctest.RandIntRange(1000, 3600000)))
	secondReadTimeout, _ := keycloak.GetDurationStringFromMilliseconds(strconv.Itoa(acctest.RandIntRange(1000, 3600000)))

	firstLdap := &keycloak.LdapUserFederation{
		RealmId:                realmName,
		Name:                   "terraform-" + acctest.RandString(10),
		Enabled:                firstEnabled,
		UsernameLDAPAttribute:  acctest.RandString(10),
		UuidLDAPAttribute:      acctest.RandString(10),
		UserObjectClasses:      []string{acctest.RandString(10), acctest.RandString(10), acctest.RandString(10)},
		ConnectionUrl:          "ldap://" + acctest.RandString(10),
		UsersDn:                acctest.RandString(10),
		BindDn:                 acctest.RandString(10),
		BindCredential:         acctest.RandString(10),
		SearchScope:            randomStringInSlice([]string{"ONE_LEVEL", "SUBTREE"}),
		ValidatePasswordPolicy: firstValidatePasswordPolicy,
		UseTruststoreSpi:       randomStringInSlice([]string{"ALWAYS", "ONLY_FOR_LDAPS", "NEVER"}),
		ConnectionTimeout:      firstConnectionTimeout,
		ReadTimeout:            firstReadTimeout,
		Pagination:             firstPagination,
		BatchSizeForSync:       acctest.RandIntRange(50, 10000),
		FullSyncPeriod:         acctest.RandIntRange(1, 3600),
		ChangedSyncPeriod:      acctest.RandIntRange(1, 3600),
		CachePolicy:            randomStringInSlice([]string{"DEFAULT", "EVICT_DAILY", "EVICT_WEEKLY", "MAX_LIFESPAN", "NO_CACHE"}),
	}

	secondLdap := &keycloak.LdapUserFederation{
		RealmId:                realmName,
		Name:                   "terraform-" + acctest.RandString(10),
		Enabled:                !firstEnabled,
		UsernameLDAPAttribute:  acctest.RandString(10),
		UuidLDAPAttribute:      acctest.RandString(10),
		UserObjectClasses:      []string{acctest.RandString(10)},
		ConnectionUrl:          "ldap://" + acctest.RandString(10),
		UsersDn:                acctest.RandString(10),
		BindDn:                 acctest.RandString(10),
		BindCredential:         acctest.RandString(10),
		SearchScope:            randomStringInSlice([]string{"ONE_LEVEL", "SUBTREE"}),
		ValidatePasswordPolicy: !firstValidatePasswordPolicy,
		UseTruststoreSpi:       randomStringInSlice([]string{"ALWAYS", "ONLY_FOR_LDAPS", "NEVER"}),
		ConnectionTimeout:      secondConnectionTimeout,
		ReadTimeout:            secondReadTimeout,
		Pagination:             !firstPagination,
		BatchSizeForSync:       acctest.RandIntRange(50, 10000),
		FullSyncPeriod:         acctest.RandIntRange(1, 3600),
		ChangedSyncPeriod:      acctest.RandIntRange(1, 3600),
		CachePolicy:            randomStringInSlice([]string{"DEFAULT", "EVICT_DAILY", "EVICT_WEEKLY", "MAX_LIFESPAN", "NO_CACHE"}),
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basicFromInterface(firstLdap),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
			{
				Config: testKeycloakLdapUserFederation_basicFromInterface(secondLdap),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_unsetTimeoutDurationStrings(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_basicWithTimeouts(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
			{
				Config: testKeycloakLdapUserFederation_basic(realmName, ldapName),
				Check:  testAccCheckKeycloakLdapUserFederationExists("keycloak_ldap_user_federation.openldap"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_editModeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	editMode := randomStringInSlice(keycloakLdapUserFederationEditModes)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("edit_mode", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected edit_mode to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("edit_mode", realmName, ldapName, editMode),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "edit_mode", editMode),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_vendorValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	vendor := randomStringInSlice(keycloakLdapUserFederationVendors)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("vendor", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected vendor to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("vendor", realmName, ldapName, vendor),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "vendor", vendor),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_searchScopeValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	searchScope := randomStringInSlice(keycloakLdapUserFederationSearchScopes)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("search_scope", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected search_scope to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("search_scope", realmName, ldapName, searchScope),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "search_scope", searchScope),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_useTrustStoreValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	useTrustStore := randomStringInSlice(keycloakLdapUserFederationTruststoreSpiSettings)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("use_truststore_spi", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected use_truststore_spi to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("use_truststore_spi", realmName, ldapName, useTrustStore),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "use_truststore_spi", useTrustStore),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_cachePolicyValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	cachePolicy := randomStringInSlice(keycloakUserFederationCachePolicies)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_basicWithAttrValidation("cache_policy", realmName, ldapName, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("expected cache_policy to be one of .+ got .+"),
			},
			{
				Config: testKeycloakLdapUserFederation_basicWithAttrValidation("cache_policy", realmName, ldapName, cachePolicy),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "cache_policy", cachePolicy),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_bindValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_noBindCredentialValidation(realmName, ldapName),
				ExpectError: regexp.MustCompile("validation error: authentication requires both BindDN and BindCredential to be set"),
			},
			{
				Config:      testKeycloakLdapUserFederation_nobindDnValidation(realmName, ldapName),
				ExpectError: regexp.MustCompile("validation error: authentication requires both BindDN and BindCredential to be set"),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_syncPeriodValidation(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)

	validSyncPeriod := acctest.RandIntRange(1, 3600)
	invalidNegativeSyncPeriod := -acctest.RandIntRange(1, 3600)
	invalidZeroSyncPeriod := 0

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakLdapUserFederation_basicWithSyncPeriod(realmName, ldapName, validSyncPeriod, invalidNegativeSyncPeriod),
				ExpectError: regexp.MustCompile(`expected .+ to be either -1 \(disabled\), or greater than zero`),
			},
			{
				Config:      testKeycloakLdapUserFederation_basicWithSyncPeriod(realmName, ldapName, invalidNegativeSyncPeriod, validSyncPeriod),
				ExpectError: regexp.MustCompile(`expected .+ to be either -1 \(disabled\), or greater than zero`),
			},
			{
				Config:      testKeycloakLdapUserFederation_basicWithSyncPeriod(realmName, ldapName, validSyncPeriod, invalidZeroSyncPeriod),
				ExpectError: regexp.MustCompile(`expected .+ to be either -1 \(disabled\), or greater than zero`),
			},
			{
				Config:      testKeycloakLdapUserFederation_basicWithSyncPeriod(realmName, ldapName, invalidZeroSyncPeriod, validSyncPeriod),
				ExpectError: regexp.MustCompile(`expected .+ to be either -1 \(disabled\), or greater than zero`),
			},
			{
				Config: testKeycloakLdapUserFederation_basicWithSyncPeriod(realmName, ldapName, validSyncPeriod, validSyncPeriod),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "full_sync_period", strconv.Itoa(validSyncPeriod)),
					resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "changed_sync_period", strconv.Itoa(validSyncPeriod)),
				),
			},
		},
	})
}

func TestAccKeycloakLdapUserFederation_bindCredential(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	ldapName := "terraform-" + acctest.RandString(10)
	firstBindCredential := acctest.RandString(10)
	secondBindCredential := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapUserFederation_bindCredential(realmName, ldapName, firstBindCredential),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "bind_credential", firstBindCredential),
			},
			{
				Config: testKeycloakLdapUserFederation_bindCredential(realmName, ldapName, secondBindCredential),
				Check:  resource.TestCheckResourceAttr("keycloak_ldap_user_federation.openldap", "bind_credential", secondBindCredential),
			},
		},
	})
}

func testAccCheckKeycloakLdapUserFederationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapUserFederationFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapUserFederationFetch(resourceName string, ldap *keycloak.LdapUserFederation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedLdap, err := getLdapUserFederationFromState(s, resourceName)
		if err != nil {
			return err
		}

		ldap.Id = fetchedLdap.Id
		ldap.RealmId = fetchedLdap.RealmId

		return nil
	}
}

func testAccCheckKeycloakLdapUserFederationDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_user_federation" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			ldap, _ := keycloakClient.GetLdapUserFederation(realm, id)
			if ldap != nil {
				return fmt.Errorf("ldap config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapUserFederationFromState(s *terraform.State, resourceName string) (*keycloak.LdapUserFederation, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldap, err := keycloakClient.GetLdapUserFederation(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap config with id %s: %s", id, err)
	}

	return ldap, nil
}

func getLdapUserFederationImportId(resourceName, bindCredential string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]

		return fmt.Sprintf("%s/%s/%s", realmId, id, bindCredential), nil
	}
}

func testKeycloakLdapUserFederation_basic(realm, ldap string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "admin"
}
	`, realm, ldap)
}

func testKeycloakLdapUserFederation_basicFromInterface(ldap *keycloak.LdapUserFederation) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                     = "%s"
	realm_id                 = "${keycloak_realm.realm.id}"

	enabled                  = %t

	username_ldap_attribute  = "%s"
	rdn_ldap_attribute       = "%s"
	uuid_ldap_attribute      = "%s"
	user_object_classes      = %s
	connection_url           = "%s"
	users_dn                 = "%s"
	bind_dn                  = "%s"
	bind_credential          = "%s"
	search_scope             = "%s"

	validate_password_policy = %t
	use_truststore_spi       = "%s"
	connection_timeout       = "%s"
	read_timeout             = "%s"
	pagination               = %t

	batch_size_for_sync      = %d
	full_sync_period         = %d
	changed_sync_period      = %d

	cache_policy             = "%s"
}
	`, ldap.RealmId, ldap.Name, ldap.Enabled, ldap.UsernameLDAPAttribute, ldap.RdnLDAPAttribute, ldap.UuidLDAPAttribute, arrayOfStringsForTerraformResource(ldap.UserObjectClasses), ldap.ConnectionUrl, ldap.UsersDn, ldap.BindDn, ldap.BindCredential, ldap.SearchScope, ldap.ValidatePasswordPolicy, ldap.UseTruststoreSpi, ldap.ConnectionTimeout, ldap.ReadTimeout, ldap.Pagination, ldap.BatchSizeForSync, ldap.FullSyncPeriod, ldap.ChangedSyncPeriod, ldap.CachePolicy)
}

func testKeycloakLdapUserFederation_basicWithAttrValidation(attr, realm, ldap, val string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	%s                      = "%s"

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "admin"
}
	`, realm, ldap, attr, val)
}

func testKeycloakLdapUserFederation_nobindDnValidation(realm, ldap string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	bind_credential         = "admin"

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
}
	`, realm, ldap)
}

func testKeycloakLdapUserFederation_noBindCredentialValidation(realm, ldap string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	bind_dn                 = "cn=admin,dc=example,dc=org"

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
}
	`, realm, ldap)
}

func testKeycloakLdapUserFederation_basicWithSyncPeriod(realm, ldap string, fullSyncPeriod, changedSyncPeriod int) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "admin"

	full_sync_period        = %d
	changed_sync_period     = %d
}
	`, realm, ldap, fullSyncPeriod, changedSyncPeriod)
}

func testKeycloakLdapUserFederation_basicWithTimeouts(realm, ldap string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "admin"

	connection_timeout      = "10s"
	read_timeout            = "5s"
}
	`, realm, ldap)
}

func testKeycloakLdapUserFederation_bindCredential(realm, ldap, bindCredential string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "%s"
}
	`, realm, ldap, bindCredential)
}

func testKeycloakLdapUserFederation_noAuth(realm, ldap string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap_no_auth" {
	name                    = "%s"
	realm_id                = "${keycloak_realm.realm.id}"

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalRole"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
}
	`, realm, ldap)
}
