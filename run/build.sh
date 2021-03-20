#!/bin/bash

set -e

uuid='2f4eab2fca750d49cc9039bac724b89b'
go run -ldflags "-X main.uuid=${uuid}" ./pkg --manifest tree -location plugin/tree.vim

eval $(go env | grep 'GOARCH\|GOOS')
export GOOS="${1:-$GOOS}"
export GOARCH="${2:-$GOARCH}"

go build -ldflags "-X main.uuid=${uuid}" -o bin/tree-${GOOS}-${GOARCH} ./pkg

