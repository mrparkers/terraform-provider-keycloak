---
page_title: "Keycloak Provider"
---

# Keycloak Provider

The Keycloak provider can be used to interact with [Keycloak](https://www.keycloak.org/).

## A note for users of the legacy Wildfly distribution

Recently, Keycloak has been updated to use Quarkus over the legacy Wildfly distribution. The only significant change here
that affects this Terraform provider is the removal of `/auth` from the default context path for the Keycloak API.

If you are using the legacy Wildfly distribution of Keycloak, you will need to set the `base_path` provider argument to
`/auth`. This can also be done by using the `KEYCLOAK_BASE_PATH` environment variable.

## Keycloak Setup

This Terraform provider can be configured to use the [client credentials](https://www.oauth.com/oauth2-servers/access-tokens/client-credentials/)
or [password](https://www.oauth.com/oauth2-servers/access-tokens/password-grant/) grant types. If you aren't
sure which to use, the client credentials grant is recommended, as it was designed for machine to machine authentication.

### Client Credentials Grant Setup (recommended)

1. Create a new client using the `openid-connect` protocol. This client can be created in the `master` realm if you would
like to manage your entire Keycloak instance, or in any other realm if you only want to manage that realm.
1. Update the client you just created:
    1. Set `Access Type` to `confidential`.
    1. Set `Standard Flow Enabled` to `OFF`.
    1. Set `Direct Access Grants Enabled` to `OFF`
    1. Set `Service Accounts Enabled` to `ON`.
1. Grant required roles for managing Keycloak via the `Service Account Roles` tab in the client you created in step 1, see [Assigning Roles](#assigning-roles) section below.

### Password Grant Setup

These steps will assume that you are using the `admin-cli` client, which is already correctly configured for this type
of authentication. Do not follow these steps if you have already followed the steps for the client credentials grant.

1. Create or identify the user whose credentials will be used for authentication.
1. Edit this user in the "Users" section of the management console and assign roles using the "Role Mappings" tab.

### Assigning Roles

There are many ways that roles can be assigned to manage Keycloak. Here are a couple of common scenarios accompanied
by suggested roles to assign. This is not an exhaustive list, and there is often more than one way to assign a particular set
of permissions.

- Managing the entire Keycloak instance: Assign the `admin` role to a user or service account within the `master` realm.
- Managing the entire `foo` realm: Assign the `realm-admin` client role from the `realm-management` client to a user or service
account within the `foo` realm.
- Managing clients for all realms within the entire Keycloak instance: Assign the `create-client` client role from each of
the realm clients to a user or service account within the `master` realm. For example, given a Keycloak instance with realms
`master`, `foo`, and `bar`, assign the `create-client` client role from the clients `master-realm`, `foo-realm`, and `bar-realm`.

## Example Usage (client credentials grant)

```hcl
provider "keycloak" {
	client_id     = "terraform"
	client_secret = "884e0f95-0f42-4a63-9b1f-94274655669e"
	url           = "http://localhost:8080"
}
```

## Example Usage (password grant)

```hcl
provider "keycloak" {
	client_id     = "admin-cli"
	username      = "keycloak"
	password      = "password"
	url           = "http://localhost:8080"
}
```

## Argument Reference

The following arguments are supported:

- `client_id` - (Required) The `client_id` for the client that was created in the "Keycloak Setup" section. Use the `admin-cli` client if you are using the password grant. Defaults to the environment variable `KEYCLOAK_CLIENT_ID`.
- `url` - (Required) The URL of the Keycloak instance, before `/auth/admin`. Defaults to the environment variable `KEYCLOAK_URL`.
- `client_secret` - (Optional) The secret for the client used by the provider for authentication via the client credentials grant. This can be found or changed using the "Credentials" tab in the client settings. Defaults to the environment variable `KEYCLOAK_CLIENT_SECRET`. This attribute is required when using the client credentials grant, and cannot be set when using the password grant.
- `client_assertion` - (Optional) use a jwt signed by a known signer to get a token exchaned for, .
- `username` - (Optional) The username of the user used by the provider for authentication via the password grant. Defaults to the environment variable `KEYCLOAK_USER`. This attribute is required when using the password grant, and cannot be set when using the client credentials grant.
- `password` - (Optional) The password of the user used by the provider for authentication via the password grant. Defaults to the environment variable `KEYCLOAK_PASSWORD`. This attribute is required when using the password grant, and cannot be set when using the client credentials grant.
- `realm` - (Optional) The realm used by the provider for authentication. Defaults to the environment variable `KEYCLOAK_REALM`, or `master` if the environment variable is not specified.
- `initial_login` - (Optional) Optionally avoid Keycloak login during provider setup, for when Keycloak itself is being provisioned by terraform. Defaults to true, which is the original method.
- `client_timeout` - (Optional) Sets the timeout of the client when addressing Keycloak, in seconds. Defaults to the environment variable `KEYCLOAK_CLIENT_TIMEOUT`, or `5` if the environment variable is not specified.
- `tls_insecure_skip_verify` - (Optional) Allows ignoring insecure certificates when set to `true`. Defaults to `false`. Disabling this security check is dangerous and should only be done in local or test environments.
- `root_ca_certificate` - (Optional) Allows x509 calls using an unknown CA certificate (for development purposes)
- `base_path` - (Optional) The base path used for accessing the Keycloak REST API.  Defaults to the environment variable `KEYCLOAK_BASE_PATH`, or an empty string if the environment variable is not specified. Note that users of the legacy distribution of Keycloak will need to set this attribute to `/auth`.
- `additional_headers` - (Optional) A map of custom HTTP headers to add to each request to the Keycloak API.



## Example JWT signer code:

```
#!/bin/bash

# Variables
KMS_KEY_ARN="$1"
ALGORITHM="RSASSA_PKCS1_V1_5_SHA_256"
CLIENT_ID="test"
AUTH_SERVER_TOKEN_ENDPOINT_URL="$2"

# Helper functions
base64url_encode() {
    openssl enc -base64 -A | tr '+/' '-_' | tr -d '='
}

generate_jti() {
    # Generates a random string for jti using /dev/urandom
    head -c 16 /dev/urandom | base64 | tr -d '/+=' | cut -c -22
}


# JWT Claims
ISSUER="$CLIENT_ID" # Issuer
SUBJECT="$CLIENT_ID" # Subject
AUDIENCE="$AUTH_SERVER_TOKEN_ENDPOINT_URL" # Audience
JTI=$(generate_jti) # JWT ID
CURRENT_TIME=$(date +%s)
EXPIRATION_TIME=$(($CURRENT_TIME + 3600)) # 1 hour from now

# Create JWT Header
HEADER=$(jq -n --arg alg "RS256" --arg typ "JWT" '{alg: $alg, typ: $typ}')
HEADER_BASE64=$(echo -n "$HEADER" | base64url_encode)

# Create JWT Payload
PAYLOAD=$(jq -n \
    --arg iss "$ISSUER" \
    --arg sub "$SUBJECT" \
    --arg aud "$AUDIENCE" \
    --arg jti "$JTI" \
    --argjson iat $CURRENT_TIME \
    --argjson exp $EXPIRATION_TIME \
    '{iss: $iss, sub: $sub, aud: $aud, jti: $jti, iat: $iat, exp: $exp}')
PAYLOAD_BASE64=$(echo -n "$PAYLOAD" | base64url_encode)

# Prepare the message to be signed
MESSAGE="$HEADER_BASE64.$PAYLOAD_BASE64"

# Use AWS KMS to sign the message
SIGNATURE_JSON=$(echo -n "$MESSAGE" | aws kms sign \
    --key-id "$KMS_KEY_ARN" \
    --message-type RAW \
    --signing-algorithm "$ALGORITHM" \
    --message fileb:///dev/stdin \
    --output json)
SIGNATURE=$(echo $SIGNATURE_JSON | jq -r '.Signature' | base64 -d | base64url_encode)

# Construct the JWT
JWT="$HEADER_BASE64.$PAYLOAD_BASE64.$SIGNATURE"

# echo "jwt = \"$JWT\"" > terraform.tfvars
# echo "$JWT" > /tmp/test

EXCHANGED_JWT=$(curl --location --request POST ${AUTH_SERVER_TOKEN_ENDPOINT_URL} \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer' \
--data-urlencode "client_assertion=${JWT}" \
--data-urlencode 'grant_type=client_credentials' | base64 | tr -d '\n')

echo "token_info_json = \"$EXCHANGED_JWT\"" > terraform.tfvars
```

you can then use it via:

```
provider "keycloak" {
  token_info_json = base64decode(var.token_info_json)
  url             = "http://localhost:8080"
  initial_login = false
}

```