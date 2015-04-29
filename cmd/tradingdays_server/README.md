gogrs - tradingdays_server
===========================

[![GoDoc](https://godoc.org/github.com/toomore/gogrs?status.svg)](https://godoc.org/github.com/toomore/gogrs/tradingdays/tradingdays_server)
[![Build Status](https://travis-ci.org/toomore/gogrs.svg?branch=master)](https://travis-ci.org/toomore/gogrs)

主要支援動態更新 CSV 檔案讀取，解決非預定開休市狀況（如：颱風假）

Install:

	go install -u github.com/toomore/gogrs/cmd/tradingdays_server

Usage:

	tradingdays_server [flags]

The flags are:

	-http
		HTTP service address (default ':59123')
	-csvcachetime
		CSV cache time.(default: 21600)

URL Path:

	/open?q={timestamp}

回傳 JSON 格式

	{
		"date": "2015-04-24T15:14:52Z",
		"open": true
	}

範例：

	http://gogrs-trd.toomore.net/open?q=1429888492
