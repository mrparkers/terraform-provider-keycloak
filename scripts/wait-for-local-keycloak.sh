#!/usr/bin/env bash

set -x

KEYCLOAK_URL="http://localhost:8080/"

printf "Waiting for local Keycloak to be ready"

until $(curl --silent --head --fail --max-time 2 ${KEYCLOAK_URL}); do
    printf '.'
    sleep 2
done

echo
