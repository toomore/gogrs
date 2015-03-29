package gogrs

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func assertType(t *testing.T, t1 interface{}, t2 interface{}) {
	if reflect.TypeOf(t1) != reflect.TypeOf(t2) {
		t.Errorf("Diff type t1(%s), t2(%s)", reflect.TypeOf(t1), reflect.TypeOf(t2))
	}
}

func TestURL(t *testing.T) {
	d := &DailyData{
		No:   "2618",
		Date: time.Now(),
	}
	assertType(t, d, &DailyData{})
}

func ExampleDailyData() {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(d)

	stockData, _ := d.GetData()
	fmt.Println(stockData)
}

func TestDailyData_GetData(*testing.T) {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(d)

	stockData, _ := d.GetData()
	d.GetData() // Test Cache.
	fmt.Println(stockData)
}

func BenchmarkGetData(b *testing.B) {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	for i := 0; i <= b.N; i++ {
		d.GetData()
	}
}

func BenchmarkGetVolumeList(b *testing.B) {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2015, 3, 27, 0, 0, 0, 0, time.Local),
	}
	d.GetData()
	for i := 0; i <= b.N; i++ {
		d.GetVolumeList()
	}
}

func ExampleDailyData_GetData_notEnoughData() {
	year, month, _ := time.Now().Date()
	d := DailyData{
		No:   "2618",
		Date: time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local),
	}

	stockData, err := d.GetData()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(stockData)
	}
}

func ExampleDailyData_Round() {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(d)

	stockData, _ := d.GetData()
	fmt.Println(stockData)

	fmt.Println(d.URL()) // 2014/12

	d.Round()
	fmt.Println(d.URL()) // 2014/11
	stockData, _ = d.GetData()
	fmt.Println(stockData)

	d.Round()
	fmt.Println(d.URL()) // 2014/10
	stockData, _ = d.GetData()
	fmt.Println(stockData)
}

func TestDailyData_Round(*testing.T) {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	d.Round()
	fmt.Println(d.URL())
}

func TestDailyData_GetDataByTimeMap(*testing.T) {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(d.GetDataByTimeMap())
}

func TestDailyData_FormatDailyData(*testing.T) {
	d := DailyData{
		No:   "2618",
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	d.GetData()
	fmt.Println(d.FormatDailyData())
}
