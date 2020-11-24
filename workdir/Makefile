include $(M_RESOURCES)/consts.mk
include $(M_RESOURCES)/defaults.mk

export

#medatada method is printing static metadata information about module
metadata: guard-M_RESOURCES
	@./runner metadata

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
#TODO should also validate state if strong requirements are met
init: guard-M_RESOURCES guard-M_SHARED
	@./runner init --shared=$(M_SHARED) --vms_count=$(M_VMS_COUNT) --public_ips=$(M_PUBLIC_IPS) --name=$(M_NAME) --vms_rsa=$(M_VMS_RSA)

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan:
	@./runner plan --shared=$(M_SHARED) --resources=$(M_RESOURCES) --client_id=$(M_ARM_CLIENT_ID)  --client_secret=$(M_ARM_CLIENT_SECRET) --subscription_id=$(M_ARM_SUBSCRIPTION_ID) --tenant_id=$(M_ARM_TENANT_ID)

new-apply:
	@./runner apply --shared=$(M_SHARED) --resources=$(M_RESOURCES) --client_id=$(M_ARM_CLIENT_ID)  --client_secret=$(M_ARM_CLIENT_SECRET) --subscription_id=$(M_ARM_SUBSCRIPTION_ID) --tenant_id=$(M_ARM_TENANT_ID)

#apply method runs module provider logic using config file
apply: guard-M_RESOURCES guard-M_SHARED new-apply update-state-after-apply terraform-output

#audit method should call logic to check if remote components are in "known" state
#TODO implement validation if remote resources are as expected, possibly with terraform plan
audit:
	#AzBI | audit | should output current state of remote components

destroy: template-tfvars terraform-destroy update-state-after-destroy

plan-destroy: template-tfvars terraform-plan-destroy

output: terraform-output

terraform-plan-json:
	#AzBI | terraform-plan-json | will show plan in json
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
		terraform show \
		-no-color \
		-json \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan

terraform-plan-destroy:
	#AzBI | terraform-plan-destroy | will prepare plan of destruction
	@cd $(M_RESOURCES)/terraform ; \
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
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		-out=$(M_SHARED)/$(M_MODULE_SHORT)/terraform-destroy.tfplan \
		$(M_RESOURCES)/terraform

terraform-destroy:
	#AzBI | terraform-destroy | will destroy using plan of destruction
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	TF_WARN_OUTPUT_ERRORS=1 \
	ARM_CLIENT_ID=$(M_ARM_CLIENT_ID) \
	ARM_CLIENT_SECRET=$(M_ARM_CLIENT_SECRET) \
	ARM_SUBSCRIPTION_ID=$(M_ARM_SUBSCRIPTION_ID) \
	ARM_TENANT_ID=$(M_ARM_TENANT_ID) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-destroy.tfplan

terraform-output:
	#AzBI | terraform-output | will prepare terraform output
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
		terraform output \
		-no-color \
		-json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate > $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json
	@yq r --prettyPrint $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json | yq r - --printMode pv '*.value' > $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml
	@yq p -i $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml $(M_MODULE_SHORT).output
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml

template-tfvars:
	#AzBI | template-tfvars | will template .tfvars.json file
	@yq read -jP $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) '$(M_MODULE_SHORT)*' > $(M_RESOURCES)/terraform/vars.tfvars.json

update-state-after-apply:
	#AzBI | update-state-after-apply | will update state file after apply
	@cp $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml
	@yq d -i $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml kind
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status applied
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/azbi-config.tmp.yml

update-state-after-destroy:
	#AzBI | update-state-after-destroy | will clean state file after destroy
	@yq d -i $(M_SHARED)/$(M_STATE_FILE_NAME) '$(M_MODULE_SHORT)'
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status destroyed

guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

$(M_SHARED)/$(M_MODULE_SHORT):
	@mkdir -p $@
