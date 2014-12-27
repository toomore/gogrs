package main

import (
	"fmt"
	"time"
)

// To get TWSE daily CSV files.
type DailyData struct {
	stock_no string
	date     time.Time
}

// To render url for csv.
func (d DailyData) Url() string {
	return fmt.Sprintf(TWSECSV, d.date.Year(), d.date.Month(), d.stock_no, RandInt())
}

// Sub one month.
func (d *DailyData) Round() {
	year, month, day := d.date.Date()
	d.date = time.Date(year, month-1, day, 0, 0, 0, 0, time.UTC)
}

func main() {
	d := &DailyData{
		stock_no: "2618",
		date:     time.Now(),
	}

	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
	d.Round()
	fmt.Println(d.Url())
}
