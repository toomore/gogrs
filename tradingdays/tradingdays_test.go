package tradingdays

import (
	"fmt"
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestReadCSV(*testing.T) {
	//readCSV()
	fmt.Println(exceptDays)
	fmt.Println(time.Date(2015, 9, 29, 0, 0, 0, 0, time.Local))
	fmt.Println(IsOpen(2015, 4, 17, time.Local))
	fmt.Println(IsOpen(2015, 4, 18, time.Local))
	fmt.Println(IsOpen(2015, 4, 20, time.Local))
	fmt.Println(IsOpen(2015, 4, 25, time.Local))
}

func BenchmarkIsOpen(t *testing.B) {
	for i := 0; i < t.N; i++ {
		IsOpen(2015, 4, 19, utils.TaipeiTimeZone)
	}
}
