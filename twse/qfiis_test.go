package twse

import (
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestQFIISTOP20_Get(t *testing.T) {
	qf := &QFIISTOP20{Date: time.Date(2015, 5, 25, 0, 0, 0, 0, utils.TaipeiTimeZone)}
	t.Log(qf.URL())
	t.Log(qf.Get())
}

func TestBFI82U_Get(t *testing.T) {
	bfi := NewBFI82U(
		time.Date(2015, 5, 25, 0, 0, 0, 0, utils.TaipeiTimeZone),
		time.Date(2015, 5, 26, 0, 0, 0, 0, utils.TaipeiTimeZone),
	)
	t.Log(bfi.URL())
	data, _ := bfi.Get()
	t.Logf("%+v", data)
}

func TestT86_Get(t *testing.T) {
	t86 := &T86{Date: time.Date(2015, 5, 25, 0, 0, 0, 0, utils.TaipeiTimeZone)}
	t.Log(t86.URL())
	if data, err := t86.Get("ALLBUT0999"); err == nil {
		for i, v := range data[:5] {
			t.Logf("%d %+v", i, v)
		}
	} else {
		t.Error("Error", err)
	}
}

func TestTWTXXU_Get(t *testing.T) {
	date := time.Date(2015, 5, 26, 0, 0, 0, 0, utils.TaipeiTimeZone)
	for _, v := range []*TWTXXU{NewTWT38U(date), NewTWT44U(date), NewTWT43U(date)} {
		t.Log(v.URL())
		if data, err := v.Get(); err == nil {
			t.Log(len(data), err)
			t.Logf("%+v\n", data[len(data)-1])
		} else {
			t.Error(err)
		}
	}
}
