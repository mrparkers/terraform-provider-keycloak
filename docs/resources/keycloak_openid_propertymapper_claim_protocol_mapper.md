---
page_title: "keycloak_openid_propertymapper_claim_protocol_mapper Resource"
---

# keycloak\_openid\_propertymapper\_claim\_protocol\_mapper Resource

Allows for creating and managing claim protocol mappers within Keycloak.

The property claim mappers allow you to define a claim with based on dynamic values to support latest keycloak apis.

Protocol mappers can be defined for a single client, or they can be defined for a client scope which can be shared between multiple different clients.

## Example Usage (Client)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client" "openid_client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client"

  name    = "client"
  enabled = true

  access_type         = "CONFIDENTIAL"
  valid_redirect_uris = [
    "http://localhost:8080/openid-callback"
  ]
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "userattribute_id_claim_mapper" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.openid_client.id
  name        = "property-mapper"

  claim_name  = "property"
  json_type     = "String"

  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-property-mapper"

  set {
    name = "user.attribute"
    value = "id"
  }
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "clientrole_claim_mapper" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.openid_client.id
  name        = "client-role-mapper"

  claim_name  = "clientrole"
  json_type     = "String"

  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-client-role-mapper"

  add_to_introspection_token = true
  add_to_id_token = true
  add_to_access_token = true
  add_to_userinfo = true
  add_to_lightweight_claim = true

  set {
    name = "multivalued"
    value = "false"
  }

  set {
    name = "usermodel.clientRoleMapping.clientId"
    value = "admin-cli"
  }

  set {
    name = "usermodel.clientRoleMapping.rolePrefix"
    value = "prefix"
  }
}
```

## Example Usage (Client Scope)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client_scope" "client_scope" {
  realm_id = keycloak_realm.realm.id
  name     = "client-scope"
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "userattribute_id_claim_mapper" {
  realm_id    = keycloak_realm.realm.id
  client_scope_id = keycloak_openid_client_scope.client_scope.id
  name        = "property-mapper"

  claim_name  = "property"
  json_type     = "String"

  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-property-mapper"

  set {
    name = "user.attribute"
    value = "id"
  }
}

resource "keycloak_openid_propertymapper_claim_protocol_mapper" "clientrole_claim_mapper" {
  realm_id    = keycloak_realm.realm.id
  client_scope_id = keycloak_openid_client_scope.client_scope.id
  name        = "client-role-mapper"

  claim_name  = "clientrole"
  json_type     = "String"

  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-client-role-mapper"

  add_to_introspection_token = true
  add_to_id_token = true
  add_to_access_token = true
  add_to_userinfo = true
  add_to_lightweight_claim = true

  set {
    name = "multivalued"
    value = "false"
  }

  set {
    name = "usermodel.clientRoleMapping.clientId"
    value = "admin-cli"
  }

  set {
    name = "usermodel.clientRoleMapping.rolePrefix"
    value = "prefix"
  }
}
```

## Argument Reference

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `claim_name` - (Required) The name of the claim to insert into a token.
- `claim_value` - (Required) The hardcoded value of the claim.
- `client_id` - (Optional) The client this protocol mapper should be attached to. Conflicts with `client_scope_id`. One of `client_id` or `client_scope_id` must be specified.
- `client_scope_id` - (Optional) The client scope this protocol mapper should be attached to. Conflicts with `client_id`. One of `client_id` or `client_scope_id` must be specified.
- `claim_value_type` - (Optional) The claim type used when serializing JSON tokens. Can be one of `String`, `JSON`, `long`, `int`, or `boolean`. Defaults to `String`.
- `add_to_id_token` - (Optional) Indicates if the property should be added as a claim to the id token. Defaults to `true`.
- `add_to_access_token` - (Optional) Indicates if the property should be added as a claim to the access token. Defaults to `true`.
- `add_to_userinfo` - (Optional) Indicates if the property should be added as a claim to the UserInfo response body. Defaults to `true`.
- `add_to_introspection_token` - (Optional) Indicates if the property should be added as a claim to the introspection token. Defaults to `true`.
- `add_to_lightweight_claim` - (Optional) Indicates if the property should be added as a lightweight claim. Defaults to `false`.
- `set` - (Block Set) Custom values to be merged with the values. (see below for nested schema)

### Nested Schema for `set`

Required:

- `name` (String)
- `value` (String)

## Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_openid_propertymapper_claim_protocol_mapper.claim_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_openid_propertymapper_claim_protocol_mapper.claim_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
