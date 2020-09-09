ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION := 0.0.1
USER := epiphany
IMAGE := azbi

#makes it easier to replace the value
M_NAME ?= mkyc-module-tests

-include $(ROOT_DIR)/azure.mk

export

.PHONY: all clean init plan apply destroy-plan destroy setup

all: init apply

clean:
	@rm -rf $(ROOT_DIR)/shared

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		init \
		M_VMS_COUNT=2 \
		M_PUBLIC_IPS=false \
		M_NAME=$(M_NAME)

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

setup: $(ROOT_DIR)/shared/azure_rsa

$(ROOT_DIR)/shared/azure_rsa:
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
