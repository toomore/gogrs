// To get TWSE daily CSV files.
package gogrs

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func (d DailyData) GetData() ([][]string, error) {
	urlpath := fmt.Sprintf("%s%s", TWSEHOST, d.Url())
	csv_data, err := http.Get(urlpath)
	if err != nil {
		fmt.Println("[err] >>> ", err)
		return nil, err
	} else {
		defer csv_data.Body.Close()
		data, _ := ioutil.ReadAll(csv_data.Body)
		data_content := strings.Split(string(data), "\n")
		for i, _ := range data_content {
			data_content[i] = strings.TrimSpace(data_content[i])
		}
		csv_reader := csv.NewReader(strings.NewReader(strings.Join(data_content[2:], "\n")))
		return csv_reader.ReadAll()
	}
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
