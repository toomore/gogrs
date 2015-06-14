// Package realtime - Fetch realtime stock data info
// 擷取盤中個股、指數即時股價資訊
//
package realtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

type msgArray []map[string]string
type unixMapData map[int64]Data

// StockRealTime start with No, Timestamp, Date.
type StockRealTime struct {
	No          string      // 股票代碼
	Date        time.Time   // 擷取時間
	UnixMapData unixMapData // 時間資料暫存
	Exchange    string      // tse, otc
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
	if utils.ExchangeMap[stock.Exchange] {
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
				utils.RandInt(),
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
	Category string // 股票類別代號
}

// Data is realtime return formated data.
type Data struct {
	BestAskPrice   []float64              // 最佳五檔賣出價資訊
	BestBidPrice   []float64              // 最佳五檔買進價資訊
	BestAskVolume  []int64                // 最佳五檔賣出量資訊
	BestBidVolume  []int64                // 最佳五檔買進量資訊
	Open           float64                // 開盤價格
	Highest        float64                // 最高價
	Lowest         float64                // 最低價
	Price          float64                // 該盤成交價格
	LimitUp        float64                // 漲停價
	LimitDown      float64                // 跌停價
	Volume         float64                // 該盤成交量
	VolumeAcc      float64                // 累計成交量
	YesterdayPrice float64                // 昨日收盤價格
	TradeTime      time.Time              // 交易時間
	Info           StockInfo              // 相關資訊
	SysInfo        map[string]interface{} // 系統回傳資訊
}

var (
	errorNetworkFail   = errors.New("Network fail: %s")
	errorNotEnoughData = errors.New("Not enough data.")
)

func (stock *StockRealTime) get() (StockBlob, error) {
	var (
		err   error
		resp  *http.Response
		value StockBlob
	)

	if resp, err = http.Get(stock.URL()); err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&value)

		if len(value.MsgArray) == 0 {
			err = errorNotEnoughData
		}
	} else {
		err = fmt.Errorf(errorNetworkFail.Error(), err)
	}

	return value, err
}

// Get return stock realtime map data.
func (stock *StockRealTime) Get() (Data, error) {
	var (
		err    error
		result Data
		value  StockBlob
	)

	if value, err = stock.get(); err == nil && len(value.MsgArray) != 0 {
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
		tlong, _ := strconv.ParseInt(value.MsgArray[0]["tlong"], 10, 64)
		result.TradeTime = time.Unix(tlong/1000, 0)

		result.Info.No = value.MsgArray[0]["c"]
		result.Info.FullName = value.MsgArray[0]["nf"]
		result.Info.Name = value.MsgArray[0]["n"]
		result.Info.Ticker = value.MsgArray[0]["ch"]
		result.Info.Exchange = value.MsgArray[0]["ex"]
		result.Info.Category = value.MsgArray[0]["i"]

		result.SysInfo = make(map[string]interface{})
		result.SysInfo = value.QueryTime

		// Record
		stock.UnixMapData[tlong/1000] = result
	}
	return result, err
}

// NewTWSE 建立一個上市股票
func NewTWSE(No string, Date time.Time) *StockRealTime {
	return &StockRealTime{
		No:          No,
		Date:        Date,
		Exchange:    "tse",
		UnixMapData: make(unixMapData),
	}
}

// NewOTC 建立一個上櫃股票
func NewOTC(No string, Date time.Time) *StockRealTime {
	return &StockRealTime{
		No:          No,
		Date:        Date,
		Exchange:    "otc",
		UnixMapData: make(unixMapData),
	}
}

// NewWeight 大盤指數
func NewWeight(Date time.Time) *StockRealTime {
	return &StockRealTime{
		No:          "t00",
		Date:        Date,
		Exchange:    "tse",
		UnixMapData: make(unixMapData),
	}
}

// NewOTCI 上櫃指數
func NewOTCI(Date time.Time) *StockRealTime {
	return &StockRealTime{
		No:          "o00",
		Date:        Date,
		Exchange:    "otc",
		UnixMapData: make(unixMapData),
	}
}

// NewFRMSA 寶島指數
func NewFRMSA(Date time.Time) *StockRealTime {
	return &StockRealTime{
		No:          "FRMSA",
		Date:        Date,
		Exchange:    "tse",
		UnixMapData: make(unixMapData),
	}
}
