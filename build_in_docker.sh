#!/bin/bash
#
# When run from inside a golang project directory, this script will build that
# project, mounted inside a golang docker container.

i=${1:-"quay.io/geonet/golang-godep"}
p=${PWD/${GOPATH}\/src\//}
w="/go/src/${p}"

docker run --rm -e "GOBIN=${w}/build/bin" -e "CGO_ENABLED=0" -e "GOOS=linux" -v $PWD:${w} -w ${w}/ -it $i godep go install -a -installsuffix cgo ./...

docker build -t pagerduty-jobs:base - < Dockerfile.base

# Docker images for apps
cd build/bin
for b in pd-*
do
	rm ../Dockerfile
	echo "FROM pagerduty-jobs:base" > ../Dockerfile
	echo "ADD bin/${b} /${b}" >> ../Dockerfile
	echo "USER nobody" >> ../Dockerfile
	echo "ENTRYPOINT [\"/${b}\"]" >> ../Dockerfile
	docker build -t quay.io/geonet/pagerduty-jobs:$b ..
done


# vim: set ts=4 sw=4 tw=0 :
