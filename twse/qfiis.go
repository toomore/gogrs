package twse

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

// QFIISTOP20 取得「外資及陸資持股比率前二十名彙總表」
type QFIISTOP20 struct {
	Date time.Time
}

// URL 擷取網址
func (q QFIISTOP20) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, fmt.Sprintf(utils.QFIISTOP20, q.Date.Year(), q.Date.Month(), q.Date.Day()))
}

// Get 擷取資料
func (q *QFIISTOP20) Get() ([][]string, error) {
	data, _ := hCache.Get(q.URL(), false)

	csvArrayContent := strings.Split(string(data), "\n")
	for i, v := range csvArrayContent[2 : len(csvArrayContent)-1] {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[:len(csvArrayContent)-3], "\n")))
	return csvReader.ReadAll()
}
