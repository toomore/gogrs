gogrs
======

[![GoDoc](https://godoc.org/github.com/toomore/gogrs?status.svg)](https://godoc.org/github.com/toomore/gogrs) [![Build Status](https://travis-ci.org/toomore/gogrs.svg?branch=master)](https://travis-ci.org/toomore/gogrs) [![Coverage Status](https://coveralls.io/repos/toomore/gogrs/badge.svg?branch=master)](https://coveralls.io/r/toomore/gogrs?branch=master) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/toomore/gogrs/master/LICENSE)

gogrs is a tool for fetching data from Taiwan Stock Exchange(TWSE) and dockerizing. gogrs now is still in development. I will try my best to speed up to completed the same function with [grs](https://github.com/toomore/grs) (Python). gogrs 是擷取台灣上市股票股價資訊工具，目前還在大量的開發中。原始工具是用 [grs](https://github.com/toomore/grs)（Python 套件），目標是將基本功能用 go 來實作。

Packages
---------

1. realtime - [擷取盤中個股、指數即時股價資訊](https://godoc.org/github.com/toomore/gogrs/realtime)
2. twse - [擷取台灣股市上市/上櫃股票資訊、上市/上櫃類股清單、外資及陸資持股比率前二十名彙總表、 三大法人買賣金額統計表、三大法人買賣超日報、自營商、投信、外資及陸資買賣超彙總表](https://godoc.org/github.com/toomore/gogrs/twse)
3. tradingdays - [股市開休市判斷（支援非國定假日：颱風假）與當日區間判斷（盤中、盤後、盤後盤）](https://godoc.org/github.com/toomore/gogrs/tradingdays)
4. utils - [套件所需的公用工具（總和、平均、序列差、持續天數、民國日期解析、簡單亂數、標準差、簡單 net/http 快取）](https://godoc.org/github.com/toomore/gogrs/utils)

Cmd
----

![gogrs cmd demo](https://s3-ap-northeast-1.amazonaws.com/toomore/gogrs/gogrs_cmd_demo_20150615.png "gogrs cmd demo")

1. gogrs_example - [簡單範例測試](https://godoc.org/github.com/toomore/gogrs/cmd/gogrs_example)
2. realtime - [擷取盤中即時資訊與大盤、上櫃、寶島指數](https://godoc.org/github.com/toomore/gogrs/cmd/realtime)
3. twsereport - [每日收盤後產生符合選股條件的報告](https://godoc.org/github.com/toomore/gogrs/cmd/twsereport)
4. twsecache - [清除 twse cache](https://godoc.org/github.com/toomore/gogrs/cmd/twsecache)
5. tradingdays_server - [提供簡單的日期查詢 API Server](https://godoc.org/github.com/toomore/gogrs/cmd/tradingdays_server)

Docker
-------

Download image.

    docker pull toomore/gogrs

`tag:latest` bind to `branch:master`, more docker [info](https://registry.hub.docker.com/u/toomore/gogrs/).
Or [minify gogrs docker](https://registry.hub.docker.com/u/toomore/gogrs-mini/).

Run `tradingdays_server`.

    docker run -d -p 80:59123 toomore/gogrs tradingdays_server

Or login run other cmd

    docker run -it toomore/gogrs

Create a ramdisk volume

    docker create -v /run/shm/:/run/shm --name ramdisk toomore/gogrs-mini

Run with ramdisk volume

    docker run -it --volumes-from ramdisk toomore/gogrs

TODO
-----

1. docker-compose scale `tradingdays_server` with unused ports.
2. In English comment.
3. 盤中預估量能。
4. 個股對應分類股資訊。
5. 顯示三大法人。

License
--------

The MIT License (MIT)

Copyright © 2015 Toomore Chiang, http://toomore.net/ <toomore0929@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

http://toomore.mit-license.org/
