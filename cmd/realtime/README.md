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
	-color
		色彩化（default: true）
	-nonstop
		自動重複，單位秒數

範例

	realtime -twse=2618,2329 -otc=8446,2719 -index

回傳內容

	2015/06/11 17:46:28 發行量加權股價指數(t00) $9302.49(-28.50) -0.31% 5440/97357 [2015-06-11 13:33:00 +0800 CST] [20150611 17:46:29]
	2015/06/11 17:46:28 櫃買指數(o00) $135.34(0.08) 0.06% 1228/24569 [2015-06-11 13:33:00 +0800 CST] [20150611 17:46:29]
	2015/06/11 17:46:28 寶島股價指數(FRMSA) $10768.05(-29.39) -0.27% 0/0 [2015-06-11 13:33:00 +0800 CST] [20150611 17:46:29]
	2015/06/11 17:46:28 華研(8446) $140.00(-6.00) -4.11% 29/351 [2015-06-11 14:30:00 +0800 CST] [20150611 17:46:29]
	2015/06/11 17:46:28 華泰(2329) $13.10(-0.40) -2.96% 217/3993 [2015-06-11 14:30:00 +0800 CST] [20150611 17:46:29]
	2015/06/11 17:46:28 燦星旅(2719) $24.30(-0.70) -2.80% 7/107 [2015-06-11 14:30:00 +0800 CST] [20150611 17:46:28]
	2015/06/11 17:46:28 長榮航(2618) $19.30(-0.70) -3.50% 1001/19166 [2015-06-11 14:30:00 +0800 CST] [20150611 17:46:29]
