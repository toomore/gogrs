gogrs - twsereport
===================

[![GoDoc](https://godoc.org/github.com/toomore/gogrs?status.svg)](https://godoc.org/github.com/toomore/gogrs/cmd/twsereport)
[![Build Status](https://travis-ci.org/toomore/gogrs.svg?branch=master)](https://travis-ci.org/toomore/gogrs)

每日收盤後產生符合選股條件的報告

Install:

	go install github.com/toomore/gogrs/cmd/twsereport

Usage:

	twsereport [flags]

The flags are:

	-twse
		上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329
	-twsecate
		上市股票類別，可使用 ',' 分隔多組代碼，例：11,15
	-otc
		上櫃股票代碼，可使用 ',' 分隔多組代碼，例：4406,8446
	-otccate
		上櫃股票類別，可使用 ',' 分隔多組代碼，例：02,14
	-ncpu
		指定 CPU 數量，預設為實際 CPU 數量
	-color
		色彩化
