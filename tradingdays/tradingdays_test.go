package tradingdays

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestIsOpen(t *testing.T) {
	fmt.Println(exceptDays)
	if IsOpen(2015, 4, 17, time.Local) != true {
		t.Error("Should be `true`")
	}
	if IsOpen(2015, 4, 18, time.Local) != false {
		t.Error("Should be `false`")
	}
	if IsOpen(2015, 4, 20, time.Local) != true {
		t.Error("Should be `true`")
	}
	if IsOpen(2015, 5, 1, time.Local) != false {
		t.Error("Should be `false`")
	}
}

func BenchmarkIsOpen(t *testing.B) {
	for i := 0; i < t.N; i++ {
		IsOpen(2015, 4, 19, utils.TaipeiTimeZone)
	}
}
