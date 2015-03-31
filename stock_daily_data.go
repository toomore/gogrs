package gogrs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type memData map[int64][][]string

// DailyData start with stock no, date.
type DailyData struct {
	No      string
	Date    time.Time
	RawData [][]string
	hasData memData
}

// URL return stock csv url path.
func (d DailyData) URL() string {
	path := fmt.Sprintf(TWSECSV, d.Date.Year(), d.Date.Month(), d.Date.Year(), d.Date.Month(), d.No, RandInt())
	return fmt.Sprintf("%s%s", TWSEHOST, path)
}

// Round will do sub one month.
func (d *DailyData) Round() {
	year, month, _ := d.Date.Date()
	d.Date = time.Date(year, month-1, 1, 0, 0, 0, 0, time.UTC)
}

// GetData return csv data in array.
func (d *DailyData) GetData() ([][]string, error) {
	if d.hasData == nil {
		d.hasData = make(memData)
	}
	if len(d.hasData[d.Date.Unix()]) == 0 {
		csvFiles, err := http.Get(d.URL())
		if err != nil {
			return nil, fmt.Errorf("Network fail: %s", err)
		}
		defer csvFiles.Body.Close()
		data, _ := ioutil.ReadAll(csvFiles.Body)
		csvArrayContent := strings.Split(string(data), "\n")
		for i := range csvArrayContent {
			csvArrayContent[i] = strings.TrimSpace(csvArrayContent[i])
		}
		if len(csvArrayContent) > 2 {
			csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[2:], "\n")))
			allData, err := csvReader.ReadAll()
			d.RawData = append(allData, d.RawData...)
			d.hasData[d.Date.Unix()] = allData
			return allData, err
		}
		return nil, errors.New("Not enough data.")
	}
	return d.hasData[d.Date.Unix()], nil
}

// GetDataByTimeMap return a map by key of time.Time
func (d DailyData) GetDataByTimeMap() map[time.Time]interface{} {
	data := make(map[time.Time]interface{})
	dailyData, _ := d.GetData()
	for _, v := range dailyData {
		data[ParseDate(v[0])] = v
	}
	return data
}

func (d DailyData) getColsList(colsNo int) []string {
	var result []string
	result = make([]string, len(d.RawData))
	for i, value := range d.RawData {
		result[i] = value[colsNo]
	}
	return result
}

func (d DailyData) getColsListFloat64(colsNo int) []float64 {
	var result []float64
	result = make([]float64, len(d.RawData))
	for i, v := range d.getColsList(colsNo) {
		result[i], _ = strconv.ParseFloat(v, 64)
	}
	return result
}

func (d DailyData) GetVolumeList() []uint64 {
	var result []uint64
	result = make([]uint64, len(d.RawData))
	for i, v := range d.getColsList(1) {
		result[i], _ = strconv.ParseUint(strings.Replace(v, ",", "", -1), 10, 64)
	}
	return result
}

func (d DailyData) GetOpenList() []float64 {
	return d.getColsListFloat64(3)
}

func (d DailyData) GetPriceList() []float64 {
	return d.getColsListFloat64(6)
}

func (d DailyData) GetRangeList() []float64 {
	return d.getColsListFloat64(7)
}

func (d DailyData) MA(days int) []float64 {
	var result []float64
	var priceList = d.GetPriceList()
	result = make([]float64, len(priceList)-days+1)
	for i, _ := range priceList[days-1:] {
		result[i] = AvgFlast64(priceList[i : i+days])
	}
	return result
}

func (d DailyData) MAV(days int) []uint64 {
	var result []uint64
	var volumeList = d.GetVolumeList()
	result = make([]uint64, len(volumeList)-days+1)
	for i, _ := range volumeList[days-1:] {
		result[i] = AvgUint64(volumeList[i : i+days])
	}
	return result
}

func (d DailyData) IsRed() bool {
	var rangeList []float64 = d.GetRangeList()
	return rangeList[len(rangeList)-1] > 0
}

// FmtDailyData is struct for daily data format.
type FmtDailyData struct {
	Date       time.Time
	Volume     uint64  //成交股數
	TotalPrice uint64  //成交金額
	Open       float64 //開盤價
	High       float64 //最高價
	Low        float64 //最低價
	Price      float64 //收盤價
	Range      float64 //漲跌價差
	Totalsale  uint64  //成交筆數
}

// FormatDailyData is format daily data.
func (d DailyData) FormatDailyData() []FmtDailyData {
	result := make([]FmtDailyData, len(d.RawData))
	var loopd FmtDailyData
	for i, v := range d.RawData {
		loopd.Date = ParseDate(v[0])

		volume, _ := strconv.ParseUint(strings.Replace(v[1], ",", "", -1), 10, 32)
		loopd.Volume = volume

		totalprice, _ := strconv.ParseUint(strings.Replace(v[2], ",", "", -1), 10, 32)
		loopd.TotalPrice = totalprice

		open, _ := strconv.ParseFloat(v[3], 64)
		loopd.Open = open

		high, _ := strconv.ParseFloat(v[4], 64)
		loopd.High = high

		low, _ := strconv.ParseFloat(v[5], 64)
		loopd.Low = low

		price, _ := strconv.ParseFloat(v[6], 64)
		loopd.Price = price

		rangeData, _ := strconv.ParseFloat(v[7], 64)
		loopd.Range = rangeData

		totalsale, _ := strconv.ParseUint(strings.Replace(v[8], ",", "", -1), 10, 64)
		loopd.Totalsale = totalsale

		result[i] = loopd
	}
	return result
}
