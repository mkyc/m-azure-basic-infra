include $(M_RESOURCES)/templates.mk
include $(M_RESOURCES)/defaults.mk

export

#medatada method is printing static metadata information about module
metadata: guard-M_RESOURCES
	#AzBI | metadata | should print component metadata
	@echo "$$M_METADATA_CONTENT"

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
init: guard-M_RESOURCES guard-M_SHARED guard-M_CONFIG_NAME guard-M_VMS_COUNT guard-M_PUBLIC_IPS
	#AzBI | init | should create default config in shared directory
	#AzBI | init | should also validate state if strong requirements are met
	@if test -f $(M_SHARED)/$(M_CONFIG_NAME); then mv $(M_SHARED)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_CONFIG_NAME).backup ; fi
	@echo "$$M_CONFIG_CONTENT" > $(M_SHARED)/$(M_CONFIG_NAME)
	@#TODO create/update state

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan: template
	#AzBI | plan | should validate if config file is correct
	#AzBI | plan | should validate how config is different from state
	@#TODO implement config validation
	@#TODO implement state comparison
	@#TODO consider parsing terraform plan output
	#AzBI | plan | will run init
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
		terraform init \
		-no-color \
		-input=false \
		$(M_RESOURCES)/terraform
	#AzBI | plan | will run plan
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
		-out=$(M_SHARED)/terraform.tfplan \
		$(M_RESOURCES)/terraform
	#AzBI | plan | will show plan in json
	@TF_DATA_DIR=$(M_SHARED)/.terraform \
	TF_IN_AUTOMATION=true \
		terraform show \
		-no-color \
		-json \
		$(M_SHARED)/terraform.tfplan

#apply method runs module provider logic using config file
apply: plan guard-M_RESOURCES guard-M_SHARED
	#AzBI | apply | should apply module logic
	#AzBI | apply | will run apply
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
		$(M_SHARED)/terraform.tfplan

#audit method should call logic to check if remote components are in "known" state
audit:
	#AzBI | audit | should output current state of remote components
	@#TODO implement validation if remote resources are as expected, possibly with terraform plan

#template method should template terraform files using embedded templates and configuration
template:
	#AzBI | template | should template terraform files using config
	#AzBI | template | will template .tfvars.json file
	@mkdir -p $(M_SHARED)/.terraform
	@yq read -jP $(M_SHARED)/$(M_CONFIG_NAME) 'azi*' > $(M_RESOURCES)/terraform/vars.tfvars.json

destroy: template
	#AzBI | destroy | destroyer of the worlds is planing ...
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
	@if test -z $$M_DO_DESTROY ; \
	then \
		echo \#AzBI \| destroy \| destroyer of the worlds is sleeping ... ; \
	else \
	    echo \#AzBI \| destroy \| destroyer of the worlds is just about to begin ... ; \
		TF_DATA_DIR=$(M_SHARED)/.terraform \
		TF_IN_AUTOMATION=true \
		ARM_CLIENT_ID=$(M_ARM_CLIENT_ID) \
		ARM_CLIENT_SECRET=$(M_ARM_CLIENT_SECRET) \
		ARM_SUBSCRIPTION_ID=$(M_ARM_SUBSCRIPTION_ID) \
		ARM_TENANT_ID=$(M_ARM_TENANT_ID) \
			terraform destroy \
			-no-color \
			-input=false \
			-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
			-auto-approve \
			-state=$(M_SHARED)/terraform.tfstate \
			$(M_RESOURCES)/terraform ; \
	fi

guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
