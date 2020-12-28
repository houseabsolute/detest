#!/bin/bash

status=0

exe=$( which golangci-lint ) 
if [ -z "$exe" ]; then
    exe="./bin/golangci-lint"
fi

golangci-lint run \
    -c ./.golangci-lint-main.yml \
    --skip-dirs pkg/detest/internal/.+

status+=$?

golangci-lint run \
    -c ./.golangci-lint-internal.yml \
    ./pkg/detest/internal/...

status+=$?

exit $status
