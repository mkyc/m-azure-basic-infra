define M_CONFIG_CONTENT
kind: $(M_MODULE_SHORT)-config
$(M_MODULE_SHORT):
  size: $(M_VMS_COUNT)
  use_public_ip: $(M_PUBLIC_IPS)
  location: "$(M_LOCATION)"
  name: "$(M_NAME)"
  address_space: ["10.0.0.0/16"]
  address_prefixes: ["10.0.1.0/24"]
  rsa_pub_path: "$(M_SHARED)/$(M_VMS_RSA).pub"
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
