package gogrs

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
		Date:      time.Now(),
	}
	r.URL()
	r.GetData()
}

func ExampleStockRealTime() {
	r := StockRealTime{
		No:        "2618",
		Timestamp: utils.RandInt(),
		Date:      time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}

	data, _ := r.GetData()
	fmt.Printf("%v", data)
}
