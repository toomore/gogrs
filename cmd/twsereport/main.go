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
	"github.com/toomore/gogrs/utils"
)

type base interface {
	MA(days int) []float64
	Len() int
	PlusData()
}

// MA 3 > 6 > 18
func check01(b base) bool {
	defer wg.Done()
	var start = b.Len()
	for {
		if b.Len() >= 18 {
			break
		}
		b.PlusData()
		if (b.Len() - start) == 0 {
			break
		}
	}
	var ma3 = b.MA(3)
	var ma6 = b.MA(6)
	var ma18 = b.MA(18)
	//log.Println(ma3[len(ma3)-1], ma6[len(ma6)-1], ma18[len(ma18)-1])
	if ma3[len(ma3)-1] > ma6[len(ma6)-1] && ma6[len(ma6)-1] > ma18[len(ma18)-1] {
		return true
	}
	return false
}

func check02(b base) bool {
	defer wg.Done()
	days, up := utils.CountCountineFloat64(utils.DeltaFloat64(b.MA(3)))
	if up && days > 1 {
		return true
	}
	return false
}

var wg sync.WaitGroup
var twseNo = flag.String("twse", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")

func main() {
	flag.Parse()
	var datalist []*twse.Data

	if *twseNo != "" {
		twselist := strings.Split(*twseNo, ",")
		datalist = make([]*twse.Data, len(twselist))

		for i, no := range twselist {
			datalist[i] = twse.NewTWSE(no, tradingdays.FindRecentlyOpened(time.Now()))
		}
	}

	if len(datalist) > 0 {
		for _, checkfunc := range []func(base) bool{check01, check02} {
			fmt.Printf("----- %v -----\n", checkfunc)
			for _, stock := range datalist {
				wg.Add(1)
				go func(checkfunc func(base) bool, stock *twse.Data) {
					runtime.Gosched()
					if checkfunc(base(stock)) {
						fmt.Printf("%s\n", stock.No)
					}
				}(checkfunc, stock)
			}
			wg.Wait()
		}
	} else {
		flag.PrintDefaults()
	}
}
