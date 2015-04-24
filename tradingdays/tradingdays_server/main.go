// 提供簡單的日期查詢 API Server.
//
/*
主要支援動態更新 CSV 檔案讀取，解決非預定開休市狀況（如：颱風假）

Install:

	go install -v github.com/toomore/gogrs/tradingdays/tradingdays_server

Usage:

	tradingdays_server [flags]

The flags are:

	-http
		HTTP service address (default ':59123')
	-csvcachetime
		CSV cache time.(default: 21600)

URL Path:

	/open?q={timestamp}

回傳 JSON 格式

	{
		"date": "2015-04-24T15:14:52Z",
		"open": true
	}

範例：

	http://gogrs-trd.toomore.net/open?q=1429888492

*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/toomore/gogrs/tradingdays"
)

// Log is show viwer log.
func Log(req *http.Request) {
	var userIP string
	if userIP = req.Header.Get("X-FORWARDED-FOR"); userIP == "" {
		userIP = req.RemoteAddr
	}
	log.Println(req.URL, userIP, req.UserAgent(), req.Form, req.Referer())
}

// Home is home page.
func Home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("查詢台灣股市是否開市<br>"))
	w.Write([]byte(fmt.Sprintf("<a href=\"/open?q=%d\">範例</a><br>", time.Now().Unix())))
	w.Write([]byte("[<a href=\"https://godoc.org/github.com/toomore/gogrs/tradingdays\">Docs</a>] [<a href=\"https://github.com/toomore/gogrs/blob/master/tradingdays/tradingdays_server/main.go\">github</a>]<br>"))

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

	if data, err := strconv.ParseInt(req.FormValue("q"), 10, 64); err != nil {
		jsonStr, _ = json.Marshal(&errorJSON{Error: "Wrong date format"})
	} else {
		if csvcachetime.InCache(defaultcachetime) != true {
			tradingdays.DownloadCSV(true)
			csvcachetime.Set()
			log.Println("DownloadCSV.")
		}

		date := time.Unix(data, 0)
		jsonStr, _ = json.Marshal(&tradeJSON{
			Date: date.UTC(),
			Open: tradingdays.IsOpen(date.Year(), date.Month(), date.Day(), date.Location())})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStr)
	Log(req)
}

type cachetime struct {
	timestamp int64
}

func (c *cachetime) Set() {
	c.timestamp = time.Now().Unix()
}

func (c *cachetime) InCache(seconds int64) bool {
	result := (time.Now().Unix() - c.timestamp) <= seconds
	if result != true {
		c.timestamp = 0
	}
	return result
}

var csvcachetime cachetime

var httpport string
var defaultcachetime int64

func init() {
	tradingdays.DownloadCSV(true)
	log.Println("Init DownloadCSV.")
	csvcachetime.Set()
	flag.StringVar(&httpport, "http", ":59123", "HTTP service address (e.g., ':59123')")
	flag.Int64Var(&defaultcachetime, "csvcachetime", 21600, "CSV cache time.")
}

func main() {
	flag.Parse()
	log.Println("http:", httpport, "csvcachetime:", defaultcachetime)
	http.HandleFunc("/", Home)
	http.HandleFunc("/open", TradeOpen)
	log.Fatal(http.ListenAndServe(httpport, nil))
}
