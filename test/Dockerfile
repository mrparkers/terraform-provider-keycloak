# this Dockerfile is only used during acceptance tests within CI

FROM quay.io/keycloak/keycloak:19.0.2

ENTRYPOINT [ "/opt/keycloak/bin/kc.sh", "start-dev", "--features=preview" ]
