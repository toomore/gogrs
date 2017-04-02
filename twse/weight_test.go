package twse

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestWeight(t *testing.T) {
	if len(Weight(time.Date(2017, 2, 22, 0, 0, 0, 0, utils.TaipeiTimeZone))) == 0 {
		t.Fatal("2017/2/22 Get no data")
	}

	if len(Weight(time.Date(2017, 3, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))) == 0 {
		t.Fatal("2017/3/5 Get no data")
	}

	now := time.Now()
	if len(Weight(time.Date(2017, now.Month()+1, 1, 0, 0, 0, 0, utils.TaipeiTimeZone))) > 0 {
		t.Fatalf("%s Get data", now)
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
	// &{Date:2017-03-01 00:00:00 +0800 Asia/Taipei Open:9751.12 High:9758.78 Low:9674.78 Close:9674.78}
}
