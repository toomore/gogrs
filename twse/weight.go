package twse

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

// Weight is TWSE Weight
func Weight(date time.Time) []*WeightData {
	if csvData, err := getWeight(date); err == nil {
		return solveWeightCSV(csvData)
	}
	return nil
}

func getWeight(date time.Time) ([]byte, error) {
	return hCache.PostForm(fmt.Sprintf("%s%s",
		utils.TWSEHOST, fmt.Sprintf("/exchangeReport/FMTQIK?response=csv&date=%d%02d%02d", date.Year(), date.Month(), date.Day())), nil)
}

// WeightData struct
type WeightData struct {
	Date       time.Time
	Volume     uint64  // 成交股數
	TotalPrice uint64  //成交金額
	Totalsale  uint64  //成交筆數
	Point      float64 //收盤價
	Range      float64 //漲跌價差
}

func solveWeightCSV(raw []byte) []*WeightData {
	splitData := strings.Split(string(raw), "\n")
	csvReader := csv.NewReader(strings.NewReader(strings.Join(splitData[2:], "\n")))
	dataMapping := make(map[time.Time][]string)

	for {
		row, err := csvReader.Read()
		switch err {
		case nil:
			if len(row) >= 5 {
				date := utils.ParseDate(row[0])
				if date.IsZero() == false {
					if _, dup := dataMapping[date]; dup {
						goto Final
					}
					dataMapping[date] = row
				}
			}
		case io.EOF:
			goto Final
		}
	}
Final:
	result := make([]*WeightData, len(dataMapping))
	num := 0
	for k, v := range dataMapping {
		result[num] = &WeightData{Date: k}
		result[num].Volume, _ = strconv.ParseUint(strings.Replace(v[1], ",", "", -1), 10, 64)
		result[num].TotalPrice, _ = strconv.ParseUint(strings.Replace(v[2], ",", "", -1), 10, 64)
		result[num].Totalsale, _ = strconv.ParseUint(strings.Replace(v[3], ",", "", -1), 10, 64)
		result[num].Point, _ = strconv.ParseFloat(strings.Replace(v[4], ",", "", -1), 64)
		result[num].Range, _ = strconv.ParseFloat(v[5], 64)
		num++
	}
	return result
}
