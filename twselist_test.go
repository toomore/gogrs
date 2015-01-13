package gogrs

import (
	"fmt"
	"testing"
	"time"
)

func ExampleTWSEList_URL() {
	l := TWSEList{
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(l.URL("15"))
	// Output: http://www.twse.com.tw/ch/trading/exchange/MI_INDEX/MI_INDEX2_print.php?genpage=genpage/Report201412/A1122014122615.php&type=csv
}
func TestTWSEList_GetData_notEnoughData(*testing.T) {
	year, month, day := time.Now().Date()
	l := TWSEList{
		Date: time.Date(year, month+1, day, 0, 0, 0, 0, time.Local),
	}
	listdata, err := l.GetData("15") //航運業
	fmt.Println(listdata, err)
}

func TestTWSEList_GetData(*testing.T) {
	l := TWSEList{
		//Date: time.Now(),
		Date: time.Date(2014, 12, 26, 0, 0, 0, 0, time.Local),
	}
	listdata, err := l.GetData("15") //航運業
	fmt.Println(listdata, err)
}
