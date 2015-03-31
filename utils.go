package gogrs

import (
	"crypto/rand"
	"math/big"
	//"math/rand"
	"regexp"
	"strconv"
	"time"
)

// RandInt return random int.
func RandInt() int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(102400))
	return result.Int64()
}

//func RandInt() int64 {
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	return r.Int63n(100000)
//}

var dateReg = regexp.MustCompile(`[\d]{2,}`)

// ParseDate is to parse "104/01/13" format.
func ParseDate(strDate string) time.Time {
	p := dateReg.FindAllString(strDate, -1)
	year, _ := strconv.Atoi(p[0])
	mon, _ := strconv.Atoi(p[1])
	day, _ := strconv.Atoi(p[2])
	return time.Date(year+1911, time.Month(mon), day, 0, 0, 0, 0, time.Local)
}

func SumFloat64(data []float64) float64 {
	var result float64
	for _, v := range data {
		result += v
	}
	return result
}

func AvgFlast64(data []float64) float64 {
	return float64(int(SumFloat64(data)*100)/len(data)) / 100
}

func SumUint64(data []uint64) uint64 {
	var result uint64
	for _, v := range data {
		result += v
	}
	return result
}

func AvgUint64(data []uint64) uint64 {
	return SumUint64(data) / uint64(len(data))
}

func ThanPastUint64(data []uint64, days int, max bool) bool {
	var dataFloat64 []float64
	dataFloat64 = make([]float64, days+1)
	for i, v := range data[len(data)-1-days:] {
		dataFloat64[i] = float64(v)
	}
	return thanPast(dataFloat64, max)
}

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

	var base float64 = data[len(data)-1]
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
