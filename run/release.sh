#!/bin/bash

if [[ "$1" == "" ]]; then
  echo "Usage: ./run/release.sh <tag>"
  exit 1
fi

set -e

tag=$1

rm -f bin/*

./run/build.sh darwin arm64
./run/build.sh darwin amd64

rm -rf run pkg go.mod go.sum

tag="v$(echo $tag | sed 's/^v//')"

git add -A
git add -f bin/tree*

git commit -m "🎉 Release $tag"

git tag $tag
git push origin $tag

git reset --hard HEAD@{1}

