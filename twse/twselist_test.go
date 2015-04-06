package twse

import (
	"fmt"
	"testing"
	"time"
)

func ExampleLists_URL() {
	l := &Lists{
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(l.URL("15"))
	// Output: http://www.twse.com.tw/ch/trading/exchange/MI_INDEX/MI_INDEX2_print.php?genpage=genpage/Report201412/A1122014122615.php&type=csv
}

func TestLists_Get_notEnoughData(*testing.T) {
	year, month, day := time.Now().Date()
	l := &Lists{
		Date: time.Date(year, month+1, day, 0, 0, 0, 0, time.Local),
	}
	listdata, err := l.Get("15") //航運業
	fmt.Println(listdata, err)
}

func TestLists_Get(*testing.T) {
	l := &Lists{
		//Date: time.Now(),
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	listdata, err := l.Get("15") //航運業
	fmt.Println(listdata, err)
}
