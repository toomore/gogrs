// Package realtime - Fetch realtime stock data info
// 擷取盤中個股、指數即時股價資訊
//
package realtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
				time.Now().UnixNano()/1000000,
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
	log.Println(stock.URL())
	req, _ := http.NewRequest("GET", stock.URL(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:40.0) Gecko/20100101 Firefox/40.0")
	req.Header.Set("Referer", fmt.Sprintf("http://mis.twse.com.tw/stock/fibest.jsp?stock=%s", stock.No))
	//req.Header.Set("Cookie", "JSESSIONID=52978A658B945F6B400D480FF6EE00E3; JSESSIONID=83C2D4470AA4F430CAC02EF7A162C623")
	//req.Header.Set("Cookie", "JSESSIONID=52978A658B945F6B400D480FF6EE00E3; JSESSIONID=D7953B63154DD91F6F13957DE3893D40")
	//req.Header.Set("")

	///misr, _ := http.NewRequest("GET", "http://mis.tse.com.tw/", nil)
	/////misr, _ := http.NewRequest("GET", "http://mis.tse.com.tw/stock/index.jsp", nil)
	///misr.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:39.0) Gecko/20100101 Firefox/39.0")
	///rr, _ := utils.HttpClient.Do(misr)
	///coo := rr.Header.Get("Set-Cookie")
	///log.Println(rr.Cookies())
	///log.Println(rr.Location())

	misr2, _ := http.NewRequest("GET", "http://mis.twse.com.tw/stock/index.jsp", nil)
	//misr2.Header.Set("Host", "mis.twse.com.tw")
	//misr2.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01 Accept-Language: zh-TW,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	//misr2.Header.Set("Accept-Encoding", "gzip, deflate")
	//misr2.Header.Set("x-misr2uested-with", "XMLHttpmisr2uest")
	misr2.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:40.0) Gecko/20100101 Firefox/40.0")
	//misr2.Header.Set("Cookie", strings.Trim(strings.Split(coo, " ")[0], ";"))
	log.Println("><><>", misr2)
	rr2, _ := utils.HttpClient.Do(misr2)
	coo2 := rr2.Header.Get("Set-Cookie")

	//log.Println(coo)
	log.Println("CCCC", coo2)
	//req.Header.Set("Cookie", strings.Trim(strings.Split(coo, " ")[0], ";"))
	//req.Header.Set("Cookie", "JSESSIONID=22159BA16824188BE08E96AF54444A4D; JSESSIONID=125DB9963B8745C20E2D77DFCADB8470")
	//req.Header.Set("Host", "mis.twse.com.tw")
	//req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01 Accept-Language: zh-TW,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	//req.Header.Set("x-requested-with", "XMLHttpRequest")
	coo2 = "JSESSIONID=A7320CE8F258CCB56C8D337C690FE32B; JSESSIONID=214C21A6EAB8BCFE696BA1974D9901A7; Path=/stock"
	req.Header.Set("Cookie", coo2)
	req.Header.Set("Connection", "keep-alive")
	log.Println(">>>>", req.Header)
	if resp, err = utils.HttpClient.Do(req); err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&value)

		if len(value.MsgArray) == 0 {
			err = errorNotEnoughData
		}
	} else {
		err = fmt.Errorf(errorNetworkFail.Error(), err)
	}

	log.Println(err)
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
