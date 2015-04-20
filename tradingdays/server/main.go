package main

import (
	"log"
	"net/http"
)

func Log(req *http.Request) {
	log.Println(req.URL, req.UserAgent(), req.Form)
}

func Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
	Log(req)
}

func TradeOpen(w http.ResponseWriter, req *http.Request) {
	log.Println(req.FormValue("q"))
	w.Write([]byte(req.FormValue("q")))
	Log(req)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/open", TradeOpen)
	log.Fatal(http.ListenAndServe(":59123", nil))
}
