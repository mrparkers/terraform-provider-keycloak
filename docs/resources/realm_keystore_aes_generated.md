---
page_title: "keycloak_realm_key_aes_generated Resources"
---

# keycloak\_realm\_key\_aes_generated Resources

Allows for creating and managing Realm keystores within Keycloak.

A realm manages a logical collection of users, credentials, roles, and groups. Users log in to realms and can be federated
from multiple sources.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
	realm = "my-realm"
}

resource "keycloak_realm_key_aes_generated" "keystore_aes_generated" {
	name      = "my-aes-generated-key"
	realm_id  = keycloak_realm.my_realm.realm

	enabled = true
	active  = true

	priority     = 100
	secret_size  = 16
}
```

## Argument Reference

- `name` - (Required) Display name of provider when linked in admin console.
- `realm_id` - (Required) The realm this keystore exists in.
- `enabled` - (Optional) When `false`, key is not accessible in this realm. Defaults to `true`.
- `active` - (Optional) When `false`, key in not used for signing. Defaults to `true`.
- `priority` - (Optional) Priority for the provider. Defaults to `0`
- `secret_size` - (Optional) Size in bytes for the generated AES Key. Size 16 is for AES-128, Size 24 for AES-192 and Size 32 for AES-256. WARN: Bigger keys then 128 bits are not allowed on some JDK implementations. Defaults to `16`.

## Import

Realm keys can be imported using realm name and keystore id, you can find it in web UI.

Example:

```bash
$ terraform import keycloak_realm_key_aes_generated.keystore_aes_generated my-realm/my-realm/618cfba7-49aa-4c09-9a19-2f699b576f0b
```
