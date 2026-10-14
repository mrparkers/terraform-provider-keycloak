---
page_title: "keycloak_organization Resource"
---

# keycloak\_organization Resource

Allow for creating and managing Organizations within Keycloak.

Attributes can also be defined on Groups.

Linkage with identity providers is managed with the identity provider resources, or, on Keycloak 26.8 or later, with
`identity_provider` blocks on the organization. Only `identity_provider` blocks can link one identity provider to
several organizations, and route each domain of an organization to its own identity provider.

## Example usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_organization" "this" {
  realm    = keycloak_realm.realm.name
  name     = "org"
  alias = "org"
  enabled = true
  
  domain {
    name = "example.com"
  }
}

resource "keycloak_oidc_identity_provider" "this" {
  realm             = keycloak_realm.realm.name
  alias             = "my-idp"
  authorization_url = "https://authorizationurl.com"
  client_id         = "clientID"
  client_secret     = "clientSecret"
  token_url         = "https://tokenurl.com"

  organization_id                 = keycloak_organization.this.id
  org_domain                      = "example.com"
  org_redirect_mode_email_matches = true
}

```

## Example usage with identity provider links (Keycloak 26.8+)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_oidc_identity_provider" "corporate" {
  realm             = keycloak_realm.realm.name
  alias             = "corporate"
  authorization_url = "https://authorizationurl.com"
  client_id         = "clientID"
  client_secret     = "clientSecret"
  token_url         = "https://tokenurl.com"
}

resource "keycloak_organization" "subsidiary" {
  realm = keycloak_realm.realm.name
  name  = "subsidiary"

  identity_provider {
    alias = keycloak_oidc_identity_provider.corporate.alias
  }

  domain {
    name                    = "subsidiary.example.com"
    identity_provider_alias = keycloak_oidc_identity_provider.corporate.alias
    auto_redirect           = true
  }
}

# The same identity provider, linked to a second organization without adding
# its users as members.
resource "keycloak_organization" "partner" {
  realm = keycloak_realm.realm.name
  name  = "partner"

  identity_provider {
    alias           = keycloak_oidc_identity_provider.corporate.alias
    auto_membership = false
  }

  domain {
    name                    = "partner.example.com"
    identity_provider_alias = keycloak_oidc_identity_provider.corporate.alias
  }
}
```

## Argument Reference

- `realm` - (Required) The realm this organization exists in.
- `name` - (Required) The name of the organization.
- `alias` - (Optional) The alias unique identifies the organization. Same as the name if not specified. The alias cannot be changed after the organization has been created.
- `description` - (Optional) The description of the organization.
- `redirect_url` - (Optional) The landing page after user completes registration or accepts an invitation to the organization. If left empty, the user will be redirected to the account console by default.
- `domain` - (Optional) A list of [domains](#domain-arguments).
- `identity_provider` - (Optional) A list of [identity providers](#identity-provider-arguments) that this resource links to the organization. Requires Keycloak 26.8 or later. Links that this resource did not make, such as those of an identity provider resource with `organization_id`, are kept. When set, this resource also manages the identity provider routing of every domain, so do not set `org_domain` on an identity provider resource for this organization.
- `attributes` - (Optional) A map representing attributes for the group. In order to add multivalued attributes, use `##` to separate the values. Max length for each value is 255 chars.

### Domain arguments

- `name` - (Required) The domain name
- `verified` - (Optional) Whether domain is verified or not. Default is false.
- `identity_provider_alias` - (Optional) The alias of an identity provider in an `identity_provider` block. Users whose email matches this domain authenticate through it.
- `auto_redirect` - (Optional) When `true`, users whose email matches this domain are redirected to the identity provider without seeing the login form. Requires `identity_provider_alias`. Default is false.

### Identity provider arguments

- `alias` - (Required) The alias of the identity provider.
- `auto_membership` - (Optional) When `true`, users who authenticate through the identity provider are added as members of the organization, if the organization has no domain or one of its domains matches their email. Default is true.
- `membership_type` - (Optional) The membership type of the users that `auto_membership` adds: `UNMANAGED`, or `MANAGED` to remove a user's account together with the membership. Keycloak allows `MANAGED` on one organization per identity provider, and only with `auto_membership`. Default is `UNMANAGED`.

## Import

Organizations can be imported using the format `{{realm_id}}/{{organization_id}}`, where `organization_id` is the unique ID that Keycloak
assigns to the organizations upon creation. This value can be found in the URI when editing this organization in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_organization.this my-realm/cec54914-b702-4c7b-9431-b407817d059a
```
