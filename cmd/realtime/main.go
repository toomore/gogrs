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
	-twsecate
		上市股票類別，可使用 ',' 分隔多組代碼，例：11,15
	-showtwsecatelist
		顯示上市分類表（default: false）
	-ncpu
		指定 CPU 數量，預設為實際 CPU 數量
	-pt
		計算花費時間

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
	"github.com/toomore/gogrs/twse"
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

var chanbuf int

var twseNo = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")
var twseCate = flag.String("twsecate", "", "上市股票類別，可使用 ',' 分隔多組代碼，例：11,15")
var showtwsecatelist = flag.Bool("showcatelist", false, "顯示上市分類表")
var otcNo = flag.String("otc", "", "上櫃股票代碼，可使用 ',' 分隔多組代碼，例：8446,2719")
var index = flag.Bool("index", false, "顯示大盤、上櫃、寶島指數（default: false）")
var ncpu = flag.Int("ncpu", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
var pt = flag.Bool("pt", false, "計算花費時間")

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*ncpu)
	chanbuf = *ncpu * 2

	queue := make(chan *realtime.StockRealTime, chanbuf)
	defer close(queue)
	var wg sync.WaitGroup

	var startTime time.Time
	if *pt {
		startTime = time.Now()
	}

	if *showtwsecatelist {
		var index = 1
		for cateNo, cateName := range twse.TWSECLASS {
			fmt.Printf("%s(%s)\t\t", cateName, cateNo)
			if index%3 == 0 {
				fmt.Println("")
			}
			index++
		}
		fmt.Println("")
	}

	if *twseNo != "" {
		for _, no := range strings.Split(*twseNo, ",") {
			wg.Add(1)
			go func(no string) {
				runtime.Gosched()
				queue <- realtime.NewTWSE(no, TaipeiNow())
			}(no)
		}
	}

	if *twseCate != "" {
		l := &twse.Lists{Date: tradingdays.FindRecentlyOpened()}
		for _, no := range strings.Split(*twseCate, ",") {
			for _, s := range l.GetCategoryList(no) {
				wg.Add(1)
				go func(s twse.StockInfo) {
					runtime.Gosched()
					queue <- realtime.NewTWSE(s.No, TaipeiNow())
				}(s)
			}
		}
	}

	if *otcNo != "" {
		for _, no := range strings.Split(*otcNo, ",") {
			wg.Add(1)
			go func(no string) {
				runtime.Gosched()
				queue <- realtime.NewOTC(no, TaipeiNow())
			}(no)
		}
	}
	if *index {
		wg.Add(3)
		for _, r := range []*realtime.StockRealTime{realtime.NewWeight(TaipeiNow()), realtime.NewOTCI(TaipeiNow()), realtime.NewFRMSA(TaipeiNow())} {
			go func(r *realtime.StockRealTime) {
				runtime.Gosched()
				queue <- r
			}(r)
		}
	}
	go func() {
		for r := range queue {
			go func(r *realtime.StockRealTime) {
				defer wg.Done()
				runtime.Gosched()
				data, _ := r.Get()
				log.Println(prettyprint(data))
			}(r)
		}
	}()
	wg.Wait()
	if *pt {
		defer fmt.Println(time.Now().Sub(startTime))
	}
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}
}
