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

func (stock StockOption) GenStockURL() string {
	return fmt.Sprintf("%s%s", TWSEURL,
		fmt.Sprintf(TWSEREAL,
			"tse",
			stock.No,
			fmt.Sprintf(
				"%d%02d%02d",
				stock.Date.Year(),
				int(stock.Date.Month()),
				stock.Date.Day(),
			),
			stock.Timestamp,
		))
}

func (stock StockOption) GetData() (value stockBlob) {
	url := stock.GenStockURL()
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	} else {
		fmt.Println(resp)
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&value)
	}
	return
}
