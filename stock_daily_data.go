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
			return nil, errors.New(fmt.Sprintf("Network fail: %s", err))
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

// FmtDailyData is struct for daily data format.
type FmtDailyData struct {
	Date       time.Time
	Volume     uint64
	TotalPrice uint64
	Open       float64
	High       float64
	Low        float64
	Price      float64
	Range      float64
	Totalsale  uint64
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
