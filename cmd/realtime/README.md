gogrs - realtime 
=================

[![GoDoc](https://godoc.org/github.com/toomore/gogrs?status.svg)](https://godoc.org/github.com/toomore/gogrs/cmd/realtime)
[![Build Status](https://travis-ci.org/toomore/gogrs.svg?branch=master)](https://travis-ci.org/toomore/gogrs)

擷取盤中即時資訊與大盤、上櫃、寶島指數.

Install:

	go install github.com/toomore/gogrs/cmd/realtime

Usage:

	realtime [flags]

The flags are:

	-twse
		上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329
	-otc
		上櫃股票代碼，可使用 ',' 分隔多組代碼，例：8446,2719
	-index
		顯示大盤、上櫃、寶島指數（default: false）
	-twsecate
		上市股票類別，可使用 ',' 分隔多組代碼，例：11,15
	-showtwsecatelist
		顯示上市分類表（default: false）
	-ncpu
		指定 CPU 數量，預設為實際 CPU 數量
	-pt
		計算花費時間
	-count
		計算此次查詢的漲跌家數（default: true）

範例

	realtime -twse=2618,2329 -otc=8446,2719 -index

回傳內容

	2015/04/29 23:43:13 櫃買指數(o00) $144.72(-1.15) 1271/23250 [2015-04-29 13:33:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 發行量加權股價指數(t00) $9853.83(-80.99) 5742/113762 [2015-04-29 13:33:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 寶島股價指數(FRMSA) $11416.04(-93.58) 0/0 [2015-04-29 13:33:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 華研(8446) $124.00(1.00) 6/293 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 長榮航(2618) $24.20(-0.55) 666/15618 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 燦星旅(2719) $24.15(0.25) 6/83 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 華泰(2329) $15.95(0.50) 648/19995 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
