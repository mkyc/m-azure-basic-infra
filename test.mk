test-default-config:
	#will run default config tests
	@bash tests/tests.sh cleanup
	@bash tests/tests.sh setup
	@bash tests/tests.sh test-default-config-suite $(IMAGE_NAME)
	@bash tests/tests.sh cleanup
	#finished default config tests

test-config-with-variables:
	#will run config with variables tests
	@bash tests/tests.sh cleanup
	@bash tests/tests.sh setup
	@bash tests/tests.sh test-config-with-variables-suite $(IMAGE_NAME)
	@bash tests/tests.sh cleanup
	#finished config with variables tests

test-plan:
	#will run plan tests
	@bash tests/tests.sh cleanup
	@bash tests/tests.sh setup
	@bash tests/tests.sh test-plan-suite $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@bash tests/tests.sh cleanup
	#finished plan tests

test-apply:
	#will run apply tests
	@bash tests/tests.sh cleanup
	@bash tests/tests.sh setup
	@bash tests/tests.sh test-apply-suite $(IMAGE_NAME) $(ARM_CLIENT_ID) $(ARM_CLIENT_SECRET) $(ARM_SUBSCRIPTION_ID) $(ARM_TENANT_ID)
	@bash tests/tests.sh cleanup
	#finished apply tests

generate-report:
	@bash tests/tests.sh generate_junit_report
