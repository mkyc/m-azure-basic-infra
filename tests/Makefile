ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

IMAGE_NAME := epiphanyplatform/azbi:0.0.1

M_NAME ?= azbi-module-tests

VMS_RSA_FILENAME = test_vms_rsa

-include service-principal.mk
include templates.mk

#Tests makefiles
include test-default-config.mk
include test-config-with-variables.mk

export

prepare-service-principal: guard-CLIENT_ID guard-CLIENT_SECRET guard-SUBSCRIPTION_ID guard-TENANT_ID
	@echo "$$SERVICE_PRINCIPAL_CONTENT" > $(ROOT_DIR)/service-principal.mk

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

setup: $(ROOT_DIR)/shared/$(VMS_RSA_FILENAME)

$(ROOT_DIR)/shared/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null

clean:
	@rm -rf $(ROOT_DIR)/shared
