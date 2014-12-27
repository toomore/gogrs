package gogrs

import (
	"fmt"
	"time"
)

func main() {
	stock := StockOption{
		no:        "2618",
		timestamp: time.Now().Unix(),
		date:      time.Date(2014, time.August, 11, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(stock.GenStockUrl())
	store_data := stock.GetData()
	fmt.Printf("%+v\n", store_data)
	fmt.Printf("%+v\n", store_data.QueryTime)
	fmt.Printf("%s\n", store_data.QueryTime["sysTime"])
}
