// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/toomore/gogrs/realtime"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

var (
	cacheTime time.Time
	count     *bool
	cyan      = color.New(color.FgCyan).SprintfFunc()
	index     *bool
	nonstop   *int64
	pt        *bool
)

// realtimeCmd represents the realtime command
var realtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "realtime info",
	Long:  `盤中即時資訊`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tradingdays.DownloadCSV(false)
		TaipeiNow()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("realtime called")
	},
}

// TaipeiNow show Taipei Now time.
func TaipeiNow() time.Time {
	if cacheTime.Year() == 1 {
		d := time.Now().UTC()
		days := d.Day()
		var index int
		for {
			if tradingdays.IsOpen(d.Year(), d.Month(), days) {
				if index == 0 {
					if tradingdays.NewTimePeriod(time.Date(d.Year(), d.Month(), days, d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())).AtBefore() {
						days--
						for {
							if tradingdays.IsOpen(d.Year(), d.Month(), days) {
								break
							}
							days--
						}
					}
				}
				break
			}
			days--
			index++
		}
		cacheTime = time.Date(d.Year(), d.Month(), days, 0, 0, 0, 0, utils.TaipeiTimeZone)
	}
	return cacheTime
}

func rtprettyprint(data realtime.Data) string {
	var (
		RangeValue  = data.Price - data.Open
		outputcolor func(string, ...interface{}) string
	)
	switch {
	case RangeValue > 0:
		outputcolor = red
	case RangeValue < 0:
		outputcolor = green
	default:
		outputcolor = white
	}
	return fmt.Sprintf("%s %s %s",
		yellowBold("%s(%s)", data.Info.Name, data.Info.No),
		outputcolor("$%.2f(%.2f) %.2f%% %.0f/%.0f",
			data.Price, data.Price-data.Open, RangeValue/data.Open*100, data.Volume, data.VolumeAcc),
		cyan("[%s] [%s %s]",
			data.TradeTime, data.SysInfo["sysDate"], data.SysInfo["sysTime"]),
	)
}

func fetch(r *realtime.StockRealTime) {
	limit <- struct{}{}
	runtime.Gosched()
	defer rtwg.Done()
	data, _ := r.Get()
	if data.TradeTime.IsZero() {
		log.Println("No data")
	} else {
		log.Println(rtprettyprint(data))
	}
	if *count {
		counter++
		if data.Price-data.Open > 0 {
			up++
		} else if data.Price-data.Open < 0 {
			down++
		}
	}
	<-limit
}

var (
	counter   int
	down      int
	startTime time.Time
	up        int
	rtwg      sync.WaitGroup
)

func main() {
	color.NoColor = !*showcolor

	runtime.GOMAXPROCS(*ncpu)

	queue := make(chan *realtime.StockRealTime, *ncpu)
	defer close(queue)

	limit = make(chan struct{}, 1)

Start:
	if *pt {
		startTime = time.Now()
	}

	if *showcatelist {
		var (
			cateTitle    = []string{"The same with TWSE/OTC", "OnlyTWSE", "OnlyOTC"}
			categoryList = twse.NewCategoryList()
			index        int
		)

		for i, cate := range []map[string]string{
			categoryList.Same(), categoryList.OnlyTWSE(), categoryList.OnlyOTC(),
		} {
			index = 1
			fmt.Println(white("---------- %s ----------", cateTitle[i]))
			for cateNo, cateName := range cate {
				fmt.Printf("%s\t", fmt.Sprintf("%s%s", green("%s", cateName), yellowBold("(%s)", cateNo)))
				if index%3 == 0 {
					fmt.Println("")
				}
				index++
			}
			fmt.Println("")
		}
	}

	if *twseNo != "" {
		for _, no := range strings.Split(*twseNo, ",") {
			rtwg.Add(1)
			go func(no string) {
				runtime.Gosched()
				queue <- realtime.NewTWSE(no, TaipeiNow())
			}(no)
		}
	}

	if *twseCate != "" {
		l := twse.NewLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, no := range strings.Split(*twseCate, ",") {
			for _, s := range l.GetCategoryList(no) {
				rtwg.Add(1)
				go func(s twse.StockInfo) {
					runtime.Gosched()
					queue <- realtime.NewTWSE(s.No, TaipeiNow())
				}(s)
			}
		}
	}

	if *otcNo != "" {
		for _, no := range strings.Split(*otcNo, ",") {
			rtwg.Add(1)
			go func(no string) {
				runtime.Gosched()
				queue <- realtime.NewOTC(no, TaipeiNow())
			}(no)
		}
	}

	if *otcCate != "" {
		o := twse.NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, no := range strings.Split(*otcCate, ",") {
			for _, s := range o.GetCategoryList(no) {
				rtwg.Add(1)
				go func(s twse.StockInfo) {
					runtime.Gosched()
					queue <- realtime.NewOTC(s.No, TaipeiNow())
				}(s)
			}
		}
	}

	if *index {
		rtwg.Add(3)
		for _, r := range []*realtime.StockRealTime{realtime.NewWeight(TaipeiNow()), realtime.NewOTCI(TaipeiNow()), realtime.NewFRMSA(TaipeiNow())} {
			go func(r *realtime.StockRealTime) {
				runtime.Gosched()
				queue <- r
			}(r)
		}
	}
	go func() {
		for r := range queue {
			go fetch(r)
		}
	}()
	rtwg.Wait()
	if *count {
		fmt.Printf("All: %d, Up: %d, Down: %d, Same: %d\n",
			counter, up, down, counter-up-down)
	}
	if *pt {
		defer fmt.Println(time.Now().Sub(startTime))
	}
	if *nonstop > 0 {
		time.Sleep(time.Duration(*nonstop) * time.Second)
		counter, up, down = 0, 0, 0
		goto Start
	}
}

func init() {
	count = realtimeCmd.Flags().BoolP("count", "", true, "計算此次查詢的漲跌家數")
	index = realtimeCmd.Flags().BoolP("index", "i", false, "顯示大盤、上櫃、寶島指數（default: false）")
	ncpu = realtimeCmd.Flags().IntP("ncpu", "n", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
	nonstop = realtimeCmd.Flags().Int64P("nonstop", "", 0, "自動重複，單位秒數")
	otcCate = realtimeCmd.Flags().StringP("otccate", "e", "", "上櫃股票類別，可使用 ',' 分隔多組代碼，例：02,14")
	otcNo = realtimeCmd.Flags().StringP("otc", "o", "", "上櫃股票代碼，可使用 ',' 分隔多組代碼，例：4406,8446")
	pt = realtimeCmd.Flags().BoolP("pt", "", false, "計算花費時間")
	showcatelist = realtimeCmd.Flags().BoolP("catelist", "l", false, "顯示上市/上櫃分類表")
	showcolor = realtimeCmd.Flags().BoolP("color", "", true, "色彩化")
	twseCate = realtimeCmd.Flags().StringP("twsecate", "c", "", "上市股票類別，可使用 ',' 分隔多組代碼，例：11,15")
	twseNo = realtimeCmd.Flags().StringP("twse", "t", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")

	RootCmd.AddCommand(realtimeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// realtimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// realtimeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
