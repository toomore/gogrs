package twse

import (
	"testing"
	"time"

	"github.com/toomore/gogrs/utils"
)

func TestQFIIS_Get(t *testing.T) {
	qf := &QFIIS{Date: time.Date(2015, 5, 25, 0, 0, 0, 0, utils.TaipeiTimeZone)}
	t.Log(qf.URL())
	t.Log(qf.Get())
}
