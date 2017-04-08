#!/bin/bash

docker run -it --rm -v $(pwd)/gogrs_bin:/gogrs_bin gogrs:latest \
    bash -c "cp /go/bin/gogrs_example /gogrs_bin;
             cp /go/bin/realtime /gogrs_bin;
             cp /go/bin/tradingdays_server /gogrs_bin;
             cp /go/bin/twsecache /gogrs_bin;
             cp /go/bin/twsereport /gogrs_bin;"

docker build -t gogrs:mini ./

sudo rm -rf ./gogrs_bin
