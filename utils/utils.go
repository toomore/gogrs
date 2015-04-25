// Package utils - 套件所需的公用工具
//
package utils

import (
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
	TWSEURL     string = "http://mis.tse.com.tw"
	TWSEHOST    string = "http://www.twse.com.tw"
	OTCHOST     string = "http://www.tpex.org.tw"
	OTCCSV      string = "/ch/stock/aftertrading/daily_trading_info/st43_download.php?d=%d/%02d&stkno=%s&r=%d" // year, mon, stock, rand
	TWSECSV     string = "/ch/trading/exchange/STOCK_DAY/STOCK_DAY_print.php?genpage=genpage/Report%d%02d/%d%02d_F3_1_8_%s.php&type=csv&r=%d"
	TWSELISTCSV string = "/ch/trading/exchange/MI_INDEX/MI_INDEX.php" // year, mon, day, type
	TWSEREAL    string = "/stock/api/getStockInfo.jsp?ex_ch=%s_%s.tw_%s&json=1&delay=0&_=%d"
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

// AvgFlast64 計算平均（float64）
func AvgFlast64(data []float64) float64 {
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

func thanSumPast(data []float64, max bool) bool {
	var result = data[len(data)-1] > SumFloat64(data[:len(data)-1])
	if max {
		return result
	}
	return !result
}
