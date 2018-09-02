#!/usr/bin/env bash

KEYCLOAK_URL="http://localhost:8080"
KEYCLOAK_USER="keycloak"
KEYCLOAK_PASSWORD="password"
KEYCLOAK_CLIENT_ID="terraform"
KEYCLOAK_CLIENT_SECRET="884e0f95-0f42-4a63-9b1f-94274655669e"

echo "Creating initial terraform client"

accessToken=$(
    curl -s --fail \
        -d "username=${KEYCLOAK_USER}" \
        -d "password=${KEYCLOAK_PASSWORD}" \
        -d "client_id=admin-cli" \
        -d "grant_type=password" \
        "${KEYCLOAK_URL}/auth/realms/master/protocol/openid-connect/token" \
        | jq -r '.access_token'
)

function post() {
    curl --fail \
        -H "Authorization: bearer ${accessToken}" \
        -H "Content-Type: application/json" \
        -d "${2}" \
        "${KEYCLOAK_URL}/auth/admin${1}"
}

function get() {
    curl --fail --silent \
        -H "Authorization: bearer ${accessToken}" \
        -H "Content-Type: application/json" \
        "${KEYCLOAK_URL}/auth/admin${1}"
}

terraformClient=$(jq -n "{
    id: \"${KEYCLOAK_CLIENT_ID}\",
    name: \"${KEYCLOAK_CLIENT_ID}\",
    secret: \"${KEYCLOAK_CLIENT_SECRET}\",
    clientAuthenticatorType: \"client-secret\",
    enabled: true,
    serviceAccountsEnabled: true
}")

post "/realms/master/clients" "${terraformClient}"

masterRealmAdminRole=$(get "/realms/master/roles" | jq -r '
    .
    | map(
        select(.name == "admin")
    )
    | .[0]
')
masterRealmAdminRoleId=$(echo ${masterRealmAdminRole} | jq -r '.id')

terraformClientServiceAccount=$(get "/realms/master/clients/${KEYCLOAK_CLIENT_ID}/service-account-user")
terraformClientServiceAccountId=$(echo ${terraformClientServiceAccount} | jq -r '.id')

# [{"id":"6a4d83a2-4cc6-43c9-81b7-fa7440431eaa","name":"admin","description":"${role_admin}","composite":true,"clientRole":false,"containerId":"master"}]

serviceAccountAdminRoleMapping=$(jq -n "[{
    clientRole: false,
    composite: true,
    containerId: \"master\",
    description: \"\${role_admin}\",
    id: \"${masterRealmAdminRoleId}\",
    name: \"admin\",
}]")

post "/realms/master/users/${terraformClientServiceAccountId}/role-mappings/realm" "${serviceAccountAdminRoleMapping}"

echo "Done"
