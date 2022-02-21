package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakRealmUserProfile_featureDisabled(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRealmUserProfileDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRealmUserProfile_featureDisabled(realmName),
				ExpectError: regexp.MustCompile("User Profile is disabled"),
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_basicEmpty(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")

	realmUserProfile := &keycloak.RealmUserProfile{}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRealmUserProfileDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_template(realmName, realmUserProfile),
				Check:  testAccCheckKeycloakRealmUserProfileExists("keycloak_realm_user_profile.realm_user_profile"),
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_basicFull(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")

	realmUserProfile := &keycloak.RealmUserProfile{
		Attributes: []*keycloak.RealmUserProfileAttribute{
			{Name: "attribute1"},
			{
				Name:        "attribute2",
				DisplayName: "attribute 2",
				Group:       "group",
				Selector:    &keycloak.RealmUserProfileSelector{Scopes: []string{"roles"}},
				Required: &keycloak.RealmUserProfileRequired{
					Roles:  []string{"user"},
					Scopes: []string{"offline_access"},
				},
				Permissions: &keycloak.RealmUserProfilePermissions{
					Edit: []string{"admin", "user"},
					View: []string{"admin", "user"},
				},
				Validations: map[string]keycloak.RealmUserProfileValidationConfig{
					"person-name-prohibited-characters": map[string]interface{}{},
					"pattern":                           map[string]interface{}{"pattern": "^[a-z]+$", "error_message": "Error!"},
				},
				Annotations: map[string]string{"foo": "bar"},
			},
		},
		Groups: []*keycloak.RealmUserProfileGroup{
			{Name: "group", DisplayDescription: "Description", DisplayHeader: "Header", Annotations: map[string]string{"foo": "bar"}},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRealmUserProfileDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_template(realmName, realmUserProfile),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", realmUserProfile,
				),
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_update(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")

	before := &keycloak.RealmUserProfile{
		Attributes: []*keycloak.RealmUserProfileAttribute{
			{Name: "attribute1"},
			{
				Name:        "attribute2",
				DisplayName: "attribute 2",
				Group:       "group",
				Selector:    &keycloak.RealmUserProfileSelector{Scopes: []string{"roles"}},
				Required: &keycloak.RealmUserProfileRequired{
					Roles:  []string{"user"},
					Scopes: []string{"offline_access"},
				},
				Permissions: &keycloak.RealmUserProfilePermissions{
					Edit: []string{"admin", "user"},
					View: []string{"admin", "user"},
				},
				Validations: map[string]keycloak.RealmUserProfileValidationConfig{
					"person-name-prohibited-characters": map[string]interface{}{},
					"pattern":                           map[string]interface{}{"pattern": "^[a-z]+$", "error_message": "Error!"},
				},
				Annotations: map[string]string{"foo": "bar"},
			},
		},
		Groups: []*keycloak.RealmUserProfileGroup{
			{Name: "group", DisplayDescription: "Description", DisplayHeader: "Header", Annotations: map[string]string{"foo": "bar"}},
		},
	}

	after := &keycloak.RealmUserProfile{
		Attributes: []*keycloak.RealmUserProfileAttribute{
			{Name: "attribute1"},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRealmUserProfileDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_template(realmName, before),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", before,
				)},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, after),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", after,
				),
			},
		},
	})
}

func testKeycloakRealmUserProfile_featureDisabled(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_realm_user_profile" "realm_user_profile" {
	realm_id = keycloak_realm.realm.id
}
`, realm)
}

func testKeycloakRealmUserProfile_template(realm string, realmUserProfile *keycloak.RealmUserProfile) string {
	tmpl, err := template.New("").Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(`
resource "keycloak_realm" "realm" {
	realm = "{{ .realm }}"
	user_profile_enabled = true
}

resource "keycloak_realm_user_profile" "realm_user_profile" {
	realm_id = keycloak_realm.realm.id

	{{- range $_, $attribute := .userProfile.Attributes }}
	attribute {
        name = "{{ $attribute.Name }}"
		{{- if $attribute.DisplayName }}
        display_name = "{{ $attribute.DisplayName }}"
		{{- end }}

		{{- if $attribute.Group }}
        group = "{{ $attribute.Group }}"
		{{- end }}

		{{- if $attribute.Selector }}
		{{- if $attribute.Selector.Scopes }}
        enabled_when_scope = ["{{ StringsJoin $attribute.Selector.Scopes "\", \"" }}"]
		{{- end }}
		{{- end }}

		{{- if $attribute.Required }}
		{{- if $attribute.Required.Roles }}
        required_for_roles = ["{{ StringsJoin $attribute.Required.Roles "\", \"" }}"]
		{{- end }}
		{{- end }}

		{{- if $attribute.Required }}
		{{- if $attribute.Required.Scopes }}
        required_for_scopes = ["{{ StringsJoin $attribute.Required.Scopes "\", \"" }}"]
		{{- end }}
		{{- end }}

		{{- if $attribute.Permissions }}
        permissions {
			{{- if $attribute.Permissions.View }}
            view = ["{{ StringsJoin $attribute.Permissions.View "\", \"" }}"]
			{{- end }}
			{{- if $attribute.Permissions.Edit }}
            edit = ["{{ StringsJoin $attribute.Permissions.Edit "\", \"" }}"]
			{{- end }}
        }
		{{- end }}

		{{- if $attribute.Validations }}
		{{ range $name, $config := $attribute.Validations }}
        validator {
            name = "{{ $name }}"
            {{- if $config }}
            config = {
                {{- range $key, $value := $config }}
                {{ $key }} = "{{ $value }}"
                {{- end }}
            }
            {{- end }}
        } 
		{{- end }}
		{{- end }}

		{{- if $attribute.Annotations }}
        annotations = {
            {{- range $key, $value := $attribute.Annotations }}
            {{ $key }} = "{{ $value }}"
            {{- end }}
        }
		{{- end }}
    }
	{{- end }}

	{{- range $_, $group := .userProfile.Groups }}
    group {
        name = "{{ $group.Name }}"

		{{- if $group.DisplayHeader }}
        display_header = "{{ $group.DisplayHeader }}"
		{{- end }}

		{{- if $group.DisplayDescription }}
        display_description = "{{ $group.DisplayDescription }}"
		{{- end }}

		{{- if $group.Annotations }}
        annotations = {
            {{- range $key, $value := $group.Annotations }}
            {{ $key }} = "{{ $value }}"
            {{- end }}
        }
		{{- end }}
    }
	{{- end }}
}
	`)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var tmplBuf bytes.Buffer
	err = tmpl.Execute(&tmplBuf, map[string]interface{}{"realm": realm, "userProfile": realmUserProfile})
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return tmplBuf.String()
}

func testAccCheckKeycloakRealmUserProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRealmUserProfileFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakRealmUserProfileStateEqual(resourceName string, realmUserProfile *keycloak.RealmUserProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		realmUserProfileFromState, err := getRealmUserProfileFromState(s, resourceName)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(realmUserProfile, realmUserProfileFromState) {
			j1, _ := json.Marshal(realmUserProfile)
			j2, _ := json.Marshal(realmUserProfileFromState)
			return fmt.Errorf("%v\nshould be equal to\n%v", string(j1), string(j2))
		}

		return nil
	}
}

func testAccCheckKeycloakRealmUserProfileDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm_user_profile" {
				continue
			}

			realm := rs.Primary.Attributes["realm_id"]

			realmUserProfile, _ := keycloakClient.GetRealmUserProfile(realm)
			if realmUserProfile != nil {
				return fmt.Errorf("user profile for realm %s", realm)
			}
		}

		return nil
	}
}

func getRealmUserProfileFromState(s *terraform.State, resourceName string) (*keycloak.RealmUserProfile, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]

	realmUserProfile, err := keycloakClient.GetRealmUserProfile(realm)
	if err != nil {
		return nil, fmt.Errorf("error getting realm user profile: %s", err)
	}

	return realmUserProfile, nil
}
