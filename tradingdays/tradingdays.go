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
func IsOpen(d time.Time) bool {
	if d.Weekday() == 0 || d.Weekday() == 6 {
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
		isopen := false
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
