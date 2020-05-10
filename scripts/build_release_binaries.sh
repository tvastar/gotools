#!/bin/bash

function build() {
    export GOARCH=$1
    export GOOS=$2
    export CGO_ENABLED=$3

    mkdir -p ./bin/"$2"_"$1"
    
    CGO_ENABLED=$3 GOARCH=$1 GOOS=$2 go build -o ./bin/"$2"_"$1"/. ./...

    pushd ./bin/"$2"_"$1"
    tar -czvf ../"$2"_"$1".tar.gz *
    popd
}

set -ex
build 386 linux 0
build amd64 linux 0
build amd64 darwin 1
build 386 windows 1
build amd64 windows 1

