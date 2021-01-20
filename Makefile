ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION ?= dev
IMAGE_REPOSITORY := epiphanyplatform/azbi

IMAGE_NAME := $(IMAGE_REPOSITORY):$(VERSION)

define SERVICE_PRINCIPAL_CONTENT
ARM_CLIENT_ID ?= $(CLIENT_ID)
ARM_CLIENT_SECRET ?= $(CLIENT_SECRET)
ARM_SUBSCRIPTION_ID ?= $(SUBSCRIPTION_ID)
ARM_TENANT_ID ?= $(TENANT_ID)
endef

-include ./service-principal.mk

export

#used for correctly setting shared folder permissions
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)

.PHONY: all

all: build

.PHONY: build test pipeline-test release prepare-service-principal

build: guard-IMAGE_NAME
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		--build-arg ARG_HOST_UID=$(HOST_UID) \
		--build-arg ARG_HOST_GID=$(HOST_GID) \
		-t $(IMAGE_NAME) \
		.

#prepare service principal variables file before running this target using `CLIENT_ID=xxx CLIENT_SECRET=yyy SUBSCRIPTION_ID=zzz TENANT_ID=vvv make prepare-service-principal`
#test targets are located in ./test.mk file
test: guard-IMAGE_REPOSITORY build
	$(eval LDFLAGS = $(shell govvv -flags -pkg github.com/epiphany-platform/m-azure-basic-infrastructure/cmd -version $(VERSION)))
	@AZURE_CLIENT_ID=$(ARM_CLIENT_ID) AZURE_CLIENT_SECRET=$(ARM_CLIENT_SECRET) AZURE_SUBSCRIPTION_ID=$(ARM_SUBSCRIPTION_ID) AZURE_TENANT_ID=$(ARM_TENANT_ID) go test -ldflags="$(LDFLAGS)" -v -timeout 30m

pipeline-test:
	$(eval LDFLAGS = $(shell govvv -flags -pkg github.com/epiphany-platform/m-azure-basic-infrastructure/cmd -version $(VERSION)))
	@go test -ldflags="$(LDFLAGS)" -v -timeout 30m

prepare-service-principal: guard-CLIENT_ID guard-CLIENT_SECRET guard-SUBSCRIPTION_ID guard-TENANT_ID
	@echo "$$SERVICE_PRINCIPAL_CONTENT" > $(ROOT_DIR)/service-principal.mk

release: guard-VERSION guard-IMAGE_NAME
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		-t $(IMAGE_NAME) \
		.

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

doctor:
	go mod tidy
	go fmt ./...
	go vet ./...
	goimports -l -w .
