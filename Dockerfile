FROM golang:latest
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
    cd /go/src/github.com/toomore/gogrs && \
    go get -v ./... && \
    go get github.com/golang/lint/golint && \
    go get github.com/mattn/goveralls && \
    go get golang.org/x/tools/cmd/goimports
