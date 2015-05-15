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

type check01 struct{}

func (check01) String() string {
	return "MA 3 > 6 > 18"
}

func (check01) CheckFunc(b ...base) bool {
	defer wg.Done()
	var d = b[0]
	var start = d.Len()
	for {
		if d.Len() >= 18 {
			break
		}
		d.PlusData()
		if (d.Len() - start) == 0 {
			break
		}
	}
	if d.Len() < 18 {
		return false
	}
	var ma3 = d.MA(3)
	var ma6 = d.MA(6)
	var ma18 = d.MA(18)
	//log.Println(ma3[len(ma3)-1], ma6[len(ma6)-1], ma18[len(ma18)-1])
	if ma3[len(ma3)-1] > ma6[len(ma6)-1] && ma6[len(ma6)-1] > ma18[len(ma18)-1] {
		return true
	}
	return false
}

type check02 struct{}

func (check02) String() string {
	return "check02"
}
func (check02) CheckFunc(b ...base) bool {
	defer wg.Done()
	if b[0].Len() < 3 {
		return false
	}
	days, up := utils.CountCountineFloat64(utils.DeltaFloat64(b[0].MA(3)))
	if up && days > 1 {
		return true
	}
	return false
}

type checkGroup interface {
	String() string
	CheckFunc(...base) bool
}

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
var ckList = make(checkGroupList, 1)

func init() {
	ckList.Add(checkGroup(check01{}))
	ckList.Add(checkGroup(check02{}))
}

func main() {
	flag.Parse()
	var datalist []*twse.Data
	var catelist []twse.StockInfo
	var twselist []string
	var catenolist []string

	if *twseCate != "" {
		l := &twse.Lists{Date: tradingdays.FindRecentlyOpened(time.Now())}
		catelist = l.GetCategoryList(*twseCate)
		catenolist = make([]string, len(catelist))
		for i, s := range catelist {
			catenolist[i] = s.No
		}
	}

	if *twseNo != "" {
		twselist = strings.Split(*twseNo, ",")
	}
	datalist = make([]*twse.Data, len(twselist)+len(catelist))

	for i, no := range append(twselist, catenolist...) {
		datalist[i] = twse.NewTWSE(no, tradingdays.FindRecentlyOpened(time.Now()))
	}

	if len(datalist) > 0 {
		for _, check := range ckList {
			fmt.Printf("----- %v -----\n", check)
			wg.Add(len(datalist))
			for _, stock := range datalist {
				go func(check checkGroup, stock *twse.Data) {
					runtime.Gosched()
					if check.CheckFunc(stock) {
						fmt.Printf("%s\n", stock.No)
					}
				}(check, stock)
			}
			wg.Wait()
		}
	} else {
		flag.PrintDefaults()
	}
}
