FROM golang:alpine
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /go/src/github.com/toomore/gogrs/

ADD ./LICENSE ./
ADD ./README.md ./
ADD ./cmd/cache.go ./cmd/cache.go
ADD ./cmd/example.go ./cmd/example.go
ADD ./cmd/filter ./cmd/filter
ADD ./cmd/gendoc.go ./cmd/gendoc.go
ADD ./cmd/realtime.go ./cmd/realtime.go
ADD ./cmd/report.go ./cmd/report.go
ADD ./cmd/root.go ./cmd/root.go
ADD ./cmd/server.go ./cmd/server.go
ADD ./doc.go ./
ADD ./goclean.sh ./
ADD ./main.go ./main.go
ADD ./realtime ./realtime
ADD ./tradingdays ./tradingdays
ADD ./twse ./twse
ADD ./utils ./utils

VOLUME ["/go/bin"]

RUN  \
    apk update && apk add gcc git musl-dev && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/*
