package twse

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

type QFIIS struct {
	Date time.Time
}

func (q QFIIS) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, fmt.Sprintf(utils.QFIISTOP20, q.Date.Year(), q.Date.Month(), q.Date.Day()))
}

func (q *QFIIS) Get() ([][]string, error) {
	data, _ := hCache.Get(q.URL(), false)

	csvArrayContent := strings.Split(string(data), "\n")
	for i, v := range csvArrayContent[2 : len(csvArrayContent)-1] {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[:len(csvArrayContent)-3], "\n")))
	var (
		allData [][]string
		err     error
	)
	if allData, err = csvReader.ReadAll(); err == nil {
		for _, v := range allData {
			for i, vv := range v {
				v[i] = strings.Replace(vv, ",", "", -1)
			}
		}
	}
	return allData, err
}
