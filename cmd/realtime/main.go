package main

import (
	"flag"
	"log"
	"time"

	"github.com/toomore/gogrs/realtime"
	"github.com/toomore/gogrs/utils"
)

func TaipeiNow() time.Time {
	d := time.Now()
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, utils.TaipeiTimeZone)
}

func prettyprint(data realtime.Data) {
	log.Printf("%s(%s) $%.2f(%.2f) %.0f/%.0f [%s] [%s %s]\n", data.Info.Name, data.Info.No,
		data.Price, data.Price-data.Open, data.Volume, data.VolumeAcc,
		data.TradeTime, data.SysInfo["sysDate"], data.SysInfo["sysTime"])
}

var twseNo = flag.String("twse", "", "TWSE NO.")
var otcNo = flag.String("otc", "", "OTC NO.")

func main() {
	flag.Parse()
	var r *realtime.StockRealTime
	if *twseNo != "" {
		r = realtime.NewTWSE(*twseNo, TaipeiNow())
		data, _ := r.Get()
		prettyprint(data)
	}
	if *otcNo != "" {
		r = realtime.NewOTC(*otcNo, TaipeiNow())
		data, _ := r.Get()
		prettyprint(data)
	}
	//flag.PrintDefaults()
}
