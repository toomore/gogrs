package gogrs

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strconv"
	"time"
)

// RandInt return random int.
func RandInt() int64 {
	result, _ := rand.Int(rand.Reader, big.NewInt(102400))
	return result.Int64()
}

var dateReg = regexp.MustCompile(`[\d]{2,}`)

// ParseDate is to parse "104/01/13" format.
func ParseDate(strDate string) time.Time {
	p := dateReg.FindAllString(strDate, -1)
	year, _ := strconv.Atoi(p[0])
	mon, _ := strconv.Atoi(p[1])
	day, _ := strconv.Atoi(p[2])
	return time.Date(year+1911, time.Month(mon), day, 0, 0, 0, 0, time.Local)
}
