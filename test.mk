define t_test
	echo === RUN $1 >> $(ROOT_DIR)/output.txt
	-cd tests ; $(MAKE) $1 > $(ROOT_DIR)/output-$1.txt 2>&1 ;								\
	EXIT_CODE=$$? ;																			\
	if [ $$EXIT_CODE == "0" ] ; then 														\
		echo "--- PASS: $1 (1.00 seconds)" >> $(ROOT_DIR)/output.txt ;						\
	else 																					\
		echo "--- FAIL: $1 (1.00 seconds)" >> $(ROOT_DIR)/output.txt ;						\
		cat $(ROOT_DIR)/output-$1.txt | awk '{print "\t"$$0}' >> $(ROOT_DIR)/output.txt ;	\
		rm $(ROOT_DIR)/output-$1.txt ; 														\
	fi ;																					\
	exit $$EXIT_CODE
	echo END TEST >> $(ROOT_DIR)/output.txt
endef

test-default-config:
	#will run default config tests
	@bash tests/tests.sh clean
	@bash tests/tests.sh setup
	@bash tests/tests.sh init-default-config $(IMAGE_NAME)
	@bash tests/tests.sh check-default-config-content $(IMAGE_NAME)
	@bash tests/tests.sh clean
	#finished default config tests

test-config-with-variables:
	#will run config with variables tests
	@bash tests/tests.sh clean
	@bash tests/tests.sh setup
	@bash tests/tests.sh init-2-machines-no-public-ips-named $(IMAGE_NAME)
	@bash tests/tests.sh check-2-machines-no-public-ips-named-rsa-config-content $(IMAGE_NAME)
	@bash tests/tests.sh clean
	#finished config with variables tests

test-plan:
	#will run plan tests
	@bash tests/tests.sh clean
	@bash tests/tests.sh setup
	@bash tests/tests.sh init-2-machines-no-public-ips-named $(IMAGE_NAME)
	@bash tests/tests.sh check-2-machines-no-public-ips-named-rsa-config-content $(IMAGE_NAME)
	@bash tests/tests.sh plan-2-machines-no-public-ips-named $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@bash tests/tests.sh check-2-machines-no-public-ips-named-rsa-plan $(IMAGE_NAME)
	@bash tests/tests.sh clean
	#finished plan tests

test-apply:
	#will run apply tests
	@bash tests/tests.sh clean
	@bash tests/tests.sh setup
	@bash tests/tests.sh init-2-machines-no-public-ips-named $(IMAGE_NAME)
	@bash tests/tests.sh check-2-machines-no-public-ips-named-rsa-config-content $(IMAGE_NAME)
	@bash tests/tests.sh plan-2-machines-no-public-ips-named $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@bash tests/tests.sh check-2-machines-no-public-ips-named-rsa-plan $(IMAGE_NAME)
	@-bash tests/tests.sh apply-2-machines-no-public-ips-named $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@-bash tests/tests.sh check-2-machines-no-public-ips-named-rsa-apply $(IMAGE_NAME)
	@-bash tests/tests.sh validate-azure-resources-presence $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@bash tests/tests.sh cleanup-after-apply $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@bash tests/tests.sh clean
	#finished apply tests
