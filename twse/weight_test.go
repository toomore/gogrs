package twse

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestWeight(t *testing.T) {
	if len(Weight(time.Date(2017, 2, 22, 0, 0, 0, 0, utils.TaipeiTimeZone))) == 0 {
		t.Error("2017/2/22 Get no data")
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
	// &{Date:2017-03-01 00:00:00 +0800 Asia/Taipei Open:9751.12 High:9758.78 Low:9674.78 Close:9674.78}
}

func TestWeightVolume(t *testing.T) {
	date := time.Date(2017, 2, 5, 0, 0, 0, 0, utils.TaipeiTimeZone)
	bytedata, nums := getWeightVolume(date)
	if nums == 0 {
		t.Error("Get no data")
	} else {
		t.Log(nums)
		var result = make([]*WeightVolumeData, nums)
		solveWeightVolumeCSV(bytedata, result)

		var result2 = make([]*WeightVolumeData, nums+5)
		solveWeightVolumeCSV(bytedata, result2)
		t.Log(result2)

		var result3 = make([]*WeightVolumeData, nums-2)
		solveWeightVolumeCSV(bytedata, result3)
		t.Log(result3)

		var result4 = make([]*WeightVolumeData, 5)
		solveWeightVolumeCSV([]byte("\n\"1\",\"2\",\"3\"\n\"1\",\"2\"\n\"1\""), result4)
		t.Log(result4)
	}

	if len(WeightVolume(date)) != 18 {
		t.Error("Data nums fail")
	}

	now := time.Now()
	if len(WeightVolume(time.Date(2017, now.Month()+1, 5, 0, 0, 0, 0, utils.TaipeiTimeZone))) > 0 {
		t.Error("Data nums fail")
	}
}

func ExampleWeightVolume() {
	for _, v := range WeightVolume(time.Date(2017, 2, 5, 0, 0, 0, 0, utils.TaipeiTimeZone)) {
		if v.Date == time.Date(2017, 2, 2, 0, 0, 0, 0, utils.TaipeiTimeZone) {
			fmt.Printf("%+v", v)
			break
		}
	}
	// output:
	// &{Date:2017-02-02 00:00:00 +0800 Asia/Taipei Volume:5119534614 TotalPrice:147930853630 Totalsale:1131546 Points:9428.97 Diff:-18.98}
}
