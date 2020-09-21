ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION ?= 0.0.1
USER := epiphanyplatform
IMAGE := azbi

IMAGE_NAME := $(USER)/$(IMAGE):$(VERSION)

define SERVICE_PRINCIPAL_CONTENT
ARM_CLIENT_ID ?= $(CLIENT_ID)
ARM_CLIENT_SECRET ?= $(CLIENT_SECRET)
ARM_SUBSCRIPTION_ID ?= $(SUBSCRIPTION_ID)
ARM_TENANT_ID ?= $(TENANT_ID)
endef


-include ./service-principal.mk
include ./test.mk

export

#used for correctly setting shared folder permissions
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)

.PHONY: build test release metadata prepare-service-principal

build: guard-VERSION guard-IMAGE guard-USER
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		--build-arg ARG_HOST_UID=$(HOST_UID) \
		--build-arg ARG_HOST_GID=$(HOST_GID) \
		-t $(IMAGE_NAME) \
		.

#prepare service principal variables file before running this target using `CLIENT_ID=xxx CLIENT_SECRET=yyy SUBSCRIPTION_ID=zzz TENANT_ID=vvv make prepare-service-principal`
#test targets are located in ./test.mk file
test: build \
	test-default-config \
	test-config-with-variables \
	test-plan \
	test-apply \
	generate-report

prepare-service-principal: guard-CLIENT_ID guard-CLIENT_SECRET guard-SUBSCRIPTION_ID guard-TENANT_ID
	@echo "$$SERVICE_PRINCIPAL_CONTENT" > $(ROOT_DIR)/tests/service-principal.mk

release: guard-VERSION guard-IMAGE guard-USER
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		-t $(IMAGE_NAME) \
		.

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
