ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

VERSION:=0.0.1
USER:=epiphany
IMAGE:=azbi

clean:
	@rm -rf $(ROOT_DIR)/shared

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		init \
		M_VMS_COUNT=4 \
		M_PUBLIC_IPS=false

setup:
	@mkdir -p $(ROOT_DIR)/shared