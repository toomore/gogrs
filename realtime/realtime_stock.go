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
	Exchange    string
}

var exchangeMap = map[string]bool{"tse": true, "otc": true}

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
	if exchangeMap[stock.Exchange] {
		return fmt.Sprintf("%s%s", utils.TWSEURL,
			fmt.Sprintf(utils.TWSEREAL,
				stock.Exchange,
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
	return ""
}

// StockInfo is base stock info.
type StockInfo struct {
	Exchange string // tse or otc
	FullName string // Full company name.
	Name     string // Stock name.
	No       string // Stock no
	Ticker   string // Ticker symbol（股票代號）
}

// Data is realtime return formated data.
type Data struct {
	BestAskPrice   []float64
	BestBidPrice   []float64
	BestAskVolume  []int64
	BestBidVolume  []int64
	Open           float64
	Highest        float64
	Lowest         float64
	Price          float64
	LimitUp        float64
	LimitDown      float64
	Volume         float64
	VolumeAcc      float64
	YesterdayPrice float64
	Info           StockInfo
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
func (stock *StockRealTime) Get() (Data, error) {
	value, err := stock.get()

	if err != nil {
		return Data{}, err
	}

	if len(value.MsgArray) != 0 {
		var result Data
		aList := strings.Split(value.MsgArray[0]["a"], "_")
		result.BestAskPrice = make([]float64, len(aList)-1)
		for i, v := range aList[:len(aList)-1] {
			result.BestAskPrice[i], _ = strconv.ParseFloat(v, 10)
		}

		bList := strings.Split(value.MsgArray[0]["b"], "_")
		result.BestBidPrice = make([]float64, len(bList)-1)
		for i, v := range bList[:len(bList)-1] {
			result.BestBidPrice[i], _ = strconv.ParseFloat(v, 10)
		}

		fList := strings.Split(value.MsgArray[0]["f"], "_")
		result.BestAskVolume = make([]int64, len(fList)-1)
		for i, v := range fList[:len(fList)-1] {
			result.BestAskVolume[i], _ = strconv.ParseInt(v, 10, 64)
		}

		gList := strings.Split(value.MsgArray[0]["g"], "_")
		result.BestBidVolume = make([]int64, len(gList)-1)
		for i, v := range gList[:len(gList)-1] {
			result.BestBidVolume[i], _ = strconv.ParseInt(v, 10, 64)
		}

		result.Open, _ = strconv.ParseFloat(value.MsgArray[0]["o"], 10)
		result.Highest, _ = strconv.ParseFloat(value.MsgArray[0]["h"], 10)
		result.Lowest, _ = strconv.ParseFloat(value.MsgArray[0]["l"], 10)
		result.Price, _ = strconv.ParseFloat(value.MsgArray[0]["z"], 10)
		result.LimitUp, _ = strconv.ParseFloat(value.MsgArray[0]["u"], 10)
		result.LimitDown, _ = strconv.ParseFloat(value.MsgArray[0]["w"], 10)
		result.Volume, _ = strconv.ParseFloat(value.MsgArray[0]["tv"], 10)
		result.VolumeAcc, _ = strconv.ParseFloat(value.MsgArray[0]["v"], 10)
		result.YesterdayPrice, _ = strconv.ParseFloat(value.MsgArray[0]["y"], 10)

		result.Info.No = value.MsgArray[0]["n"]
		result.Info.FullName = value.MsgArray[0]["nf"]
		result.Info.No = value.MsgArray[0]["n"]
		result.Info.Ticker = value.MsgArray[0]["ch"]
		result.Info.Exchange = value.MsgArray[0]["ex"]

		return result, nil
	}

	return Data{}, nil
}
