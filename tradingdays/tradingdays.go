package tradingdays

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/toomore/gogrs/utils"
)

// IsOpen 判斷是否為開休日
func IsOpen(year int, month time.Month, day int, loc *time.Location) bool {
	d := time.Date(year, month, day, 0, 0, 0, 0, loc)
	if openornot, ok := exceptDays[d.Unix()]; ok {
		return openornot
	}
	if d.In(utils.TaipeiTimeZone).Weekday() == 0 || d.In(utils.TaipeiTimeZone).Weekday() == 6 {
		return false
	}
	return true
}

var exceptDays map[int64]bool
var timeLayout = "2006/1/2"

func readCSV() {
	data, _ := ioutil.ReadFile("./list.csv")
	processCSV(bytes.NewReader(data))
}

// DownloadCSV 更新開休市表
func DownloadCSV(replace bool) {
	resp, _ := http.Get("https://s3-ap-northeast-1.amazonaws.com/toomore/gogrs/list.csv")
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		if replace {
			exceptDays = make(map[int64]bool)
		}
		processCSV(bytes.NewReader(data))
	}
}

func processCSV(data io.Reader) {
	csvdata := csv.NewReader(data)

	for {
		record, err := csvdata.Read()
		if err == io.EOF {
			break
		}
		if t, err := time.ParseInLocation(timeLayout, record[0], utils.TaipeiTimeZone); err == nil {
			var isopen bool
			if record[1] == "1" {
				isopen = true
			}
			exceptDays[t.Unix()] = isopen
		}
	}
}

func init() {
	exceptDays = make(map[int64]bool)
	readCSV()
}
