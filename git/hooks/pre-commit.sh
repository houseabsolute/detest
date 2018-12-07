#!/bin/bash

set -eo pipefail

status=0

go generate ./...
if (( $? != 0 )); then
    status+=1
fi

./dev/bin/run-gometalinter-helper.sh -commit-hook
if (( $? != 0 )); then
    status+=2
fi

exit $status
