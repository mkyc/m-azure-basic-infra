ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

AZBI := epiphanyplatform/azbi:latest

#makes it easier to replace the value
M_NAME ?= azbi-module-example

VMS_RSA_FILENAME = vms_rsa

-include $(ROOT_DIR)/azure.mk

export

.PHONY: all clean metadata init plan apply output destroy-plan destroy setup

all: init plan apply

clean:
	@rm -rf $(ROOT_DIR)/shared

metadata:
	@docker run --rm \
		-t $(AZBI) \
		metadata

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		init \
		--name=$(M_NAME) \
		--vms_rsa=$(VMS_RSA_FILENAME)

plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		plan

apply: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		apply

output: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AZBI) \
		output

destroy-plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		plan \
		--destroy

destroy: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e CLIENT_ID=$$ARM_CLIENT_ID \
		-e CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		-e SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		-e TENANT_ID=$$ARM_TENANT_ID \
		-t $(AZBI) \
		destroy

setup: $(ROOT_DIR)/shared/$(VMS_RSA_FILENAME)

$(ROOT_DIR)/shared/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
