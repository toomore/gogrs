package twse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	iconv "github.com/djimenez/iconv-go"
	"github.com/toomore/gogrs/utils"
)

// TWSECLASS is a class list of TWSE.
var TWSECLASS = map[string]string{
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

// TWSEList is to get TWSE list.
type TWSEList struct {
	Date time.Time
}

// URL is to render urlpath.
func (l TWSEList) URL(strNo string) string {
	year, month, day := l.Date.Date()
	return fmt.Sprintf("%s%s", utils.TWSEHOST, fmt.Sprintf(utils.TWSELISTCSV, year, month, year, month, day, strNo))
}

// Get is to get csv data.
func (l TWSEList) Get(strNo string) ([][]string, error) {
	data, err := http.Get(l.URL(strNo))
	if err != nil {
		return nil, fmt.Errorf("Network fail: %s", err)
	}
	defer data.Body.Close()
	dataContentBig5, _ := ioutil.ReadAll(data.Body)
	dataContent, _ := iconv.ConvertString(string(dataContentBig5), "big5", "utf-8")
	csvArrayContent := strings.Split(dataContent, "\n")
	if len(csvArrayContent) > 5 {
		//for i := range csvArrayContent {
		//	fmt.Println(i, csvArrayContent[i])
		//}
		csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[2:len(csvArrayContent)-5], "\n")))
		return csvReader.ReadAll()
	}
	return nil, errors.New("Not enough data.")
}
