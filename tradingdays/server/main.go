package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/toomore/gogrs/tradingdays"
)

// Log is show viwer log.
func Log(req *http.Request) {
	log.Println(req.URL, req.UserAgent(), req.Form)
}

// Home is home page.
func Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
	Log(req)
}

type tradeJSON struct {
	Date time.Time `json:"date"`
	Open bool      `json:"open"`
}

type errorJSON struct {
	Error string `json:"error"`
}

// TradeOpen is "./open" page.
func TradeOpen(w http.ResponseWriter, req *http.Request) {
	var jsonStr []byte
	data, err := strconv.ParseInt(req.FormValue("q"), 10, 64)
	if err != nil {
		jsonStr, _ = json.Marshal(&errorJSON{Error: "Wrong date format"})
	} else {
		date := time.Unix(data, 0)
		jsonStr, _ = json.Marshal(&tradeJSON{
			Date: date.UTC(),
			Open: tradingdays.IsOpen(date.Year(), date.Month(), date.Day(), date.Location())})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStr)
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
