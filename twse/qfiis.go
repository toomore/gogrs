package twse

import (
	"fmt"
	"time"

	"github.com/toomore/gogrs/utils"
)

type QFIIS struct {
	Date time.Time
}

func (q QFIIS) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEHOST, fmt.Sprintf(utils.QFIISTOP20, q.Date.Year(), q.Date.Month(), q.Date.Day()))
}
