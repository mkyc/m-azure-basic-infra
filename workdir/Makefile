include $(M_RESOURCES)/templates.mk
include $(M_RESOURCES)/consts.mk
include $(M_RESOURCES)/defaults.mk

export

#medatada method is printing static metadata information about module
metadata: guard-M_RESOURCES
	#AzBI | metadata | should print component metadata
	@echo "$$M_METADATA_CONTENT"

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
#TODO should also validate state if strong requirements are met
init: guard-M_RESOURCES guard-M_SHARED ensure-state-file template-config-file initialize-state-file display-config-file

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan: guard-M_RESOURCES guard-M_SHARED validate-config validate-state setup template-tfvars module-plan terraform-init terraform-plan

#apply method runs module provider logic using config file
apply: guard-M_RESOURCES guard-M_SHARED module-plan terraform-apply update-state-after-apply

#audit method should call logic to check if remote components are in "known" state
#TODO implement validation if remote resources are as expected, possibly with terraform plan
audit:
	#AzBI | audit | should output current state of remote components

destroy: setup template-tfvars terraform-destroy update-state-after-destroy

plan-destroy: setup template-tfvars terraform-plan-destroy

terraform-init:
	#AzBI | terraform-init | will run terraform init
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
		terraform init \
		-no-color \
		-input=false \
		$(M_RESOURCES)/terraform

#TODO consider parsing terraform plan output
terraform-plan:
	#AzBI | terraform-plan | will run plan
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
	ARM_CLIENT_ID=$(M_ARM_CLIENT_ID) \
	ARM_CLIENT_SECRET=$(M_ARM_CLIENT_SECRET) \
	ARM_SUBSCRIPTION_ID=$(M_ARM_SUBSCRIPTION_ID) \
	ARM_TENANT_ID=$(M_ARM_TENANT_ID) \
		terraform plan \
		-no-color \
		-input=false \
		-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
		-state=$(M_SHARED)/terraform.tfstate \
		-out=$(M_SHARED)/terraform-apply.tfplan \
		$(M_RESOURCES)/terraform

terraform-plan-json:
	#AzBI | terraform-plan-json | will show plan in json
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
		terraform show \
		-no-color \
		-json \
		$(M_SHARED)/terraform-apply.tfplan

terraform-apply:
	#AzBI | terraform-apply | should apply module logic
	#AzBI | terraform-apply | will run terraform apply
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
	ARM_CLIENT_ID=$(M_ARM_CLIENT_ID) \
	ARM_CLIENT_SECRET=$(M_ARM_CLIENT_SECRET) \
	ARM_SUBSCRIPTION_ID=$(M_ARM_SUBSCRIPTION_ID) \
	ARM_TENANT_ID=$(M_ARM_TENANT_ID) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/terraform.tfstate \
		$(M_SHARED)/terraform-apply.tfplan

terraform-plan-destroy:
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
	ARM_CLIENT_ID=$(M_ARM_CLIENT_ID) \
	ARM_CLIENT_SECRET=$(M_ARM_CLIENT_SECRET) \
	ARM_SUBSCRIPTION_ID=$(M_ARM_SUBSCRIPTION_ID) \
	ARM_TENANT_ID=$(M_ARM_TENANT_ID) \
		terraform plan \
		-destroy \
		-no-color \
		-input=false \
		-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
		-state=$(M_SHARED)/terraform.tfstate \
		-out=$(M_SHARED)/terraform-destroy.tfplan \
		$(M_RESOURCES)/terraform

terraform-destroy:
	#AzBI | terraform-destroy | terraform-destroy is just about to begin ...
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
	ARM_CLIENT_ID=$(M_ARM_CLIENT_ID) \
	ARM_CLIENT_SECRET=$(M_ARM_CLIENT_SECRET) \
	ARM_SUBSCRIPTION_ID=$(M_ARM_SUBSCRIPTION_ID) \
	ARM_TENANT_ID=$(M_ARM_TENANT_ID) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/terraform.tfstate \
		$(M_SHARED)/terraform-destroy.tfplan

setup:
	#AzBI | setup | ensure crucial directories are there
	@mkdir -p $(M_SHARED)/.terraform

template-tfvars:
	#AzBI | template-tfvars | will template .tfvars.json file
	@yq read -jP $(M_SHARED)/$(M_CONFIG_NAME) '$(M_MODULE_SHORT)*' > $(M_RESOURCES)/terraform/vars.tfvars.json

template-config-file:
	#AzBI | template-config-file | will template config file (and backup previous if exists)
	@if test -f $(M_SHARED)/$(M_CONFIG_NAME); then mv $(M_SHARED)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_CONFIG_NAME).backup ; fi
	@echo "$$M_CONFIG_CONTENT" > $(M_SHARED)/$(M_CONFIG_NAME)

display-config-file:
	#AzBI | display-config-file | config file content is:
	@cat $(M_SHARED)/$(M_CONFIG_NAME)

ensure-state-file: $(M_SHARED)/$(M_STATE_FILE_NAME)
	#AzBI | ensure-state-file | checks if state file exists

initialize-state-file:
	#AzBI | initialize-state-file | will initialize state file
	@echo "$$M_STATE_INITIAL" > $(M_SHARED)/azbi-state.tmp
	@yq m -i -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/azbi-state.tmp
	@rm $(M_SHARED)/azbi-state.tmp

update-state-after-apply:
	#AzBI | update-state-after-apply | will update state file after apply
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_CONFIG_NAME)
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status applied

update-state-after-destroy:
	#AzBI | update-state-after-destroy | will clean state file after destroy
	@yq d -i $(M_SHARED)/$(M_STATE_FILE_NAME) '$(M_MODULE_SHORT)'
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status destroyed

#TODO check if there is state file
#TODO check if there is config
assert-init-completed:
	#AzBI | assert-init-completed | will check if all initialization steps are completed

#TODO validate if config is correct
#TODO consider https://github.com/santhosh-tekuri/jsonschema as it's small
validate-config:
	#AzBI | validate-config | will perform config validation

#TODO validate if state file is correct
#TODO consider https://github.com/santhosh-tekuri/jsonschema as it's small
validate-state:
	#AzBI | validate-state | will perform state file validation

module-plan:
	#AzBI | module-plan | will perform module plan
	@yq m -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_CONFIG_NAME) > $(M_SHARED)/azbi-future-state.tmp
	@- yq compare $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/azbi-future-state.tmp
	@rm $(M_SHARED)/azbi-future-state.tmp

guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

$(M_SHARED)/$(M_STATE_FILE_NAME):
	@mkdir -p $(dir $@)
	@touch $@
