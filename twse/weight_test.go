package twse

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestWeight(t *testing.T) {
	case1 := Weight(time.Date(2017, 2, 22, 0, 0, 0, 0, utils.TaipeiTimeZone))
	if len(case1) == 0 {
		t.Error("2017/2/22 Get no data")
	}

	for i, v := range case1 {
		t.Logf("[%d] %+v", i, v)
	}

	if len(Weight(time.Date(2017, 3, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))) == 0 {
		t.Error("2017/3/5 Get no data")
	}

	now := time.Now()
	if len(Weight(time.Date(2017, now.Month()+1, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))) > 0 {
		t.Errorf("%s Get data", now)
	}
}

func ExampleWeight() {
	result := Weight(time.Date(2017, 3, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))
	for _, v := range result {
		if v.Date == time.Date(2017, 3, 1, 0, 0, 0, 0, utils.TaipeiTimeZone) {
			fmt.Printf("%+v", v)
			break
		}
	}
	// output:
	// &{Date:2017-03-01 00:00:00 +0800 Asia/Taipei Volume:4352913007 TotalPrice:99915928382 Totalsale:981872 Point:9674.78 Range:-75.69}
}

func Benchmark_solveWeightCSV(b *testing.B) {
	csvData, _ := getWeight(time.Date(2017, 2, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))
	for i := 0; i <= b.N; i++ {
		solveWeightCSV(csvData)
	}
}
