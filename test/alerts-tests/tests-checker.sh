#!/bin/bash

expected_tests_count=2

rules=()
readarray -t rules < <(yq '.groups[].rules[].alert' ./rules.yaml)

tests=()
readarray -t tests < <(yq '.tests[].alert_rule_test[].alertname' ./*-tests.yaml)

errorrules=()
errorcount=()
i=0

for item in "${rules[@]}"; do
    count=0

    for j in "${tests[@]}"; do
        if [[ "$j" == "$item" ]]; then
            ((count++))
        fi
    done

    echo "Passed: $item, Tests: $count / $expected_tests_count"

    if [[ "$count" -lt $expected_tests_count ]]; then
        errorrules[i]="$item"
        errorcount[i]="$count"
        ((i++))
    fi
done

if [[ "$i" -gt 0 ]]; then
    echo "--------------------------------"
    echo "Failed: Alert rules have less than ${expected_tests_count} tests (minimum ${expected_tests_count} tests per rule needed):"
    for k in "${!errorrules[@]}"; do
        echo "Failed: ${errorrules[k]}, Tests: ${errorcount[k]} / ${expected_tests_count}"
    done
    exit 1
else
    echo "--------------------------------"
    echo "Passed: All alert rules have required tests"
    exit 0
fi
