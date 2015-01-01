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
	d := DailyData{No: "2618", Date: time.Now()}
	fmt.Println(d)

	r := StockRealTime{
		No:        "2618",
		Timestamp: time.Now().Unix(),
		Date:      time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}

	data := r.GetData()
	fmt.Printf("%v", data)

	fmt.Println("----- Test -----\n\n")
	//fmt.Println(d.GetData())
	stock_data, _ := d.GetData()
	fmt.Println(stock_data)

	fmt.Println(d.URL())

	d.Round()
	fmt.Println(d.URL())
	stock_data, _ = d.GetData()
	fmt.Println(stock_data)

	d.Round()
	fmt.Println(d.URL())
	stock_data, _ = d.GetData()
	fmt.Println(stock_data)
}

func ExampleStockRealTime() {
	r := StockRealTime{
		No:        "2618",
		Timestamp: time.Now().Unix(),
		Date:      time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}

	data := r.GetData()
	fmt.Printf("%v", data)
}
