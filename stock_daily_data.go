// To get TWSE daily CSV files.
package gogrs

import (
	"fmt"
	"time"
)

// To get TWSE daily CSV files.
type DailyData struct {
	No   string
	Date time.Time
}

// To render url for csv.
func (d DailyData) Url() string {
	return fmt.Sprintf(TWSECSV, d.Date.Year(), d.Date.Month(), d.Date.Year(), d.Date.Month(), d.No, RandInt())
}

// Sub one month.
func (d *DailyData) Round() {
	year, month, day := d.Date.Date()
	d.Date = time.Date(year, month-1, day, 0, 0, 0, 0, time.UTC)
}

func (d DailyData) GetData() string {
	urlpath := fmt.Sprintf("%s%s", TWSEHOST, d.Url())
	return urlpath
}

//func main() {
//	d := &DailyData{
//		No: "2618",
//		Date:     time.Now(),
//	}
//
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//	d.Round()
//	fmt.Println(d.Url())
//}
