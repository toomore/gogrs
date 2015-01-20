package gogrs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
			fmt.Println("[err] >>> ", err)
			return nil, err
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

func (d DailyData) GetDataMap() map[time.Time]interface{} {
	data := make(map[time.Time]interface{})
	dailyData, _ := d.GetData()
	reg := regexp.MustCompile(`[\d]{2,}`)
	for _, v := range dailyData {
		p := reg.FindAllString(v[0], -1)
		year, _ := strconv.Atoi(p[0])
		mon, _ := strconv.Atoi(p[1])
		day, _ := strconv.Atoi(p[2])
		data[time.Date(year+1911, time.Month(mon), day, 0, 0, 0, 0, time.Local)] = v
	}
	return data
}

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

var dateReg = regexp.MustCompile(`[\d]{2,}`)

func parseDate(strDate string) time.Time {
	p := dateReg.FindAllString(strDate, -1)
	year, _ := strconv.Atoi(p[0])
	mon, _ := strconv.Atoi(p[1])
	day, _ := strconv.Atoi(p[2])
	return time.Date(year+1911, time.Month(mon), day, 0, 0, 0, 0, time.Local)
}

func (d DailyData) FormatDailyData() {
	for _, v := range d.RawData {
		var d FmtDailyData
		volume, _ := strconv.ParseUint(strings.Replace(v[1], ",", "", -1), 10, 32)
		d.Volume = volume

		totalprice, _ := strconv.ParseUint(strings.Replace(v[2], ",", "", -1), 10, 32)
		d.TotalPrice = totalprice

		open, _ := strconv.ParseFloat(v[3], 64)
		d.Open = open

		high, _ := strconv.ParseFloat(v[4], 64)
		d.High = high

		low, _ := strconv.ParseFloat(v[5], 64)
		d.Low = low

		price, _ := strconv.ParseFloat(v[6], 64)
		d.Price = price

		range_, _ := strconv.ParseFloat(v[7], 64)
		d.Range = range_

		totalsale, _ := strconv.ParseUint(strings.Replace(v[8], ",", "", -1), 10, 64)
		d.Totalsale = totalsale

		fmt.Println(d)
	}
}
