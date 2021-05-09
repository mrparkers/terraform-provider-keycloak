---
page_title: "keycloak_attribute_importer_identity_provider_mapper Resource"
---

# keycloak\_attribute\_importer\_identity\_provider\_mapper Resource

Allows for creating and managing an attribute importer identity provider mapper within Keycloak.

The attribute importer mapper can be used to map attributes from externally defined users to attributes or properties of the imported Keycloak user:
- For the OIDC identity provider, this will map a claim on the ID or access token to an attribute for the imported Keycloak user.
- For the SAML identity provider, this will map a SAML attribute found within the assertion to an attribute for the imported Keycloak user.
- For social identity providers, this will map a JSON field from the user profile to an attribute for the imported Keycloak user.

~> If you are using Keycloak 10 or higher, you will need to specify the `extra_config` argument in order to define a `syncMode` for the mapper.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_oidc_identity_provider" "oidc" {
  realm             = keycloak_realm.realm.id
  alias             = "oidc"
  authorization_url = "https://example.com/auth"
  token_url         = "https://example.com/token"
  client_id         = "example_id"
  client_secret     = "example_token"
  default_scopes    = "openid random profile"
}

resource "keycloak_attribute_importer_identity_provider_mapper" "oidc" {
  realm                   = keycloak_realm.realm.id
  name                    = "email-attribute-importer"
  claim_name              = "my-email-claim"
  identity_provider_alias = keycloak_oidc_identity_provider.oidc.alias
  user_attribute          = "email"

  # extra_config with syncMode is required in Keycloak 10+
  extra_config = {
    syncMode = "INHERIT"
  }
}
```

## Argument Reference

The following arguments are supported:

- `realm` - (Required) The name of the realm.
- `name` - (Required) The name of the mapper.
- `identity_provider_alias` - (Required) The alias of the associated identity provider.
- `user_attribute` - (Required) The user attribute or property name to store the mapped result.
- `attribute_name` - (Optional) For SAML based providers, this is the name of the attribute to search for in the assertion. Conflicts with `attribute_friendly_name`.
- `attribute_friendly_name` - (Optional) For SAML based providers, this is the friendly name of the attribute to search for in the assertion. Conflicts with `attribute_name`.
- `claim_name` - (Optional) For OIDC based providers, this is the name of the claim to use.
- `extra_config` - (Optional) Key/value attributes to add to the identity provider mapper model that is persisted to Keycloak. This can be used to extend the base model with new Keycloak features.

## Import

Identity provider mappers can be imported using the format `{{realm_id}}/{{idp_alias}}/{{idp_mapper_id}}`, where `idp_alias` is the identity provider alias, and `idp_mapper_id` is the unique ID that Keycloak
assigns to the mapper upon creation. This value can be found in the URI when editing this mapper in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_attribute_importer_identity_provider_mapper.test_mapper my-realm/my-mapper/f446db98-7133-4e30-b18a-3d28fde7ca1b
```
