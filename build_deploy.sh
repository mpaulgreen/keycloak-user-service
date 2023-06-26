#!/bin/bash


DOCKER_CONF="$PWD/.docker"
# shellcheck disable=SC2035
CONTAINER_ENGINE="$(shell which podman >/dev/null 2>&1 && echo podman || echo docker)"

mkdir -p "$DOCKER_CONF"
$CONTAINER_ENGINE --config="$DOCKER_CONF" login -u="$QUAY_USER" -p="$QUAY_TOKEN" quay.io

# build images
make docker-build docker-push
