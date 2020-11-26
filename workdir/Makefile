include $(M_RESOURCES)/consts.mk
include $(M_RESOURCES)/defaults.mk

export

#medatada method is printing static metadata information about module
metadata:
	@./runner metadata

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
#TODO should also validate state if strong requirements are met
init:
	@./runner init --shared=$(M_SHARED) --vms_count=$(M_VMS_COUNT) --public_ips=$(M_PUBLIC_IPS) --name=$(M_NAME) --vms_rsa=$(M_VMS_RSA)

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan:
	@./runner plan --shared=$(M_SHARED) --resources=$(M_RESOURCES) --client_id=$(M_ARM_CLIENT_ID)  --client_secret=$(M_ARM_CLIENT_SECRET) --subscription_id=$(M_ARM_SUBSCRIPTION_ID) --tenant_id=$(M_ARM_TENANT_ID)

#apply method runs module provider logic using config file
apply:
	@./runner apply --shared=$(M_SHARED) --resources=$(M_RESOURCES) --client_id=$(M_ARM_CLIENT_ID)  --client_secret=$(M_ARM_CLIENT_SECRET) --subscription_id=$(M_ARM_SUBSCRIPTION_ID) --tenant_id=$(M_ARM_TENANT_ID)

#audit method should call logic to check if remote components are in "known" state
#TODO implement validation if remote resources are as expected, possibly with terraform plan
audit:
	@./runner audit

destroy:
	@./runner destroy --shared=$(M_SHARED) --resources=$(M_RESOURCES) --client_id=$(M_ARM_CLIENT_ID)  --client_secret=$(M_ARM_CLIENT_SECRET) --subscription_id=$(M_ARM_SUBSCRIPTION_ID) --tenant_id=$(M_ARM_TENANT_ID)

plan-destroy:
	@./runner plan --destroy --shared=$(M_SHARED) --resources=$(M_RESOURCES) --client_id=$(M_ARM_CLIENT_ID)  --client_secret=$(M_ARM_CLIENT_SECRET) --subscription_id=$(M_ARM_SUBSCRIPTION_ID) --tenant_id=$(M_ARM_TENANT_ID)

output:
	@./runner output --shared=$(M_SHARED) --resources=$(M_RESOURCES)
