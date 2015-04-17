package tradingdays

import "time"

func IsOpen(d time.Time) bool {
	if d.Weekday() == 0 || d.Weekday() == 6 {
		return false
	}
	return true
}
