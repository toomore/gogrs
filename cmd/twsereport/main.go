// 每日收盤後產生符合選股條件的報告.
//
/*
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
	-showcatelist
		顯示上市/上櫃分類表（default: false）
	-ncpu
		指定 CPU 數量，預設為實際 CPU 數量
	-color
		色彩化（default: true）

可以重新調整自己的條件組合，目前預設的為：

	1. MA 3 > 6 > 18
	2. 量大於前三天 K 線收紅
	3. 量或價走平 45 天
	4. (MA3 < MA6) > MA18 and MA3UP(1)
	5. 三日內最大量 K 線收紅 收在 MA18 之上
	6. 漲幅 7% 以上
	7. 多方力道 > 0.75

*/
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
)

type checkGroupList []checkGroup

func (c *checkGroupList) Add(f checkGroup) {
	if (*c)[0] == nil {
		(*c)[0] = f
	} else {
		*c = append(*c, f)
	}
}

var (
	wg           sync.WaitGroup
	twseNo       = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")
	twseCate     = flag.String("twsecate", "", "上市股票類別，可使用 ',' 分隔多組代碼，例：11,15")
	otcNo        = flag.String("otc", "", "上櫃股票代碼，可使用 ',' 分隔多組代碼，例：4406,8446")
	otcCate      = flag.String("otccate", "", "上櫃股票類別，可使用 ',' 分隔多組代碼，例：02,14")
	showcatelist = flag.Bool("showcatelist", false, "顯示上市/上櫃分類表")
	showcolor    = flag.Bool("color", true, "色彩化")
	ncpu         = flag.Int("ncpu", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
	ckList       = make(checkGroupList, 1)
	white        = color.New(color.FgWhite, color.Bold).SprintfFunc()
	red          = color.New(color.FgRed, color.Bold).SprintfFunc()
	green        = color.New(color.FgGreen, color.Bold).SprintfFunc()
	yellow       = color.New(color.FgYellow).SprintfFunc()
	yellowBold   = color.New(color.FgYellow, color.Bold).SprintfFunc()
	blue         = color.New(color.FgBlue).SprintfFunc()
)

func init() {
	runtime.GOMAXPROCS(*ncpu)
}

func prettyprint(stock *twse.Data, check checkGroup) string {
	var (
		Open        = stock.GetOpenList()[len(stock.GetOpenList())-1]
		Price       = stock.GetPriceList()[len(stock.GetPriceList())-1]
		RangeValue  = stock.GetDailyRangeList()[len(stock.GetDailyRangeList())-1]
		Volume      = stock.GetVolumeList()[len(stock.GetVolumeList())-1] / 1000
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

	return fmt.Sprintf("%s %s %s %s%s %s %s",
		yellow("[%s]", check),
		blue("%s", stock.RawData[stock.Len()-1][0]),
		outputcolor("%s %s", stock.No, stock.Name),
		outputcolor("$%.2f", Price),
		outputcolor("(%.2f)", RangeValue),
		outputcolor("%.2f%%", RangeValue/Open*100),
		outputcolor("%d", Volume),
	)
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var (
		datalist     []*twse.Data
		l            *twse.Lists
		o            *twse.OTCLists
		otccatelist  []string
		otcdelta     int
		otclist      []string
		twsecatelist []string
		twselist     []string
	)

	color.NoColor = !*showcolor

	if *showcatelist {
		var (
			cateTitle    = []string{"The same with TWSE/OTC", "OnlyTWSE", "OnlyOTC"}
			categoryList = twse.NewCategoryList()
			index        int
		)

		for i, cate := range []map[string]string{
			categoryList.Same(), categoryList.OnlyTWSE(), categoryList.OnlyOTC(),
		} {
			index = 1
			fmt.Println(white("---------- %s ----------", cateTitle[i]))
			for cateNo, cateName := range cate {
				fmt.Printf("%s\t", fmt.Sprintf("%s%s", green("%s", cateName), yellowBold("(%s)", cateNo)))
				if index%3 == 0 {
					fmt.Println("")
				}
				index++
			}
			fmt.Println("")
		}
	}

	if *twseCate != "" {
		l = twse.NewLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, v := range strings.Split(*twseCate, ",") {
			for _, s := range l.GetCategoryList(v) {
				twsecatelist = append(twsecatelist, s.No)
			}
		}
	}

	if *otcCate != "" {
		o = twse.NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, v := range strings.Split(*otcCate, ",") {
			for _, s := range o.GetCategoryList(v) {
				otccatelist = append(otccatelist, s.No)
			}
		}
	}

	if *twseNo != "" {
		twselist = strings.Split(*twseNo, ",")
	}
	if *otcNo != "" {
		otclist = strings.Split(*otcNo, ",")
	}
	datalist = make([]*twse.Data, len(twselist)+len(twsecatelist)+len(otclist)+len(otccatelist))

	for i, no := range append(twselist, twsecatelist...) {
		datalist[i] = twse.NewTWSE(no, tradingdays.FindRecentlyOpened(time.Now()))
	}
	otcdelta = len(twselist) + len(twsecatelist)
	for i, no := range append(otclist, otccatelist...) {
		datalist[i+otcdelta] = twse.NewOTC(no, tradingdays.FindRecentlyOpened(time.Now()))
	}

	if len(datalist) > 0 {
		for _, check := range ckList {
			fmt.Println(yellowBold("----- %v -----", check))
			wg.Add(len(datalist))
			for _, stock := range datalist {
				go func(check checkGroup, stock *twse.Data) {
					defer wg.Done()
					runtime.Gosched()
					if check.CheckFunc(stock) {
						fmt.Println(prettyprint(stock, check))
					}
				}(check, stock)
			}
			wg.Wait()
		}
	}
}
