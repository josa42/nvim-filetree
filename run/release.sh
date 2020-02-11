#!/bin/bash

if [[ "$1" == "" ]]; then
  echo "Usage: ./run/release.sh <tag>"
  exit 1
fi

set -e

tag=$1

./run/build.sh

rm -rf run pkg

tag="v$(echo $tag | sed 's/^v//')"

git add -A
git add -f bin/tree

git commit -m "🎉 Release $tag"

git tag $tag
git push origin $tag

git reset --hard HEAD@{1}

