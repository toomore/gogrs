package filter

import (
	"testing"
	"time"

	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
)

var stock = twse.NewTWSE("2618", tradingdays.FindRecentlyOpened(time.Now()))

func TestCheckGroup(t *testing.T) {
	for i, v := range AllList {
		t.Log(i, v.No(), v.String(), v.Mindata(), v.CheckFunc(stock))
	}
}

func TestCheckGroup_String(t *testing.T) {
	for i, v := range AllList {
		t.Log(i, v.No(), v)
	}
}
