gogrs-mini
===========

Minify gogrs docker, only add **binary files** into docker image.

## Build toomore/gogrs:mini

    sh ./make.sh

> Note: Will create temp folder:`gogrs_bin`, and will be removed at end.

## Run default cmd

    docker run -it --rm toomore/gogrs:mini twsereport --help

## Default cmd
  - [gogrs_example](https://godoc.org/github.com/toomore/gogrs/cmd/gogrs_example)
  - [tradingdays_server](https://godoc.org/github.com/toomore/gogrs/cmd/tradingdays_server)
  - [twsecache](https://godoc.org/github.com/toomore/gogrs/cmd/twsecache)
  - [twsereport](https://godoc.org/github.com/toomore/gogrs/cmd/twsereport)
