package main

import (
	"encoding/json"
	"flag"
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

	if csvcachetime.InCache(defaultcachetime) != true {
		tradingdays.DownloadCSV(true)
		csvcachetime.Set()
		log.Println("DownloadCSV.")
	}

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
