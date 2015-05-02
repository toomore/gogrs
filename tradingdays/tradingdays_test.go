package tradingdays

import (
	"fmt"
	"testing"
	"time"
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
