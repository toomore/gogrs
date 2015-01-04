package gogrs

import (
	"fmt"
	"time"
)

type TWSEList struct {
	Date time.Time
}

func (l TWSEList) baseURL() string {
	year, month, day := l.Date.Date()
	return fmt.Sprintf("%s%s", TWSEHOST, fmt.Sprintf(TWSELISTCSV, year, month, year, month, day))
}

func (l TWSEList) URL(strNo string) string {
	return fmt.Sprintf(l.baseURL, strNo)
}
