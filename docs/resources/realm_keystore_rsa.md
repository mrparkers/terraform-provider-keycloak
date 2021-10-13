---
page_title: "keycloak_realm_keystore_rsa Resources"
---

# keycloak\_realm\_keystore\_rsa Resources

Allows for creating and managing `rsa` Realm keystores within Keycloak.

A realm keystore manages generated key pairs that are used by Keycloak to perform cryptographic signatures and encryption.

## Example Usage

```hcl
resource "keycloak_realm" "my_realm" {
	realm = "my-realm"
}

resource "keycloak_realm_keystore_rsa" "keystore_rsa" {
	name      = "my-rsa-key"
	realm_id  = keycloak_realm.my_realm.realm

	enabled = true
	active  = true

	private_key = "<your rsa private key>"
	certificate = "<your certificate>"

	priority  = 100
	algorithm = "RS256"
	keystore_size  = 2048
}
```

## Argument Reference

- `name` - (Required) Display name of provider when linked in admin console.
- `realm_id` - (Required) The realm this keystore exists in.
- `private_key` - (Required) Private RSA Key encoded in PEM format.
- `certificate` - (Required) X509 Certificate encoded in PEM format.
- `enabled` - (Optional) When `false`, key is not accessible in this realm. Defaults to `true`.
- `active` - (Optional) When `false`, key in not used for signing. Defaults to `true`.
- `priority` - (Optional) Priority for the provider. Defaults to `0`
- `algorithm` - (Optional) Intended algorithm for the key. Defaults to `RS256`
- `keystore_size` - (Optional) Size for the generated keys. Defaults to `2048`.

## Import

Realm keys can be imported using realm name and keystore id, you can find it in web UI.

Example:

```bash
$ terraform import keycloak_realm_keystore_rsa.keystore_rsa my-realm/my-realm/618cfba7-49aa-4c09-9a19-2f699b576f0b
```
