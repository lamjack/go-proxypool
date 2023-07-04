REGISTRY_HOST=docker.io
USERNAME=lamjack
NAMESPACE=lamjack
NAME=poolproxy

IMAGE=$(REGISTRY_HOST)/$(NAMESPACE)/$(NAME)
DOCKER_CONFIG_PATH=~/.docker/$(REGISTRY_HOST)_$(NAMESPACE)
VERSION:=$(shell date +'%y.%m.%d')-$(shell git rev-parse --short HEAD)

SHELL=/bin/bash

DOCKER_BUILD_CONTEXT=.
DOCKER_FILE_PATH=./Dockerfile
DOCKER_BUILD_ARGS=--platform linux/amd64

all: docker-build tag-docker do-push

docker-build:
	@echo -e "\033[32m docker build version:$(VERSION) \033[0m"
	docker build $(DOCKER_BUILD_ARGS) -t $(IMAGE):$(VERSION) $(DOCKER_BUILD_CONTEXT) -f $(DOCKER_FILE_PATH)

tag-docker:
	@echo create tag latest for API
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	@echo create tag $(VERSION) for API
	docker tag $(IMAGE):$(VERSION) $(IMAGE):$(VERSION)

do-push:
	@echo publish API latest to $(REGISTRY_HOST)
	docker --config $(DOCKER_CONFIG_PATH) push $(IMAGE):latest
	@echo publish $(IMAGE_API):$(VERSION) to $(REGISTRY_HOST)
	docker --config $(DOCKER_CONFIG_PATH) push $(IMAGE):$(VERSION)
	@echo -e "\033[32m $(IMAGE):$(VERSION) \033[0m"

registry-login:
	docker --config $(DOCKER_CONFIG_PATH) login --username=$(USERNAME) $(REGISTRY_HOST)