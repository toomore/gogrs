package tradingdays

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func ExampleIsOpen() {
	fmt.Println(IsOpen(2015, 4, 17))
	fmt.Println(IsOpen(2015, 5, 1))
	// output:
	// true
	// false
}

func TestFindRecentlyOpened(t *testing.T) {
	result := FindRecentlyOpened(time.Now())
	if time.Now().UTC().Sub(result).Hours() > 24*7 {
		t.Error("Should not more than 7 days")
	}

	var date = time.Date(2015, 5, 11, 8, 30, 0, 0, utils.TaipeiTimeZone)
	result2 := FindRecentlyOpened(date)
	if result2.Unix() != time.Date(2015, 5, 8, 0, 0, 0, 0, time.UTC).Unix() {
		t.Error("Should be at 2015/5/8")
	}

	date = time.Date(2015, 5, 10, 8, 30, 0, 0, utils.TaipeiTimeZone)
	result2 = FindRecentlyOpened(date)
	if result2.Unix() != time.Date(2015, 5, 8, 0, 0, 0, 0, time.UTC).Unix() {
		t.Error("Should be at 2015/5/8")
	}

	date0601 := time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
	date0602 := time.Date(2015, 6, 2, 0, 0, 0, 0, time.UTC).Unix()
	test06011 := FindRecentlyOpened(time.Date(2015, 6, 2, 7, 0, 0, 0, utils.TaipeiTimeZone)).Unix()
	test06012 := FindRecentlyOpened(time.Date(2015, 6, 2, 8, 0, 0, 0, utils.TaipeiTimeZone)).Unix()
	test06013 := FindRecentlyOpened(time.Date(2015, 6, 2, 9, 0, 0, 0, utils.TaipeiTimeZone)).Unix()
	test06014 := FindRecentlyOpened(time.Date(2015, 6, 2, 10, 0, 0, 0, utils.TaipeiTimeZone)).Unix()
	test06021 := FindRecentlyOpened(time.Date(2015, 6, 2, 14, 30, 0, 0, utils.TaipeiTimeZone)).Unix()
	if !(date0601 == test06011 && date0601 == test06012 && date0601 == test06013 && date0601 == test06014) {
		t.Error("Should be at 2015/6/1")
	}
	if !(date0602 == test06021) {
		t.Error("Should be at 2015/6/2")
	}
}

func TestIsOpen(t *testing.T) {
	DownloadCSV(true)
	if IsOpen(2015, 4, 17) != true {
		t.Error("Should be `true`")
	}
	if IsOpen(2015, 4, 18) != false {
		t.Error("Should be `false`")
	}
	if IsOpen(2015, 4, 20) != true {
		t.Error("Should be `true`")
	}
	if IsOpen(2015, 5, 1) != false {
		t.Error("Should be `false`")
	}
}

func TestTimePerid(t *testing.T) {
	var tp1 = NewTimePeriod(time.Date(2015, 5, 8, 0, 0, 0, 0, utils.TaipeiTimeZone))
	var tp2 = NewTimePeriod(time.Date(2015, 5, 8, 10, 0, 0, 0, utils.TaipeiTimeZone))
	var tp3 = NewTimePeriod(time.Date(2015, 5, 8, 14, 0, 0, 0, utils.TaipeiTimeZone))
	var tp4 = NewTimePeriod(time.Date(2015, 5, 8, 20, 0, 0, 0, utils.TaipeiTimeZone))

	if tp1.AtBefore() != true {
		t.Error("Should be `true`")
	}

	if tp1.AtOpen() != false {
		t.Error("Should be `false`")
	}

	if tp1.AtAfterOpen() != false {
		t.Error("Should be `false`")
	}

	if tp1.AtClose() != false {
		t.Error("Should be `false`")
	}

	if tp2.AtOpen() != true {
		t.Error("Should be `true`")
	}

	if tp2.AtBefore() != false {
		t.Error("Should be `false`")
	}

	if tp3.AtAfterOpen() != true {
		t.Error("Should be `true`")
	}

	if tp4.AtClose() != true {
		t.Error("Should be `true`")
	}
}

func TestDownloadCSV(*testing.T) {
	DownloadCSV(true)
	DownloadCSV(true)
	DownloadCSV(true)
	DownloadCSV(true)
	DownloadCSV(true)
	DownloadCSV(true)
}

func BenchmarkTimePeriod(b *testing.B) {
	var tp = NewTimePeriod(time.Date(2015, 5, 8, 20, 0, 0, 0, utils.TaipeiTimeZone))
	for i := 0; i < b.N; i++ {
		tp.AtBefore()
		tp.AtOpen()
		tp.AtAfterOpen()
		tp.AtClose()
	}
}

func ExampleNewTimePeriod() {
	var tp = NewTimePeriod(time.Date(2015, 5, 8, 20, 0, 0, 0, utils.TaipeiTimeZone))
	fmt.Println(tp.AtBefore(), tp.AtOpen(), tp.AtAfterOpen(), tp.AtClose())
	// output:
	// false false false true
}

func BenchmarkIsOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsOpen(2015, 4, 19)
	}
}

func BenchmarkDownloadCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DownloadCSV(true)
	}
}

func BenchmarkFindRecentlyOpened(b *testing.B) {
	var now = time.Now()
	for i := 0; i < b.N; i++ {
		FindRecentlyOpened(now)
	}
}

func init() {
	log.Println("Testing init")
	DownloadCSV(true)
}
