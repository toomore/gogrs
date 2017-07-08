FROM alpine:latest
MAINTAINER Toomore Chiang <toomore0929@gmail.com>

WORKDIR /bin

ADD ./gogrs_bin ./

RUN apk update && apk add ca-certificates bash-completion && \
    cd /usr/share/bash-completion/completions && gogrs doc -b && \
    rm -rf /var/cache/apk/* /var/lib/apk/* /etc/apk/cache/*

CMD ["/bin/bash", "-l"]
