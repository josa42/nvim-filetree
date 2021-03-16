#!/bin/bash

set -e

export GOOS="$1"
export GOARCH="$2"

eval $(go env | grep 'GOARCH\|GOOS')

uuid='2f4eab2fca750d49cc9039bac724b89b'

go build -ldflags "-X main.uuid=${uuid}" -o bin/tree-${GOOS}-${GOARCH} ./pkg
./bin/tree-${GOOS}-${GOARCH} --manifest tree -location plugin/tree.vim


