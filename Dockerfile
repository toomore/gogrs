FROM golang:wheezy 
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

RUN go get golang.org/x/tools/cmd/cover && go get golang.org/x/tools/cmd/vet && \
    go get golang.org/x/tools/cmd/goimports && go get github.com/golang/lint/golint && \
    go get github.com/toomore/gogrs
RUN cd /go/src/github.com/toomore/gogrs && go get -v ./... && sh ./goclean.sh
