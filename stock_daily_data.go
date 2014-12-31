// Package gogrs is try to get stock data from TWSE.
package gogrs

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// DailyData start with stock no, date.
type DailyData struct {
	No   string
	Date time.Time
}

// URL return stock csv url path.
func (d DailyData) URL() string {
	return fmt.Sprintf(TWSECSV, d.Date.Year(), d.Date.Month(), d.Date.Year(), d.Date.Month(), d.No, RandInt())
}

// Round will do sub one month.
func (d *DailyData) Round() {
	year, month, day := d.Date.Date()
	d.Date = time.Date(year, month-1, day, 0, 0, 0, 0, time.UTC)
}

// GetData return csv data in array.
func (d DailyData) GetData() ([][]string, error) {
	urlpath := fmt.Sprintf("%s%s", TWSEHOST, d.URL())
	csvFiles, err := http.Get(urlpath)
	if err != nil {
		fmt.Println("[err] >>> ", err)
		return nil, err
	}
	defer csvFiles.Body.Close()
	data, _ := ioutil.ReadAll(csvFiles.Body)
	csvArrayContent := strings.Split(string(data), "\n")
	for i := range csvArrayContent {
		csvArrayContent[i] = strings.TrimSpace(csvArrayContent[i])
	}
	csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[2:], "\n")))
	return csvReader.ReadAll()
}

//func main() {
//	d := &DailyData{
//		No: "2618",
//		Date:     time.Now(),
//	}
//
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//	d.Round()
//	fmt.Println(d.URL())
//}
