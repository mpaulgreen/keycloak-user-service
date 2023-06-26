BINARY_NAME=keycloak-user-service

# DEV TEMPLATE SETTINGS FOR LOCAL ENG WORK
DEV_IMAGE_TAG := latest
DEV_IMAGE_NAME := keycloak-user-service
DEV_REGISTRY ?= quay.io
DEV_REGISTRY_REPO ?= ecosystem-appeng
DEV_IMAGE ?= $(DEV_REGISTRY)/$(DEV_REGISTRY_REPO)/${DEV_IMAGE_NAME}

# PRODUCTION TEMPLATE SETTINGS FOR BUILDING CONTAINER
IMAGE_TAG := $(shell git rev-parse --short=8 HEAD)
IMG_NAME ?= keycloak-user-service
REGISTRY ?= quay.io
REGISTRY_REPO ?= app-sre
IMAGE ?= $(REGISTRY)/$(REGISTRY_REPO)/${IMG_NAME}

CONTAINER_ENGINE ?= $(shell which podman >/dev/null 2>&1 && echo podman || echo docker)
ifneq (,$(wildcard $(CURDIR)/.docker))
	DOCKER_CONF := $(CURDIR)/.docker
else
	DOCKER_CONF := $(HOME)/.docker
endif

.PHONY: all

all: test clean build

## Build:
build: ## Build your project and put the output binary in out/bin/
	mkdir -p out/bin
	go build -o out/bin/$(BINARY_NAME) .

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out

## Test:
tests: ## Run the tests of the project
	eval "export KEYCLOAK_BACKEND_URL=http://localhost:8080 && export ADMIN_USER=admin && export ADMIN_PASSWORD=admin && export GRANT_TYPE=password && export CLIENT_ID=admin-cli && export KEYCLOAK_REALM=master && export KEYCLOAK_USERS_RESOURCE_URI=admin/realms/master/users && export USER_SERVICE_PORT=8000 && export DISABLE_KEYCLOAK_CERT_VERIFICATION=true && go test ./test/..."

run-local-keycloak:
	@$(CONTAINER_ENGINE) run -p 8080:8080 --env KEYCLOAK_ADMIN=admin --env KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak  start-dev

run-local:
	. local-exec-env-vars.sh && go run .

## Docker:
docker-build: ## Use the dockerfile to build the container
	@$(CONTAINER_ENGINE) build --pull -t $(IMAGE):latest .
	@$(CONTAINER_ENGINE) tag $(IMAGE):latest $(IMAGE):$(IMAGE_TAG)

docker-push: ## push the image
	@$(CONTAINER_ENGINE) --config=$(DOCKER_CONF) push $(IMAGE):latest
	@$(CONTAINER_ENGINE) --config=$(DOCKER_CONF) push $(IMAGE):$(IMAGE_TAG)

dev-container-build:
	@$(CONTAINER_ENGINE) build --ulimit nofile=16384:65536 --pull -t $(DEV_IMAGE):latest .
	@$(CONTAINER_ENGINE) tag $(DEV_IMAGE):latest $(DEV_IMAGE):$(DEV_IMAGE_TAG)

dev-container-push:
	@$(CONTAINER_ENGINE) push $(DEV_IMAGE):latest
	@$(CONTAINER_ENGINE) push $(DEV_IMAGE):$(DEV_IMAGE_TAG)