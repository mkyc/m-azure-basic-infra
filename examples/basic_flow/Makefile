ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION := 0.0.1
USER := epiphany
IMAGE := azbi

.PHONY: all clean init apply setup

all: init apply

clean:
	@rm -rf $(ROOT_DIR)/shared

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		init \
		M_VMS_COUNT=4 \
		M_PUBLIC_IPS=false

apply: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(USER)/$(IMAGE):$(VERSION) \
		apply

setup:
	@mkdir -p $(ROOT_DIR)/shared
