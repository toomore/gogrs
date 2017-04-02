package twse

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"regexp"
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

var filterWeight = regexp.MustCompile(`<!------ start of 交易資訊 > 指數資料 >加權股價指數歷史資料 --->\n(.+)\n<!------ end of  交易資訊 > 指數資料 >加權股價指數歷史資料 --->`)

func getWeight(date time.Time) ([]byte, error) {
	data := make(url.Values)
	year := fmt.Sprintf("%d", date.Year()-1911)
	mon := fmt.Sprintf("%02d", date.Month())
	data.Add("myear", year)
	data.Add("mmon", mon)
	resp, _ := hCache.PostForm("http://www.tse.com.tw/ch/trading/indices/MI_5MINS_HIST/MI_5MINS_HIST.php", data)
	if filterWeight.Match(resp) {
		for _, result := range filterWeight.FindAllSubmatch(resp, -1) {
			data := make(url.Values)
			data.Add("html", base64.StdEncoding.EncodeToString(result[1]))
			data.Add("dirname", fmt.Sprintf("MI_5MINS_HIST%s%s", year, mon))
			return hCache.PostForm("http://www.tse.com.tw/ch/trading/indices/MI_5MINS_HIST/MI_5MINS_HIST_print.php?language=ch&save=csv", data)
		}
	}
	return nil, errorNotEnoughData
}

// WeightData struct
type WeightData struct {
	Date  time.Time
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func atoiMust(s string) float64 {
	result, _ := strconv.ParseFloat(strings.Replace(strings.Trim(s, " "), ",", "", -1), 64)
	return result
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
		result[num] = &WeightData{
			Date:  k,
			Open:  atoiMust(v[1]),
			High:  atoiMust(v[2]),
			Low:   atoiMust(v[3]),
			Close: atoiMust(v[4]),
		}
		num++
	}
	return result
}
