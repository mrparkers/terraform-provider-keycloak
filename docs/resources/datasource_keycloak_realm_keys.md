# datasource keycloak_realm_keys

Use this data source to get the keys of a realm. Keys can be filtered by algorithm and status.

Remarks:

- A key must meet all filter criteria
- This datasource may return more than one value.
- If no key matches the filter criteria, then an error is returned.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

data "keycloak_realm_keys" "keys" {
  realm_id  = keycloak_realm.realm
  algorithms = ["AES", "RS256"]
  status = ["ACTIVE", "PASSIVE"]
}

# show certificate of first key:
output "certificate" {
  value = data.keycloak_realm_keys.realm.keys[0].certificate
}

```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm of which the keys are retrieved.
- `algorithms` - (Optional) When specified, keys are filtered by algorithm (values for algorithm: `HS256`, `RS256`,`AES`, ...)
- `status` - (Optional) When specified, keys are filtered by status (values for status: `ACTIVE`, `DISABLED` and `PASSIVE`)
