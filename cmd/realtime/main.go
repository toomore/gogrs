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

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var twseNo = flag.String("twse", "", "TWSE NO.")
var otcNo = flag.String("otc", "", "OTC NO.")
var index = flag.Bool("index", false, "Show Index.")

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
	result := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(queue))
	for _, r := range queue {
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
	//flag.PrintDefaults()
}
