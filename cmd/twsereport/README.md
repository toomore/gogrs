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

可以重新調整自己的條件組合，目前預設的為：

    1. MA 3 > 6 > 18
    2. 量大於前三天 K 線收紅
    3. 量或價走平 45 天
    4. (MA3 < MA6) > MA18 and MA3UP(1)
    5. 三日內最大量 K 線收紅 收在 MA18 之上
    6. 漲幅 7% 以上
    7. 多方力道 > 0.75
