package tradingdays

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"time"

	"github.com/toomore/gogrs/utils"
)

// IsOpen is check open or not.
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
var timeLayout = "2006/1/02"

func readCSV() {
	data, _ := ioutil.ReadFile("./list.csv")
	csvdata := csv.NewReader(bytes.NewReader(data))

	for {
		record, err := csvdata.Read()
		if err == io.EOF {
			break
		}
		//fmt.Println(record[1])
		t, err := time.ParseInLocation(timeLayout, record[0], utils.TaipeiTimeZone)
		var isopen bool
		if record[1] == "1" {
			isopen = true
		}
		exceptDays[t.Unix()] = isopen
		//fmt.Println(t, err, t.Unix(), isopen)
		//fmt.Println(exceptDays)
	}
}

func init() {
	exceptDays = make(map[int64]bool)
	readCSV()
}
