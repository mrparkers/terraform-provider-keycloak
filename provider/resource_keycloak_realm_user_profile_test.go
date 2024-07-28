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

var (
	default_userprofile_attributs = []*keycloak.RealmUserProfileAttribute{
		{
			Name:        "username",
			DisplayName: "${username}",
			Multivalued: false,

			Permissions: &keycloak.RealmUserProfilePermissions{
				Edit: []string{"admin", "user"},
				View: []string{"admin", "user"},
			},
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"length":                            map[string]interface{}{"min": "3", "max": "255"},
				"person-name-prohibited-characters": map[string]interface{}{},
				"up-username-not-idn-homograph":     map[string]interface{}{},
			},
		},
		{
			Name:        "email",
			DisplayName: "${email}",
			Multivalued: false,

			Permissions: &keycloak.RealmUserProfilePermissions{
				Edit: []string{"admin", "user"},
				View: []string{"admin", "user"},
			},
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"email":  map[string]interface{}{},
				"length": map[string]interface{}{"max": "255"},
			},
			Required: &keycloak.RealmUserProfileRequired{
				Roles: []string{"user"},
			},
		},
		{
			Name:        "firstName",
			DisplayName: "${firstName}",
			Multivalued: false,

			Permissions: &keycloak.RealmUserProfilePermissions{
				Edit: []string{"admin", "user"},
				View: []string{"admin", "user"},
			},
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"length":                            map[string]interface{}{"max": "255"},
				"person-name-prohibited-characters": map[string]interface{}{},
			},
			Required: &keycloak.RealmUserProfileRequired{
				Roles: []string{"user"},
			},
		},
		{
			Name:        "lastName",
			DisplayName: "${lastName}",
			Multivalued: false,

			Permissions: &keycloak.RealmUserProfilePermissions{
				Edit: []string{"admin", "user"},
				View: []string{"admin", "user"},
			},
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"length":                            map[string]interface{}{"max": "255"},
				"person-name-prohibited-characters": map[string]interface{}{},
			},
			Required: &keycloak.RealmUserProfileRequired{
				Roles: []string{"user"},
			},
		},
	}
	default_userprofile_groups = []*keycloak.RealmUserProfileGroup{
		{
			Name:               "user-metadata",
			DisplayHeader:      "User metadata",
			DisplayDescription: "Attributes, which refer to user metadata",
		},
	}
)

func TestAccKeycloakRealmUserProfile_Importer(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-realm-acc")
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_UnmanagedAttributePolicy(realmName, "ENABLED"),
				Check:  testAccCheckKeycloakRealmUserProfileUnmanagedAttributePolicy("keycloak_realm_user_profile.realm_user_profile", "ENABLED"),
			},
			{
				ResourceName:      "keycloak_realm_user_profile.realm_user_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_UnmanagedAttributePolicy(t *testing.T) {

	realmName := acctest.RandomWithPrefix("tf-acc")
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_UnmanagedAttributePolicy(realmName, "DISABLED"),
				Check:  testAccCheckKeycloakRealmUserProfileUnmanagedAttributePolicy("keycloak_realm_user_profile.realm_user_profile", "DISABLED"),
			},
			{
				Config: testKeycloakRealmUserProfile_UnmanagedAttributePolicy(realmName, "ENABLED"),
				Check:  testAccCheckKeycloakRealmUserProfileUnmanagedAttributePolicy("keycloak_realm_user_profile.realm_user_profile", "ENABLED"),
			},
			{
				Config: testKeycloakRealmUserProfile_UnmanagedAttributePolicy(realmName, "ADMIN_VIEW"),
				Check:  testAccCheckKeycloakRealmUserProfileUnmanagedAttributePolicy("keycloak_realm_user_profile.realm_user_profile", "ADMIN_VIEW"),
			},
			{
				Config: testKeycloakRealmUserProfile_UnmanagedAttributePolicy(realmName, "ADMIN_EDIT"),
				Check:  testAccCheckKeycloakRealmUserProfileUnmanagedAttributePolicy("keycloak_realm_user_profile.realm_user_profile", "ADMIN_EDIT"),
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_basicFull(t *testing.T) {

	skipIfVersionIsLessThanOrEqualTo(testCtx, t, keycloakClient, keycloak.Version_14)

	realmName := acctest.RandomWithPrefix("tf-acc")

	attributes := default_userprofile_attributs
	attributes = append(attributes, &keycloak.RealmUserProfileAttribute{Name: "attribute1"})
	attributes = append(attributes, &keycloak.RealmUserProfileAttribute{
		Name:        "attribute2",
		DisplayName: "attribute 2",
		Multivalued: false,
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
			"pattern":                           map[string]interface{}{"pattern": "\"^[a-z]+$\"", "error_message": "\"Error!\""},
		},
		Annotations: map[string]interface{}{
			"foo":               "\"bar\"",
			"inputOptionLabels": "{\"a\":\"b\"}",
		},
	})

	groups := default_userprofile_groups
	groups = append(groups, &keycloak.RealmUserProfileGroup{
		Name:               "group",
		DisplayDescription: "Description",
		DisplayHeader:      "Header",
		Annotations: map[string]interface{}{
			"foo":  "\"bar\"",
			"test": "{\"a2\":\"b2\"}",
		},
	})

	realmUserProfile := &keycloak.RealmUserProfile{
		Attributes: attributes,
		Groups:     groups,
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
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

func TestAccKeycloakRealmUserProfile_group(t *testing.T) {

	skipIfVersionIsLessThanOrEqualTo(testCtx, t, keycloakClient, keycloak.Version_14)

	realmName := acctest.RandomWithPrefix("tf-acc")

	withoutGroup := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{Name: "attribute"}),
		Groups:     default_userprofile_groups,
	}

	withGroup := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{Name: "attribute"}),
		Groups:     append(default_userprofile_groups, &keycloak.RealmUserProfileGroup{Name: "group"}),
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withoutGroup),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withoutGroup,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withGroup),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withGroup,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withoutGroup),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withoutGroup,
				),
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_attributeValidator(t *testing.T) {

	skipIfVersionIsLessThanOrEqualTo(testCtx, t, keycloakClient, keycloak.Version_14)

	realmName := acctest.RandomWithPrefix("tf-acc")

	withoutValidator := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{Name: "attribute"}),
		Groups:     default_userprofile_groups,
	}

	withInitialConfig := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"length":  map[string]interface{}{"min": "5", "max": "10"},
				"options": map[string]interface{}{"options": "[\"cgu\"]"},
			},
		}),
		Groups: default_userprofile_groups,
	}

	withNewConfig := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"length": map[string]interface{}{"min": "6", "max": "10"},
			},
		}),
		Groups: default_userprofile_groups,
	}

	withNewValidator := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
			Validations: map[string]keycloak.RealmUserProfileValidationConfig{
				"person-name-prohibited-characters": map[string]interface{}{},
				"length":                            map[string]interface{}{"min": "6", "max": "10"},
			},
		}),
		Groups: default_userprofile_groups,
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withoutValidator),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withoutValidator,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withInitialConfig),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withInitialConfig,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withNewConfig),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withNewConfig,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withNewValidator),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withNewValidator,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withNewConfig),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withNewConfig,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withoutValidator),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withoutValidator,
				),
			},
		},
	})
}

func TestAccKeycloakRealmUserProfile_attributePermissions(t *testing.T) {

	skipIfVersionIsLessThanOrEqualTo(testCtx, t, keycloakClient, keycloak.Version_14)

	realmName := acctest.RandomWithPrefix("tf-acc")

	withoutPermissions := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
		}),
		Groups: default_userprofile_groups,
	}

	viewAttributeMissing := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
			Permissions: &keycloak.RealmUserProfilePermissions{
				Edit: []string{"admin", "user"},
			},
		}),
		Groups: default_userprofile_groups,
	}

	editAttributeMissing := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
			Permissions: &keycloak.RealmUserProfilePermissions{
				View: []string{"admin", "user"},
			},
		}),
		Groups: default_userprofile_groups,
	}

	bothAttributesMissing := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name:        "attribute",
			Permissions: &keycloak.RealmUserProfilePermissions{},
		}),
		Groups: default_userprofile_groups,
	}

	withRightPermissions := &keycloak.RealmUserProfile{
		Attributes: append(default_userprofile_attributs, &keycloak.RealmUserProfileAttribute{
			Name: "attribute",
			Permissions: &keycloak.RealmUserProfilePermissions{
				Edit: []string{"admin", "user"},
				View: []string{"admin", "user"},
			},
		}),
		Groups: default_userprofile_groups,
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withoutPermissions),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withoutPermissions,
				),
			},
			{
				Config:      testKeycloakRealmUserProfile_template(realmName, viewAttributeMissing),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testKeycloakRealmUserProfile_template(realmName, editAttributeMissing),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testKeycloakRealmUserProfile_template(realmName, bothAttributesMissing),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withRightPermissions),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withRightPermissions,
				),
			},
			{
				Config: testKeycloakRealmUserProfile_template(realmName, withoutPermissions),
				Check: testAccCheckKeycloakRealmUserProfileStateEqual(
					"keycloak_realm_user_profile.realm_user_profile", withoutPermissions,
				),
			},
		},
	})
}

func testKeycloakRealmUserProfile_UnmanagedAttributePolicy(realm string, unmanagedAttributePolicy string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_realm_user_profile" "realm_user_profile" {
	realm_id = keycloak_realm.realm.id

	unmanaged_attribute_policy = "%s"
	
	attribute {
		name= "username"
		display_name= "$${username}"
	
		validator {
			name = "length"
			config = {
				min= "3"
				max= "255"
			}
		}
		validator {
			name = "person-name-prohibited-characters"
		}
		validator {
			name = "up-username-not-idn-homograph"
		}
	
		permissions {
			view = ["admin", "user"]
			edit = ["admin", "user"]
		}
	
		multivalued = false
	}
	
	attribute {
		name= "email"
		display_name= "$${email}"
	
		validator {
			name = "email"
		}
		validator {
			name = "length"
			config = {
				max= "255"
			}
		}
	
		required_for_roles = [ "user" ]
	
		permissions {      
			view = ["admin", "user"]
			edit = ["admin", "user"]
		}
	
		multivalued = false
	}
	
	attribute {
		name= "firstName"
		display_name= "$${firstName}"
	
		validator {
			name = "length"
			config = {
				max= "255"
			}
		}
		validator {
			name = "person-name-prohibited-characters"
		}
	
		required_for_roles = [ "user" ]
	
		permissions {      
			view = ["admin", "user"]      
			edit = ["admin", "user"]
		}
	
		multivalued = false
	}
	
	attribute {
		name= "lastName"
		display_name= "$${lastName}"
	
		validator {
			name = "length"
			config = {
				max= "255"
			}
		}
		validator {
			name = "person-name-prohibited-characters"
		}
	
		required_for_roles = [ "user" ]
	
		permissions {      
			view = ["admin", "user"]      
			edit = ["admin", "user"]
		}
	
		multivalued = false
	}
	
	group {
		name               = "user-metadata"
		display_header      = "User metadata"
		display_description = "Attributes, which refer to user metadata"
	}
}
`, realm, unmanagedAttributePolicy)
}

func testKeycloakRealmUserProfile_template(realm string, realmUserProfile *keycloak.RealmUserProfile) string {
	tmpl, err := template.New("").Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(`
resource "keycloak_realm" "realm" {
	realm 	   = "{{ .realm }}"
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

		multivalued = "{{ $attribute.Multivalued }}"

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
                {{ $key }} = jsonencode ( {{ $value }} )
                {{- end }}
            }
            {{- end }}
        }
		{{- end }}
		{{- end }}

		{{- if $attribute.Annotations }}
        annotations = {
            {{- range $key, $value := $attribute.Annotations }}
            {{ $key }} = jsonencode ( {{ $value }} )
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
            {{ $key }} = jsonencode ( {{ $value }} )
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

	return strings.Replace(tmplBuf.String(), "${", "$${", -1)
}

func testAccCheckKeycloakRealmUserProfileUnmanagedAttributePolicy(resourceName string, expectedUnmanagedAttributePolicy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		realm_user_profile, err := getRealmUserProfileFromState(s, resourceName)
		if err != nil {
			return err
		}

		if expectedUnmanagedAttributePolicy == "DISABLED" {
			if len(realm_user_profile.UnmanagedAttributePolicy) != 0 {
				return fmt.Errorf("expected UnmanagedAttributePolicy value empty (%s), but was %s", expectedUnmanagedAttributePolicy, realm_user_profile.UnmanagedAttributePolicy)
			}
		} else if realm_user_profile.UnmanagedAttributePolicy != expectedUnmanagedAttributePolicy {
			return fmt.Errorf("expected UnmanagedAttributePolicy value %s, but was %s", expectedUnmanagedAttributePolicy, realm_user_profile.UnmanagedAttributePolicy)
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
			return fmt.Errorf("\n%v\nshould be equal to\n%v", string(j1), string(j2))
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

	realmUserProfile, err := keycloakClient.GetRealmUserProfile(testCtx, realm)
	if err != nil {
		return nil, fmt.Errorf("error getting realm user profile: %s", err)
	}

	return realmUserProfile, nil
}
