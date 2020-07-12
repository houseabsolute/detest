#!/bin/bash

status=0

golangci-lint run \
    -c ./.golangci-lint-main.yml \
    --skip-dirs pkg/detest/internal/.+

status+=$?

golangci-lint run \
    -c ./.golangci-lint-internal.yml \
    ./pkg/detest/internal/...

status+=$?

exit $status
