#!/bin/bash
#
# When run from inside a golang project directory, this script will build that
# project, mounted inside a golang docker container.

i=${1:-"quay.io/geonet/golang-godep"}
p=${PWD/*\/gocode\/src\//} # I have my gocode in ~/gocode/src/...
w="/go/src/${p}"

#docker run --rm -e "GOBIN=${w}/bin" -e "CGO_ENABLED=0" -e "GOOS=linux" -v "/etc/passwd:/etc/passwd:ro" -u ${USER} -v $PWD:${w} -w ${w}/ -it $i godep go install -a -installsuffix cgo ./... 

docker run --rm -e "GOBIN=${w}/bin" -e "CGO_ENABLED=0" -e "GOOS=linux" -v $PWD:${w} -w ${w}/ -it $i godep go install -a -installsuffix cgo ./... 
# vim: set ts=4 sw=4 tw=0 :
