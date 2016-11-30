package twse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/toomore/gogrs/utils"
)

type unixMapData map[int64][][]string

var (
	errorNetworkFail   = errors.New("Network fail: %s")
	errorNotEnoughData = errors.New("not enough data")
)

// Data start with stock no, date.
type Data struct {
	No             string
	Name           string
	Date           time.Time
	RawData        [][]string
	UnixMapData    unixMapData
	exchange       string
	openList       []float64
	highList       []float64
	lowList        []float64
	priceList      []float64
	rangeList      []float64
	dailyRangeList []float64
	volumeList     []uint64
	datelist       []time.Time
}

// NewTWSE 建立一個 TWSE 上市股票
func NewTWSE(No string, Date time.Time) *Data {
	return &Data{
		No:          No,
		Date:        Date,
		exchange:    "tse",
		UnixMapData: make(unixMapData),
	}
}

// NewOTC 建立一個 OTC 上櫃股票
func NewOTC(No string, Date time.Time) *Data {
	return &Data{
		No:          No,
		Date:        Date,
		exchange:    "otc",
		UnixMapData: make(unixMapData),
	}
}

// URL return stock csv url path.
func (d Data) URL() string {
	switch d.exchange {
	case "tse":
		return fmt.Sprintf("%s%s", utils.TWSEHOST, utils.TWSECSV)
	case "otc":
		return fmt.Sprintf("%s%s",
			utils.OTCHOST,
			fmt.Sprintf(utils.OTCCSV, d.Date.Year()-1911, d.Date.Month(), d.No))
	}

	return ""
}

// Round will do sub one month.
func (d *Data) Round() {
	year, month, _ := d.Date.Date()
	d.Date = time.Date(year, month-1, 1, 0, 0, 0, 0, d.Date.Location())
}

// PlusData will do Round() and Get().
func (d *Data) PlusData() {
	d.Round()
	d.Get()
}

func (d *Data) clearCache() {
	d.rangeList = nil
	d.dailyRangeList = nil
	d.openList = nil
	d.highList = nil
	d.lowList = nil
	d.priceList = nil
	d.volumeList = nil
	d.datelist = nil
}

var noName = regexp.MustCompile(`^([A-Z0-9]+)(.+)`)

// Get return csv data in array.
func (d *Data) Get() ([][]string, error) {
	if len(d.UnixMapData[d.Date.Unix()]) == 0 {
		var data []byte
		var err error
		switch d.exchange {
		case "tse":
			data, err = hCache.PostForm(
				d.URL(),
				url.Values{"download": {"csv"}, "query_year": {strconv.Itoa(d.Date.Year())}, "query_month": {strconv.Itoa(int(d.Date.Month()))}, "CO_ID": {d.No}})
		case "otc":
			data, err = hCache.Get(d.URL(), true)
		}
		if err != nil {
			return nil, fmt.Errorf(errorNetworkFail.Error(), err)
		}
		csvArrayContent := strings.Split(string(data), "\n")
		for i := range csvArrayContent {
			csvArrayContent[i] = strings.TrimSpace(csvArrayContent[i])
		}
		var csvReader *csv.Reader
		if (d.exchange == "tse" && len(csvArrayContent) > 2) || (d.exchange == "otc" && len(csvArrayContent) > 5) {
			if d.exchange == "tse" {
				if d.Name == "" {
					groups := noName.FindStringSubmatch(strings.Split(csvArrayContent[0], " ")[1])
					d.No, d.Name = groups[1], groups[2]
				}
				csvReader = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[2:], "\n")))
			} else if d.exchange == "otc" {
				if d.Name == "" {
					d.Name = strings.Split(csvArrayContent[2], ":")[1]
				}
				csvReader = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[5:len(csvArrayContent)-1], "\n")))
			}
			allData, err := csvReader.ReadAll()
			d.RawData = append(allData, d.RawData...)
			d.UnixMapData[d.Date.Unix()] = allData
			d.clearCache()
			return allData, err
		}
		return nil, errorNotEnoughData
	}
	return d.UnixMapData[d.Date.Unix()], nil
}

// GetByTimeMap return a map by key of time.Time
func (d Data) GetByTimeMap() map[time.Time]interface{} {
	var (
		dailyData [][]string
		data      map[time.Time]interface{}
		err       error
	)
	if dailyData, err = d.Get(); err == nil {
		data = make(map[time.Time]interface{})
		for _, v := range dailyData {
			data[utils.ParseDate(v[0])] = v
		}
	}
	return data
}

func (d Data) getColsList(colsNo int) []string {
	var result []string
	result = make([]string, len(d.RawData))
	for i, value := range d.RawData {
		result[i] = value[colsNo]
	}
	return result
}

func (d Data) getColsListFloat64(colsNo int) []float64 {
	var result []float64
	result = make([]float64, len(d.RawData))
	for i, v := range d.getColsList(colsNo) {
		result[i], _ = strconv.ParseFloat(strings.Replace(v, ",", "", -1), 64)
	}
	return result
}

// GetVolumeList 取得 成交股數 序列
func (d *Data) GetVolumeList() []uint64 {
	if d.volumeList == nil {
		var (
			unit   uint64 = 1
			result []uint64
		)
		if d.exchange == "otc" {
			unit = 1000
		}
		result = make([]uint64, len(d.RawData))
		for i, v := range d.getColsList(1) {
			result[i], _ = strconv.ParseUint(strings.Replace(v, ",", "", -1), 10, 64)
			result[i] = result[i] * unit
		}
		d.volumeList = result
	}
	return d.volumeList
}

// GetOpenList 取得 開盤價 序列
func (d *Data) GetOpenList() []float64 {
	if d.openList == nil {
		d.openList = d.getColsListFloat64(3)
	}
	return d.openList
}

// GetHighList 取得 最高價 序列
func (d *Data) GetHighList() []float64 {
	if d.highList == nil {
		d.highList = d.getColsListFloat64(4)
	}
	return d.highList
}

// GetLowList 取得 最低價 序列
func (d *Data) GetLowList() []float64 {
	if d.lowList == nil {
		d.lowList = d.getColsListFloat64(5)
	}
	return d.lowList
}

// GetPriceList 取得 收盤價 序列
func (d *Data) GetPriceList() []float64 {
	if d.priceList == nil {
		d.priceList = d.getColsListFloat64(6)
	}
	return d.priceList
}

// GetRangeList 取得 漲跌價差 序列（為當日收盤價與前一日收盤價比較）
func (d *Data) GetRangeList() []float64 {
	if d.rangeList == nil {
		d.rangeList = d.getColsListFloat64(7)
	}
	return d.rangeList
}

// GetDailyRangeList 計算收盤與開盤價差
func (d *Data) GetDailyRangeList() []float64 {
	if d.dailyRangeList == nil {
		d.dailyRangeList = utils.CalDiffFloat64(d.GetPriceList(), d.GetOpenList())
	}
	return d.dailyRangeList
}

// GetDateList 取出資料時間序列
func (d *Data) GetDateList() []time.Time {
	if d.datelist == nil {
		datelistdata := d.getColsList(0)
		d.datelist = make([]time.Time, len(datelistdata))
		for i, v := range datelistdata {
			d.datelist[i] = utils.ParseDate(v)
		}
	}
	return d.datelist
}

// MA 計算 收盤價 的移動平均
func (d Data) MA(days int) []float64 {
	var (
		priceList = d.GetPriceList()
		result    []float64
	)
	result = make([]float64, len(priceList)-days+1)
	for i := range priceList[days-1:] {
		result[i] = utils.AvgFloat64(priceList[i : i+days])
	}
	return result
}

// MABR 計算 收盤價移動平均 的乖離
func (d Data) MABR(days1, days2 int) []float64 {
	return utils.CalDiffFloat64(d.MA(days1), d.MA(days2))
}

// MAV 計算 成交股數 的移動平均
func (d Data) MAV(days int) []uint64 {
	var (
		result     []uint64
		volumeList = d.GetVolumeList()
	)
	result = make([]uint64, len(volumeList)-days+1)
	for i := range volumeList[days-1:] {
		result[i] = utils.AvgUint64(volumeList[i : i+days])
	}
	return result
}

// MAVBR 計算 成交股數移動平均 的乖離
func (d Data) MAVBR(days1, days2 int) []int64 {
	var (
		MAV1    = d.MAV(days1)
		MAV2    = d.MAV(days2)
		length  int
		result1 []int64
		result2 []int64
	)
	if len(MAV1) <= len(MAV2) {
		length = len(MAV1)
	} else {
		length = len(MAV2)
	}
	result1 = make([]int64, length)
	result2 = make([]int64, length)
	for i := 1; i <= length; i++ {
		result1[length-i] = int64(MAV1[len(MAV1)-i])
		result2[length-i] = int64(MAV2[len(MAV2)-i])
	}

	return utils.CalDiffInt64(result1, result2)
}

// IsRed 計算是否收紅 K
func (d Data) IsRed() bool {
	var openList = d.GetOpenList()
	var priceList = d.GetPriceList()
	return priceList[len(priceList)-1] > openList[len(openList)-1]
}

// IsThanYesterday 計算漲跌與昨日收盤相比
func (d Data) IsThanYesterday() bool {
	var rangeList = d.GetRangeList()
	return rangeList[len(rangeList)-1] > 0
}

// Len 返回資料序列長度
func (d Data) Len() int {
	return len(d.RawData)
}

// FmtData is struct for daily data format.
type FmtData struct {
	Date       time.Time
	Volume     uint64  //成交股數
	TotalPrice uint64  //成交金額
	Open       float64 //開盤價
	High       float64 //最高價
	Low        float64 //最低價
	Price      float64 //收盤價
	Range      float64 //漲跌價差
	Totalsale  uint64  //成交筆數
}

// FormatData is format daily data.
func (d Data) FormatData() []FmtData {
	var (
		loopd  FmtData
		result []FmtData
	)
	result = make([]FmtData, len(d.RawData))
	for i, v := range d.RawData {
		loopd.Date = utils.ParseDate(v[0])
		loopd.Volume, _ = strconv.ParseUint(strings.Replace(v[1], ",", "", -1), 10, 64)
		loopd.TotalPrice, _ = strconv.ParseUint(strings.Replace(v[2], ",", "", -1), 10, 64)
		loopd.Open, _ = strconv.ParseFloat(v[3], 64)
		loopd.High, _ = strconv.ParseFloat(v[4], 64)
		loopd.Low, _ = strconv.ParseFloat(v[5], 64)
		loopd.Price, _ = strconv.ParseFloat(v[6], 64)
		loopd.Range, _ = strconv.ParseFloat(v[7], 64)
		loopd.Totalsale, _ = strconv.ParseUint(strings.Replace(v[8], ",", "", -1), 10, 64)
		result[i] = loopd
	}
	return result
}

var hCache *utils.HTTPCache

func init() {
	hCache = utils.NewHTTPCache(utils.GetOSRamdiskPath(), "cp950")
}
