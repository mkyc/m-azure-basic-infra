#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

function selfcheck() {
  local required_binaries=(which docker az yq ssh)
  local failed=1  # false
  local binary
  for binary in ${required_binaries[@]}; do
    if ! which $binary >/dev/null 2>&1; then
      failed=0  # true
      echo "FATAL: $binary is missing from PATH"
    fi
  done
  if [[ $failed -eq 0 ]]; then
    exit 1
  fi
}

function usage() {
  echo "usage:
    $0 cleanup
    $0 setup
    $0 generate_junit_report
    $0 test-default-config-suite [image_name]
    $0 test-config-with-variables-suite [image_name]
    $0 test-plan-suite [image_name] [ARM_CLIENT_ID] [ARM_CLIENT_SECRET] [ARM_SUBSCRIPTION_ID] [ARM_TENANT_ID]
    $0 test-apply-suite [image_name] [ARM_CLIENT_ID] [ARM_CLIENT_SECRET] [ARM_SUBSCRIPTION_ID] [ARM_TENANT_ID]
    "
}

function test-default-config-suite() {
  #$1 is IMAGE_NAME
  start_suite test-default-config

  local r=0
  run_test init-default-config "$r" "$1"          && r=0 || r=$?
  run_test check-default-config-content "$r" "$1" && r=0 || r=$?

  stop_suite test-default-config "$r"
}

function test-config-with-variables-suite() {
  #$1 is IMAGE_NAME
  start_suite test-config-with-variables

  local r=0
  run_test init-2-machines-no-public-ips-named "$r" "$1"                     && r=0 || r=$?
  run_test check-2-machines-no-public-ips-named-rsa-config-content "$r" "$1" && r=0 || r=$?

  stop_suite test-config-with-variables "$r"
}

function test-plan-suite() {
  #$1 is IMAGE_NAME
  #$2 is ARM_CLIENT_ID
  #$3 is ARM_CLIENT_SECRET
  #$4 is ARM_SUBSCRIPTION_ID
  #$5 is ARM_TENANT_ID
  start_suite test-plan

  local r=0
  run_test init-2-machines-no-public-ips-named "$r" "$1"                     && r=0 || r=$?
  run_test check-2-machines-no-public-ips-named-rsa-config-content "$r" "$1" && r=0 || r=$?
  run_test plan-2-machines-no-public-ips-named "$r" "$1 $2 $3 $4 $5"         && r=0 || r=$?
  run_test check-2-machines-no-public-ips-named-rsa-plan "$r" "$1"           && r=0 || r=$?

  stop_suite test-plan "$r"
}

function test-apply-suite() {
  #$1 is IMAGE_NAME
  #$2 is ARM_CLIENT_ID
  #$3 is ARM_CLIENT_SECRET
  #$4 is ARM_SUBSCRIPTION_ID
  #$5 is ARM_TENANT_ID
  start_suite test-apply

  local r=0
  run_test init-2-machines-no-public-ips-named "$r" "$1"                          && r=0 || r=$?
  run_test check-2-machines-no-public-ips-named-rsa-config-content "$r" "$1"      && r=0 || r=$?
  run_test plan-2-machines-no-public-ips-named "$r" "$1 $2 $3 $4 $5"              && r=0 || r=$?
  run_test check-2-machines-no-public-ips-named-rsa-plan "$r" "$1"                && r=0 || r=$?
  run_test apply-2-machines-no-public-ips-named "$r" "$1 $2 $3 $4 $5"             && r=0 || r=$?
  run_test check-2-machines-no-public-ips-named-rsa-apply "$r" "$1"               && r=0 || r=$?
  run_test plan-2-machines-no-public-ips-enable-public-ips "$r" "$1 $2 $3 $4 $5"  && r=0 || r=$?
  run_test apply-2-machines-no-public-ips-enable-public-ips "$r" "$1 $2 $3 $4 $5" && r=0 || r=$?
  run_test validate-azure-resources-presence "$r" "$1 $2 $3 $4 $5"                && r=0 || r=0
  run_test validate-ssh-connectivity "$r"                                         && r=0 || r=0
  run_test cleanup-after-apply "$r" "$1 $2 $3 $4 $5"                              && r=0 || r=$?

  stop_suite test-apply "$r"
}

function init-default-config() {
  echo "#	will initialize config with \"docker run ... init\" command"
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    init
}

function check-default-config-content() {
  echo "#	will test if file ./shared/azbi/azbi-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/azbi-config.yml; then exit 1; fi
  echo "# will test if file ./shared/azbi/azbi-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azbi/azbi-config.yml "$TESTS_DIR"/mocks/default-config/config.yml
}

function init-2-machines-no-public-ips-named() {
  echo "#	will initialize config with \"docker run ... init M_VMS_COUNT=2 M_PUBLIC_IPS=false M_NAME=azbi-module-tests M_VMS_RSA=test_vms_rsa command\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    init \
    M_VMS_COUNT=2 \
    M_PUBLIC_IPS=false \
    M_NAME=azbi-module-tests \
    M_VMS_RSA=test_vms_rsa
}

function check-2-machines-no-public-ips-named-rsa-config-content() {
  echo "#	will test if file ./shared/azbi/azbi-config.yml exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/azbi-config.yml; then exit 1; fi
  echo "#	will test if file ./shared/azbi/azbi-config.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/azbi/azbi-config.yml "$TESTS_DIR"/mocks/config-with-variables/config.yml
}

function plan-2-machines-no-public-ips-named() {
  echo "#	will plan with \"docker run ... plan M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    plan \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

function check-2-machines-no-public-ips-named-rsa-plan() {
  echo "#	will test if file ./shared/state.yml exists"
  if ! test -f "$TESTS_DIR"/shared/state.yml; then exit 1; fi
  echo "#	will test if file ./shared/state.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/state.yml "$TESTS_DIR"/mocks/plan/state.yml
  echo "#	will test if file ./shared/azbi/terraform-apply.tfplan exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/terraform-apply.tfplan; then exit 1; fi
  echo "#	will test if file ./shared/azbi/terraform-apply.tfplan size is greater than 0"
  local filesize=$(du "$TESTS_DIR"/shared/azbi/terraform-apply.tfplan | cut -f1)
  if [[ ! $filesize -gt 0 ]]; then exit 1; fi
}

function apply-2-machines-no-public-ips-named() {
  echo "#	will apply with \"docker run ... apply M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    apply \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

function check-2-machines-no-public-ips-named-rsa-apply() {
  echo "#	will test if file ./shared/state.yml exists"
  if ! test -f "$TESTS_DIR"/shared/state.yml; then exit 1; fi
  echo "#	will test if file ./shared/state.yml has expected content"
  cmp -b "$TESTS_DIR"/shared/state.yml "$TESTS_DIR"/mocks/apply/state.yml
  echo "#	will test if file ./shared/azbi/terraform.tfstate exists"
  if ! test -f "$TESTS_DIR"/shared/azbi/terraform.tfstate; then exit 1; fi
  echo "#	will test if file ./shared/azbi/terraform.tfstate size is greater than 0"
  local filesize=$(du "$TESTS_DIR"/shared/azbi/terraform.tfstate | cut -f1)
  if [[ ! $filesize -gt 0 ]]; then exit 1; fi
}

function plan-2-machines-no-public-ips-enable-public-ips() {
  echo "#   will enable public ips inside azbi-config.yml"
  yq w --inplace "$TESTS_DIR"/shared/azbi/azbi-config.yml azbi.use_public_ip true
  echo "#	will plan with \"docker run ... plan M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    plan \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

function apply-2-machines-no-public-ips-enable-public-ips() {
  echo "#	will apply with \"docker run ... apply M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    apply \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

function validate-azure-resources-presence() {
  echo "#	will do az login"
  az login --service-principal --username "$2" --password "$3" --tenant "$5" -o none
  echo "#	will test if there is expected resource group in subscription"
  local group_id=$(az group show --subscription "$4" --name azbi-module-tests-rg --query id)
  if [[ -z $group_id ]]; then exit 1; fi
  echo "#	will test if there is expected amount of machines in resource group"
  local vms_count=$(az vm list --subscription "$4" --resource-group azbi-module-tests-rg -o yaml | yq r - --length)
  if [[ $vms_count -ne 2 ]]; then exit 1; fi
  local nsg_nic_count=$(az network nsg show --subscription "$4" --resource-group azbi-module-tests-rg --name vm-nic-nsg-0 -o yaml | yq r - networkInterfaces --length)
  if [[ $nsg_nic_count -ne 2 ]]; then exit 1; fi
}

function validate-ssh-connectivity() {
  local ssh_options='-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectionAttempts=10 -o ControlMaster=no -o BatchMode=yes'
  local ip_address
  unset IFS && yq r "$TESTS_DIR"/shared/state.yml 'azbi.output[public_ips.value].*' | while read ip_address; do
    if ! ssh -q -F /dev/null $ssh_options -i "$TESTS_DIR"/shared/test_vms_rsa -nT "operations@$ip_address" -- uname -a; then exit 1; fi
  done
}

function cleanup-after-apply() {
  echo "#	will apply with \"docker run ... plan-destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    plan-destroy \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
  echo "#	will apply with \"docker run ... destroy M_ARM_CLIENT_ID=... M_ARM_CLIENT_SECRET=... M_ARM_SUBSCRIPTION_ID=... M_ARM_TENANT_ID=...\""
  docker run --rm \
    -v "$MOUNT_DIR"/shared:/shared \
    -t "$1" \
    destroy \
    M_ARM_CLIENT_ID="$2" \
    M_ARM_CLIENT_SECRET="$3" \
    M_ARM_SUBSCRIPTION_ID="$4" \
    M_ARM_TENANT_ID="$5"
}

selfcheck

TESTS_DIR=/tests-share
MOUNT_DIR=/tmp/tests-share

# shellcheck disable=SC1090
source "$(dirname "$0")/suite.sh"

case $1 in
test-default-config-suite)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  test-default-config-suite "$2"
  ;;
test-config-with-variables-suite)
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi
  test-config-with-variables-suite "$2"
  ;;
test-plan-suite)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  test-plan-suite "$2" "$3" "$4" "$5" "$6"
  ;;
test-apply-suite)
  if [[ $# -ne 6 ]]; then
    usage
    exit 1
  fi
  test-apply-suite "$2" "$3" "$4" "$5" "$6"
  ;;
cleanup)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  cleanup
  ;;
setup)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  setup
  ;;
generate_junit_report)
  if [[ $# -ne 1 ]]; then
    usage
    exit 1
  fi
  generate_junit_report
  ;;
*)
  usage
  exit 1
  ;;
esac
