package gogrs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//STOCKPATH = '/stock/api/getStockInfo.jsp?ex_ch=%(exchange)s_%(no)s.tw_%(date)s&json=1&delay=%(delay)s&_=%(timestamp)s'

type StockOption struct {
	No        string
	Timestamp int64
	Date      time.Time
}

type stockBlob struct {
	Rtcode    string
	UserDelay int
	Rtmessage string
	Referer   string
	MsgArray  []map[string]string
	QueryTime map[string]interface{}
}

func (stock StockOption) GenStockUrl() string {
	return fmt.Sprintf(
		"%sstock/api/getStockInfo.jsp?ex_ch=%s_%s.tw_%s&json=1&delay=0&_=%d",
		TWSEURL,
		"tse",
		stock.No,
		fmt.Sprintf(
			"%d%02d%02d",
			stock.Date.Year(),
			int(stock.Date.Month()),
			stock.Date.Day(),
		),
		stock.Timestamp,
	)
}

func (stock StockOption) GetData() (value stockBlob) {
	url := stock.GenStockUrl()
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&value)
	return
}
