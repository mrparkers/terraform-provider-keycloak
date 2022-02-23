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
resource "keycloak_realm" "my_realm" {
	realm = "my-realm"

  attributes = {
		userProfileEnabled = true
	}
}

resource keycloak_realm_user_profile userprofile {
	realm_id  = keycloak_realm.my_realm.realm

  attribute {
    name = "field1"
    display_name = "Field 1"
    group = "group1"

    enabled_when_scope = ["offline_access"]

    required_for_roles = ["user"]
    required_for_scopes = ["offline_access"]

    permissions {
      view = ["admin", "user"]
      edit = ["admin", "user"]
    }

    validator {
      name = "person-name-prohibited-characters"
    } 

    validator {
      name = "pattern"
      config = {
        pattern = "^[a-z]+$"
        error_message = "Nope"
      }
    }

    annotations = {
      foo = "bar"
    }
  }

  attribute {
    name = "field2"
  }

  group {
    name = "group1"
    display_header = "Group 1"
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

- `realm_id` - (Required) The name of the realm the user profile apply to.
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

- `name` - (Required) The name of the validator.
- `config` - (Optional) A map defining the configuration of the validator.

### Group Arguments

- `name` - (Required) The name of the group.
- `display_header` - (Optional) The display header of the group.
- `display_description` - (Optional) The display description of the group.
- `annotations` - (Optional) A map of annotations for the group.

## Import

This resource currently does not support importing.
