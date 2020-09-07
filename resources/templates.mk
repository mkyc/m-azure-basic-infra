define M_METADATA
labels:
  version: $(M_VERSION)
  name: Azure Basic Infrastructure
  short: AzBI
  kind: infrastructure
  provider: azure
  provides-vms: true
  provides-pubips: true
endef
export M_METADATA
