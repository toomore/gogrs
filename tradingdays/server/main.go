package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func Log(req *http.Request) {
	log.Println(req.URL, req.UserAgent(), req.Form)
}

func Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
	Log(req)
}

func TradeOpen(w http.ResponseWriter, req *http.Request) {
	data, err := strconv.ParseInt(req.FormValue("q"), 10, 64)
	if err != nil {
		w.Write([]byte("Wrong data format."))
	} else {
		date := time.Unix(data, 0)
		w.Write([]byte(date.String()))
	}
	Log(req)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/open", TradeOpen)
	log.Fatal(http.ListenAndServe(":59123", nil))
}
