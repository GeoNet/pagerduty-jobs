#!/bin/bash -eu

# Builds Docker images for the arg list.  These must be project directories
# where this script is executed.
#
# Builds a statically linked executable and adds it to the container.
# Adds the assets dir from each project to the container e.g., origin/assets
# It is not an error for the assets dir to not exist.
# Any assets needed by the application should be read from the assets dir
# relative to the executable.
#
# usage: ./build.sh project [project]

if [ $# -eq 0 ]; then
  echo Error: please supply a project to build. Usage: ./build.sh project [project]
  exit 1
fi

# code will be compiled in this container
BUILDER_IMAGE='quay.io/geonet/golang:1.13.1-alpine'
RUNNER_IMAGE='quay.io/geonet/alpine:3.10'

VERSION='git-'$(git rev-parse --short HEAD)

for i in "$@"; do
 
  mkdir -p cmd/$i/assets
  dockerfile="Dockerfile"
  if test -f "cmd/${i}/Dockerfile"; then
    dockerfile="cmd/${i}/Dockerfile"
  fi


  docker build \
    --build-arg=BUILD="$i" \
    --build-arg=RUNNER_IMAGE="$RUNNER_IMAGE" \
    --build-arg=BUILDER_IMAGE="$BUILDER_IMAGE" \
    --build-arg=GIT_COMMIT_SHA="$VERSION" \
    --build-arg=ASSET_DIR="./cmd/$i/assets" \
    -t "quay.io/geonet/${i}:${VERSION}" \
    -f $dockerfile .

done

# vim: set ts=4 sw=4 tw=0 et:
