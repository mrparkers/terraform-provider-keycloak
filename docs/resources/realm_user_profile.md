---
page_title: "keycloak_realm_user_profile Resource"
---

# keycloak_realm_user_profile Resource

Allows for managing Realm User Profiles within Keycloak.

A user profile defines a schema for representing user attributes and how they are managed within a realm.
This is a preview feature, hence not fully supported and disabled by default.
To enable it, start the server with one of the following flags:
- WildFly distribution: `-Dkeycloak.profile.feature.declarative_user_profile=enabled`
- Quarkus distribution: `--features=preview` or `--features=declarative-user-profile`

The realm linked to the `keycloak_realm_user_profile` resource must have the user profile feature enabled.
It can be done via the administration UI, or by setting the `userProfileEnabled` realm attribute to `true`.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm = "my-realm"

  attributes = {
    userProfileEnabled = true
  }
}

resource "keycloak_realm_user_profile" "userprofile" {
  realm_id = keycloak_realm.my_realm.id

  attribute {
    name         = "field1"
    display_name = "Field 1"
    group        = "group1"

    enabled_when_scope = ["offline_access"]

    required_for_roles  = ["user"]
    required_for_scopes = ["offline_access"]

    permissions {
      view = ["admin", "user"]
      edit = ["admin", "user"]
    }

    validator {
      person_name_prohibited_characters {}
      pattern {
        pattern       = "^[a-z]+$"
        error_message = "Nope"
      }
      length {
        min           = 1
        max           = 10
        trim_disabled = false
      }
    }

    annotations = {
      foo = "bar"
    }
  }

  attribute {
    name = "field2"
  }

  attribute {
    name = "field3"
    display_name = "Field 3"

    validator {
      options {
        options = ["option1", "option2"]
      }
      email {}
    }
  
  }


  attribute {
    name = "field4"

    validator {
      double {
        max = 5.5
        min = 1.5
      }
    }
  }

  group {
    name                = "group1"
    display_header      = "Group 1"
    display_description = "A first group"

    annotations = {
      foo = "bar"
    }
  }

  group {
    name = "group2"
  }
}
```

## Argument Reference

- `realm_id` - (Required) The ID of the realm the user profile applies to.
- `attribute` - (Optional) An ordered list of [attributes](#attribute-arguments).
- `group` - (Optional) A list of [groups](#group-arguments).

### Attribute Arguments

- `name` - (Required) The name of the attribute.
- `display_name` - (Optional) The display name of the attribute.
- `group` - (Optional) The group that the attribute belong to.
- `enabled_when_scope` - (Optional) A list of scopes. The attribute will only be enabled when these scopes are requested by clients.
- `required_for_roles` - (Optional) A list of roles for which the attribute will be required.
- `required_for_scopes` - (Optional) A list of scopes for which the attribute will be required.
- `permissions` - (Optional) The [permissions](#permissions-arguments) configuration information.
- `validator` - (Optional) A list of [validators](#validator-arguments) for the attribute.
- `annotations` - (Optional) A map of annotations for the attribute.

#### Permissions Arguments

- `edit` - (Optional) A list of profiles that will be able to edit the attribute. One of `admin`, `user`.
- `view` - (Optional) A list of profiles that will be able to view the attribute. One of `admin`, `user`.

#### Validator Arguments

- `length` - (Optional) Check the length of a string value [length](#length).
- `integer` - (Optional) Check if the value is an integer and within a lower and upper range [integer](#integer).
- `double` - (Optional) Check if the value is a double and within a lower and upper range [double](#double).
- `uri` - (Optional) Check if the value is a valid URI. Has no arguments.
- `pattern` - (Optional) Check if the value matches a specific RegEx pattern [pattern](#pattern)
- `email` - (Optional) Check if the value has a valid e-mail format. Has no arguments.
- `local_date` - (Optional) Check if the value has a valid format based on the realm and/or user locale. Has no arguments.
- `person_name_prohibited_characters` - (Optional) Check if the value is a valid person name as an additional barrier for attacks such as script injection. The validation is based on a default RegEx pattern that blocks characters not common in person names [person_name_prohibited_characters](#personnameprohibitedcharacters).
- `username_prohibited_characters` - (Optional) Check if the value is a valid username as an additional barrier for attacks such as script injection. The validation is based on a default RegEx pattern that blocks characters not common in usernames [username_prohibited_characters](#usernameprohibitedcharacters).
- `options` - (Optional) Check if the value is from the defined set of allowed values. Useful to validate values entered through select and multiselect fields[options](#options).

#### Length

- `min` - (Required) Minimum character length.
- `max`- (Required) Maximum character length.
- `trim_disabled` - (Optional) Whether value is trimmed prior to validation.

#### Integer

- `min` - (Required) Defines lower range.
- `max`- (Required) Defines upper range.

#### Double

- `min` - (Required) Defines lower range.
- `max`- (Required) Defines upper range.

#### Pattern

- `pattern`- (Required) The RegEx pattern to use when validating values.
- `error_message`- (Optional) The key of the error message in i18n bundle. If not set a generic message is used.

#### Person_name_prohibited_characters

- `error_message`- (Optional) The key of the error message in i18n bundle. If not set a generic message is used.

#### User_name_prohibited_characters

- `error_message`- (Optional) The key of the error message in i18n bundle. If not set a generic message is used.

#### Options

- `options`- (Required) Array of strings containing allowed values.

### Group Arguments

- `name` - (Required) The name of the group.
- `display_header` - (Optional) The display header of the group.
- `display_description` - (Optional) The display description of the group.
- `annotations` - (Optional) A map of annotations for the group.

## Import

This resource currently does not support importing.
