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

RUN  \
    apk update && apk add bash git build-base && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    go get github.com/golang/lint/golint && \
    go get golang.org/x/tools/cmd/goimports && \
    go get -v ./...
