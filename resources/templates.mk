define M_METADATA_CONTENT
labels:
  version: $(M_VERSION)
  name: Azure Basic Infrastructure
  short: AzBI
  kind: infrastructure
  provider: azure
  provides-vms: true
  provides-pubips: true
endef

define M_CONFIG_CONTENT
azi:
  size: $(M_VMS_COUNT)
  provide-public-IPs: $(M_PUBLIC_IPS)
endef

define M_TFVARS_CONTENT
location = $(M_LOCATION)
rg_name = $(M_RG_NAME)
endef
