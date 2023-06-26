#!/bin/sh
export KEYCLOAK_REALM=master
export KEYCLOAK_USERS_RESOURCE_URI=admin/realms/${KEYCLOAK_REALM}/users
export USER_SERVICE_PORT=8000
export DISABLE_KEYCLOAK_CERT_VERIFICATION=true

# Overwrite the environment specific variables here
case $1 in

  docker)
    export KEYCLOAK_BACKEND_URL=http://0.0.0.0:8080
    ;;

  openshift-local)
    export KEYCLOAK_BACKEND_URL=http://keycloak-postgres-operator.apps-crc.testing:8443
    ;;

  openshift-dev)
    export KEYCLOAK_BACKEND_URL=https://keycloak.fips-test.svc.cluster.local:8443
    ;;

  *)
    export KEYCLOAK_BACKEND_URL=http://localhost:8080
    ;;
esac

