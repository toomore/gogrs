package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/toomore/gogrs/realtime"
	"github.com/toomore/gogrs/utils"
)

// TaipeiNow show Taipei Now time.
func TaipeiNow() time.Time {
	d := time.Now()
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, utils.TaipeiTimeZone)
}

func prettyprint(data realtime.Data) string {
	return fmt.Sprintf("%s(%s) $%.2f(%.2f) %.0f/%.0f [%s] [%s %s]",
		data.Info.Name, data.Info.No,
		data.Price, data.Price-data.Open, data.Volume, data.VolumeAcc,
		data.TradeTime, data.SysInfo["sysDate"], data.SysInfo["sysTime"])
}

var twseNo = flag.String("twse", "", "TWSE NO.")
var otcNo = flag.String("otc", "", "OTC NO.")
var index = flag.Bool("index", false, "Show Index.")

func main() {
	flag.Parse()
	var r *realtime.StockRealTime
	if *twseNo != "" {
		r = realtime.NewTWSE(*twseNo, TaipeiNow())
		data, _ := r.Get()
		log.Println(prettyprint(data))
	}
	if *otcNo != "" {
		r = realtime.NewOTC(*otcNo, TaipeiNow())
		data, _ := r.Get()
		log.Println(prettyprint(data))
	}
	if *index {
		var wg sync.WaitGroup
		indexList := [3]*realtime.StockRealTime{realtime.NewWeight(TaipeiNow()),
			realtime.NewOTCI(TaipeiNow()), realtime.NewFRMSA(TaipeiNow())}
		result := make(chan string)
		wg.Add(3)
		for _, r := range indexList {
			go func(r *realtime.StockRealTime) {
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
	}
	//flag.PrintDefaults()
}
