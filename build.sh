#!/usr/bin/env bash

if [[ ! -f build.sh ]]; then
echo 'build must be run within its container folder' 1>&2
exit 1
fi

dist=go_main

export GOPATH=$GOPATH:${PWD}
echo $GOPATH

go build -o ./bin/${dist} main
