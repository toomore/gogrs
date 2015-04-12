package twse

import (
	"fmt"
	"time"

	"github.com/toomore/gogrs/utils"
)

func ExampleLists_Get() {
	l := &Lists{
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	listdata, _ := l.Get("15") //航運業
	fmt.Println(listdata[0])
	// output:
	// [2208 台船 729,340 324 12,048,156 16.45 16.60 16.45 16.45  0.00 16.45 67 16.50 58 41.13]
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

func ExampleLists_URL() {
	l := &Lists{
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	fmt.Println(l.URL("15"))
	// output: http://www.twse.com.tw/ch/trading/exchange/MI_INDEX/MI_INDEX2_print.php?genpage=genpage/Report201412/A1122014122615.php&type=csv
}
