package twse

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/utils"
)

func TestLists_Get_Rawdata(*testing.T) {
	l := NewLists(time.Date(2014, 12, 23, 0, 0, 0, 0, utils.TaipeiTimeZone))
	//listdata, err := l.Get("MS")
	//fmt.Println(l.categoryRawData, "\n\n", listdata, err)
	//l.FmtData
	l.Get("MS")
	l.Get("ms")
}

func TestLists_Get_categoryNoList(t *testing.T) {
	l := NewLists(time.Date(2015, 4, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
	l.Get("15") //航運業
	l.Get("01") //水泥業
	t.Log(l.FmtData["2618"])
	t.Log(l.FmtData)
	t.Log(l.categoryRawData)
	t.Log(l.categoryNoList)
	t.Log(l.GetCategoryList("15"))
	t.Log("ALLBUT0999:", len(l.GetCategoryList("ALLBUT0999")))
	t.Log("ALL:", len(l.GetCategoryList("ALL")))
	ll := NewLists(time.Date(2015, 5, 22, 0, 0, 0, 0, utils.TaipeiTimeZone))
	if len(ll.GetCategoryList("ALLBUT0999")) < 10 {
		t.Error("應該沒那麼少")
	}
	if len(ll.GetCategoryList("ALL")) < 10 {
		t.Error("應該沒那麼少")
	}
}

func TestOTCLists(t *testing.T) {
	o := NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))
	t.Log(o.GetCategoryList("04"))
	t.Log(o.Get("04"))
}

func TestCategoryList(t *testing.T) {
	categoryList := NewCategoryList()
	t.Log(categoryList.Same())
	t.Log(categoryList.OnlyTWSE())
	t.Log(categoryList.OnlyOTC())
}

func TestBaseLists(t *testing.T) {
	otc := NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))
	twselist := NewLists(tradingdays.FindRecentlyOpened(time.Now()))
	for _, v := range []BaseLists{otc, twselist} {
		if data, err := v.Get("02"); err == nil {
			t.Log(len(data))
		} else {
			t.Error(err)
		}
	}
}

func ExampleLists_GetCategoryList() {
	l := NewLists(time.Date(2015, 4, 27, 0, 0, 0, 0, utils.TaipeiTimeZone))
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
	l := NewLists(time.Date(2015, 4, 9, 0, 0, 0, 0, utils.TaipeiTimeZone))
	l.Get("15") //航運業
	fmt.Printf("%+v", l.FmtData["2618"])
	// output:
	// {No:2618 Name:長榮航 Volume:46670950 TotalPrice:1136982254 Open:24 High:24.65 Low:24 Price:24 Range:0.55 Totalsale:11117 LastBuyPrice:24 LastBuyVolume:2027 LastSellPrice:24.1 LastSellVolume:10 PERatio:0 IssuedShares:0}
}

func ExampleLists_Get() {
	l := NewLists(time.Date(2014, 12, 26, 0, 0, 0, 0, utils.TaipeiTimeZone))
	listdata, _ := l.Get("15") //航運業
	fmt.Println(listdata[0])
	// output:
	// [2208 台船 729,340 324 12,048,156 16.45 16.60 16.45 16.45   0.00 16.45 67 16.50 58 41.13 ]
}

func ExampleLists_Get_notEnoughData() {
	year, month, day := time.Now().Date()
	l := NewLists(time.Date(year, month+1, day, 0, 0, 0, 0, utils.TaipeiTimeZone))
	_, err := l.Get("15") //航運業
	fmt.Println(err)
	// output:
	// Not enough data
}
