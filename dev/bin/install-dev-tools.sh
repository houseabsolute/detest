#!/bin/bash

set -eo pipefail

function run () {
    echo $1
    eval $1
}

function set_bindir () {
    BINDIR="$GOPATH/bin"
}

function install_go_tools () {
    run "./dev/bin/download-golangci-lint.sh -b $BINDIR v1.12.3"
    # Built with `godownloader --source raw --repo golang/dep --exe dep --nametpl 'dep-{{ .Os }}-{{ .Arch }}' > ./dev/bin/download-dep.sh`
    run "./dev/bin/download-dep.sh -b $BINDIR v0.5.0"
}

if [ "$1" == "-v" ]; then
    set -x
fi

set_bindir
install_go_tools

exit 0
