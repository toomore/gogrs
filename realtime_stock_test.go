package gogrs

import (
	"testing"
	"time"
)

func TestStockRealTime(*testing.T) {
	r := &StockRealTime{
		No:        "2618",
		Timestamp: RandInt(),
		Date:      time.Now(),
	}
	r.URL()
	r.GetData()
}
