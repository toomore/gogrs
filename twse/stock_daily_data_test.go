package twse

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func assertType(t *testing.T, t1 interface{}, t2 interface{}) {
	if reflect.TypeOf(t1) != reflect.TypeOf(t2) {
		t.Errorf("Diff type t1(%s), t2(%s)", reflect.TypeOf(t1), reflect.TypeOf(t2))
	}
}

func TestURL(t *testing.T) {
	var d = NewTWSE("2618", time.Now())
	assertType(t, d, &Data{})
}

func ExampleData() {
	var d = NewTWSE("2618", time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))
	fmt.Println(d.Date)

	stockData, _ := d.Get()
	fmt.Println(stockData[0])
	// output:
	// 2014-12-26 00:00:00 +0800 Asia/Taipei
	// [103/12/01 64,418,143 1,350,179,448 20.20 21.40 20.20 21.35 +1.35 13,249]
}

func TestData_Get(*testing.T) {
	var d = NewTWSE("2618", time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))

	d.Get()
	d.Get() // Test Cache.
}

var twse = NewTWSE("2329", time.Date(2015, 03, 20, 0, 0, 0, 0, utils.TaipeiTimeZone))
var otc = NewOTC("8446", time.Date(2015, 03, 20, 0, 0, 0, 0, utils.TaipeiTimeZone))

func TestGetList(*testing.T) {
	for _, stock := range []*Data{twse, otc} {
		stock.URL()
		stock.Get()
		//fmt.Println(stock.RawData)
		stock.MA(6)
		stock.MAV(6)
		stock.GetPriceList()
		utils.ThanPastFloat64(stock.GetPriceList(), 3, true)
		utils.ThanPastFloat64(stock.GetPriceList(), 3, false)
		stock.GetVolumeList()
		utils.ThanPastUint64(stock.GetVolumeList(), 3, true)
		utils.ThanPastUint64(stock.GetVolumeList(), 3, false)
		stock.GetRangeList()
		stock.GetOpenList()
		stock.IsRed()
	}
}

func TestMABR(t *testing.T) {
	twse.Get()
	var sample1mabr = twse.MABR(3, 6)
	var sample1ma3 = twse.MA(3)
	var sample1ma6 = twse.MA(6)
	if sample1mabr[len(sample1mabr)-1] != sample1ma3[len(sample1ma3)-1]-sample1ma6[len(sample1ma6)-1] {
		t.Error("Should be the sample")
	}

	var sample2mavbr = twse.MAVBR(3, 6)
	var sample2mav3 = twse.MAV(3)
	var sample2mav6 = twse.MAV(6)
	if sample2mavbr[len(sample2mavbr)-1] != int64(sample2mav3[len(sample2mav3)-1]-sample2mav6[len(sample2mav6)-1]) {
		t.Error("Should be the sample")
	}
}

func BenchmarkMA(b *testing.B) {
	twse.Get()
	for i := 0; i <= b.N; i++ {
		twse.MA(3)
	}
}

func BenchmarkMABR(b *testing.B) {
	twse.Get()
	for i := 0; i <= b.N; i++ {
		twse.MABR(3, 6)
	}
}

func BenchmarkMAV(b *testing.B) {
	twse.Get()
	for i := 0; i <= b.N; i++ {
		twse.MAV(3)
	}
}

func BenchmarkMAVBR(b *testing.B) {
	twse.Get()
	for i := 0; i <= b.N; i++ {
		twse.MAVBR(3, 6)
	}
}

func BenchmarkGet(b *testing.B) {
	var d = NewTWSE("2618", time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))
	for i := 0; i <= b.N; i++ {
		d.Get()
	}
}

func BenchmarkGetVolumeList(b *testing.B) {
	var d = NewTWSE("2618", time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
	d.Get()
	for i := 0; i <= b.N; i++ {
		d.GetVolumeList()
	}
}

func BenchmarkGetPriceList(b *testing.B) {
	var d = NewTWSE("2618", time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
	d.Get()
	for i := 0; i <= b.N; i++ {
		d.GetPriceList()
	}
}

// 新增一個 TWSE 上市股票
func Example_newTWSE() {
	var stock = NewTWSE("2618", time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
	stock.Get()
	fmt.Println(stock.RawData[0])
	// output:
	// [104/03/02 13,384,378 305,046,992 23.00 23.05 22.50 22.90 -0.10 3,793]
}

// 新增一個 OTC 上櫃股票
func Example_newOTC() {
	var stock = NewOTC("8446", time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
	stock.Get()
	fmt.Println(stock.RawData[0])
	// output:
	// [104/03/02 354 33,018 92.00 94.90 90.80 92.60 3.50 299]
}

func ExampleData_Get_notEnoughData() {
	year, month, _ := time.Now().Date()
	var d = NewTWSE("2618", time.Date(year, month+1, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))

	stockData, err := d.Get()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(stockData)
	}
	// output:
	// Not enough data.
}

func ExampleData_PlusData() {
	var stock = NewTWSE("2618", time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
	stock.Get() // 2015/3
	fmt.Println(stock.Date)
	stock.PlusData() // 2015/2
	fmt.Println(stock.Date)
	// output:
	// 2015-03-27 00:00:00 +0800 Asia/Taipei
	// 2015-02-01 00:00:00 +0800 Asia/Taipei
}

func ExampleData_Round() {
	var d = NewTWSE("2618", time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))

	fmt.Println(d.Date) // 2014/12

	d.Round()
	fmt.Println(d.Date) // 2014/11

	d.Round()
	fmt.Println(d.Date) // 2014/10
	// output:
	// 2014-12-26 00:00:00 +0800 Asia/Taipei
	// 2014-11-01 00:00:00 +0800 Asia/Taipei
	// 2014-10-01 00:00:00 +0800 Asia/Taipei
}

func TestData_Round(t *testing.T) {
	var now = time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone)
	var past = time.Date(2015, 2, 1, 0, 0, 0, 0, utils.TaipeiTimeZone)
	var d = NewTWSE("2618", now)

	t.Log(d.Date)
	if d.Date == past {
		t.Fatal(d.Date, past)
	}
	d.Round()
	t.Log(d.Date)
	if d.Date != past {
		t.Fatal(d.Date, past)
	}
}

func TestData_PlusData(t *testing.T) {
	var now = time.Date(2015, 3, 27, 0, 0, 0, 0, utils.TaipeiTimeZone)
	var d = NewTWSE("2618", now)
	d.PlusData()
	var d2 = NewTWSE("2618", time.Date(2015, 2, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))
	d2.Get()
	for i := range d2.RawData {
		if d.RawData[i][0] != d2.RawData[i][0] {
			t.Fatal("Data not difference.")
			t.Log(d.RawData, d2.RawData)
		}
	}
}

func TestData_GetByTimeMap(*testing.T) {
	var d = NewTWSE("2618", time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))
	d.GetByTimeMap()
}

func TestData_FormatData(*testing.T) {
	var d = NewTWSE("2618", time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))
	d.Get()
	d.FormatData()
}
