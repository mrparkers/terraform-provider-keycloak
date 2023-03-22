# This Dockerfile is responsible for building images that are used for testing in CI
# This is needed because GitHub Actions doesn't support overriding the entrypoint of a service

FROM quay.io/keycloak/keycloak:${KEYCLOAK_VERSION}

ENTRYPOINT [ "/opt/keycloak/bin/kc.sh", "start-dev", "--features=preview" ]
