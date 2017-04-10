FROM golang:alpine
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /go/src/github.com/toomore/gogrs/

ADD ./cmd ./cmd
ADD ./realtime ./realtime
ADD ./tradingdays ./tradingdays
ADD ./twse ./twse
ADD ./utils ./utils

ADD ./LICENSE ./
ADD ./README.md ./
ADD ./doc.go ./
ADD ./goclean.sh ./

VOLUME ["/go/bin"]

RUN  \
    apk update && apk add gcc git musl-dev && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/*
