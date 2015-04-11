package realtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestStockRealTime(*testing.T) {
	r := &StockRealTime{
		No: "2618",
		//Date:      time.Now(),
		Date:     time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone),
		Exchange: "tse",
	}

	r.URL()
	v, _ := r.Get()
	fmt.Println(v.BestAskPrice)
	fmt.Println(v.BestBidPrice)
	fmt.Println(v.BestAskVolume)
	fmt.Println(v.BestBidVolume)
	fmt.Println(v)
	fmt.Println("UnixMapData", r.UnixMapData)
}

func TestStockRealTimeOTC(*testing.T) {
	r := &StockRealTime{
		No: "8446",
		//Date:      time.Now(),
		Date:     time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone),
		Exchange: "otc",
	}

	r.URL()
	v, _ := r.Get()
	fmt.Println(v.BestAskPrice)
	fmt.Println(v.BestBidPrice)
	fmt.Println(v.BestAskVolume)
	fmt.Println(v.BestBidVolume)
	fmt.Println(v)
	fmt.Println("UnixMapData", r.UnixMapData)
}

func TestStockRealTimeIndexs(*testing.T) {
	var date = time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone)

	weight := NewWeight(date)
	otc := NewOTC(date)
	farmsa := NewFRMSA(date)
	fmt.Println(weight.Get())
	fmt.Println(otc.Get())
	fmt.Println(farmsa.Get())
}

func BenchmarkGet(b *testing.B) {
	r := &StockRealTime{
		No: "2618",
		//Date:      time.Now(),
		Date:     time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone),
		Exchange: "tse",
	}

	for i := 0; i <= b.N; i++ {
		r.Get()
	}
}

// 擷取 長榮航(2618) 上市即時盤股價資訊
func ExampleStockRealTime_Get_twse() {
	r := StockRealTime{
		No:       "2618",
		Date:     time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone),
		Exchange: "tse",
	}

	data, _ := r.Get()
	fmt.Printf("%v", data)
}

// 擷取 華研(8446) 上櫃即時盤股價資訊
func ExampleStockRealTime_Get_otc() {
	r := StockRealTime{
		No:       "8446",
		Date:     time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone),
		Exchange: "otc",
	}

	data, _ := r.Get()
	fmt.Printf("%v", data.Info)
	// output:
	// {otc 華研國際音樂股份有限公司  華研 8446.tw}
}

func ExampleStockRealTime_NewWeight() {
	weight := NewWeight(time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))
	data, _ := weight.Get()
	fmt.Println(data.Info)
	// output:
	// {tse   發行量加權股價指數 t00.tw}
}

func ExampleStockRealTime_NewOTC() {
	otc := NewOTC(time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))
	data, _ := otc.Get()
	fmt.Println(data.Info)
	// output:
	// {otc   櫃檯指數 o00.tw}
}

func ExampleStockRealTime_NewFRMSA() {
	farmsa := NewFRMSA(time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))
	data, _ := farmsa.Get()
	fmt.Println(data.Info)
	// output:
	// {tse   寶島股價指數 FRMSA.tw}
}
