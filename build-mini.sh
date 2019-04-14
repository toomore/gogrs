#!/bin/bash

docker run -it --rm -v $(pwd)/gogrs_bin:/gogrs_bin toomore/gogrs:latest \
    sh -c "GO111MODULE=on go get -v ./...;
           cp /go/bin/gogrs /gogrs_bin;"

docker build -t toomore/gogrs:mini -f ./Dockerfile-mini ./

sudo rm -rf ./gogrs_bin
