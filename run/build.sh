#!/bin/bash

set -e

go build -o bin/tree ./pkg
./bin/tree --manifest tree -location plugin/tree.vim


