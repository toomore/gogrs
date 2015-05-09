package tradingdays

import (
	"fmt"
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
	result := FindRecentlyOpened()
	if time.Now().UTC().Sub(result).Hours() > 24*7 {
		t.Error("Should not more than 7 days")
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

func BenchmarkTimePeriod(b *testing.B) {
	var tp = NewTimePeriod(time.Date(2015, 5, 8, 20, 0, 0, 0, utils.TaipeiTimeZone))
	for i := 0; i < b.N; i++ {
		tp.AtBefore()
		tp.AtOpen()
		tp.AtAfterOpen()
		tp.AtClose()
	}
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
	for i := 0; i < b.N; i++ {
		FindRecentlyOpened()
	}
}

func init() {
	DownloadCSV(true)
}
