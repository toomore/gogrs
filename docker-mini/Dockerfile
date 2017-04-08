FROM alpine:latest
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /bin

ADD ./gogrs_bin ./

RUN apk update && apk add ca-certificates && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
