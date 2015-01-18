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
