package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (stock StockOption) GetData() (value stockBlob) {
	url := stock.GenStockUrl()
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&value)
	return
}

type jsonBlob struct {
	Args    map[string]interface{}
	Headers map[string]interface{}
}

type QueryTimeBlob struct {
	sysTime           string
	sessionLatestTime int
	sysDate           string
}

type stockBlob struct {
	Rtcode    string
	UserDelay int
	Rtmessage string
	Referer   string
	MsgArray  []map[string]string
	QueryTime map[string]interface{}
}

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
