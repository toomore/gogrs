// Package realtime - Fetch realtime stock data info http://mis.tse.com.tw/
// 擷取盤中即時股價資訊
//
package realtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

//STOCKPATH = '/stock/api/getStockInfo.jsp?ex_ch=%(exchange)s_%(no)s.tw_%(date)s&json=1&delay=%(delay)s&_=%(timestamp)s'

type msgArray []map[string]string
type unixMapData map[int64]msgArray

// StockRealTime start with No, Timestamp, Date.
type StockRealTime struct {
	No          string
	Timestamp   int64
	Date        time.Time
	UnixMapData unixMapData
}

// StockBlob return map data.
type StockBlob struct {
	Rtcode    string
	UserDelay int
	Rtmessage string
	Referer   string
	MsgArray  msgArray
	QueryTime map[string]interface{}
}

// URL return realtime url path.
func (stock StockRealTime) URL() string {
	return fmt.Sprintf("%s%s", utils.TWSEURL,
		fmt.Sprintf(utils.TWSEREAL,
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

type Date struct {
	BestAskPrice []float64
}

func (stock *StockRealTime) get() (StockBlob, error) {
	var value StockBlob
	url := stock.URL()
	resp, err := http.Get(url)
	if err != nil {
		return value, fmt.Errorf("Network fail: %s", err)
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&value)

	if len(value.MsgArray) != 0 {
		unixTime, _ := strconv.ParseInt(value.MsgArray[0]["tlong"], 10, 64)
		if stock.UnixMapData == nil {
			stock.UnixMapData = make(unixMapData)
		}

		// Should format data.
		stock.UnixMapData[unixTime/1000] = value.MsgArray
	}

	return value, nil
}

// Get return stock realtime map data.
func (stock *StockRealTime) Get() (Date, error) {
	value, err := stock.get()

	if err != nil {
		return Date{}, err
	}

	if len(value.MsgArray) != 0 {
		var result Date
		aList := strings.Split(value.MsgArray[0]["a"], "_")
		result.BestAskPrice = make([]float64, len(aList)-1)
		for i, v := range aList[:len(aList)-1] {
			result.BestAskPrice[i], _ = strconv.ParseFloat(v, 10)
		}
		return result, nil
	}

	return Date{}, nil
}
