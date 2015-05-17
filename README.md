gogrs
======

[![GoDoc](https://godoc.org/github.com/toomore/gogrs?status.svg)](https://godoc.org/github.com/toomore/gogrs) [![Build Status](https://travis-ci.org/toomore/gogrs.svg?branch=master)](https://travis-ci.org/toomore/gogrs) [![Coverage Status](https://coveralls.io/repos/toomore/gogrs/badge.svg?branch=master)](https://coveralls.io/r/toomore/gogrs?branch=master)

gogrs now is still in development. I will try my best to speed up to completed the same function with [grs](https://github.com/toomore/grs) (Python). gogrs 是擷取台灣上市股票股價資訊工具，目前還在大量的開發中。原始工具是用 [grs](https://github.com/toomore/grs)（Python 套件），目標是將基本功能用 go 來實作。

Packages
---------

1. realtime - [擷取盤中個股、指數即時股價資訊](https://godoc.org/github.com/toomore/gogrs/realtime)
2. twse - [擷取台灣股市上市、上櫃股票資訊](https://godoc.org/github.com/toomore/gogrs/twse)
3. tradingdays - [股市開休市判斷（支援非國定假日：颱風假）與當日區間判斷（盤中、盤後、盤後盤）](https://godoc.org/github.com/toomore/gogrs/tradingdays)
4. utils - [套件所需的公用工具](https://godoc.org/github.com/toomore/gogrs/utils)

Cmd
----

1. gogrs_example - [簡單範例測試](https://godoc.org/github.com/toomore/gogrs/cmd/gogrs_example)
2. realtime - [擷取盤中即時資訊與大盤、上櫃、寶島指數](https://godoc.org/github.com/toomore/gogrs/cmd/realtime)
3. twsereport - [每日收盤後產生符合選股條件的報告](https://godoc.org/github.com/toomore/gogrs/cmd/twsereport)
4. tradingdays_server - [提供簡單的日期查詢 API Server](https://godoc.org/github.com/toomore/gogrs/cmd/tradingdays_server)
