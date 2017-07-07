// Copyright © 2017 Toomore Chiang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
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
	w.Write([]byte("[<a href=\"https://godoc.org/github.com/toomore/gogrs/tradingdays\">Docs</a>] [<a href=\"https://github.com/toomore/gogrs/blob/master/cmd/tradingdays_server/main.go\">github</a>]<br>"))

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
		if csvcachetime.InCache(*defaultttl) != true {
			tradingdays.DownloadCSV(true)
			csvcachetime.Set()
			log.Println("DownloadCSV.")
		}

		date := time.Unix(data, 0)
		jsonStr, _ = json.Marshal(&tradeJSON{
			Date: date.UTC(),
			Open: tradingdays.IsOpen(date.Year(), date.Month(), date.Day())})
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

var (
	csvcachetime cachetime
	httpport     *string
	defaultttl   *int64
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run tradingdays server",
	Long:  `run tradingdays server`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tradingdays.DownloadCSV(true)
		log.Println("Init DownloadCSV.")
		csvcachetime.Set()

	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("http:", *httpport, "CSVCacheTime:", *defaultttl)
		http.HandleFunc("/", Home)
		http.HandleFunc("/open", TradeOpen)
		log.Fatal(http.ListenAndServe(*httpport, nil))
	},
}

func init() {
	httpport = serverCmd.Flags().StringP("port", "p", ":59123", "HTTP Port")
	defaultttl = serverCmd.Flags().Int64P("ttl", "t", 21600, "Cache time")

	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
