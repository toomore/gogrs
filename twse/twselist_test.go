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

func TestLists_Get_categoryNoList(t *testing.T) {
	l := &Lists{
		Date: time.Date(2015, 4, 27, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	l.Get("15") //航運業
	l.Get("01") //水泥業
	t.Log(l.FmtData["2618"])
	t.Log(l.FmtData)
	t.Log(l.categoryRawData)
	t.Log(l.categoryNoList)
	t.Log(l.GetCategoryList("15"))
	t.Log(l.Get("ALLBUT0999"))
	t.Log(l.GetCategoryList("ALLBUT0999"))
}

func ExampleLists_GetCategoryList() {
	l := &Lists{
		Date: time.Date(2015, 4, 27, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	categoryList := l.GetCategoryList("15")
	for _, v := range categoryList {
		if v.No == "2618" {
			fmt.Printf("%+v", v)
			break
		}
	}
	// output:
	// {No:2618 Name:長榮航}
}

func ExampleLists_Get_fmtData() {
	l := &Lists{
		Date: time.Date(2015, 4, 9, 0, 0, 0, 0, utils.TaipeiTimeZone),
	}
	l.Get("15") //航運業
	fmt.Printf("%+v", l.FmtData["2618"])
	// output:
	// {No:2618 Name:長榮航 Volume:46670950 TotalPrice:1136982254 Open:24 High:24.65 Low:24 Price:24 Range:0.55 Totalsale:11117 LastBuyPrice:24 LastBuyVolume:2027 LastSellPrice:24.1 LastSellVolume:10 PERatio:0}
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
