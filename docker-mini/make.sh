#!/bin/bash

BASE=$(pwd)

cd $GOPATH/src/github.com/toomore/gogrs
go get -v ./...

cd $GOPATH/bin
cp ./twsereport $BASE
cp ./realtime $BASE
cp ./twsecache $BASE
cp ./tradingdays_server $BASE

cd $BASE
docker build -t toomore/gogrs-mini .
