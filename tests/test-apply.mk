#Default config tests

apply-2-machines-no-pubips-named:
	#	will apply with "docker run ... apply M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(IMAGE_NAME) \
		apply \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

check-2-machines-no-pubips-named-rsa-apply:
	#	will test if file ./shared/state.yml exists
	@if ! test -f $(ROOT_DIR)/shared/state.yml; then exit 1 ; fi
	#	will test if file ./shared/state.yml has expected content
	@cmp -b $(ROOT_DIR)/shared/state.yml $(ROOT_DIR)/mocks/apply/state.yml
	#	will test if file ./shared/azbi/terraform.tfstate exists
	@if ! test -f $(ROOT_DIR)/shared/azbi/terraform.tfstate; then exit 1 ; fi
	#	will test if file ./shared/azbi/terraform.tfstate size is greater than 0
	@FSIZE=$$(ls -l $$ROOT_DIR/shared/azbi/terraform.tfstate | awk '{print $$5}') ;\
	if [ ! $$FSIZE -gt 0 ] ; then exit 1 ; fi

validate-azure-resources-presence:
	#	will do az login
	@az login --service-principal --username $$ARM_CLIENT_ID --password $$ARM_CLIENT_SECRET --tenant $$ARM_TENANT_ID -o none
	#	will test if there is expected resource group in subscription
	@GID=$$(az group show --subscription $$ARM_SUBSCRIPTION_ID --name azbi-module-tests-rg --query id) ;\
	if test -z $$GID ; then exit 1 ; fi
	#	will test if there is expected amount of machines in resource group
	@VMSC=$$(az vm list --subscription $$ARM_SUBSCRIPTION_ID --resource-group azbi-module-tests-rg -o yaml | yq r - --length) ;\
	if [ $$VMSC -ne 2 ] ; then exit 1 ; fi

cleanup-after-apply:
	#	will apply with "docker run ... plan-destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(IMAGE_NAME) \
		plan-destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID
	#	will apply with "docker run ... destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(IMAGE_NAME) \
		destroy \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID
