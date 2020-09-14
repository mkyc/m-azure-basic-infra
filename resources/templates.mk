define M_METADATA_CONTENT
labels:
  version: $(M_VERSION)
  name: Azure Basic Infrastructure
  short: $(M_MODULE_SHORT)
  kind: infrastructure
  provider: azure
  provides-vms: true
  provides-pubips: true
endef

define M_CONFIG_CONTENT
kind: $(M_MODULE_SHORT)-config
$(M_MODULE_SHORT):
  size: $(M_VMS_COUNT)
  use_public_ip: $(M_PUBLIC_IPS)
  location: "$(M_LOCATION)"
  name: "$(M_NAME)"
  address_space:
  - "10.0.0.0/16"
  subnet_cidrs:
  - "10.0.1.0/24"
  subnet_names:
  - "$(M_NAME)-sn1"
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
