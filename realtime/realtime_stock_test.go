package realtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestStockRealTime(*testing.T) {
	r := &StockRealTime{
		No:        "2618",
		Timestamp: utils.RandInt(),
		//Date:      time.Now(),
		Date: time.Date(2015, 4, 1, 0, 0, 0, 0, time.Local),
	}
	r.URL()
	v, _ := r.Get()
	fmt.Println(v.BestAskPrice)
	fmt.Println(v.BestBidPrice)
	fmt.Println(v.BestAskVolume)
	fmt.Println(v.BestBidVolume)
	fmt.Println("UnixMapData", r.UnixMapData)
}

func ExampleStockRealTime() {
	r := StockRealTime{
		No:        "2618",
		Timestamp: utils.RandInt(),
		Date:      time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}

	data, _ := r.Get()
	fmt.Printf("%v", data)
}
