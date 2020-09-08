ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION := 0.0.1
USER := epiphany
IMAGE := azbi

.PHONY: all clean init plan apply setup

all: init apply

clean:
	@rm -rf $(ROOT_DIR)/shared

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		init \
		M_VMS_COUNT=4 \
		M_PUBLIC_IPS=false \
		M_RG_NAME="mkyc-module-tests"

plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		plan \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

apply: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		apply \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy-plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID \
        M_DO_DESTROY=true

setup:
	@mkdir -p $(ROOT_DIR)/shared
