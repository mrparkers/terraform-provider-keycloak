package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakGroup_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroup_basic(realmName, groupName),
				Check:  testAccCheckKeycloakGroupExists("keycloak_group.group"),
			},
			{
				ResourceName:        "keycloak_group.group",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakGroup_createAfterManualDestroy(t *testing.T) {
	var group = &keycloak.Group{}

	realmName := "terraform-" + acctest.RandString(10)
	groupName := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroup_basic(realmName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group"),
					testAccCheckKeycloakGroupFetch("keycloak_group.group", group),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteGroup(group.RealmId, group.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakGroup_basic(realmName, groupName),
				Check:  testAccCheckKeycloakGroupExists("keycloak_group.group"),
			},
		},
	})
}

func TestAccKeycloakGroup_updateGroupName(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	groupNameBefore := "terraform-group-" + acctest.RandString(10)
	groupNameAfter := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroup_basic(realmName, groupNameBefore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group"),
					resource.TestCheckResourceAttr("keycloak_group.group", "name", groupNameBefore),
				),
			},
			{
				Config: testKeycloakGroup_basic(realmName, groupNameAfter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group"),
					resource.TestCheckResourceAttr("keycloak_group.group", "name", groupNameAfter),
				),
			},
		},
	})
}

func TestAccKeycloakGroup_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)

	group := "terraform-group-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroup_updateRealmBefore(realmOne, realmTwo, group),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group"),
					testAccCheckKeycloakGroupBelongsToRealm("keycloak_group.group", realmOne),
				),
			},
			{
				Config: testKeycloakGroup_updateRealmAfter(realmOne, realmTwo, group),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists("keycloak_group.group"),
					testAccCheckKeycloakGroupBelongsToRealm("keycloak_group.group", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakGroup_nested(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	parentGroupName := "terraform-parent-group-" + acctest.RandString(10)
	firstChildGroupName := "terraform-child-group-" + acctest.RandString(10)
	secondChildGroupName := "terraform-child-group-" + acctest.RandString(10)

	parentGroupResource := "keycloak_group.parent_group"
	firstChildGroupResource := "keycloak_group.first_child_group"
	secondChildGroupResource := "keycloak_group.second_child_group"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGroup_nested(realmName, parentGroupName, firstChildGroupName, secondChildGroupName, firstChildGroupResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists(parentGroupResource),
					testAccCheckKeycloakGroupExists(firstChildGroupResource),
					testAccCheckKeycloakGroupExists(secondChildGroupResource),

					resource.TestCheckNoResourceAttr(parentGroupResource, "parent_id"),
					resource.TestCheckResourceAttrPair(firstChildGroupResource, "parent_id", parentGroupResource, "id"),
					resource.TestCheckResourceAttrPair(secondChildGroupResource, "parent_id", firstChildGroupResource, "id"),
				),
			},
			{
				Config: testKeycloakGroup_nested(realmName, parentGroupName, firstChildGroupName, secondChildGroupName, parentGroupResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakGroupExists(parentGroupResource),
					testAccCheckKeycloakGroupExists(firstChildGroupResource),
					testAccCheckKeycloakGroupExists(secondChildGroupResource),

					resource.TestCheckNoResourceAttr(parentGroupResource, "parent_id"),
					resource.TestCheckResourceAttrPair(firstChildGroupResource, "parent_id", parentGroupResource, "id"),
					resource.TestCheckResourceAttrPair(secondChildGroupResource, "parent_id", parentGroupResource, "id"),
				),
			},
		},
	})
}

func testAccCheckKeycloakGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getGroupFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakGroupFetch(resourceName string, group *keycloak.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedGroup, err := getGroupFromState(s, resourceName)
		if err != nil {
			return err
		}

		group.Id = fetchedGroup.Id
		group.RealmId = fetchedGroup.RealmId

		return nil
	}
}

func testAccCheckKeycloakGroupBelongsToRealm(resourceName, realm string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		group, err := getGroupFromState(s, resourceName)
		if err != nil {
			return err
		}

		if group.RealmId != realm {
			return fmt.Errorf("expected group with id %s to have realm_id of %s, but got %s", group.Id, realm, group.RealmId)
		}

		return nil
	}
}

func testAccCheckKeycloakGroupDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_group" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			group, _ := keycloakClient.GetGroup(realm, id)
			if group != nil {
				return fmt.Errorf("group with id %s still exists", id)
			}
		}

		return nil
	}
}

func getGroupFromState(s *terraform.State, resourceName string) (*keycloak.Group, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	group, err := keycloakClient.GetGroup(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting group with id %s: %s", id, err)
	}

	return group, nil
}

func testKeycloakGroup_basic(realm, group string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
	`, realm, group)
}

func testKeycloakGroup_updateRealmBefore(realmOne, realmTwo, group string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm_1.id}"
}
	`, realmOne, realmTwo, group)
}

func testKeycloakGroup_updateRealmAfter(realmOne, realmTwo, group string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_group" "group" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm_2.id}"
}
	`, realmOne, realmTwo, group)
}

func testKeycloakGroup_nested(realm, parentGroup, firstChildGroup, secondChildGroup, secondChildGroupParent string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_group" "parent_group" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_group" "first_child_group" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	parent_id = "${keycloak_group.parent_group.id}"
}

resource "keycloak_group" "second_child_group" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	parent_id = "${%s.id}"
}
	`, realm, parentGroup, firstChildGroup, secondChildGroup, secondChildGroupParent)
}
