package main

import (
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request) {
	log.Println(req.FormValue("q"))
	w.Write([]byte("Hello World"))
}

func main() {
	http.HandleFunc("/", Home)
	log.Fatal(http.ListenAndServe(":59123", nil))
}
