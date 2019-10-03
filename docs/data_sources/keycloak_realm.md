# keycloak_realm data source

This data source can be used to fetch properties of a Keycloak realm for
usage with other resources.

### Example Usage

```hcl
data "keycloak_realm" "realm" {
    realm   = "my-realm"
}

# use the data source

resource "keycloak_role" "group" {
    realm_id = "${data.keycloak_realm.id}"
    name     = "group"
}

```

### Argument Reference

The following arguments are supported:

- `realm` - (Required) The realm name.

### Attributes Reference

See the docs for the [`keycloak_realm` resource](../resources/keycloak_realm.md) for details on the exported attributes.
