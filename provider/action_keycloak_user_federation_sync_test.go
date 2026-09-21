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

func TestAccKeycloakUserFederationSyncAction_ldap(t *testing.T) {
	t.Parallel()

	ldapName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
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
				Config: testKeycloakCustomUserFederation_basic(name, "custom") + `
action "keycloak_user_federation_sync" "custom" {
	config {
		realm_id           = data.keycloak_realm.realm.id
		user_federation_id = keycloak_custom_user_federation.custom.id
		mode               = "changed_users"
	}
}

resource "terraform_data" "trigger" {
	lifecycle {
		action_trigger {
			events  = [after_create]
			actions = [action.keycloak_user_federation_sync.custom]
		}
	}
}
`,
				Check: testAccCheckKeycloakCustomUserFederationExists("keycloak_custom_user_federation.custom"),
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
	connection_url          = "ldap://openldap"
	users_dn                = "ou=users,dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "adminpassword"

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
	`, testAccRealmUserFederation.Realm, ldap, mode)
}
