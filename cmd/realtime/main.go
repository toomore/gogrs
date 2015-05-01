// 擷取盤中即時資訊與大盤、上櫃、寶島指數.
//
/*

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

範例

	/realtime -twse=2618,2329 -otc=8446,2719 -index=true

回傳內容

	2015/04/29 23:43:13 櫃買指數(o00) $144.72(-1.15) 1271/23250 [2015-04-29 13:33:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 發行量加權股價指數(t00) $9853.83(-80.99) 5742/113762 [2015-04-29 13:33:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 寶島股價指數(FRMSA) $11416.04(-93.58) 0/0 [2015-04-29 13:33:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 華研(8446) $124.00(1.00) 6/293 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 長榮航(2618) $24.20(-0.55) 666/15618 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 燦星旅(2719) $24.15(0.25) 6/83 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]
	2015/04/29 23:43:13 華泰(2329) $15.95(0.50) 648/19995 [2015-04-29 14:30:00 +0800 CST] [20150429 23:43:13]

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/toomore/gogrs/realtime"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/utils"
)

var cacheTime time.Time

// TaipeiNow show Taipei Now time.
func TaipeiNow() time.Time {
	if cacheTime.Year() == 1 {
		d := time.Now().UTC()
		days := d.Day()
		for {
			if tradingdays.IsOpen(d.Year(), d.Month(), days) {
				break
			}
			days--
		}
		if d.Before(time.Date(d.Year(), d.Month(), days, 1, 0, 0, 0, time.FixedZone("UTC", 0))) {
			if d.Hour() == 0 && d.Minute() <= 59 {
				days--
			}
		}
		cacheTime = time.Date(d.Year(), d.Month(), days, 0, 0, 0, 0, utils.TaipeiTimeZone)
	}
	return cacheTime
}

func prettyprint(data realtime.Data) string {
	return fmt.Sprintf("%s(%s) $%.2f(%.2f) %.0f/%.0f [%s] [%s %s]",
		data.Info.Name, data.Info.No,
		data.Price, data.Price-data.Open, data.Volume, data.VolumeAcc,
		data.TradeTime, data.SysInfo["sysDate"], data.SysInfo["sysTime"])
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var twseNo = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")
var otcNo = flag.String("otc", "", "上櫃股票代碼，可使用 ',' 分隔多組代碼，例：8446,2719")
var index = flag.Bool("index", false, "顯示大盤、上櫃、寶島指數（default: false）")

func main() {
	flag.Parse()
	queue := []*realtime.StockRealTime{}
	if *twseNo != "" {
		for _, no := range strings.Split(*twseNo, ",") {
			queue = append(queue, realtime.NewTWSE(no, TaipeiNow()))
		}
	}
	if *otcNo != "" {
		for _, no := range strings.Split(*otcNo, ",") {
			queue = append(queue, realtime.NewOTC(no, TaipeiNow()))
		}
	}
	if *index {
		queue = append(queue, []*realtime.StockRealTime{realtime.NewWeight(TaipeiNow()),
			realtime.NewOTCI(TaipeiNow()), realtime.NewFRMSA(TaipeiNow())}...)
	}
	result := make(chan string, runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(len(queue))
	for _, r := range queue {
		go func(r *realtime.StockRealTime) {
			runtime.Gosched()
			data, _ := r.Get()
			result <- prettyprint(data)
		}(r)
	}
	go func() {
		for v := range result {
			wg.Done()
			log.Println(v)
		}
	}()
	wg.Wait()
	if len(queue) == 0 {
		flag.PrintDefaults()
	}
}
