package main

import (
	"fmt"
	"time"
)

func DailyDataUrl(stock string) (url string) {
	t := time.Now()
	url = fmt.Sprintf(TWSECSV, t.Year(), t.Month(), stock, t.UnixNano())
	return
}

func main() {
	fmt.Println(DailyDataUrl("2618"))
}
