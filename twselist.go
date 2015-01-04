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

type TWSEList struct {
	Date time.Time
}

func (l TWSEList) baseURL() string {
	year, month, day := l.Date.Date()
	return fmt.Sprintf("%s%s", TWSEHOST, fmt.Sprintf(TWSELISTCSV, year, month, year, month, day))
}

func (l TWSEList) URL(strNo string) string {
	return fmt.Sprintf(l.baseURL(), strNo)
}

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
