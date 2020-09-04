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

#medatada method is printing static metadata information about module
metadata: guard-M_RESOURCES
	@echo \# AzBI \| metadata \| should print component metadata
	@echo "$$M_METADATA"

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
init: guard-M_RESOURCES guard-M_SHARED guard-M_CONFIG
	@echo \# AzBI \| init \| should create default config in shared directory
	@echo \# AzBI \| init \| should also validate state if strong requirements are met
	@cp ${M_RESOURCES}/default-config.yml ${M_SHARED}/${M_CONFIG}
	#TODO read additional parameters
	#TODO create/update state

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan: template
	@echo \# AzBI \| plan \| should validate if config file is correct
	@echo \# AzBI \| plan \| should validate how config is different from state
	#TODO implement config validation
	#TODO implement state comparison
	#TODO consider parsing terraform plan output

#apply method runs module provider logic using config file
apply: plan
	@echo \# AzBI \| apply \| should apply module logic
	#TODO run terraform apply

#audit method should call logic to check if remote components are in "known" state
audit:
	@echo \# AzBI \| audit \| should output current state of remote components
	#TODO implement validation if remote resources are as expected, possibly with terraform plan

#template method should template terraform files using embedded templates and configuration
template:
	@echo \# AzBI \| template \| should template terraform files using config
	#TODO template terraform files

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
