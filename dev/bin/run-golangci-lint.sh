#!/bin/bash

status=0

golangci-lint run \
    -c ./.golangci-lint-main.yml \
    --skip-dirs internal/.+

status+=$?

golangci-lint run \
    -c ./.golangci-lint-internal.yml \
    ./internal/...

status+=$?

exit $status
