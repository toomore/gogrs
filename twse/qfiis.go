package twse

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

// BaseSellBuy 買進賣出合計
type BaseSellBuy struct {
	No    string
	Name  string
	Buy   int64 // 買進
	Sell  int64 // 賣出
	Total int64 // 合計
}

// QFIISTOP20 取得「外資及陸資持股比率前二十名彙總表」
type QFIISTOP20 struct {
	Date time.Time
}

// URL 擷取網址
func (q QFIISTOP20) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, utils.QFIISTOP20)
}

// Get 擷取資料
func (q QFIISTOP20) Get() ([][]string, error) {
	var (
		err  error
		data []byte
	)
	if data, err = hCache.PostForm(q.URL(), url.Values{"download": {"csv"},
		"qdate": {fmt.Sprintf("%d/%02d/%02d", q.Date.Year()-1911, q.Date.Month(), q.Date.Day())}}); err == nil {
		if len(data) > 0 {
			csvArrayContent := strings.Split(string(data), "\n")[2:]
			for i, v := range csvArrayContent {
				csvArrayContent[i] = strings.Replace(v, "=", "", -1)
			}
			return csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll()
		}
	}
	return nil, err
}

// BFI82U 取得「三大法人買賣金額統計表」
type BFI82U struct {
	Begin time.Time
	End   time.Time
}

// NewBFI82U 三大法人買賣金額統計表
func NewBFI82U(begin, end time.Time) *BFI82U {
	return &BFI82U{Begin: begin, End: end}
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
func (b BFI82U) Get() ([]BaseSellBuy, error) {
	var (
		csvdata [][]string
		data    []byte
		err     error
		result  []BaseSellBuy
	)
	if data, err = hCache.Get(b.URL(), false); err != nil {
		return nil, err
	}
	if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(strings.Split(string(data), "\n")[2:], "\n"))).ReadAll(); err == nil {
		result = make([]BaseSellBuy, len(csvdata))
		for i, v := range csvdata {
			result[i].Name = v[0]
			result[i].Buy, _ = strconv.ParseInt(strings.Replace(v[1], ",", "", -1), 10, 64)
			result[i].Sell, _ = strconv.ParseInt(strings.Replace(v[2], ",", "", -1), 10, 64)
			result[i].Total, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
		}
	}
	return result, err
}

// T86 取得「三大法人買賣超日報(股)」
type T86 struct {
	Date time.Time
}

// URL 擷取網址
func (t T86) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, utils.T86)
}

// T86Data 各欄位資料
type T86Data struct {
	No     string
	Name   string
	FII    BaseSellBuy // 外資
	SIT    BaseSellBuy // 投信
	DProp  BaseSellBuy // 自營商(自行買賣)
	DHedge BaseSellBuy // 自營商(避險)
	Diff   int64       // 三大法人買賣超股數
}

// Get 擷取資料
func (t T86) Get(cate string) ([]T86Data, error) {
	data, err := hCache.PostForm(t.URL(), url.Values{
		"download": {"csv"},
		"qdate":    {fmt.Sprintf("%d/%02d/%02d", t.Date.Year()-1911, t.Date.Month(), t.Date.Day())},
		"select2":  {cate},
		"sorting":  {"by_issue"}})
	if err != nil {
		return nil, err
	}
	csvArrayContent := strings.Split(string(data), "\n")[2:]
	for i, v := range csvArrayContent {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}
	var result []T86Data
	if csvdata, err := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[:len(csvArrayContent)-8], "\n"))).ReadAll(); err == nil {
		result = make([]T86Data, len(csvdata))
		for i, v := range csvdata {
			if len(v) >= 11 {
				result[i].No = v[0]
				result[i].Name = v[1]

				result[i].FII.Buy, _ = strconv.ParseInt(strings.Replace(v[2], ",", "", -1), 10, 64)
				result[i].FII.Sell, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
				result[i].FII.Total = result[i].FII.Buy - result[i].FII.Sell

				result[i].SIT.Buy, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
				result[i].SIT.Sell, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)
				result[i].SIT.Total = result[i].SIT.Buy - result[i].SIT.Sell

				result[i].DProp.Buy, _ = strconv.ParseInt(strings.Replace(v[6], ",", "", -1), 10, 64)
				result[i].DProp.Sell, _ = strconv.ParseInt(strings.Replace(v[7], ",", "", -1), 10, 64)
				result[i].DProp.Total = result[i].DProp.Buy - result[i].DProp.Sell

				result[i].DHedge.Buy, _ = strconv.ParseInt(strings.Replace(v[8], ",", "", -1), 10, 64)
				result[i].DHedge.Sell, _ = strconv.ParseInt(strings.Replace(v[9], ",", "", -1), 10, 64)
				result[i].DHedge.Total = result[i].DHedge.Buy - result[i].DHedge.Sell

				result[i].Diff, _ = strconv.ParseInt(strings.Replace(v[10], ",", "", -1), 10, 64)
			}
		}
	} else {
		return nil, err
	}
	return result, err
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
func (t TWTXXU) Get() ([][]BaseSellBuy, error) {
	var (
		csvdata  [][]string
		data     []byte
		datalist int
		err      error
		path     = t.URL()
		r        *url.URL
		result   [][]BaseSellBuy
	)

	if r, err = url.Parse(path); err == nil {
		rawquery, _ := url.ParseQuery(r.RawQuery)
		if data, err = hCache.PostForm(path, rawquery); err != nil {
			return nil, err
		}

		var csvArrayContent = strings.Split(string(data), "\n")
		switch t.fund {
		case "TWT43U":
			csvArrayContent = csvArrayContent[3 : len(csvArrayContent)-4]
			datalist = 3
		case "TWT44U", "TWT38U":
			csvArrayContent = csvArrayContent[2 : len(csvArrayContent)-6]
			datalist = 1
		}

		for i, v := range csvArrayContent {
			csvArrayContent[i] = strings.Replace(v, "=", "", -1)
		}

		if csvdata, err = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent, "\n"))).ReadAll(); err == nil {
			result = make([][]BaseSellBuy, len(csvdata))
			for i, v := range csvdata {
				result[i] = make([]BaseSellBuy, datalist)
				result[i][0].Name = v[0]
				result[i][0].No = v[1]
				result[i][0].Buy, _ = strconv.ParseInt(strings.Replace(v[2], ",", "", -1), 10, 64)
				result[i][0].Sell, _ = strconv.ParseInt(strings.Replace(v[3], ",", "", -1), 10, 64)
				result[i][0].Total, _ = strconv.ParseInt(strings.Replace(v[4], ",", "", -1), 10, 64)
				if datalist > 1 {
					result[i][1].Name = v[0]
					result[i][1].No = v[1]
					result[i][1].Buy, _ = strconv.ParseInt(strings.Replace(v[5], ",", "", -1), 10, 64)
					result[i][1].Sell, _ = strconv.ParseInt(strings.Replace(v[6], ",", "", -1), 10, 64)
					result[i][1].Total, _ = strconv.ParseInt(strings.Replace(v[7], ",", "", -1), 10, 64)
					result[i][2].Name = v[0]
					result[i][2].No = v[1]
					result[i][2].Buy, _ = strconv.ParseInt(strings.Replace(v[8], ",", "", -1), 10, 64)
					result[i][2].Sell, _ = strconv.ParseInt(strings.Replace(v[9], ",", "", -1), 10, 64)
					result[i][2].Total, _ = strconv.ParseInt(strings.Replace(v[10], ",", "", -1), 10, 64)
				}
			}
		}
	}
	return result, err
}
