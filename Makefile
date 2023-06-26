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