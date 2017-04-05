package twse

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
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

// WeightVolumeData is struct
type WeightVolumeData struct {
	Date       time.Time //日期
	Volume     uint64    //成交股數
	TotalPrice uint64    //成交金額
	Totalsale  uint64    //成交筆數
	Points     float64   //發行量加權股價指數
	Diff       float64   //漲跌點數
}

// WeightVolume is TWSE Weight Volume
func WeightVolume(date time.Time) []*WeightVolumeData {
	byteData, nums := getWeightVolume(date)
	if nums > 0 {
		result := make([]*WeightVolumeData, nums)
		solveWeightVolumeCSV(byteData, result)
		return result
	}
	return nil
}

func getWeightVolume(date time.Time) ([]byte, int) {
	data := make(url.Values)
	data.Add("download", "csv")
	data.Add("query_year", fmt.Sprintf("%d", date.Year()))
	data.Add("query_month", fmt.Sprintf("%d", date.Month()))
	resp, _ := hCache.PostForm("http://www.tse.com.tw/ch/trading/exchange/FMTQIK/FMTQIK.php", data)
	dataSplitbytes := bytes.Split(resp, []byte("\n"))
	if len(dataSplitbytes) > 6 {
		return bytes.Join(dataSplitbytes[2:len(dataSplitbytes)-4], []byte("\n")),
			len(dataSplitbytes) - 6
	}
	return nil, 0
}

func solveWeightVolumeCSV(raw []byte, result []*WeightVolumeData) {
	csvReader := csv.NewReader(bytes.NewReader(raw))
	for i := range result {
		row, err := csvReader.Read()
		switch err {
		case nil:
			result[i] = &WeightVolumeData{}
			result[i].Date = utils.ParseDate(row[0])
			result[i].Volume, _ = strconv.ParseUint(strings.Replace(row[1], ",", "", -1), 10, 64)
			result[i].TotalPrice, _ = strconv.ParseUint(strings.Replace(row[2], ",", "", -1), 10, 64)
			result[i].Totalsale, _ = strconv.ParseUint(strings.Replace(row[3], ",", "", -1), 10, 64)
			result[i].Points, _ = strconv.ParseFloat(strings.Replace(row[4], ",", "", -1), 64)
			result[i].Diff, _ = strconv.ParseFloat(strings.Replace(row[5], ",", "", -1), 64)
		case io.EOF:
			return
		default:
			log.Println(err)
			return
		}
	}
}
