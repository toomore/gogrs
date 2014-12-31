/*
Package gogrs is try to get stock data from TWSE.

Main example.

	import (
		"fmt"
		"time"

		"github.com/toomore/gogrs"
	)

	func main() {
		d := gogrs.DailyData{No: "2618", Date: time.Now()}
		fmt.Println(d)

		r := gogrs.StockRealTime{
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
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
		d.Round()
		fmt.Println(d.URL())
	}
*/
package gogrs
