#!/usr/bin/env bash

usage() {
  echo "usage:
    $0 clean
    $0 init-default-config [image_name]
    $0 check-default-config-content [image_name]
    $0 init-2-machines-no-public-ips-named [image_name]
    $0 check-2-machines-no-public-ips-named-rsa-config-content [image_name]
    $0 plan-2-machines-no-public-ips-named [image_name] ARM_CLIENT_ID ARM_CLIENT_SECRET ARM_SUBSCRIPTION_ID ARM_TENANT_ID
    $0 check-2-machines-no-public-ips-named-rsa-plan [image_name]
    $0 apply-2-machines-no-public-ips-named [image_name] ARM_CLIENT_ID ARM_CLIENT_SECRET ARM_SUBSCRIPTION_ID ARM_TENANT_ID
    $0 check-2-machines-no-public-ips-named-rsa-apply [image_name]
    $0 validate-azure-resources-presence [image_name] ARM_CLIENT_ID ARM_CLIENT_SECRET ARM_SUBSCRIPTION_ID ARM_TENANT_ID
    $0 cleanup-after-apply [image_name] ARM_CLIENT_ID ARM_CLIENT_SECRET ARM_SUBSCRIPTION_ID ARM_TENANT_ID
    "
}

clean() {
  rm -rf "$TESTS_DIR"/shared
}

setup() {
  mkdir -p "$TESTS_DIR"/shared/
  if [[ ! -f "$TESTS_DIR"/shared/test_vms_rsa ]]
  then
    ssh-keygen -t rsa -b 4096 -f "$TESTS_DIR"/shared/test_vms_rsa -N '' >/dev/null 2>&1
  fi
}

init-default-config() {
  echo "#	will initialize config with \"docker run ... init\" command"
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$IMAGE_NAME" \
    init
}

check-default-config-content() {
  echo "#	will test if file ./shared/azbi/azbi-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/azbi-config.yml; then exit 1; fi
  echo "# will test if file ./shared/azbi/azbi-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azbi/azbi-config.yml "$TESTS_DIR"/mocks/default-config/config.yml
}

init-2-machines-no-public-ips-named() {
  echo "#	will initialize config with \"docker run ... init M_VMS_COUNT=2 M_PUBLIC_IPS=false M_NAME=azbi-module-tests M_VMS_RSA=test_vms_rsa command\""
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$IMAGE_NAME" \
    init \
    M_VMS_COUNT=2 \
    M_PUBLIC_IPS=false \
    M_NAME=azbi-module-tests \
    M_VMS_RSA=test_vms_rsa
}

check-2-machines-no-public-ips-named-rsa-config-content() {
  echo "#	will test if file ./shared/azbi/azbi-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/azbi-config.yml; then exit 1; fi
  echo "#	will test if file ./shared/azbi/azbi-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azbi/azbi-config.yml "$TESTS_DIR"/mocks/config-with-variables/config.yml
}

plan-2-machines-no-public-ips-named() {
  echo "#	will plan with \"docker run ... plan M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$IMAGE_NAME" \
    plan \
    M_ARM_CLIENT_ID="$1" \
    M_ARM_CLIENT_SECRET="$2" \
    M_ARM_SUBSCRIPTION_ID="$3" \
    M_ARM_TENANT_ID="$4"
}

check-2-machines-no-public-ips-named-rsa-plan() {
  echo "#	will test if file ./shared/state.yml exists"
  if ! test -f "$TESTS_DIR"/shared/state.yml; then exit 1; fi
  echo "#	will test if file ./shared/state.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/state.yml "$TESTS_DIR"/mocks/plan/state.yml
  echo "#	will test if file ./shared/azbi/terraform-apply.tfplan exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/terraform-apply.tfplan; then exit 1; fi
  echo "#	will test if file ./shared/azbi/terraform-apply.tfplan size is greater than 0"
  filesize=$(du "$TESTS_DIR"/shared/azbi/terraform-apply.tfplan | cut -f1)
  if [[ ! $filesize -gt 0 ]]; then exit 1; fi
}

apply-2-machines-no-public-ips-named() {
  echo "#	will apply with \"docker run ... apply M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$IMAGE_NAME" \
    apply \
    M_ARM_CLIENT_ID="$1" \
    M_ARM_CLIENT_SECRET="$2" \
    M_ARM_SUBSCRIPTION_ID="$3" \
    M_ARM_TENANT_ID="$4"
}

check-2-machines-no-public-ips-named-rsa-apply() {
  echo "#	will test if file ./shared/state.yml exists"
  if ! test -f "$TESTS_DIR"/shared/state.yml; then exit 1; fi
  echo "#	will test if file ./shared/state.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/state.yml "$TESTS_DIR"/mocks/apply/state.yml
  echo "#	will test if file ./shared/azbi/terraform.tfstate exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/terraform.tfstate; then exit 1; fi
  echo "#	will test if file ./shared/azbi/terraform.tfstate size is greater than 0"
  filesize=$(du "$TESTS_DIR"/shared/azbi/terraform.tfstate | cut -f1)
  if [[ ! $filesize -gt 0 ]]; then exit 1; fi
}

validate-azure-resources-presence() {
  echo "#	will do az login"
  az login --service-principal --username "$1" --password "$2" --tenant "$4" -o none
  echo "#	will test if there is expected resource group in subscription"
  group_id=$(az group show --subscription "$3" --name azbi-module-tests-rg --query id)
  if [[ -z $group_id ]]; then exit 1; fi
  echo "#	will test if there is expected amount of machines in resource group"
  vms_count=$(az vm list --subscription "$3" --resource-group azbi-module-tests-rg -o yaml | yq r - --length)
  if [[ $vms_count -ne 2 ]]; then exit 1; fi
}

cleanup-after-apply() {
  echo "#	will apply with \"docker run ... plan-destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$IMAGE_NAME" \
    plan-destroy \
    M_ARM_CLIENT_ID="$1" \
    M_ARM_CLIENT_SECRET="$2" \
    M_ARM_SUBSCRIPTION_ID="$3" \
    M_ARM_TENANT_ID="$4"
  echo "#	will apply with \"docker run ... destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$TESTS_DIR"/shared:/shared \
    -t "$IMAGE_NAME" \
    destroy \
    M_ARM_CLIENT_ID="$1" \
    M_ARM_CLIENT_SECRET="$2" \
    M_ARM_SUBSCRIPTION_ID="$3" \
    M_ARM_TENANT_ID="$4"
}

TESTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
IMAGE_NAME=$2

case $1 in
clean)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  clean
  ;;
setup)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  setup
  ;;
init-default-config)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  init-default-config
  ;;
check-default-config-content)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  check-default-config-content
  ;;
init-2-machines-no-public-ips-named)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  init-2-machines-no-public-ips-named
  ;;
check-2-machines-no-public-ips-named-rsa-config-content)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  check-2-machines-no-public-ips-named-rsa-config-content
  ;;
plan-2-machines-no-public-ips-named)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  plan-2-machines-no-public-ips-named "$3" "$4" "$5" "$6"
  ;;
check-2-machines-no-public-ips-named-rsa-plan)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  check-2-machines-no-public-ips-named-rsa-plan
  ;;
apply-2-machines-no-public-ips-named)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  apply-2-machines-no-public-ips-named "$3" "$4" "$5" "$6"
  ;;
check-2-machines-no-public-ips-named-rsa-apply)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  check-2-machines-no-public-ips-named-rsa-apply
  ;;
validate-azure-resources-presence)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  validate-azure-resources-presence "$3" "$4" "$5" "$6"
  ;;
cleanup-after-apply)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  cleanup-after-apply "$3" "$4" "$5" "$6"
  ;;
*)
  usage
  exit 1
  ;;
esac
