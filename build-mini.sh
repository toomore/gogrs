#!/bin/bash

docker run -it --rm -v $(pwd)/gogrs_bin:/gogrs_bin toomore/gogrs:latest \
    sh -c "go get -v ./...;
           cp /go/bin/gogrs /gogrs_bin;"

docker build -t toomore/gogrs:mini -f ./Dockerfile-mini ./

sudo rm -rf ./gogrs_bin
