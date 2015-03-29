package main

import (
	"fmt"
	"time"

	"github.com/toomore/gogrs"
)

var stock2618 = &gogrs.DailyData{
	No:   "2618",
	Date: time.Date(2015, 03, 26, 0, 0, 0, 0, time.Local),
}

func main() {
	stock2618.GetData()
	fmt.Println(stock2618.RawData)
	fmt.Println(stock2618.GetVolumeList())
}
