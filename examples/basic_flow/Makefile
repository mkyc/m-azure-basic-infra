ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

AZBI := epiphanyplatform/azbi:0.0.1

#makes it easier to replace the value
M_NAME ?= azbi-module-tests

VMS_RSA_FILENAME = vms_rsa

-include $(ROOT_DIR)/azure.mk

export

.PHONY: all clean init plan apply output destroy-plan destroy setup

all: init plan apply

clean:
	@rm -rf $(ROOT_DIR)/shared

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		init \
		M_VMS_COUNT=2 \
		M_PUBLIC_IPS=false \
		M_NAME=$(M_NAME) \
		M_VMS_RSA=$(VMS_RSA_FILENAME)

plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		plan \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

apply: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		apply \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

output: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		output \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy-plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		plan-destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

destroy: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

setup: $(ROOT_DIR)/shared/$(VMS_RSA_FILENAME)

$(ROOT_DIR)/shared/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
