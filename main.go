package main

import (
	"fmt"
	"time"
)

const TWSEURL string = "http://mis.tse.com.tw/"

//STOCKPATH = '/stock/api/getStockInfo.jsp?ex_ch=%(exchange)s_%(no)s.tw_%(date)s&json=1&delay=%(delay)s&_=%(timestamp)s'

type StockOption struct {
	no        string
	timestamp int64
	date      time.Time
}

func (stock StockOption) GenStockUrl() string {
	return fmt.Sprintf(
		"%sstock/api/getStockInfo.jsp?ex_ch=%s_%s.tw_%s&json=1&delay=0&_=%d",
		TWSEURL,
		"tse",
		stock.no,
		fmt.Sprintf(
			"%d%02d%02d",
			stock.date.Year(),
			int(stock.date.Month()),
			stock.date.Day(),
		),
		stock.timestamp,
	)
}

type jsonBlob struct {
	Args    map[string]interface{}
	Headers map[string]interface{}
}

func main() {
	//url := "http://httpbin.org/get?name=Toomore"
	//resp, _ := http.Get(url)
	//defer resp.Body.Close()
	//resp_data, _ := ioutil.ReadAll(resp.Body)
	//var json_blob jsonBlob
	//json.Unmarshal(resp_data, &json_blob)
	//fmt.Printf("%+v", json_blob)
	//fmt.Println(json_blob.Args["name"], TWSEURL)
	option := StockOption{
		no:        "2618",
		timestamp: time.Now().Unix(),
		date:      time.Date(2014, time.August, 11, 0, 0, 0, 0, time.Local),
	}
	fmt.Println(option)
	fmt.Println(option.GenStockUrl())
}
