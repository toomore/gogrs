package twse

import (
	"encoding/csv"
	"fmt"
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
func (q *QFIISTOP20) Get() ([][]string, error) {
	data, _ := hCache.Get(q.URL(), false)

	csvArrayContent := strings.Split(string(data), "\n")
	for i, v := range csvArrayContent[2:] {
		csvArrayContent[i] = strings.Replace(v, "=", "", -1)
	}

	csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[:len(csvArrayContent)-3], "\n")))
	return csvReader.ReadAll()
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
