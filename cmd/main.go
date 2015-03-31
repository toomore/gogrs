package main

import (
	"fmt"
	"time"

	"github.com/toomore/gogrs"
)

var stock = &gogrs.DailyData{
	No:   "2329",
	Date: time.Date(2015, 03, 20, 0, 0, 0, 0, time.Local),
}

func main() {
	stock.GetData()
	fmt.Println(stock.RawData)
	fmt.Println(stock.MA(6))
	fmt.Println(stock.MAV(6))
	fmt.Println(stock.GetPriceList())
	fmt.Println(gogrs.ThanPastFloat64(stock.GetPriceList(), 3, true))
	fmt.Println(gogrs.ThanPastFloat64(stock.GetPriceList(), 3, false))
	fmt.Println(stock.GetVolumeList())
	fmt.Println(gogrs.ThanPastUint64(stock.GetVolumeList(), 3, true))
	fmt.Println(gogrs.ThanPastUint64(stock.GetVolumeList(), 3, false))
	fmt.Println(stock.GetRangeList())
	fmt.Println(stock.IsRed())
}
