package tradingdays

import (
	"fmt"
	"testing"
)

func ExampleIsOpen() {
	fmt.Println(IsOpen(2015, 4, 17))
	fmt.Println(IsOpen(2015, 5, 1))
	// output:
	// true
	// false
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

func BenchmarkIsOpen(t *testing.B) {
	for i := 0; i < t.N; i++ {
		IsOpen(2015, 4, 19)
	}
}

func BenchmarkDownloadCSV(t *testing.B) {
	for i := 0; i < t.N; i++ {
		DownloadCSV(true)
	}
}
