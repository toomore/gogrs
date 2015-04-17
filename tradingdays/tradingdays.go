package tradingdays

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

func IsOpen(d time.Time) bool {
	if d.Weekday() == 0 || d.Weekday() == 6 {
		return false
	}
	return true
}

func readCSV() {
	data, _ := ioutil.ReadFile("./list.csv")
	csvdata := csv.NewReader(bytes.NewReader(data))
	for {
		record, err := csvdata.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(record[1])
	}
}
