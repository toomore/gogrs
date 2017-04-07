// Package tradingdays - Check the day is open or not and at which period time
// 股市開休市判斷（支援非國定假日：颱風假）與當日區間判斷（盤中、盤後、盤後盤）
//
package tradingdays

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/toomore/gogrs/utils"
)

// IsOpen 判斷是否為開休日
func IsOpen(year int, month time.Month, day int) bool {
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if openornot, ok := exceptDays[d.Unix()]; ok {
		return openornot
	}
	if d.In(utils.TaipeiTimeZone).Weekday() == 0 || d.In(utils.TaipeiTimeZone).Weekday() == 6 {
		return false
	}
	return true
}

// FindRecentlyOpened 回傳最近一個開市時間（UTC 0）
func FindRecentlyOpened(date time.Time) time.Time {
	var (
		d     = date.In(utils.TaipeiTimeZone)
		days  = d.Day()
		index int
		tp    *TimePeriod
	)

	for {
		if IsOpen(d.Year(), d.Month(), days) {
			if index == 0 {
				tp = NewTimePeriod(time.Date(d.Year(), d.Month(), days, d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), utils.TaipeiTimeZone))
				if tp.AtBefore() || tp.AtOpen() {
					days--
					for {
						if IsOpen(d.Year(), d.Month(), days) {
							break
						}
						days--
					}
				}
			}
			return time.Date(d.Year(), d.Month(), days, 0, 0, 0, 0, time.UTC)
		}
		days--
		index++
	}

}

var (
	csvEtag    string
	exceptDays map[int64]bool
)

func readCSV() {
	_, filename, _, _ := runtime.Caller(1)
	data, _ := ioutil.ReadFile(path.Join(path.Dir(filename), "list.csv"))
	processCSV(bytes.NewReader(data))
}

// DownloadCSV 更新開休市表
//
// 從 https://s3-ap-northeast-1.amazonaws.com/toomore/gogrs/list.csv
// 下載表更新，主要發生在非國定假日，如：颱風假。
func DownloadCSV(replace bool) {
	req, _ := http.NewRequest("GET", utils.S3CSV, nil)
	req.Header.Set("If-None-Match", csvEtag)
	resp, _ := utils.HTTPClient.Do(req)
	if resp.StatusCode != http.StatusNotModified && resp.StatusCode == http.StatusOK {
		csvEtag = resp.Header.Get("Etag")
		updateCSVData(resp, replace)
	}
}

func updateCSVData(resp *http.Response, replace bool) {
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		if replace {
			exceptDays = make(map[int64]bool)
		}
		processCSV(bytes.NewReader(data))
	}
}

func processCSV(data io.Reader) {
	csvdata := csv.NewReader(data)

	for {
		record, err := csvdata.Read()
		if err == io.EOF {
			break
		}
		if t, err := time.ParseInLocation("2006/1/2", record[0], time.UTC); err == nil {
			var isopen bool
			if record[1] == "1" {
				isopen = true
			}
			exceptDays[t.Unix()] = isopen
		}
	}
}

// TimePeriod 簡單計算時間是否在當日的開盤前、盤中、盤後盤、收盤的時間
/*
相關時間：
	- 開盤前	00:00 ~ 09:00
	- 盤中	09:00 ~ 13:30
	- 盤後盤	13:30 ~ 14:30
	- 收盤	14:30 ~ 24:00

證交所會在 13:30 後產生一次報表，14:30 盤後盤零股交易後還會再產生一次報表。
*/
type TimePeriod struct {
	target      int64
	zero        int64 // 00:00
	start       int64 // 09:00
	firstclose  int64 // 13:30
	secondclose int64 // 14:40
	night       int64 // 24:00
}

// NewTimePeriod 建立一個時間區間判斷
func NewTimePeriod(date time.Time) *TimePeriod {
	var d = &lazyTime{date: date.In(utils.TaipeiTimeZone)}
	return &TimePeriod{
		target:      d.date.Unix(),
		zero:        d.time(0, 0).Unix(),
		start:       d.time(9, 0).Unix(),
		firstclose:  d.time(13, 30).Unix(),
		secondclose: d.time(14, 30).Unix(),
		night:       d.time(24, 0).Unix(),
	}
}

// AtBefore 檢查是否為 開盤前（0000-0900） 期間
func (t TimePeriod) AtBefore() bool {
	if t.target >= t.zero && t.target < t.start {
		return true
	}

	return false
}

// AtOpen 檢查是否為 開盤（0900-1330） 期間
func (t TimePeriod) AtOpen() bool {
	if t.target >= t.start && t.target < t.firstclose {
		return true
	}
	return false
}

// AtAfterOpen 檢查是否為 盤後盤（1330-1430） 期間
func (t TimePeriod) AtAfterOpen() bool {
	if t.target >= t.firstclose && t.target < t.secondclose {
		return true
	}
	return false
}

// AtClose 檢查是否為 收盤（1430-2400）期間
func (t TimePeriod) AtClose() bool {
	if t.target >= t.secondclose && t.target < t.night {
		return true
	}
	return false
}

type lazyTime struct {
	date time.Time
}

func (t lazyTime) time(hour, min int) time.Time {
	return time.Date(t.date.Year(), t.date.Month(), t.date.Day(), hour, min, 0, 0, t.date.Location())
}

func init() {
	exceptDays = make(map[int64]bool)
	readCSV()
}
