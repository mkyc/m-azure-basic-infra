#!/usr/bin/env bash

function run_test() {
  # $1 is function name
  # $2 is result of previous test
  # $3 are all target function arguments
  fn_name=$1
  shift
  previous_r=$1
  shift
  if [[ "$previous_r" -eq "0" ]]; then
    start_test "$fn_name"
    # shellcheck disable=SC2068
    $fn_name $@ >>"$TESTS_DIR"/output/"$fn_name"-tests.out 2>&1
    result=$?
    stop_test "$fn_name" $result
    return $result
  else
    start_test "$fn_name"
    skip_test "$fn_name"
    return "$previous_r"
  fi
}

function start_suite() {
  #$1 is test name
  suite_start_time=$(microtime)
}

function stop_suite() {
  #$1 is test name
  #$2 is $?
  if [[ $2 == "0" ]]; then
    duration=$(echo "$(microtime) - $suite_start_time" | bc)
    execution_time=$(printf "%.2fs" "$duration")
    printf "PASS\nok\t\t%s %s\n" "$1" "$execution_time" >>"$TESTS_DIR"/tests.out
  else
    duration=$(echo "$(microtime) - $suite_start_time" | bc)
    execution_time=$(printf "%.2fs" "$duration")
    printf "FAIL\nexit status %s\nFAIL\t\t%s %s\n" "$2" "$1" "$execution_time" >>"$TESTS_DIR"/tests.out
  fi
}

function start_test() {
  #$1 is test name
  echo "=== RUN $1" >>"$TESTS_DIR"/tests.out
  test_start_time=$(microtime)
}

function stop_test() {
  #$1 is test name
  #$2 is $?
  if [[ $2 == "0" ]]; then
    pass_test "$1"
  else
    fail_test "$1"
  fi
}

function pass_test() {
  #$1 is test name
  duration=$(echo "$(microtime) - $test_start_time" | bc)
  execution_time=$(printf "%.2f seconds" "$duration")
  echo "--- PASS: $1 ($execution_time)" >>"$TESTS_DIR"/tests.out
}

function fail_test() {
  #$1 is test name
  duration=$(echo "$(microtime) - $test_start_time" | bc)
  execution_time=$(printf "%.2f seconds" "$duration")
  echo "--- FAIL: $1 ($execution_time)" >>"$TESTS_DIR"/tests.out
  awk <"$TESTS_DIR"/output/"$1"-tests.out '{print "\t\t"$0}' >>"$TESTS_DIR"/tests.out
}

function skip_test() {
  #$1 is test name
  duration=$(echo "$(microtime) - $test_start_time" | bc)
  execution_time=$(printf "%.2f seconds" "$duration")
  echo "--- SKIP: $1 ($execution_time)" >>"$TESTS_DIR"/tests.out
  printf "\tprevious test failed\n" >>"$TESTS_DIR"/tests.out
}

function cleanup() {
  rm -rf "$TESTS_DIR"/shared
  rm -rf "$TESTS_DIR"/output
}

function setup() {
  mkdir -p "$TESTS_DIR"/shared/
  mkdir -p "$TESTS_DIR"/output/
  if [[ ! -f "$TESTS_DIR"/shared/test_vms_rsa ]]; then
    ssh-keygen -t rsa -b 4096 -f "$TESTS_DIR"/shared/test_vms_rsa -N '' >/dev/null 2>&1
  fi
}

function generate_junit_report() {
  go-junit-report <"$TESTS_DIR"/tests.out >"$TESTS_DIR"/repot.xml
  rm "$TESTS_DIR"/tests.out
}

function microtime() {
  #this is due to macOS BSD'ish date command
  python -c 'import time; print time.time()'
}
