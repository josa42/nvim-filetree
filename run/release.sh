#!/bin/bash

if [[ "$1" == "" ]]; then
  echo "Usage: ./run/release.sh <tag>"
  exit 1
fi

set -e

tag=$1

./run/build.sh

ls -a \
  | grep -v '^bin$' \
  | grep -v '^plugin$' \
  | grep -v '^syntax$' \
  | grep -v '^LICENSE$' \
  | grep -v '^README.md$' \
  | grep -v '^\.git$' \
  | grep -v '^\.\.$' \
  | grep -v '^\.$' \
  | xargs -n1 rm -rf

git add -A
git commit -m "Release $tag"
git tag "$tag"
git reset --hard HEAD@{1}
