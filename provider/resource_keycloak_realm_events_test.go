package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakRealmEvents_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmEventsDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmEvents_basic(realmName),
				Check:  testAccCheckKeycloakRealmEventsExists("keycloak_realm_events.realm_events"),
			},
		},
	})
}

func TestAccKeycloakRealmEvents_update(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	before := &keycloak.RealmEventsConfig{
		AdminEventsDetailsEnabled: true,
		AdminEventsEnabled:        true,
		EnabledEventTypes:         []string{"LOGIN", "LOGOUT"},
		EventsEnabled:             true,
		EventsExpiration:          1234,
		EventsListeners:           []string{"jboss-logging"},
	}

	after := &keycloak.RealmEventsConfig{
		AdminEventsDetailsEnabled: false,
		AdminEventsEnabled:        false,
		EnabledEventTypes:         []string{"LOGIN"},
		EventsEnabled:             false,
		EventsExpiration:          12345,
		EventsListeners:           []string{"jboss-logging", "example-listener"},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmEventsDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmEvents_basicFromInterface(realmName, before),
				Check:  testAccCheckKeycloakRealmEventsExists("keycloak_realm_events.realm_events"),
			},
			{
				Config: testKeycloakRealmEvents_basicFromInterface(realmName, after),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRealmEventsExists("keycloak_realm_events.realm_events"),
					func(state *terraform.State) error {
						realmEventsConfig, err := getRealmEventsFromState(state, "keycloak_realm_events.realm_events")
						if err != nil {
							return err
						}

						if realmEventsConfig.AdminEventsDetailsEnabled {
							return fmt.Errorf("exptected admin_event_details_enabled to be false")
						}

						if realmEventsConfig.AdminEventsEnabled {
							return fmt.Errorf("exptected admin_events_enabled to be false")
						}

						if realmEventsConfig.EventsEnabled {
							return fmt.Errorf("exptected events_enabled to be false")
						}

						if realmEventsConfig.EventsExpiration != 12345 {
							return fmt.Errorf("exptected events_expiration to be 12345")
						}

						if len(realmEventsConfig.EnabledEventTypes) != 1 {
							return fmt.Errorf("exptected to enabled_event_types to contain exactly one element")
						}

						if len(realmEventsConfig.EventsListeners) != 2 {
							return fmt.Errorf("exptected to event_listeners to contain exactly two element elements")
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccKeycloakRealmEvents_unsetEnabledEventTypes(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	before := &keycloak.RealmEventsConfig{
		AdminEventsDetailsEnabled: true,
		AdminEventsEnabled:        true,
		EnabledEventTypes:         []string{"LOGIN", "LOGOUT"},
		EventsEnabled:             true,
		EventsExpiration:          1234,
		EventsListeners:           []string{"jboss-logging"},
	}

	after := &keycloak.RealmEventsConfig{
		AdminEventsDetailsEnabled: true,
		AdminEventsEnabled:        true,
		EnabledEventTypes:         []string{},
		EventsEnabled:             true,
		EventsExpiration:          1234,
		EventsListeners:           []string{"jboss-logging"},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmEventsDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmEvents_basicFromInterface(realmName, before),
				Check:  testAccCheckKeycloakRealmEventsExists("keycloak_realm_events.realm_events"),
			},
			{
				Config: testKeycloakRealmEvents_basicFromInterface(realmName, after),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakRealmEventsExists("keycloak_realm_events.realm_events"),
					func(state *terraform.State) error {
						realmEventsConfig, err := getRealmEventsFromState(state, "keycloak_realm_events.realm_events")
						if err != nil {
							return err
						}

						if len(realmEventsConfig.EnabledEventTypes) != 67 {
							return fmt.Errorf("exptected to enabled_event_types to contain all (67) event types, but it contains %d", len(realmEventsConfig.EnabledEventTypes))
						}

						return nil
					},
				),
			},
		},
	})
}

func getRealmEventsFromState(s *terraform.State, resourceName string) (*keycloak.RealmEventsConfig, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]

	realmEventsConfig, err := keycloakClient.GetRealmEventsConfig(realm)
	if err != nil {
		return nil, fmt.Errorf("error getting realm events config: %s", err)
	}

	return realmEventsConfig, nil
}

func testAccCheckKeycloakRealmEventsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRealmEventsFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakRealmEventsDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_user_attribute_mapper" {
				continue
			}

			realm := rs.Primary.Attributes["realm_id"]
			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			realmEventsConfig, err := keycloakClient.GetRealmEventsConfig(realm)
			if err != nil {
				return err
			}

			if len(realmEventsConfig.EnabledEventTypes) < 1 {
				return fmt.Errorf("Expected enabled_event_types to be greater than zero after destroy")
			}

		}

		return nil
	}
}

func testKeycloakRealmEvents_basic(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_events" "realm_events" {
  realm_id = "${keycloak_realm.realm.id}"

  admin_events_enabled         = true
  admin_events_details_enabled = true
  events_enabled               = true
  events_expiration            = 1234

  enabled_event_types = [
	"LOGIN",
	"LOGOUT",
  ]

  events_listeners = [
    "jboss-logging",
	"example-listener",
  ]
}
	`, realm)
}

func testKeycloakRealmEvents_basicFromInterface(realm string, realmEventsConfig *keycloak.RealmEventsConfig) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_realm_events" "realm_events" {
  realm_id = "${keycloak_realm.realm.id}"

  admin_events_enabled         = %t
  admin_events_details_enabled = %t
  events_enabled               = %t
  events_expiration            = %d

  enabled_event_types = %s

  events_listeners = %s
}
	`, realm, realmEventsConfig.AdminEventsEnabled, realmEventsConfig.AdminEventsDetailsEnabled, realmEventsConfig.EventsEnabled, realmEventsConfig.EventsExpiration, arrayOfStringsForTerraformResource(realmEventsConfig.EnabledEventTypes), arrayOfStringsForTerraformResource(realmEventsConfig.EventsListeners))
}
