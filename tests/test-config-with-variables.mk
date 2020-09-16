#Default config tests

init-2-machines-no-pubips-named: setup
	#	will initialize config with "docker run ... init M_VMS_COUNT=2 M_PUBLIC_IPS=false M_NAME=azbi-module-tests M_VMS_RSA=test_vms_rsa command
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(IMAGE_NAME) \
		init \
		M_VMS_COUNT=2 \
		M_PUBLIC_IPS=false \
		M_NAME=azbi-module-tests \
		M_VMS_RSA=test_vms_rsa

check-2-machines-no-pubips-named-rsa-config-content:
	#	will test if file ./shared/azbi/azbi-config.yml exists
	@if ! test -f $(ROOT_DIR)/shared/azbi/azbi-config.yml; then exit 1 ; fi
	#	will test if file ./shared/azbi/azbi-config.yml has expected content
	@cmp -b $(ROOT_DIR)/shared/azbi/azbi-config.yml $(ROOT_DIR)/mocks/config-with-variables-test.yml
