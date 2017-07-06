// Package utils - 套件所需的公用工具（總和、平均、序列差、持續天數、
// 民國日期解析、簡單亂數、標準差、簡單 net/http 快取）.
//
package utils

import (
	"math"
	"regexp"
	"strconv"
	"time"
)

// ExchangeMap is simple to check tse or otc.
var ExchangeMap = map[string]bool{"tse": true, "otc": true}

// TaipeiTimeZone is for time.Data() setting.
var TaipeiTimeZone = time.FixedZone("Asia/Taipei", 8*3600)

// TWSE base url.
const (
	TWSEURL     string = "http://mis.twse.com.tw"
	TWSEHOST    string = "http://www.twse.com.tw"
	OTCHOST     string = "http://www.tpex.org.tw"
	HOME        string = "/stock/index.jsp"
	OTCCSV      string = "/ch/stock/aftertrading/daily_trading_info/st43_download.php?d=%d/%02d&stkno=%s&r=%%d"           // year, mon, stock, rand
	OTCLISTCSV  string = "/web/stock/aftertrading/otc_quotes_no1430/stk_wn1430_download.php?l=zh-tw&d=%s&se=%s&s=0,asc,0" // date, cate
	TWSECSV     string = "/exchangeReport/STOCK_DAY?response=csv&date=%d%02d%02d&stockNo=%s"
	TWSELISTCSV string = "/exchangeReport/MI_INDEX?response=csv&date=%d%02d%02d&type=%s" // year, mon, day, type
	TWSEREAL    string = "/stock/api/getStockInfo.jsp?ex_ch=%s_%s.tw&json=1&delay=0&_=%d"
	QFIISTOP20  string = "/fund/MI_QFIIS_sort_20?response=csv&date=%d%02d%02d"   // yyyymmdd
	BFI82U      string = "/fund/BFI82U?response=csv&dayDate=%d%02d%02d&type=day" // yyyymmdd
	T86         string = "/fund/T86?response=csv&date=%d%02d%02d&selectType=ALL"
	TWTXXU      string = "/fund/%s?response=csv&date=%d%02d%02d"
	S3CSV       string = "https://s3-ap-northeast-1.amazonaws.com/toomore/gogrs/list.csv"
)

// RandInt return random int.
func RandInt() int {
	return time.Now().Nanosecond()
}

//func RandInt() int64 {
//	result, _ := rand.Int(rand.Reader, big.NewInt(102400))
//	return result.Int64()
//}

//func RandInt() int64 {
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	return r.Int63n(100000)
//}

var dateReg = regexp.MustCompile(`([\d]{2,})/([\d]{1,2})/([\d]{1,2})`)

// ParseDate is to parse "104/01/13" format.
func ParseDate(strDate string) time.Time {
	p := dateReg.FindStringSubmatch(strDate)
	if len(p) == 0 {
		return time.Time{}
	}
	year, _ := strconv.Atoi(p[1])
	mon, _ := strconv.Atoi(p[2])
	day, _ := strconv.Atoi(p[3])
	return time.Date(year+1911, time.Month(mon), day, 0, 0, 0, 0, TaipeiTimeZone)
}

// SumFloat64 計算總和（float64）
func SumFloat64(data []float64) float64 {
	var result float64
	for _, v := range data {
		result += v
	}
	return result
}

// AvgFloat64 計算平均（float64）
func AvgFloat64(data []float64) float64 {
	return float64(int(SumFloat64(data)*100)/len(data)) / 100
}

// SumUint64 計算總和（uint64）
func SumUint64(data []uint64) uint64 {
	var result uint64
	for _, v := range data {
		result += v
	}
	return result
}

// AvgUint64 計算平均（uint64）
func AvgUint64(data []uint64) uint64 {
	return SumUint64(data) / uint64(len(data))
}

// ThanPastUint64 計算最後一個數值是否為過去幾天最大或最小（uint64）
func ThanPastUint64(data []uint64, days int, max bool) bool {
	var dataFloat64 = make([]float64, days+1)
	for i, v := range data[len(data)-1-days:] {
		dataFloat64[i] = float64(v)
	}
	return thanPast(dataFloat64, max)
}

// ThanPastFloat64 計算最後一個數值是否為過去幾天最大或最小（float64）
func ThanPastFloat64(data []float64, days int, max bool) bool {
	return thanPast(data[len(data)-1-days:], max)
}

func thanPast(data []float64, max bool) bool {
	//var dataFloat64 []float64
	//dataFloat64 = make([]float64, days+1)
	//for i, v := range data[len(data)-1-days : len(data)-1] {
	//	switch v.(type) {
	//	case int64:
	//		dataFloat64[i] = float64(v.(int64))
	//	case float64:
	//		dataFloat64[i] = v.(float64)
	//	}
	//}

	var base = data[len(data)-1]
	var condition func(b float64) bool

	if max {
		condition = func(b float64) bool { return base > b }
	} else {
		condition = func(b float64) bool { return base < b }
	}

	for _, v := range data[:len(data)-1] {
		if condition(v) {
			continue
		} else {
			return false
		}
	}
	return true
}

// ThanSumPastFloat64 計算最後一個數值是否為過去幾天的總和大或小（float64）
func ThanSumPastFloat64(data []float64, days int, max bool) bool {
	return thanSumPast(data[len(data)-1-days:], max)
}

// ThanSumPastUint64 計算最後一個數值是否為過去幾天的總和大或小（uint64）
func ThanSumPastUint64(data []uint64, days int, max bool) bool {
	var dataFloat64 = make([]float64, days+1)
	for i, v := range data[len(data)-1-days:] {
		dataFloat64[i] = float64(v)
	}
	return thanSumPast(dataFloat64, max)
}

func thanSumPast(data []float64, max bool) bool {
	var result = data[len(data)-1] > SumFloat64(data[:len(data)-1])
	if max {
		return result
	}
	return !result
}

// CountCountineFloat64 計算最後一個數值為正或負值的持續天數（float64）
func CountCountineFloat64(data []float64) (int, bool) {
	var condition func(x float64) bool
	if data[len(data)-1] > 0 {
		condition = func(x float64) bool { return x > 0 }
	} else {
		condition = func(x float64) bool { return x < 0 }
	}
	counter := 1
	for i := 1; i <= len(data)-1; i++ {
		if condition(data[len(data)-1-i]) {
			counter++
		} else {
			break
		}
	}
	return counter, data[len(data)-1] > 0
}

// CalDiffFloat64 計算兩序列的差（listA - listB）
func CalDiffFloat64(listA, listB []float64) []float64 {
	var length int
	var result []float64
	if len(listA) <= len(listB) {
		length = len(listA)
	} else {
		length = len(listB)
	}
	result = make([]float64, length)
	for i := 1; i <= length; i++ {
		result[length-i] = listA[len(listA)-i] - listB[len(listB)-i]
	}
	return result
}

// CalDiffInt64 計算兩序列的差（listA - listB）
func CalDiffInt64(listA, listB []int64) []int64 {
	var length int
	var result []int64
	if len(listA) <= len(listB) {
		length = len(listA)
	} else {
		length = len(listB)
	}
	result = make([]int64, length)
	for i := 1; i <= length; i++ {
		result[length-i] = listA[len(listA)-i] - listB[len(listB)-i]
	}
	return result
}

// CalLHPower 計算兩序列與終點的比例
func CalLHPower(final, low, high []float64) []float64 {
	var result []float64
	result = make([]float64, len(final))
	for i := range result {
		result[i] = (final[i] - low[i]) / (high[i] - low[i])
	}
	return result
}

// SDUint64 計算標準差
func SDUint64(list []uint64) float64 {
	var data = make([]float64, len(list))
	for i := 0; i < len(list); i++ {
		data[i] = float64(list[i])
	}
	return SD(data)
}

// SD 計算標準差
func SD(list []float64) float64 {
	var avg = AvgFloat64(list)
	var data = make([]float64, len(list))
	for i := 0; i < len(data); i++ {
		data[i] = list[i] - avg
		data[i] *= data[i]
	}
	return math.Sqrt(AvgFloat64(data))
}

// DeltaFloat64 計算數列間距差
func DeltaFloat64(data []float64) []float64 {
	var result = make([]float64, len(data)-1)
	for i := 0; i < (len(data) - 1); i++ {
		result[i] = data[i+1] - data[i]
	}
	return result
}

// DeltaInt64 計算數列間距差
func DeltaInt64(data []int64) []int64 {
	var result = make([]int64, len(data)-1)
	for i := 0; i < (len(data) - 1); i++ {
		result[i] = data[i+1] - data[i]
	}
	return result
}
