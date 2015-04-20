package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/toomore/gogrs/tradingdays"
)

func Log(req *http.Request) {
	log.Println(req.URL, req.UserAgent(), req.Form)
}

func Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
	Log(req)
}

type tradeJSON struct {
	Date time.Time `json:"date"`
	Open bool      `json:"open"`
}

func TradeOpen(w http.ResponseWriter, req *http.Request) {
	data, err := strconv.ParseInt(req.FormValue("q"), 10, 64)
	if err != nil {
		w.Write([]byte("Wrong data format."))
	} else {
		date := time.Unix(data, 0)
		json_str, _ := json.Marshal(&tradeJSON{
			Date: date.UTC(),
			Open: tradingdays.IsOpen(date.Year(), date.Month(), date.Day(), date.Location())})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_str)
	}
	Log(req)
}

func init() {
	tradingdays.DownloadCSV(true)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/open", TradeOpen)
	log.Fatal(http.ListenAndServe(":59123", nil))
}
