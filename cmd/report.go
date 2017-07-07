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
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/toomore/gogrs/cmd/filter"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
)

var (
	wg           sync.WaitGroup
	twseNo       *string
	twseCate     *string
	otcNo        *string
	otcCate      *string
	showcatelist *bool
	showcolor    *bool
	ncpu         *int
	white        = color.New(color.FgWhite, color.Bold).SprintfFunc()
	red          = color.New(color.FgRed, color.Bold).SprintfFunc()
	green        = color.New(color.FgGreen, color.Bold).SprintfFunc()
	yellow       = color.New(color.FgYellow).SprintfFunc()
	yellowBold   = color.New(color.FgYellow, color.Bold).SprintfFunc()
	blue         = color.New(color.FgBlue).SprintfFunc()
	limit        = make(chan struct{}, 1)
)

func prettyprint(stock *twse.Data, check filter.CheckGroup) string {
	var (
		Open        = stock.GetOpenList()[len(stock.GetOpenList())-1]
		Price       = stock.GetPriceList()[len(stock.GetPriceList())-1]
		RangeValue  = stock.GetDailyRangeList()[len(stock.GetDailyRangeList())-1]
		Volume      = stock.GetVolumeList()[len(stock.GetVolumeList())-1] / 1000
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

	return fmt.Sprintf("%s %s %s %s%s %s %s",
		yellow("[%s]", check),
		blue("%s", stock.RawData[stock.Len()-1][0]),
		outputcolor("%s %s", stock.No, stock.Name),
		outputcolor("$%.2f", Price),
		outputcolor("(%.2f)", RangeValue),
		outputcolor("%.2f%%", RangeValue/Open*100),
		outputcolor("%d", Volume),
	)
}

func gocheck(check filter.CheckGroup, stock *twse.Data) {
	defer wg.Done()
	if check.CheckFunc(stock) {
		fmt.Println(prettyprint(stock, check))
	}
	<-limit
}

func run() int {
	var (
		datalist     []*twse.Data
		l            *twse.Lists
		o            *twse.OTCLists
		otccatelist  []string
		otcdelta     int
		otclist      []string
		twsecatelist []string
		twselist     []string
	)

	color.NoColor = !*showcolor

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

	if *twseCate != "" {
		l = twse.NewLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, v := range strings.Split(*twseCate, ",") {
			for _, s := range l.GetCategoryList(v) {
				twsecatelist = append(twsecatelist, s.No)
			}
		}
	}

	if *otcCate != "" {
		o = twse.NewOTCLists(tradingdays.FindRecentlyOpened(time.Now()))
		for _, v := range strings.Split(*otcCate, ",") {
			for _, s := range o.GetCategoryList(v) {
				otccatelist = append(otccatelist, s.No)
			}
		}
	}

	if *twseNo != "" {
		twselist = strings.Split(*twseNo, ",")
	}
	if *otcNo != "" {
		otclist = strings.Split(*otcNo, ",")
	}
	datalist = make([]*twse.Data, len(twselist)+len(twsecatelist)+len(otclist)+len(otccatelist))

	for i, no := range append(twselist, twsecatelist...) {
		datalist[i] = twse.NewTWSE(no, tradingdays.FindRecentlyOpened(time.Now()))
	}
	otcdelta = len(twselist) + len(twsecatelist)
	for i, no := range append(otclist, otccatelist...) {
		datalist[i+otcdelta] = twse.NewOTC(no, tradingdays.FindRecentlyOpened(time.Now()))
	}

	if len(datalist) > 0 {
		for _, check := range filter.AllList {
			fmt.Println(yellowBold("----- %v -----", check))
			wg.Add(len(datalist))
			for _, stock := range datalist {
				limit <- struct{}{}
				go gocheck(check, stock)
			}
			wg.Wait()
		}
	}
	return len(datalist)
}

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "daily report",
	Long:  `show daily report`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tradingdays.DownloadCSV(true)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if run() == 0 {
			cmd.Help()
		}
	},
}

func init() {
	ncpu = reportCmd.Flags().IntP("ncpu", "n", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
	otcCate = reportCmd.Flags().StringP("otccate", "e", "", "上櫃股票類別，可使用 ',' 分隔多組代碼，例：02,14")
	otcNo = reportCmd.Flags().StringP("otc", "o", "", "上櫃股票代碼，可使用 ',' 分隔多組代碼，例：4406,8446")
	showcatelist = reportCmd.Flags().BoolP("catelist", "l", false, "顯示上市/上櫃分類表")
	showcolor = reportCmd.Flags().BoolP("color", "", true, "色彩化")
	twseCate = reportCmd.Flags().StringP("twsecate", "c", "", "上市股票類別，可使用 ',' 分隔多組代碼，例：11,15")
	twseNo = reportCmd.Flags().StringP("twse", "t", "", "上市股票代碼，可使用 ',' 分隔多組代碼，例：2618,2329")

	RootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
