package filter

import (
	"testing"
	"time"

	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

var stocklist = []*twse.Data{
	twse.NewTWSE("1201", time.Date(2016, 1, 8, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("1201", time.Date(2016, 2, 14, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("1201", time.Date(2017, 3, 21, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("1201", time.Date(2017, 4, 12, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2412", time.Date(2015, 7, 21, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2412", time.Date(2015, 9, 21, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2618", time.Date(2015, 3, 8, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2618", time.Date(2016, 1, 8, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2618", time.Date(2016, 2, 14, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2618", time.Date(2017, 3, 21, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("2618", time.Date(2017, 4, 7, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	twse.NewTWSE("4938", time.Date(2016, 3, 21, 0, 0, 0, 0, utils.TaipeiTimeZone)),
	//twse.NewTWSE("4938", time.Date(2017, time.Now().Month()+1, 21, 0, 0, 0, 0, utils.TaipeiTimeZone)),
}

func TestCheckGroup(t *testing.T) {
	for _, stock := range stocklist {
		for i, v := range AllList {
			t.Log(i, v.No(), v.String(), v.Mindata(), v.CheckFunc(stock))
		}
	}
}

func TestCheckGroup_String(t *testing.T) {
	for i, v := range AllList {
		t.Log(i, v.No(), v)
	}
}
