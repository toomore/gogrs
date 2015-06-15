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
	-otccate
		上櫃股票類別，可使用 ',' 分隔多組代碼，例：04,11
	-showcatelist
		顯示上市分類表（default: false）
	-ncpu
		指定 CPU 數量，預設為實際 CPU 數量
	-pt
		計算花費時間
	-count
		計算此次查詢的漲跌家數（default: true）
	-color
		色彩化（default: true）

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

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/toomore/gogrs/realtime"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

var cacheTime time.Time

func init() {
	TaipeiNow()
}

// TaipeiNow show Taipei Now time.
func TaipeiNow() time.Time {
	if cacheTime.Year() == 1 {
		d := time.Now().UTC()
		days := d.Day()
		var index int
		for {
			if tradingdays.IsOpen(d.Year(), d.Month(), days) {
				if index == 0 {
					if tradingdays.NewTimePeriod(time.Date(d.Year(), d.Month(), days, d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())).AtBefore() {
						days--
						for {
							if tradingdays.IsOpen(d.Year(), d.Month(), days) {
								break
							}
							days--
						}
					}
				}
				break
			}
			days--
			index++
		}
		cacheTime = time.Date(d.Year(), d.Month(), days, 0, 0, 0, 0, utils.TaipeiTimeZone)
	}
	return cacheTime
}

var (
	chanbuf      int
	twseNo       = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")
	twseCate     = flag.String("twsecate", "", "上市股票類別，可使用 ',' 分隔多組代碼，例：11,15")
	showcatelist = flag.Bool("showcatelist", false, "顯示上市分類表")
	otcNo        = flag.String("otc", "", "上櫃股票代碼，可使用 ',' 分隔多組代碼，例：8446,2719")
	otcCate      = flag.String("otccate", "", "上櫃股票類別，可使用 ',' 分隔多組代碼，例：04,11")
	index        = flag.Bool("index", false, "顯示大盤、上櫃、寶島指數（default: false）")
	ncpu         = flag.Int("ncpu", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
	pt           = flag.Bool("pt", false, "計算花費時間")
	count        = flag.Bool("count", true, "計算此次查詢的漲跌家數")
	showcolor    = flag.Bool("color", true, "色彩化")
	white        = color.New(color.FgWhite, color.Bold).SprintfFunc()
	red          = color.New(color.FgRed, color.Bold).SprintfFunc()
	green        = color.New(color.FgGreen, color.Bold).SprintfFunc()
	yellowBold   = color.New(color.FgYellow, color.Bold).SprintfFunc()
	cyan         = color.New(color.FgCyan).SprintfFunc()
)

func prettyprint(data realtime.Data) string {
	var (
		RangeValue  = data.Price - data.Open
		outputcolor func(string, ...interface{}) string
	)
	switch {
	case RangeValue > 0:
		outputcolor = red
	case RangeValue < 0:
		outputcolor = green
	default:
		outputcolor = white
	}
	return fmt.Sprintf("%s %s %s",
		yellowBold("%s(%s)", data.Info.Name, data.Info.No),
		outputcolor("$%.2f(%.2f) %.2f%% %.0f/%.0f",
			data.Price, data.Price-data.Open, RangeValue/data.Open*100, data.Volume, data.VolumeAcc),
		cyan("[%s] [%s %s]",
			data.TradeTime, data.SysInfo["sysDate"], data.SysInfo["sysTime"]),
	)
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	color.NoColor = !*showcolor

	runtime.GOMAXPROCS(*ncpu)
	chanbuf = *ncpu * 2

	queue := make(chan *realtime.StockRealTime, chanbuf)
	defer close(queue)

	var (
		counter   int
		down      int
		startTime time.Time
		up        int
		wg        sync.WaitGroup
	)

	if *pt {
		startTime = time.Now()
	}

	if *showcatelist {
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
		l := twse.NewLists(tradingdays.FindRecentlyOpened(time.Now()))
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

	if *otcCate != "" {
		o := twse.NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, no := range strings.Split(*otcCate, ",") {
			for _, s := range o.GetCategoryList(no) {
				wg.Add(1)
				go func(s twse.StockInfo) {
					runtime.Gosched()
					queue <- realtime.NewOTC(s.No, TaipeiNow())
				}(s)
			}
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
				if *count {
					counter++
					if data.Price-data.Open > 0 {
						up++
					} else if data.Price-data.Open < 0 {
						down++
					}
				}
			}(r)
		}
	}()
	wg.Wait()
	if *count {
		fmt.Printf("All: %d, Up: %d, Down: %d, Same: %d\n",
			counter, up, down, counter-up-down)
	}
	if *pt {
		defer fmt.Println(time.Now().Sub(startTime))
	}
}
