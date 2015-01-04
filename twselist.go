package gogrs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// TWSEList is to get TWSE list.
type TWSEList struct {
	Date time.Time
}

// URL is to render urlpath.
func (l TWSEList) URL(strNo string) string {
	year, month, day := l.Date.Date()
	return fmt.Sprintf("%s%s", TWSEHOST, fmt.Sprintf(TWSELISTCSV, year, month, year, month, day, strNo))
}

// GetData is to get csv data.
func (l TWSEList) GetData(strNo string) ([][]string, error) {
	data, _ := http.Get(l.URL(strNo))
	defer data.Body.Close()
	dataContent, _ := ioutil.ReadAll(data.Body)
	csvArrayContent := strings.Split(string(dataContent), "\n")
	if len(csvArrayContent) > 5 {
		//for i := range csvArrayContent {
		//	fmt.Println(i, csvArrayContent[i])
		//}
		csvReader := csv.NewReader(strings.NewReader(strings.Join(csvArrayContent[2:len(csvArrayContent)-5], "\n")))
		return csvReader.ReadAll()
	}
	return nil, errors.New("Not enough data.")
}
