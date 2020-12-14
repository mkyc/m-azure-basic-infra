ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

#makes it easier to replace the value
M_NAME ?= azbi-module-example

VMS_RSA_FILENAME = vms_rsa

-include "$(ROOT_DIR)/azure.mk"

export

.PHONY: all clean init plan apply output destroy-plan destroy setup

all: init plan apply

#TODO that all needs fixing

#clean:
#	@rm -rf $(M_SHARED)
#
#init: setup
#	@cd $(M_WORKDIR) && $(MAKE) init \
#		M_VMS_COUNT=2 \
#		M_PUBLIC_IPS=false \
#		M_NAME=$(M_NAME) \
#		M_VMS_RSA=$(VMS_RSA_FILENAME)
#
#plan: setup
#	@cd $(M_WORKDIR) && $(MAKE) plan \
#		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
#		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
#		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
#		M_ARM_TENANT_ID=$$ARM_TENANT_ID
#
#apply: setup
#	@cd $(M_WORKDIR) && $(MAKE) apply \
#		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
#		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
#		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
#		M_ARM_TENANT_ID=$$ARM_TENANT_ID
#
#output: setup
#	@cd $(M_WORKDIR) && $(MAKE) output \
#		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
#		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
#		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
#		M_ARM_TENANT_ID=$$ARM_TENANT_ID
#
#destroy-plan: setup
#	@cd $(M_WORKDIR) && $(MAKE) plan-destroy \
#		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
#		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
#		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
#		M_ARM_TENANT_ID=$$ARM_TENANT_ID
#
#destroy: setup
#	@cd $(M_WORKDIR) && $(MAKE) destroy \
#		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
#		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
#		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
#		M_ARM_TENANT_ID=$$ARM_TENANT_ID
#
#setup: $(M_SHARED)/$(VMS_RSA_FILENAME)
#
#$(M_SHARED)/$(VMS_RSA_FILENAME):
#	@mkdir -p $(dir $@)
#	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
