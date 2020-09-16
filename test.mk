test-default-config:
	#will run default config tests
	@cd tests ; $(MAKE) init-default-config
	@cd tests ; $(MAKE) check-default-config-content
	@cd tests ; $(MAKE) clean
	#finished default config tests
