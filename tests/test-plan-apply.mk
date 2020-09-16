#Default config tests

plan-2-machines-no-pubips-named:
	#	will plan with "docker run ... plan M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(IMAGE_NAME) \
		plan \
		M_ARM_CLIENT_ID=$$ARM_CLIENT_ID \
		M_ARM_CLIENT_SECRET=$$ARM_CLIENT_SECRET \
		M_ARM_SUBSCRIPTION_ID=$$ARM_SUBSCRIPTION_ID \
		M_ARM_TENANT_ID=$$ARM_TENANT_ID

check-2-machines-no-pubips-named-rsa-plan:
	#	will test if file ./shared/state.yml exists
	@if ! test -f $(ROOT_DIR)/shared/azbi/azbi-config.yml; then exit 1 ; fi
	#	will test if file ./shared/state.yml has expected content
	@cmp -b $(ROOT_DIR)/shared/state.yml $(ROOT_DIR)/mocks/plan/state.yml
	#	will test if file ./shared/azbi/terraform-apply.tfplan exists
	@if ! test -f $(ROOT_DIR)/shared/azbi/terraform-apply.tfplan; then exit 1 ; fi
	#	will test if file ./shared/azbi/terraform-apply.tfplan size is greater than 0
	@FSIZE=$$(ls -l $$ROOT_DIR/shared/azbi/terraform-apply.tfplan | awk '{print $$5}') ;\
	if [ ! $$FSIZE -gt 0 ] ; then exit 1 ; fi
