package twse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	iconv "github.com/djimenez/iconv-go"
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

// Lists is to get TWSE list.
type Lists struct {
	Date            time.Time
	FmtData         map[string]FmtListData
	categoryRawData map[string][][]string
	categoryNoList  map[string][]string
}

// Get is to get csv data.
func (l *Lists) Get(category string) ([][]string, error) {
	if TWSECLASS[category] == "" {
		return nil, errors.New("Not support.")
	}

	year, month, day := l.Date.Date()
	data, err := http.PostForm(fmt.Sprintf("%s%s", utils.TWSEHOST, utils.TWSELISTCSV),
		url.Values{"download": {"csv"}, "selectType": {category},
			"qdate": {fmt.Sprintf("%d/%02d/%02d", year-1911, month, day)}})
	defer data.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("Network fail: %s", err)
	}

	dataContentBig5, _ := ioutil.ReadAll(data.Body)
	dataContent, _ := iconv.ConvertString(string(dataContentBig5), "big5", "utf-8")
	csvArrayContent := strings.Split(dataContent, "\n")

	if category == "MS" {
		if len(csvArrayContent) > 6 {
			csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[4:51], "\n")))
			returnData, err := csvReader.ReadAll()
			return returnData, err
		}
	} else {
		if len(csvArrayContent) > 9 {
			csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[4:len(csvArrayContent)-7], "\n")))
			returnData, err := csvReader.ReadAll()
			if err == nil {
				if l.categoryRawData == nil {
					l.categoryRawData = make(map[string][][]string)
				}
				l.categoryRawData[category] = returnData
				l.formatData(category)
			}
			return returnData, err
		}
	}
	return nil, errors.New("Not enough data.")
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
}

func (l *Lists) formatData(categoryNo string) {
	if l.FmtData == nil {
		l.FmtData = make(map[string]FmtListData)
	}

	if l.categoryNoList == nil {
		l.categoryNoList = make(map[string][]string)
	}

	if _, ok := l.categoryNoList[categoryNo]; !ok {
		l.categoryNoList[categoryNo] = make([]string, len(l.categoryRawData[categoryNo]))
	}

	for i, v := range l.categoryRawData[categoryNo] {
		var data FmtListData
		data.No = strings.Trim(v[0], " ")
		data.Name = strings.Trim(v[1], " ")
		data.Volume, _ = strconv.ParseUint(v[2], 10, 32)
		data.Totalsale, _ = strconv.ParseUint(v[3], 10, 32)
		data.TotalPrice, _ = strconv.ParseUint(v[4], 10, 32)
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
		l.categoryNoList[categoryNo][i] = data.No
	}
}
