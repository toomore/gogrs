package twse

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestLists_Get_Rawdata(*testing.T) {
	l := &Lists{
		Date: time.Date(2014, 12, 23, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	//listdata, err := l.Get("MS")
	//fmt.Println(l.categoryRawData, "\n\n", listdata, err)
	//l.FmtData
	l.Get("MS")
	l.Get("ms")
}

func ExampleLists_Get_fmtData() {
	l := &Lists{
		Date: time.Date(2015, 4, 9, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	l.Get("15") //航運業
	fmt.Println(l.FmtData["2618"])
	// output:
	// {2618 長榮航 46670950 1136982254 24 24.65 24 24 0.55 11117 24 2027 24.1 10 0}
}

func ExampleLists_Get() {
	l := &Lists{
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	listdata, _ := l.Get("15") //航運業
	fmt.Println(listdata[0])
	// output:
	// [2208   台船   729340 324 12048156 16.45 16.6 16.45 16.45   0 16.45 67 16.5 58 41.13]
}

func ExampleLists_Get_notEnoughData() {
	year, month, day := time.Now().Date()
	l := &Lists{
		Date: time.Date(year, month+1, day, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	_, err := l.Get("15") //航運業
	fmt.Println(err)
	// output:
	// Not enough data.
}
