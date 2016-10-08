package realtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/utils"
)

func TestStockRealTime(t *testing.T) {
	r := NewTWSE("2618", tradingdays.FindRecentlyOpened(time.Now()))

	t.Log(r.URL())
	r.Get()
	t.Logf("%+v", r)
}

func TestStockRealTime_noData(t *testing.T) {
	r := NewTWSE("26188", tradingdays.FindRecentlyOpened(time.Now()))

	_, err := r.Get()
	if err != errorNotEnoughData {
		t.Error("Should be \"No Data.\"")
	}
}

func TestStockRealTime_URL(t *testing.T) {
	r := NewTWSE("2618", time.Date(2015, 4, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))
	r.URL()
}

func TestStockRealTimeOTC(t *testing.T) {
	r := NewOTC("8446", tradingdays.FindRecentlyOpened(time.Now()))
	r.URL()
	t.Log(r.Get())
}

func TestStockRealTimeIndexs(*testing.T) {
	var date = tradingdays.FindRecentlyOpened(time.Now())

	weight := NewWeight(date)
	otc := NewOTCI(date)
	farmsa := NewFRMSA(date)
	weight.Get()
	otc.Get()
	farmsa.Get()
}

func BenchmarkGet(b *testing.B) {
	r := NewTWSE("2618", tradingdays.FindRecentlyOpened(time.Now()))

	for i := 0; i <= b.N; i++ {
		r.Get()
	}
}

// 擷取 長榮航(2618) 上市即時盤股價資訊
func ExampleStockRealTime_Get_twse() {
	r := NewTWSE("2618", tradingdays.FindRecentlyOpened(time.Now()))

	data, _ := r.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:tse FullName:長榮航空股份有限公司 Name:長榮航 No:2618 Ticker:2618.tw Category:15}
}

// 擷取 華研(8446) 上櫃即時盤股價資訊
func ExampleStockRealTime_Get_otc() {
	r := NewOTC("8446", tradingdays.FindRecentlyOpened(time.Now()))

	data, _ := r.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:otc FullName:華研國際音樂股份有限公司 Name:華研 No:8446 Ticker:8446.tw Category:32}
}

func ExampleNewWeight() {
	weight := NewWeight(tradingdays.FindRecentlyOpened(time.Now()))
	data, _ := weight.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:tse FullName: Name:發行量加權股價指數 No:t00 Ticker:t00.tw Category:tidx.tw}
}

func ExampleNewOTCI() {
	otc := NewOTCI(tradingdays.FindRecentlyOpened(time.Now()))
	data, _ := otc.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:otc FullName: Name:櫃買指數 No:o00 Ticker:o00.tw Category:oidx.tw}
}

func ExampleNewFRMSA() {
	farmsa := NewFRMSA(tradingdays.FindRecentlyOpened(time.Now()))
	data, _ := farmsa.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:tse FullName: Name:寶島股價指數 No:FRMSA Ticker:FRMSA.tw Category:tidx.tw}
}

func ExampleNewTWSE() {
	twse := NewTWSE("2618", tradingdays.FindRecentlyOpened(time.Now()))
	data, _ := twse.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:tse FullName:長榮航空股份有限公司 Name:長榮航 No:2618 Ticker:2618.tw Category:15}
}

func ExampleNewOTC() {
	otc := NewOTC("8446", tradingdays.FindRecentlyOpened(time.Now()))
	data, _ := otc.Get()
	fmt.Printf("%+v", data.Info)
	// output:
	// {Exchange:otc FullName:華研國際音樂股份有限公司 Name:華研 No:8446 Ticker:8446.tw Category:32}
}
