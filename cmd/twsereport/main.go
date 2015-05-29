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
	-ncpu
		指定 CPU 數量，預設為實際 CPU 數量

*/
package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

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

var wg sync.WaitGroup
var twseNo = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")
var twseCate = flag.String("twsecate", "", "上市股票類別，可使用 ',' 分隔多組代碼，例：11,15")
var ncpu = flag.Int("ncpu", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
var ckList = make(checkGroupList, 1)

func init() {
	runtime.GOMAXPROCS(*ncpu)
}

func prettyprint(stock *twse.Data, check checkGroup) string {
	return fmt.Sprintf("[%s] %s %s %s $%.2f(%.2f) %d",
		check,
		stock.RawData[stock.Len()-1][0],
		stock.No, stock.Name,
		stock.GetPriceList()[len(stock.GetPriceList())-1],
		stock.GetRangeList()[len(stock.GetRangeList())-1],
		stock.GetVolumeList()[len(stock.GetVolumeList())-1]/1000,
	)
}

func main() {
	flag.Parse()
	var datalist []*twse.Data
	var catelist []twse.StockInfo
	var twselist []string
	var catenolist []string

	if *twseCate != "" {
		l := &twse.Lists{Date: tradingdays.FindRecentlyOpened(time.Now())}

		for _, v := range strings.Split(*twseCate, ",") {
			catelist = l.GetCategoryList(v)
			for _, s := range catelist {
				catenolist = append(catenolist, s.No)
			}
		}
	}

	if *twseNo != "" {
		twselist = strings.Split(*twseNo, ",")
	}
	datalist = make([]*twse.Data, len(twselist)+len(catenolist))

	for i, no := range append(twselist, catenolist...) {
		datalist[i] = twse.NewTWSE(no, tradingdays.FindRecentlyOpened(time.Now()))
	}

	if len(datalist) > 0 {
		for _, check := range ckList {
			fmt.Printf("----- %v -----\n", check)
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
	} else {
		flag.PrintDefaults()
	}
}
