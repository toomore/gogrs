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

// TWSECLASS is a class list of TWSE.
var TWSECLASS = map[string]string{
	"MS":         "大盤統計資訊",
	"0049":       "封閉式基金",
	"0099P":      "ETF",
	"01":         "水泥工業",
	"019919T":    "受益證券",
	"02":         "食品工業",
	"03":         "塑膠工業",
	"04":         "紡織纖維",
	"05":         "電機機械",
	"06":         "電器電纜",
	"07":         "化學生技醫療",
	"08":         "玻璃陶瓷",
	"09":         "造紙工業",
	"0999":       "認購權證", //(不含牛證)
	"0999B":      "熊證",
	"0999C":      "牛證",
	"0999G9":     "認股權憑證",
	"0999GA":     "附認股權特別股",
	"0999GD":     "附認股權公司債",
	"0999P":      "認售權證", //(不含熊證)
	"0999X":      "可展延牛證",
	"0999Y":      "可展延熊證",
	"10":         "鋼鐵工業",
	"11":         "橡膠工業",
	"12":         "汽車工業",
	"13":         "電子工業",
	"14":         "建材營造",
	"15":         "航運業",
	"16":         "觀光事業",
	"17":         "金融保險",
	"18":         "貿易百貨",
	"19":         "綜合",
	"20":         "其他",
	"21":         "化學工業",
	"22":         "生技醫療業",
	"23":         "油電燃氣業",
	"24":         "半導體業",
	"25":         "電腦及週邊設備業",
	"26":         "光電業",
	"27":         "通信網路業",
	"28":         "電子零組件業",
	"29":         "電子通路業",
	"30":         "資訊服務業",
	"31":         "其他電子業",
	"9299":       "存託憑證",
	"ALL":        "全部",
	"ALLBUT0999": "全部(不含權證、牛熊證、可展延牛熊證)",
	"CB":         "可轉換公司債",
}

// TWSECLASS is a class list of TWSE.
var OTCCLASS = map[string]string{
	"02": "食品工業",
	"03": "塑膠工業",
	"04": "紡織纖維",
	"05": "電機機械",
	"06": "電器電纜",
	"08": "玻璃陶瓷",
	"10": "鋼鐵工業",
	"11": "橡膠工業",
	"14": "建材營造",
	"15": "航運業",
	"16": "觀光事業",
	"17": "金融保險",
	"18": "貿易百貨",
	"20": "其他",
	"21": "化學工業",
	"22": "生技醫療業",
	"23": "油電燃氣業",
	"24": "半導體業",
	"25": "電腦及週邊設備業",
	"26": "光電業",
	"27": "通信網路業",
	"28": "電子零組件業",
	"29": "電子通路業",
	"30": "資訊服務業",
	"31": "其他電子業",
	"32": "文化創意業",
	"80": "管理股票",
	"AA": "受益證券",
	"AL": "所有證券",
	"BC": "牛證熊證",
	"EE": "上櫃指數股票型基金(ETF)",
	"EW": "所有證券(不含權證、牛熊證)",
	"GG": "認股權憑證",
	"TD": "台灣存託憑證(TDR)",
	"WW": "認購售權證",
}

// CategoryList to show TWSECLASS, OTCCLASS.
type CategoryList struct {
	same     map[string]string
	onlyTWSE map[string]string
	onlyOTC  map[string]string
}

// NewCategoryList to new one.
func NewCategoryList() *CategoryList {
	return &CategoryList{
		same:     nil,
		onlyTWSE: nil,
		onlyOTC:  nil,
	}
}

// Same return TWSECLASS and OTCCLASS.
func (c *CategoryList) Same() map[string]string {
	var result = c.same
	if result == nil {
		result = make(map[string]string)
		for i, v := range TWSECLASS {
			if _, ok := OTCCLASS[i]; ok {
				result[i] = v
			}
		}
	}
	return result
}

// OnlyTWSE return TWSECLASS but OTCCLASS
func (c *CategoryList) OnlyTWSE() map[string]string {
	var result = c.onlyTWSE
	if result == nil {
		result = make(map[string]string)
		same := c.Same()
		for i, v := range TWSECLASS {
			if _, ok := same[i]; !ok {
				result[i] = v
			}
		}
	}
	return result
}

// OnlyOTC return OTCCLASS but TWSECLASS
func (c *CategoryList) OnlyOTC() map[string]string {
	var result = c.onlyOTC
	if result == nil {
		result = make(map[string]string)
		same := c.Same()
		for i, v := range OTCCLASS {
			if _, ok := same[i]; !ok {
				result[i] = v
			}
		}
	}
	return result
}

// StockInfo is simple stock info for no, name.
type StockInfo struct {
	No   string
	Name string
}

// Lists is to get TWSE list.
type Lists struct {
	Date            time.Time
	FmtData         map[string]FmtListData
	categoryRawData map[string][][]string
	categoryNoList  map[string][]StockInfo
}

var errorNotSupport = errors.New("Not support.")

// NewLists new a Lists.
func NewLists(t time.Time) *Lists {
	return &Lists{
		Date:            t,
		FmtData:         make(map[string]FmtListData),
		categoryRawData: make(map[string][][]string),
		categoryNoList:  make(map[string][]StockInfo),
	}
}

// Get is to get TWSE csv data.
func (l *Lists) Get(category string) ([][]string, error) {
	if TWSECLASS[category] == "" {
		return nil, errorNotSupport
	}

	year, month, day := l.Date.Date()
	data, err := hCache.PostForm(fmt.Sprintf("%s%s", utils.TWSEHOST, utils.TWSELISTCSV),
		url.Values{"download": {"csv"}, "selectType": {category},
			"qdate": {fmt.Sprintf("%d/%02d/%02d", year-1911, month, day)}})

	if err != nil {
		return nil, fmt.Errorf(errorNetworkFail.Error(), err)
	}

	csvArrayContent := strings.Split(string(data), "\n")

	var csvReader *csv.Reader
	switch category {
	case "MS":
		if len(csvArrayContent) > 6 {
			csvReader = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[4:51], "\n")))
		}
	case "ALLBUT0999", "ALL":
		if len(csvArrayContent) > 155 {
			re := regexp.MustCompile("^=?[\"]{1}[0-9A-Z]{4,}")
			var pickdata []string
			for _, v := range csvArrayContent {
				if re.MatchString(v) {
					if v[0] == 61 {
						pickdata = append(pickdata, v[1:])
					} else {
						pickdata = append(pickdata, v)
					}
				}
			}
			csvReader = csv.NewReader(strings.NewReader(strings.Join(pickdata, "\n")))
		}
	default:
		if len(csvArrayContent) > 9 {
			csvReader = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[4:len(csvArrayContent)-7], "\n")))
		}
	}
	if csvReader != nil {
		returnData, err := csvReader.ReadAll()
		switch category {
		default:
			if err == nil {
				l.categoryRawData[category] = returnData
				l.formatData(category)
			}
		case "MS":
		}
		return returnData, err
	}
	return nil, errorNotEnoughData
}

// GetCategoryList 取得分類的股票代碼與名稱列表
func (l Lists) GetCategoryList(category string) []StockInfo {
	if _, ok := l.categoryNoList[category]; !ok {
		l.Get(category)
	}
	return l.categoryNoList[category]
}

// FmtListData 格式化個股的資料資訊
type FmtListData struct {
	No             string
	Name           string
	Volume         uint64  //成交股數
	TotalPrice     uint64  //成交金額
	Open           float64 //開盤價
	High           float64 //最高價
	Low            float64 //最低價
	Price          float64 //收盤價
	Range          float64 //漲跌價差
	Totalsale      uint64  //成交筆數
	LastBuyPrice   float64 //最後揭示買價
	LastBuyVolume  uint64  //最後揭示買量
	LastSellPrice  float64 //最後揭示賣價
	LastSellVolume uint64  //最後揭示賣量
	PERatio        float64 //本益比
	IssuedShares   uint64  //發行股數
}

func (l *Lists) formatData(categoryNo string) {
	if _, ok := l.categoryNoList[categoryNo]; !ok {
		l.categoryNoList[categoryNo] = make([]StockInfo, len(l.categoryRawData[categoryNo]))
	}

	for i, v := range l.categoryRawData[categoryNo] {
		var data FmtListData
		data.No = strings.Trim(v[0], " ")
		data.Name = strings.Trim(v[1], " ")
		data.Volume, _ = strconv.ParseUint(strings.Replace(v[2], ",", "", -1), 10, 64)
		data.Totalsale, _ = strconv.ParseUint(v[3], 10, 32)
		data.TotalPrice, _ = strconv.ParseUint(strings.Replace(v[4], ",", "", -1), 10, 64)
		data.Open, _ = strconv.ParseFloat(v[5], 64)
		data.High, _ = strconv.ParseFloat(v[6], 64)
		data.Low, _ = strconv.ParseFloat(v[7], 64)
		data.Price, _ = strconv.ParseFloat(v[8], 64)
		data.Range, _ = strconv.ParseFloat(fmt.Sprintf("%s%s", v[9], v[10]), 64)
		data.LastBuyPrice, _ = strconv.ParseFloat(v[11], 64)
		data.LastBuyVolume, _ = strconv.ParseUint(v[12], 10, 32)
		data.LastSellPrice, _ = strconv.ParseFloat(v[13], 64)
		data.LastSellVolume, _ = strconv.ParseUint(v[14], 10, 32)
		data.PERatio, _ = strconv.ParseFloat(v[15], 64)

		l.FmtData[data.No] = data
		l.categoryNoList[categoryNo][i] = StockInfo{No: data.No, Name: data.Name}
	}
}

// OTCLists is to get OTC list.
type OTCLists struct {
	Date            time.Time
	FmtData         map[string]FmtListData
	categoryRawData map[string][][]string
	categoryNoList  map[string][]StockInfo
}

// NewOTCLists new a Lists.
func NewOTCLists(date time.Time) *OTCLists {
	return &OTCLists{
		Date:            date,
		FmtData:         make(map[string]FmtListData),
		categoryRawData: make(map[string][][]string),
		categoryNoList:  make(map[string][]StockInfo),
	}
}

// Get is to get OTC csv data.
func (o *OTCLists) Get(category string) ([][]string, error) {
	var (
		csvArrayContent []string
		csvReader       *csv.Reader
		data            []byte
		err             error
		rawData         [][]string
		url             string
	)

	url = fmt.Sprintf("%s%s", utils.OTCHOST, fmt.Sprintf(utils.OTCLISTCSV, fmt.Sprintf("%d/%02d/%02d", o.Date.Year()-1911, o.Date.Month(), o.Date.Day()), category))

	if data, err = hCache.Get(url, false); err == nil {
		csvArrayContent = strings.Split(string(data), "\n")
		if len(csvArrayContent) > 5 {
			csvReader = csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[4:len(csvArrayContent)-1], "\n")))
			if rawData, err = csvReader.ReadAll(); err == nil {
				o.categoryRawData[category] = rawData
				o.formatData(category)
				return rawData, nil
			}
		}
	}
	return nil, err
}

// GetCategoryList 取得分類的股票代碼與名稱列表
func (o OTCLists) GetCategoryList(category string) []StockInfo {
	if _, ok := o.categoryNoList[category]; !ok {
		o.Get(category)
	}
	return o.categoryNoList[category]
}

func (o *OTCLists) formatData(categoryNo string) {
	if _, ok := o.categoryNoList[categoryNo]; !ok {
		o.categoryNoList[categoryNo] = make([]StockInfo, len(o.categoryRawData[categoryNo]))
	}

	for i, v := range o.categoryRawData[categoryNo] {
		var data FmtListData
		data.No = strings.Trim(v[0], " ")
		data.Name = strings.Trim(v[1], " ")
		data.Volume, _ = strconv.ParseUint(strings.Replace(v[7], ",", "", -1), 10, 64)
		data.Totalsale, _ = strconv.ParseUint(strings.Replace(v[9], ",", "", -1), 10, 64)
		data.TotalPrice, _ = strconv.ParseUint(strings.Replace(v[8], ",", "", -1), 10, 64)
		data.Open, _ = strconv.ParseFloat(v[4], 64)
		data.High, _ = strconv.ParseFloat(v[5], 64)
		data.Low, _ = strconv.ParseFloat(v[6], 64)
		data.Price, _ = strconv.ParseFloat(v[2], 64)
		data.Range, _ = strconv.ParseFloat(strings.Replace(v[3], " ", "", -1), 64)
		data.LastBuyPrice, _ = strconv.ParseFloat(v[10], 64)
		data.LastSellPrice, _ = strconv.ParseFloat(v[11], 64)
		data.IssuedShares, _ = strconv.ParseUint(strings.Replace(v[12], ",", "", -1), 10, 64)

		o.FmtData[data.No] = data
		o.categoryNoList[categoryNo][i] = StockInfo{No: data.No, Name: data.Name}
	}
}
