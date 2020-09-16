#Default config tests

init-default-config: setup
	#	will initialize config with "docker run ... init" command
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(IMAGE_NAME) \
		init

check-default-config-content:
	#	will test if file ./shared/azbi/azbi-config.yml exists
	@if ! test -f $(ROOT_DIR)/shared/azbi/azbi-config.yml; then exit 1 ; fi
	#	will test if file ./shared/azbi/azbi-config.yml has expected content
	@cmp -b $(ROOT_DIR)/shared/azbi/azbi-config.yml $(ROOT_DIR)/mocks/default-config/config.yml
