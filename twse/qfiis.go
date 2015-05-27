package twse

import (
	"encoding/csv"
	"fmt"
	"net/url"
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
func (q QFIISTOP20) Get() ([][]string, error) {
	data, _ := hCache.Get(q.URL(), false)

	csvArrayContent := strings.Split(string(data), "\n")[2:]
	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	return csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll()
}

// BFI82U 取得「三大法人買賣金額統計表」
type BFI82U struct {
	Begin time.Time
	End   time.Time
}

// URL 擷取網址
func (b BFI82U) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.BFI82U,
			b.Begin.Year(), b.Begin.Month(), b.Begin.Day(),
			b.End.Year(), b.End.Month(), b.End.Day(),
		))
}

// Get 擷取資料
func (b BFI82U) Get() ([][]string, error) {
	data, _ := hCache.Get(b.URL(), false)
	return csv.NewReader(strings.NewReader(strings.Join(strings.Split(string(data), "\n")[2:], "\n"))).ReadAll()
}

// T86 取得「三大法人買賣超日報(股)」
type T86 struct {
	Date time.Time
}

// URL 擷取網址
func (t T86) URL(cate string) string {
	var ym = fmt.Sprintf("%d%02d", t.Date.Year(), t.Date.Month())
	var ymd = fmt.Sprintf("%s%d", ym, t.Date.Day())
	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.T86,
			ym, ymd, cate, ymd))
}

// Get 擷取資料
func (t T86) Get(cate string) ([][]string, error) {
	data, _ := hCache.Get(t.URL(cate), false)
	csvArrayContent := strings.Split(string(data), "\n")[2:]
	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}
	return csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll()
}

// TWTXXU 產生 自營商、投信、外資及陸資買賣超彙總表
type TWTXXU struct {
	Date time.Time
	fund string
}

// NewTWT43U 自營商買賣超彙總表
func NewTWT43U(date time.Time) *TWTXXU {
	return &TWTXXU{Date: date, fund: "TWT43U"}
}

// NewTWT44U 投信買賣超彙總表
func NewTWT44U(date time.Time) *TWTXXU {
	return &TWTXXU{Date: date, fund: "TWT44U"}
}

// NewTWT38U 外資及陸資買賣超彙總表
func NewTWT38U(date time.Time) *TWTXXU {
	return &TWTXXU{Date: date, fund: "TWT38U"}
}

// URL 擷取網址
func (t TWTXXU) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST,
		fmt.Sprintf(utils.TWTXXU,
			t.fund, t.fund, t.Date.Year()-1911, t.Date.Month(), t.Date.Day()))
}

// Get 擷取資料
func (t TWTXXU) Get() ([][]string, error) {
	var (
		r      *url.URL
		err    error
		result [][]string
	)

	path := t.URL()

	if r, err = url.Parse(path); err == nil {
		rawquery, _ := url.ParseQuery(r.RawQuery)
		data, _ := hCache.PostForm(path, rawquery)

		var csvArrayContent = strings.Split(string(data), "\n")
		switch t.fund {
		case "TWT43U":
			csvArrayContent = csvArrayContent[3 : len(csvArrayContent)-4]
		case "TWT44U", "TWT38U":
			csvArrayContent = csvArrayContent[2 : len(csvArrayContent)-5]
		}

		for i, v := range csvArrayContent {
			csvArrayContent[i] = strings.Replace(v, "=", "", -1)
		}

		result, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll()
	}
	return result, err
}
