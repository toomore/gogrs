package main

import (
	"log"
	"time"

	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

type base interface {
	MA(days int) []float64
}

func check01(b base) bool {
	days, up := utils.CountCountineFloat64(b.MA(3))
	log.Println(days, up, b.MA(3))
	if up && days > 1 {
		return true
	}
	return false
}

func main() {
	stock := twse.NewTWSE("2329", time.Now())
	stock.Get()
	log.Println(check01(stock))
}
