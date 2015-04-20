package main

import (
	"fmt"
	"time"

	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

var stock = twse.NewTWSE("2329", time.Date(2015, 3, 20, 0, 0, 0, 0, time.Local))

// ShowAll is lazy to show all XD.
func ShowAll(stock *twse.Data) {
	fmt.Println(stock.RawData)
	fmt.Println(stock.MA(6))
	fmt.Println(stock.MAV(6))
	fmt.Println(stock.GetPriceList())
	fmt.Println(utils.ThanPastFloat64(stock.GetPriceList(), 3, true))
	fmt.Println(utils.ThanPastFloat64(stock.GetPriceList(), 3, false))
	fmt.Println(stock.GetVolumeList())
	fmt.Println(utils.ThanPastUint64(stock.GetVolumeList(), 3, true))
	fmt.Println(utils.ThanPastUint64(stock.GetVolumeList(), 3, false))
	fmt.Println(stock.GetRangeList())
	fmt.Println(stock.IsRed())
}

func main() {
	stock.Get()
	ShowAll(stock)
	fmt.Println("-----------------------------")
	stock.PlusData()
	ShowAll(stock)
	fmt.Println("-----------------------------")
	fmt.Println(tradingdays.IsOpen(2015, 5, 1, utils.TaipeiTimeZone))
}
