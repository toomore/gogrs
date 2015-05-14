package main

import (
	"log"
	"time"

	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

type base interface {
	MA(days int) []float64
	Len() int
	PlusData()
}

// MA 3 > 6 > 18
func check01(b base) bool {
	var start = b.Len()
	for {
		if b.Len() >= 18 {
			break
		}
		log.Println("In PlusData")
		b.PlusData()
		if (b.Len() - start) == 0 {
			break
		}
	}
	var ma3 = b.MA(3)
	var ma6 = b.MA(6)
	var ma18 = b.MA(18)
	//log.Println(ma3[len(ma3)-1], ma6[len(ma6)-1], ma18[len(ma18)-1])
	if ma3[len(ma3)-1] > ma6[len(ma6)-1] && ma6[len(ma6)-1] > ma18[len(ma18)-1] {
		return true
	}
	return false
}

func check02(b base) bool {
	days, up := utils.CountCountineFloat64(utils.DeltaFloat64(b.MA(3)))
	log.Println(days, up, b.MA(3))
	if up && days > 1 {
		return true
	}
	return false
}

func main() {
	stock := twse.NewTWSE("2329", tradingdays.FindRecentlyOpened(time.Now()))
	log.Println(check01(base(stock)))
	log.Println(check02(base(stock)))
}
