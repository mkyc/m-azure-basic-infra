ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

IMAGE_NAME := epiphanyplatform/azbi:0.0.1

M_NAME ?= azbi-module-tests

VMS_RSA_FILENAME = vms_rsa

-include $(ROOT_DIR)/service-principal.mk
include $(ROOT_DIR)/templates.mk

export

prepare-service-principal: guard-CLIENT_ID guard-CLIENT_SECRET guard-SUBSCRIPTION_ID guard-TENANT_ID
	@echo "$$SERVICE_PRINCIPAL_CONTENT" > $(ROOT_DIR)/service-principal.mk

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
