#!/bin/bash

set -eo pipefail

gometalinter-helper \
    -ignore ./.gitignore \
    $@
