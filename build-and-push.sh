#!/bin/bash

umask 0022

REPO="github.com"
ORG="GeoNet"

name=$(basename ${PWD})
p="${REPO}/${ORG}/${name}"
w="/go/src/${p}"

# Build all executables in the golang-godep container.  Output statically linked binaries to ${tmp_dir}
docker run --rm \
    -e CGO_ENABLED=0 -e GOOS=linux \
    -v ${PWD}:${w}:ro \
    -v ${PWD}/bin:/go/bin \
    -w ${w} \
    golang:1.6.0-alpine \
    go install -a -ldflags "${BUILD}" -installsuffix cgo ./...

docker build -t ${name}:base - < Dockerfile.base

# Docker images for apps
cd bin
for b in pd-*
do
    tmp_dir=$(mktemp -d)
    chmod 1777 ${tmp_dir}

    cp ${b} ${tmp_dir}

    cat > ${tmp_dir}/Dockerfile << EOF
FROM ${name}:base
ADD ${b} /
USER nobody
ENTRYPOINT ["/${b}"]
EOF

    echo "'docker build --rm=true -t quay.io/geonet/${name}:${b} ${tmp_dir}'"
    docker build --rm=true -t quay.io/geonet/${name}:${b} ${tmp_dir}

    rm -rf ${tmp_dir}
done

docker push quay.io/geonet/${name}

exit 0

# vim: set ts=4 sw=4 tw=0 et:
