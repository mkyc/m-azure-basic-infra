test-default-config:
	#will run default config tests
	@cd tests ; $(MAKE) clean
	@cd tests ; $(MAKE) init-default-config
	@cd tests ; $(MAKE) check-default-config-content
	@cd tests ; $(MAKE) clean
	#finished default config tests

test-config-with-variables:
	#will run config with variables tests
	@cd tests ; $(MAKE) clean
	@cd tests ; $(MAKE) init-2-machines-no-pubips-named
	@cd tests ; $(MAKE) check-2-machines-no-pubips-named-rsa-config-content
	@cd tests ; $(MAKE) clean
	#finished config with variables tests

test-plan:
	#will run plan tests
	@cd tests ; $(MAKE) clean
	@cd tests ; $(MAKE) init-2-machines-no-pubips-named
	@cd tests ; $(MAKE) check-2-machines-no-pubips-named-rsa-config-content
	@cd tests ; $(MAKE) plan-2-machines-no-pubips-named
	@cd tests ; $(MAKE) check-2-machines-no-pubips-named-rsa-plan
	@cd tests ; $(MAKE) clean
	#finished plan tests