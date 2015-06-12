FROM debian:jessie
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ADD ./twsereport /usr/local/bin/
ADD ./twsecache /usr/local/bin/
ADD ./realtime /usr/local/bin/
ADD ./tradingdays_server /usr/local/bin/
