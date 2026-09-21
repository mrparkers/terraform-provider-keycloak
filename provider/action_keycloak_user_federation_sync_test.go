package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

const (
	testLdapConnectionUrl  = "ldap://openldap"
	testLdapBindDn         = "cn=admin,dc=example,dc=org"
	testLdapBindCredential = "adminpassword"
)

func TestAccKeycloakUserFederationSyncAction_ldap(t *testing.T) {
	t.Parallel()

	ldapName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfLdapServerIsNotAvailable(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// actions are supported since Terraform 1.14
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: testAccCheckKeycloakLdapUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserFederationSyncAction_ldap(ldapName, "full"),
				Check:  testAccCheckKeycloakUserFederationWasSynced("keycloak_ldap_user_federation.openldap"),
			},
		},
	})
}

// the action is not limited to LDAP, it works for every user storage provider which supports synchronization.
// The custom user federation example imports a user on sync, which records the kind of the last sync. This makes
// the sync observable without the need for an LDAP server. Note that this test does not use the realm of the LDAP
// tests, since user lookups fail on older Keycloak versions if the realm contains an unreachable LDAP user federation.
func TestAccKeycloakUserFederationSyncAction_custom(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: testAccCheckKeycloakCustomUserFederationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserFederationSyncAction_custom(name, "dummy", "full"),
				Check:  testAccCheckKeycloakCustomUserFederationWasSynced(name, "full"),
			},
			{
				// the update of the user federation triggers the action again
				Config: testKeycloakUserFederationSyncAction_custom(name, "dummy-updated", "changed_users"),
				Check:  testAccCheckKeycloakCustomUserFederationWasSynced(name, "changed"),
			},
		},
	})
}

func TestAccKeycloakUserFederationSyncAction_invalidMode(t *testing.T) {
	t.Parallel()

	ldapName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakUserFederationSyncAction_ldap(ldapName, "invalid"),
				ExpectError: regexp.MustCompile("invalid mode"),
			},
		},
	})
}

// The sync of an LDAP user federation requires an LDAP server, which is available in the local environment (see
// docker-compose.yml), but not on CI. A sync which does not depend on an LDAP server is covered by
// TestAccKeycloakUserFederationSyncAction_custom.
func skipIfLdapServerIsNotAvailable(t *testing.T) {
	err := keycloakClient.TestLdapAuthentication(testCtx, testAccRealmUserFederation.Realm, testLdapConnectionUrl, testLdapBindDn, testLdapBindCredential)
	if err != nil {
		t.Skipf("skipping: Keycloak cannot reach the LDAP server %s: %s", testLdapConnectionUrl, err)
	}
}

// the sync triggered by the action already imported all LDAP users, so another full sync must only report updates.
func testAccCheckKeycloakUserFederationWasSynced(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		result, err := keycloakClient.SyncUserFederation(testCtx, rs.Primary.Attributes["realm_id"], rs.Primary.ID, keycloak.UserFederationSyncFull)
		if err != nil {
			return err
		}

		if result.Added != 0 || result.Updated == 0 {
			return fmt.Errorf("expected users to be already imported by the sync action, but got: %+v", result)
		}

		// the request timeout can be overridden per request
		_, err = keycloakClient.SyncUserFederation(keycloak.WithRequestTimeout(testCtx, time.Nanosecond), rs.Primary.Attributes["realm_id"], rs.Primary.ID, keycloak.UserFederationSyncFull)
		if err == nil || !strings.Contains(err.Error(), "Client.Timeout") {
			return fmt.Errorf("expected sync request to run into the overridden request timeout, but got: %v", err)
		}

		return nil
	}
}

func testAccCheckKeycloakCustomUserFederationWasSynced(name, syncMode string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		username := name + "-synced"

		user, err := keycloakClient.GetUserByUsername(testCtx, testAccRealm.Realm, username)
		if err != nil {
			return err
		}

		if user == nil {
			return fmt.Errorf("expected user %s to be imported by the sync action", username)
		}

		// the custom user federation example records the kind of the last sync as last name
		if user.LastName != syncMode {
			return fmt.Errorf("expected user %s to be imported by a %s sync, but got: %q", username, syncMode, user.LastName)
		}

		return nil
	}
}

func testKeycloakUserFederationSyncAction_custom(name, dummyConfig, mode string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_custom_user_federation" "custom" {
	name        = "%s"
	realm_id    = data.keycloak_realm.realm.id
	provider_id = "custom"

	enabled     = true

	config = {
		dummyConfig      = "%s"
		importUserOnSync = "true"
	}

	lifecycle {
		action_trigger {
			events  = [after_create, after_update]
			actions = [action.keycloak_user_federation_sync.custom]
		}
	}
}

action "keycloak_user_federation_sync" "custom" {
	config {
		realm_id           = data.keycloak_realm.realm.id
		user_federation_id = keycloak_custom_user_federation.custom.id
		mode               = "%s"
	}
}
	`, testAccRealm.Realm, name, dummyConfig, mode)
}

func testKeycloakUserFederationSyncAction_ldap(ldap, mode string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "%s"
	realm_id                = data.keycloak_realm.realm.id

	enabled                 = true
	import_enabled          = true
	edit_mode               = "READ_ONLY"

	username_ldap_attribute = "uid"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryUUID"
	user_object_classes     = [
		"inetOrgPerson"
	]
	connection_url          = "%s"
	users_dn                = "ou=users,dc=example,dc=org"
	bind_dn                 = "%s"
	bind_credential         = "%s"

	lifecycle {
		action_trigger {
			events  = [after_create, after_update]
			actions = [action.keycloak_user_federation_sync.openldap]
		}
	}
}

action "keycloak_user_federation_sync" "openldap" {
	config {
		realm_id           = data.keycloak_realm.realm.id
		user_federation_id = keycloak_ldap_user_federation.openldap.id
		mode               = "%s"
		client_timeout     = 60
	}
}
	`, testAccRealmUserFederation.Realm, ldap, testLdapConnectionUrl, testLdapBindDn, testLdapBindCredential, mode)
}
